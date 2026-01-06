package user

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
	mockUserRepo := new(MockUserRepo)
	mockUserAuthRepo := new(MockUserAuthRepo)
	cfg := &config.Config{
		JWTSecret:           "secret",
		JWTExpHours:         1,
		RefreshTokenSecret:  "refresh_secret",
		RefreshTokenExpDays: 7,
	}
	service := NewAuthService(mockUserRepo, mockUserAuthRepo, cfg)

	req := &domain.RegisterRequest{
		Username: "TestUser",
		Email:    "test@example.com",
		Phone:    "08123456789",
		Password: "password123",
	}

	t.Run("Success", func(t *testing.T) {
		mockUserRepo.On("FindByEmail", req.Email).Return(nil, nil).Once()
		mockUserRepo.On("FindByUsername", req.Username).Return(nil, nil).Once()
		mockUserRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil).Once()
		mockUserAuthRepo.On("StoreRefreshToken", mock.AnythingOfType("uint"), mock.AnythingOfType("string"), mock.AnythingOfType("time.Time"), "TestAgent", "127.0.0.1").Return(nil).Once()

		resp, err := service.Register(req, "TestAgent", "127.0.0.1")

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Token)
		assert.NotEmpty(t, resp.RefreshToken)
		assert.Contains(t, resp.User.ProfilePicture, "ui-avatars.com")
	})

	t.Run("Fail_EmailExists", func(t *testing.T) {
		existingUser := &domain.User{ID: 1, Email: req.Email}
		mockUserRepo.On("FindByEmail", req.Email).Return(existingUser, nil).Once()

		resp, err := service.Register(req, "TestAgent", "127.0.0.1")

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, messages.ErrEmailAlreadyExists, err.Error())
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
	mockUser := &domain.User{ID: 1, Email: "test@example.com", PasswordHash: string(hash), Username: "TestUser"}

	t.Run("Fail_NotFound", func(t *testing.T) {
		mockUserRepo := new(MockUserRepo)
		mockUserAuthRepo := new(MockUserAuthRepo)
		service := NewAuthService(mockUserRepo, mockUserAuthRepo, cfg)

		mockUserAuthRepo.On("CheckLockStatus", req.Email).Return(false, "", nil).Once()
		mockUserRepo.On("FindByEmail", req.Email).Return(nil, nil).Once()
		mockUserAuthRepo.On("RecordLoginFailure", req.Email).Return(nil).Once()

		resp, err := service.Login(req, "TestAgent", "127.0.0.1")
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Success", func(t *testing.T) {
		mockUserRepo := new(MockUserRepo)
		mockUserAuthRepo := new(MockUserAuthRepo)
		service := NewAuthService(mockUserRepo, mockUserAuthRepo, cfg)

		mockUserAuthRepo.On("CheckLockStatus", req.Email).Return(false, "", nil).Once()
		mockUserRepo.On("FindByEmail", req.Email).Return(mockUser, nil).Once()
		mockUserAuthRepo.On("StoreRefreshToken", mockUser.ID, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time"), "TestAgent", "127.0.0.1").Return(nil).Once()

		resp, err := service.Login(req, "TestAgent", "127.0.0.1")

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Token)
		assert.NotEmpty(t, resp.RefreshToken)
		assert.Equal(t, mockUser.Email, resp.User.Email)
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
	mockUser := &domain.User{ID: 1, Username: "refresh_user"}

	t.Run("Success", func(t *testing.T) {
		mockUserRepo := new(MockUserRepo)
		mockUserAuthRepo := new(MockUserAuthRepo)
		service := NewAuthService(mockUserRepo, mockUserAuthRepo, cfg)

		mockUserAuthRepo.On("ValidateRefreshToken", mock.AnythingOfType("string")).Return(mockUser.ID, nil).Once()
		mockUserRepo.On("FindByID", mockUser.ID).Return(mockUser, nil).Once()
		mockUserAuthRepo.On("RevokeRefreshToken", mock.AnythingOfType("string")).Return(nil).Once()
		mockUserAuthRepo.On("StoreRefreshToken", mockUser.ID, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time"), "TestAgent", "127.0.0.1").Return(nil).Once()

		resp, err := service.RefreshSession(mockRefreshToken, "TestAgent", "127.0.0.1")

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Token)
		assert.NotEmpty(t, resp.RefreshToken)
		assert.NotEqual(t, mockRefreshToken, resp.RefreshToken)
	})

	t.Run("Fail_InvalidToken", func(t *testing.T) {
		mockUserRepo := new(MockUserRepo)
		mockUserAuthRepo := new(MockUserAuthRepo)
		service := NewAuthService(mockUserRepo, mockUserAuthRepo, cfg)

		mockUserAuthRepo.On("ValidateRefreshToken", mock.AnythingOfType("string")).Return(uint(0), errors.New("invalid token")).Once()

		resp, err := service.RefreshSession(mockRefreshToken, "TestAgent", "127.0.0.1")

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestAuthService_Logout(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:           "secret",
		JWTExpHours:         1,
		RefreshTokenSecret:  "refresh_secret",
		RefreshTokenExpDays: 7,
	}

	t.Run("Success", func(t *testing.T) {
		mockUserRepo := new(MockUserRepo)
		mockUserAuthRepo := new(MockUserAuthRepo)
		service := NewAuthService(mockUserRepo, mockUserAuthRepo, cfg)

		mockUserAuthRepo.On("RevokeAllUserTokens", uint(1)).Return(nil).Once()

		err := service.Logout(1)

		assert.NoError(t, err)
	})
}
