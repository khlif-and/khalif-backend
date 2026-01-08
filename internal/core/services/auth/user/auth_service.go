package user

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"
	"khalif-backend/internal/infrastructure/email"
	"khalif-backend/internal/platform/config"
	"khalif-backend/internal/platform/logger"
	"khalif-backend/pkg/messages"
	"khalif-backend/pkg/utils"

	"go.uber.org/zap"
)

type authService struct {
	userRepo     ports.UserRepository
	authRepo     ports.UserAuthRepository
	emailService email.EmailService
	cfg          *config.Config
}

func NewAuthService(userRepo ports.UserRepository, authRepo ports.UserAuthRepository, emailService email.EmailService, cfg *config.Config) ports.UserAuthService {
	return &authService{
		userRepo:     userRepo,
		authRepo:     authRepo,
		emailService: emailService,
		cfg:          cfg,
	}
}

// generateOTP generates a 6-digit OTP code
func generateOTP() (string, error) {
	max := big.NewInt(999999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// Register creates a new unactivated user and sends OTP email
func (s *authService) Register(req *domain.RegisterRequest) error {
	existing, _ := s.userRepo.FindByEmail(req.Email)
	if existing != nil {
		return errors.New(messages.ErrEmailAlreadyExists)
	}
	existingUser, _ := s.userRepo.FindByUsername(req.Username)
	if existingUser != nil {
		return errors.New(messages.ErrUsernameExists)
	}

	hashedPwd, err := utils.HashPassword(req.Password)
	if err != nil {
		return errors.New(messages.ErrInternalServer)
	}

	profilePic := req.ProfilePicture
	if profilePic == "" {
		initials := utils.GenerateInitials(req.Username)
		ambientColor := utils.GenerateAmbientColor(req.Username)
		colorHex := ambientColor[1:]
		profilePic = fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=%s&color=fff", initials, colorHex)
	}

	newUser := &domain.User{
		Username:       req.Username,
		Email:          req.Email,
		Phone:          req.Phone,
		PasswordHash:   hashedPwd,
		ProfilePicture: profilePic,
		IsActivated:    false, // User must verify OTP
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return err
	}

	// Generate and send OTP
	otpCode, err := generateOTP()
	if err != nil {
		logger.Log.Error("Failed to generate OTP", zap.Error(err))
		return errors.New(messages.ErrInternalServer)
	}

	// Store OTP with 10 minute expiry
	expiresAt := time.Now().Add(10 * time.Minute)
	if err := s.authRepo.StoreOTP(newUser.ID, otpCode, expiresAt); err != nil {
		logger.Log.Error("Failed to store OTP", zap.Error(err))
		return errors.New(messages.ErrInternalServer)
	}

	// Send OTP email via Brevo
	if err := s.emailService.SendOTP(newUser.Email, newUser.Username, otpCode); err != nil {
		logger.Log.Error("Failed to send OTP email", zap.Error(err))
		// Don't return error - user can request resend
		logger.Log.Warn("OTP email failed but registration continued", zap.String("email", newUser.Email))
	}

	logger.Log.Info("User registered, OTP sent", zap.String("email", newUser.Email), zap.String("otp", otpCode))
	return nil
}

// VerifyOTP verifies the OTP code and activates the user account
func (s *authService) VerifyOTP(req *domain.VerifyOTPRequest, userAgent, ipAddress string) (*domain.LoginResponse, error) {
	// Validate OTP
	user, err := s.authRepo.ValidateOTP(req.Email, req.OTP)
	if err != nil {
		logger.Log.Warn("Invalid OTP attempt", zap.String("email", req.Email), zap.Error(err))
		return nil, errors.New(messages.ErrInvalidOTP)
	}

	// Mark OTP as used
	if err := s.authRepo.MarkOTPUsed(user.ID); err != nil {
		logger.Log.Error("Failed to mark OTP as used", zap.Error(err))
	}

	// Activate user
	if err := s.userRepo.ActivateUser(user.ID); err != nil {
		logger.Log.Error("Failed to activate user", zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	// Update user object
	user.IsActivated = true

	// Generate tokens (auto-login)
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, "user", s.cfg.JWTSecret, int(s.cfg.JWTExpHours))
	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}

	refreshTokenString, refreshTokenHash, err := utils.GenerateRefreshToken(user.ID, s.cfg.RefreshTokenSecret, int(s.cfg.RefreshTokenExpDays))
	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}

	err = s.authRepo.StoreRefreshToken(user.ID, refreshTokenHash,
		time.Now().Add(time.Hour*24*time.Duration(s.cfg.RefreshTokenExpDays)), userAgent, ipAddress)

	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}

	// Send welcome email
	go func() {
		if err := s.emailService.SendWelcome(user.Email, user.Username); err != nil {
			logger.Log.Error("Failed to send welcome email", zap.String("email", user.Email), zap.Error(err))
		}
	}()

	logger.Log.Info("User verified and activated", zap.String("email", user.Email))

	return &domain.LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshTokenString,
		User:         user,
	}, nil
}

