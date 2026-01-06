package admin

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"
	"khalif-backend/internal/platform/config"
	"khalif-backend/internal/platform/logger"
	"khalif-backend/pkg/messages"
	"khalif-backend/pkg/utils"

	"go.uber.org/zap"
)

type authService struct {
	adminRepo ports.AdminRepository
	authRepo  ports.AuthRepository
	cfg       *config.Config
}

func NewAuthService(adminRepo ports.AdminRepository, authRepo ports.AuthRepository, cfg *config.Config) ports.AuthService {
	return &authService{
		adminRepo: adminRepo,
		authRepo:  authRepo,
		cfg:       cfg,
	}
}

func (s *authService) Register(req *domain.RegisterRequest, userAgent, ipAddress string) (*domain.LoginResponse, error) {
	count, err := s.adminRepo.Count()
	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}
	if count >= 2 {
		return nil, errors.New(messages.ErrAdminLimitReached)
	}

	existing, _ := s.adminRepo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New(messages.ErrEmailAlreadyExists)
	}
	existingUser, _ := s.adminRepo.FindByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New(messages.ErrUsernameExists)
	}

	hashedPwd, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}

	profilePic := req.ProfilePicture
	if profilePic == "" {
		initials := utils.GenerateInitials(req.Username)
		ambientColor := utils.GenerateAmbientColor(req.Username)
		colorHex := ambientColor[1:]
		profilePic = fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=%s&color=fff", initials, colorHex)
	}

	newAdmin := &domain.Admin{
		Username:       req.Username,
		Email:          req.Email,
		Phone:          req.Phone,
		PasswordHash:   hashedPwd,
		ProfilePicture: profilePic,
	}

	if err := s.adminRepo.Create(newAdmin); err != nil {
		return nil, err
	}

	accessToken, err := utils.GenerateAccessToken(newAdmin.ID, newAdmin.Username, "admin", s.cfg.JWTSecret, int(s.cfg.JWTExpHours))
	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}

	refreshTokenString, refreshTokenHash, err := utils.GenerateRefreshToken(newAdmin.ID, s.cfg.RefreshTokenSecret, int(s.cfg.RefreshTokenExpDays))
	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}

	err = s.authRepo.StoreRefreshToken(newAdmin.ID, refreshTokenHash, 
		time.Now().Add(time.Hour * 24 * time.Duration(s.cfg.RefreshTokenExpDays)), userAgent, ipAddress)
	
	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}

	return &domain.LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshTokenString,
		Admin:        newAdmin,
	}, nil
}

func (s *authService) Login(req *domain.LoginRequest, userAgent, ipAddress string) (*domain.LoginResponse, error) {
	isLocked, msg, err := s.authRepo.CheckLockStatus(req.Email)

	if err != nil {
		logger.Log.Error("DB error checking lock status", zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}
	if isLocked {
		return nil, errors.New(msg)
	}

	admin, err := s.adminRepo.FindByEmail(req.Email)
	if admin == nil || err != nil {
		logger.Log.Info("User not found or DB error in Login", zap.Error(err), zap.String("email", req.Email))
		s.authRepo.RecordLoginFailure(req.Email)
		return nil, errors.New(messages.ErrInvalidCredentials)
	}

	if !utils.CheckPasswordHash(req.Password, admin.PasswordHash) {
		logger.Log.Info("Password mismatch in Login", zap.String("email", req.Email))
		s.authRepo.RecordLoginFailure(req.Email)
		return nil, errors.New(messages.ErrInvalidCredentials)
	}

	accessToken, err := utils.GenerateAccessToken(admin.ID, admin.Username, "admin", s.cfg.JWTSecret, int(s.cfg.JWTExpHours))
	if err != nil {
		logger.Log.Error("Failed to generate access token", zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	refreshTokenString, refreshTokenHash, err := utils.GenerateRefreshToken(admin.ID, s.cfg.RefreshTokenSecret, int(s.cfg.RefreshTokenExpDays))
	if err != nil {
		logger.Log.Error("Failed to generate refresh token", zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	err = s.authRepo.StoreRefreshToken(admin.ID, refreshTokenHash, 
		time.Now().Add(time.Hour * 24 * time.Duration(s.cfg.RefreshTokenExpDays)), userAgent, ipAddress)
	
	if err != nil {
		logger.Log.Error("Failed to store refresh token", zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	return &domain.LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshTokenString,
		Admin:        admin,
	}, nil
}

func (s *authService) RefreshSession(refreshToken, userAgent, ipAddress string) (*domain.LoginResponse, error) {
	hasher := sha256.New()
	hasher.Write([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hasher.Sum(nil))

	adminID, err := s.authRepo.ValidateRefreshToken(tokenHash)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	admin, err := s.adminRepo.FindByID(adminID)
	if err != nil || admin == nil {
		return nil, errors.New(messages.ErrUserNotFound)
	}

	err = s.authRepo.RevokeRefreshToken(tokenHash)
	if err != nil {
		return nil, err
	}

	newAccessToken, err := utils.GenerateAccessToken(admin.ID, admin.Username, "admin", s.cfg.JWTSecret, int(s.cfg.JWTExpHours))
	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}

	newRefreshToken, newRefreshHash, err := utils.GenerateRefreshToken(admin.ID, s.cfg.RefreshTokenSecret, int(s.cfg.RefreshTokenExpDays))
	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}

	err = s.authRepo.StoreRefreshToken(admin.ID, newRefreshHash, 
		time.Now().Add(time.Hour * 24 * time.Duration(s.cfg.RefreshTokenExpDays)), userAgent, ipAddress)

	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}

	return &domain.LoginResponse{
		Token:        newAccessToken,
		RefreshToken: newRefreshToken,
		Admin:        admin,
	}, nil
}

func (s *authService) GetMe(adminID uint) (*domain.Admin, error) {
	return s.adminRepo.FindByID(adminID)
}

func (s *authService) Logout(adminID uint) error {
	return s.authRepo.RevokeAllTokens(adminID)
}
