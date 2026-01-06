package admin

import (
	"khalif-backend/internal/core/domain"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockAdminRepo struct {
	mock.Mock
}

func (m *MockAdminRepo) Create(admin *domain.Admin) error {
	args := m.Called(admin)
	return args.Error(0)
}

func (m *MockAdminRepo) FindByEmail(email string) (*domain.Admin, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Admin), args.Error(1)
}

func (m *MockAdminRepo) FindByID(id uint) (*domain.Admin, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Admin), args.Error(1)
}

func (m *MockAdminRepo) FindByUsername(username string) (*domain.Admin, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Admin), args.Error(1)
}

func (m *MockAdminRepo) Update(admin *domain.Admin) error {
	args := m.Called(admin)
	return args.Error(0)
}

func (m *MockAdminRepo) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

type MockAuthRepo struct {
	mock.Mock
}

func (m *MockAuthRepo) CheckLockStatus(email string) (bool, string, error) {
	args := m.Called(email)
	return args.Bool(0), args.String(1), args.Error(2)
}

func (m *MockAuthRepo) RecordLoginFailure(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockAuthRepo) StoreRefreshToken(userID uint, tokenHash string, expiresAt time.Time, userAgent, ipAddress string) error {
	args := m.Called(userID, tokenHash, expiresAt, userAgent, ipAddress)
	return args.Error(0)
}

func (m *MockAuthRepo) ValidateRefreshToken(tokenHash string) (uint, error) {
	args := m.Called(tokenHash)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockAuthRepo) RevokeRefreshToken(tokenHash string) error {
	args := m.Called(tokenHash)
	return args.Error(0)
}

func (m *MockAuthRepo) RevokeAllTokens(adminID uint) error {
	args := m.Called(adminID)
	return args.Error(0)
}