// ResendOTP generates and sends a new OTP to the user
func (s *authService) ResendOTP(email string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil || user == nil {
		// Don't reveal if email exists
		return nil
	}

	if user.IsActivated {
		return errors.New(messages.ErrAccountAlreadyActivated)
	}

	// Invalidate old OTPs
	if err := s.authRepo.InvalidateOldOTPs(user.ID); err != nil {
		logger.Log.Error("Failed to invalidate old OTPs", zap.Error(err))
	}

	// Generate new OTP
	otpCode, err := generateOTP()
	if err != nil {
		logger.Log.Error("Failed to generate OTP", zap.Error(err))
		return errors.New(messages.ErrInternalServer)
	}

	// Store new OTP
	expiresAt := time.Now().Add(10 * time.Minute)
	if err := s.authRepo.StoreOTP(user.ID, otpCode, expiresAt); err != nil {
		logger.Log.Error("Failed to store OTP", zap.Error(err))
		return errors.New(messages.ErrInternalServer)
	}

	// Send OTP email
	if err := s.emailService.SendOTP(user.Email, user.Username, otpCode); err != nil {
		logger.Log.Error("Failed to send OTP email", zap.Error(err))
		return errors.New(messages.ErrInternalServer)
	}

	logger.Log.Info("OTP resent", zap.String("email", user.Email), zap.String("otp", otpCode))
	return nil
}

// Login authenticates a user and returns tokens
func (s *authService) Login(req *domain.LoginRequest, userAgent, ipAddress string) (*domain.LoginResponse, error) {
	isLocked, msg, err := s.authRepo.CheckLockStatus(req.Email)
	if err != nil {
		logger.Log.Error("DB error checking lock status", zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}
	if isLocked {
		return nil, errors.New(msg)
	}

	user, err := s.userRepo.FindByEmail(req.Email)
	if user == nil || err != nil {
		s.authRepo.RecordLoginFailure(req.Email)
		return nil, errors.New(messages.ErrInvalidCredentials)
	}

	// Check if account is activated
	if !user.IsActivated {
		return nil, errors.New(messages.ErrAccountNotActivated)
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		s.authRepo.RecordLoginFailure(req.Email)
		return nil, errors.New(messages.ErrInvalidCredentials)
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, "user", s.cfg.JWTSecret, int(s.cfg.JWTExpHours))
	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}

	refreshTokenString, refreshTokenHash, err := utils.GenerateRefreshToken(user.ID, s.cfg.RefreshTokenSecret, int(s.cfg.RefreshTokenExpDays))
	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}

	err = s.authRepo.StoreRefreshToken(user.ID, refreshTokenHash,
		time.Now().Add(time.Hour*24*time.Duration(s.cfg.RefreshTokenExpDays)), userAgent, ipAddress)

	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}

	return &domain.LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshTokenString,
		User:         user,
	}, nil
}

