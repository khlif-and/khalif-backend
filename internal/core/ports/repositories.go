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

type AudioRepository interface {
	Create(audio *domain.Audio) error
	FindByID(id uint) (*domain.Audio, error)
	FindByUUID(uuid string) (*domain.Audio, error)
	FindAll(page, limit int) ([]domain.Audio, int64, error)
	FindByMoodCategoryID(moodCategoryID uint, page, limit int) ([]domain.Audio, int64, error)
	Update(audio *domain.Audio) error
	Delete(id uint) error
	IncrementListeningCount(id uint) error
}

type MoodCategoryRepository interface {
	Create(mood *domain.MoodCategory) error
	FindByID(id uint) (*domain.MoodCategory, error)
	FindByUUID(uuid string) (*domain.MoodCategory, error)
	FindByName(name string) (*domain.MoodCategory, error)
	FindAll() ([]domain.MoodCategory, int64, error)
	Update(mood *domain.MoodCategory) error
	Delete(id uint) error
}

type UstadzRepository interface {
	Create(ustadz *domain.Ustadz) error
	FindByID(id uint) (*domain.Ustadz, error)
	FindByUUID(uuid string) (*domain.Ustadz, error)
	FindByName(name string) (*domain.Ustadz, error)
	FindAll() ([]domain.Ustadz, int64, error)
	Update(ustadz *domain.Ustadz) error
	Delete(id uint) error
}

type LikeRepository interface {
	Create(like *domain.Like) error
	FindByUUID(uuid string) (*domain.Like, error)
	FindByUserAndAudio(userID, audioID uint) (*domain.Like, error)
	FindByUserID(userID uint) ([]domain.Like, int64, error)
	Delete(id uint) error
	IncrementAudioLikeCount(audioID uint) error
	DecrementAudioLikeCount(audioID uint) error
}
