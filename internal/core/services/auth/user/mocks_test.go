package user

import (
	"khalif-backend/internal/core/domain"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepo) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) FindByID(id uint) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) FindByUsername(username string) (*domain.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

type MockUserAuthRepo struct {
	mock.Mock
}

func (m *MockUserAuthRepo) CheckLockStatus(email string) (bool, string, error) {
	args := m.Called(email)
	return args.Bool(0), args.String(1), args.Error(2)
}

func (m *MockUserAuthRepo) RecordLoginFailure(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockUserAuthRepo) StoreRefreshToken(userID uint, tokenHash string, expiresAt time.Time, userAgent, ipAddress string) error {
	args := m.Called(userID, tokenHash, expiresAt, userAgent, ipAddress)
	return args.Error(0)
}

func (m *MockUserAuthRepo) ValidateRefreshToken(tokenHash string) (uint, error) {
	args := m.Called(tokenHash)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockUserAuthRepo) RevokeRefreshToken(tokenHash string) error {
	args := m.Called(tokenHash)
	return args.Error(0)
}

func (m *MockUserAuthRepo) RevokeAllUserTokens(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}