func (s *authService) RefreshSession(refreshToken, userAgent, ipAddress string) (*domain.LoginResponse, error) {
	hasher := sha256.New()
	hasher.Write([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hasher.Sum(nil))

	userID, err := s.authRepo.ValidateRefreshToken(tokenHash)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return nil, errors.New(messages.ErrUserNotFound)
	}

	err = s.authRepo.RevokeRefreshToken(tokenHash)
	if err != nil {
		return nil, err
	}

	newAccessToken, err := utils.GenerateAccessToken(user.ID, user.Username, "user", s.cfg.JWTSecret, int(s.cfg.JWTExpHours))
	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}

	newRefreshToken, newRefreshHash, err := utils.GenerateRefreshToken(user.ID, s.cfg.RefreshTokenSecret, int(s.cfg.RefreshTokenExpDays))
	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}

	err = s.authRepo.StoreRefreshToken(user.ID, newRefreshHash,
		time.Now().Add(time.Hour*24*time.Duration(s.cfg.RefreshTokenExpDays)), userAgent, ipAddress)

	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}

	return &domain.LoginResponse{
		Token:        newAccessToken,
		RefreshToken: newRefreshToken,
		User:         user,
	}, nil
}

func (s *authService) GetMe(userID uint) (*domain.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *authService) Logout(userID uint) error {
	return s.authRepo.RevokeAllUserTokens(userID)
}

// ForgotPassword sends a password reset token to the user's email
func (s *authService) ForgotPassword(email string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil || user == nil {
		// Don't reveal if email exists - return success anyway
		return nil
	}

	// Invalidate old tokens
	if err := s.authRepo.InvalidateOldPasswordResetTokens(user.ID); err != nil {
		logger.Log.Error("Failed to invalidate old reset tokens", zap.Error(err))
	}

	// Generate reset token (32 bytes = 64 hex chars)
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		logger.Log.Error("Failed to generate reset token", zap.Error(err))
		return errors.New(messages.ErrInternalServer)
	}
	resetToken := hex.EncodeToString(tokenBytes)

	// Hash the token for storage
	hasher := sha256.New()
	hasher.Write([]byte(resetToken))
	tokenHash := hex.EncodeToString(hasher.Sum(nil))

	// Store token with 30 minute expiry
	expiresAt := time.Now().Add(30 * time.Minute)
	if err := s.authRepo.StorePasswordResetToken(user.ID, tokenHash, expiresAt); err != nil {
		logger.Log.Error("Failed to store reset token", zap.Error(err))
		return errors.New(messages.ErrInternalServer)
	}

	// Send email with reset token
	if err := s.emailService.SendPasswordReset(user.Email, user.Username, resetToken); err != nil {
		logger.Log.Error("Failed to send password reset email", zap.Error(err))
		return errors.New(messages.ErrInternalServer)
	}

	logger.Log.Info("Password reset email sent", zap.String("email", user.Email))
	return nil
}

// ResetPassword resets the user's password using the token
func (s *authService) ResetPassword(token, newPassword string) error {
	// Hash the token to compare with stored hash
	hasher := sha256.New()
	hasher.Write([]byte(token))
	tokenHash := hex.EncodeToString(hasher.Sum(nil))

	// Validate token
	userID, err := s.authRepo.ValidatePasswordResetToken(tokenHash)
	if err != nil {
		logger.Log.Warn("Invalid password reset attempt", zap.Error(err))
		return errors.New(messages.ErrInvalidToken)
	}

	// Hash new password
	hashedPwd, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New(messages.ErrInternalServer)
	}

	// Get user and update password
	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return errors.New(messages.ErrUserNotFound)
	}

	user.PasswordHash = hashedPwd
	if err := s.userRepo.Update(user); err != nil {
		logger.Log.Error("Failed to update password", zap.Error(err))
		return errors.New(messages.ErrInternalServer)
	}

	// Mark token as used
	if err := s.authRepo.MarkPasswordResetTokenUsed(tokenHash); err != nil {
		logger.Log.Error("Failed to mark token as used", zap.Error(err))
	}

	// Revoke all existing sessions for security
	if err := s.authRepo.RevokeAllUserTokens(userID); err != nil {
		logger.Log.Error("Failed to revoke user tokens", zap.Error(err))
	}

	logger.Log.Info("Password reset successful", zap.Uint("userID", userID))
	return nil
}
