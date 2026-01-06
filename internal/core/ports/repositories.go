package ports

import (
	"time"
	"khalif-backend/internal/core/domain"
)

type AdminRepository interface {
	Create(admin *domain.Admin) error
	FindByEmail(email string) (*domain.Admin, error)
	FindByID(id uint) (*domain.Admin, error)
	FindByUsername(username string) (*domain.Admin, error)
	Update(admin *domain.Admin) error
	Count() (int64, error)
}

type AuthRepository interface {
	CheckLockStatus(email string) (bool, string, error)
	RecordLoginFailure(email string) error

	StoreRefreshToken(userID uint, tokenHash string, expiresAt time.Time, userAgent, ipAddress string) error
	ValidateRefreshToken(tokenHash string) (uint, error)
	RevokeRefreshToken(tokenHash string) error
	RevokeAllTokens(adminID uint) error
}

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	Update(user *domain.User) error
}

type UserAuthRepository interface {
	CheckLockStatus(email string) (bool, string, error)
	RecordLoginFailure(email string) error

	StoreRefreshToken(userID uint, tokenHash string, expiresAt time.Time, userAgent, ipAddress string) error
	ValidateRefreshToken(tokenHash string) (uint, error)
	RevokeRefreshToken(tokenHash string) error
	RevokeAllUserTokens(userID uint) error
}
