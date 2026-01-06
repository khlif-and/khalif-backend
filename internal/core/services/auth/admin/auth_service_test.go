package admin

import (
	"errors"
	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/platform/config"
	"khalif-backend/pkg/messages"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register(t *testing.T) {
	mockAdminRepo := new(MockAdminRepo)
	mockAuthRepo := new(MockAuthRepo)
	cfg := &config.Config{
		JWTSecret:           "secret",
		JWTExpHours:         1,
		RefreshTokenSecret:  "refresh_secret",
		RefreshTokenExpDays: 7,
	}
	service := NewAuthService(mockAdminRepo, mockAuthRepo, cfg)

	req := &domain.RegisterRequest{
		Username: "TestUser",
		Email:    "test@example.com",
		Phone:    "08123456789",
		Password: "password123",
	}

	t.Run("Success", func(t *testing.T) {
		mockAdminRepo.On("Count").Return(int64(0), nil).Once()
		mockAdminRepo.On("FindByEmail", req.Email).Return(nil, nil).Once()
		mockAdminRepo.On("FindByUsername", req.Username).Return(nil, nil).Once()
		mockAdminRepo.On("Create", mock.AnythingOfType("*domain.Admin")).Return(nil).Once()
		mockAuthRepo.On("StoreRefreshToken", mock.AnythingOfType("uint"), mock.AnythingOfType("string"), mock.AnythingOfType("time.Time"), "TestAgent", "127.0.0.1").Return(nil).Once()

		resp, err := service.Register(req, "TestAgent", "127.0.0.1")

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Token)
		assert.NotEmpty(t, resp.RefreshToken)
		assert.Contains(t, resp.Admin.ProfilePicture, "ui-avatars.com")
	})

	t.Run("Fail_MaxAdmin", func(t *testing.T) {
		mockAdminRepo.On("Count").Return(int64(2), nil).Once()

		resp, err := service.Register(req, "TestAgent", "127.0.0.1")

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, messages.ErrAdminLimitReached, err.Error())
	})
}

func TestAuthService_Login(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:           "secret",
		JWTExpHours:         1,
		RefreshTokenSecret:  "refresh_secret",
		RefreshTokenExpDays: 7,
	}

	req := &domain.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockAdmin := &domain.Admin{ID: 1, Email: "test@example.com", PasswordHash: string(hash), Username: "TestUser"}

	t.Run("Fail_NotFound", func(t *testing.T) {
		mockAdminRepo := new(MockAdminRepo)
		mockAuthRepo := new(MockAuthRepo)
		service := NewAuthService(mockAdminRepo, mockAuthRepo, cfg)

		mockAuthRepo.On("CheckLockStatus", req.Email).Return(false, "", nil).Once()
		mockAdminRepo.On("FindByEmail", req.Email).Return(nil, nil).Once()
		mockAuthRepo.On("RecordLoginFailure", req.Email).Return(nil).Once()

		resp, err := service.Login(req, "TestAgent", "127.0.0.1")
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Success", func(t *testing.T) {
		mockAdminRepo := new(MockAdminRepo)
		mockAuthRepo := new(MockAuthRepo)
		service := NewAuthService(mockAdminRepo, mockAuthRepo, cfg)

		mockAuthRepo.On("CheckLockStatus", req.Email).Return(false, "", nil).Once()
		mockAdminRepo.On("FindByEmail", req.Email).Return(mockAdmin, nil).Once()
		mockAuthRepo.On("StoreRefreshToken", mockAdmin.ID, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time"), "TestAgent", "127.0.0.1").Return(nil).Once()

		resp, err := service.Login(req, "TestAgent", "127.0.0.1")

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Token)
		assert.NotEmpty(t, resp.RefreshToken)
		assert.Equal(t, mockAdmin.Email, resp.Admin.Email)
	})
}

func TestAuthService_RefreshSession(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:           "secret",
		JWTExpHours:         1,
		RefreshTokenSecret:  "refresh_secret",
		RefreshTokenExpDays: 7,
	}

	mockRefreshToken := "valid_refresh_token"
	mockAdmin := &domain.Admin{ID: 1, Username: "refresh_user"}

	t.Run("Success", func(t *testing.T) {
		mockAdminRepo := new(MockAdminRepo)
		mockAuthRepo := new(MockAuthRepo)
		service := NewAuthService(mockAdminRepo, mockAuthRepo, cfg)

		mockAuthRepo.On("ValidateRefreshToken", mock.AnythingOfType("string")).Return(mockAdmin.ID, nil).Once()
		mockAdminRepo.On("FindByID", mockAdmin.ID).Return(mockAdmin, nil).Once()
		mockAuthRepo.On("RevokeRefreshToken", mock.AnythingOfType("string")).Return(nil).Once()
		mockAuthRepo.On("StoreRefreshToken", mockAdmin.ID, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time"), "TestAgent", "127.0.0.1").Return(nil).Once()

		resp, err := service.RefreshSession(mockRefreshToken, "TestAgent", "127.0.0.1")

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Token)
		assert.NotEmpty(t, resp.RefreshToken)
		assert.NotEqual(t, mockRefreshToken, resp.RefreshToken)
	})

	t.Run("Fail_InvalidToken", func(t *testing.T) {
		mockAdminRepo := new(MockAdminRepo)
		mockAuthRepo := new(MockAuthRepo)
		service := NewAuthService(mockAdminRepo, mockAuthRepo, cfg)

		mockAuthRepo.On("ValidateRefreshToken", mock.AnythingOfType("string")).Return(uint(0), errors.New("invalid token")).Once()

		resp, err := service.RefreshSession(mockRefreshToken, "TestAgent", "127.0.0.1")

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
