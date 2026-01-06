package user

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
	userRepo ports.UserRepository
	authRepo ports.UserAuthRepository
	cfg      *config.Config
}

func NewAuthService(userRepo ports.UserRepository, authRepo ports.UserAuthRepository, cfg *config.Config) ports.UserAuthService {
	return &authService{
		userRepo: userRepo,
		authRepo: authRepo,
		cfg:      cfg,
	}
}

func (s *authService) Register(req *domain.RegisterRequest, userAgent, ipAddress string) (*domain.LoginResponse, error) {
	existing, _ := s.userRepo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New(messages.ErrEmailAlreadyExists)
	}
	existingUser, _ := s.userRepo.FindByUsername(req.Username)
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

	newUser := &domain.User{
		Username:       req.Username,
		Email:          req.Email,
		Phone:          req.Phone,
		PasswordHash:   hashedPwd,
		ProfilePicture: profilePic,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	accessToken, err := utils.GenerateAccessToken(newUser.ID, newUser.Username, "user", s.cfg.JWTSecret, int(s.cfg.JWTExpHours))
	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}

	refreshTokenString, refreshTokenHash, err := utils.GenerateRefreshToken(newUser.ID, s.cfg.RefreshTokenSecret, int(s.cfg.RefreshTokenExpDays))
	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}

	err = s.authRepo.StoreRefreshToken(newUser.ID, refreshTokenHash, 
		time.Now().Add(time.Hour * 24 * time.Duration(s.cfg.RefreshTokenExpDays)), userAgent, ipAddress)
	
	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}

	return &domain.LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshTokenString,
		User:         newUser,
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

	user, err := s.userRepo.FindByEmail(req.Email)
	if user == nil || err != nil {
		s.authRepo.RecordLoginFailure(req.Email)
		return nil, errors.New(messages.ErrInvalidCredentials)
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
		time.Now().Add(time.Hour * 24 * time.Duration(s.cfg.RefreshTokenExpDays)), userAgent, ipAddress)
	
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
		time.Now().Add(time.Hour * 24 * time.Duration(s.cfg.RefreshTokenExpDays)), userAgent, ipAddress)

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
