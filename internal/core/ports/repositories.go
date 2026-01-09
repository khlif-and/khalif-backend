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
	ActivateUser(userID uint) error
}

type UserAuthRepository interface {
	CheckLockStatus(email string) (bool, string, error)
	RecordLoginFailure(email string) error

	StoreRefreshToken(userID uint, tokenHash string, expiresAt time.Time, userAgent, ipAddress string) error
	ValidateRefreshToken(tokenHash string) (uint, error)
	RevokeRefreshToken(tokenHash string) error
	RevokeAllUserTokens(userID uint) error

	// OTP methods
	StoreOTP(userID uint, otpCode string, expiresAt time.Time) error
	ValidateOTP(email, otpCode string) (*domain.User, error)
	MarkOTPUsed(userID uint) error
	InvalidateOldOTPs(userID uint) error

	// Password Reset methods
	StorePasswordResetToken(userID uint, tokenHash string, expiresAt time.Time) error
	ValidatePasswordResetToken(tokenHash string) (uint, error)
	MarkPasswordResetTokenUsed(tokenHash string) error
	InvalidateOldPasswordResetTokens(userID uint) error
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
	// Listening History with SP
	RecordListening(userID, audioID uint) (alreadyListened bool, newCount int64, err error)
	GetUserListeningHistory(userID uint, page, limit int) ([]domain.ListeningHistory, int64, error)
	// Radio
	GetRadioQueue(seedAudio *domain.Audio, limit int) ([]domain.Audio, error)
}

type MoodCategoryRepository interface {
	Create(mood *domain.MoodCategory) error
	FindByID(id uint) (*domain.MoodCategory, error)
	FindByUUID(uuid string) (*domain.MoodCategory, error)
	FindByName(name string) (*domain.MoodCategory, error)
	FindAll(page, limit int) ([]domain.MoodCategory, int64, error)
	Update(mood *domain.MoodCategory) error
	Delete(id uint) error
}

type UstadzRepository interface {
	Create(ustadz *domain.Ustadz) error
	FindByID(id uint) (*domain.Ustadz, error)
	FindByUUID(uuid string) (*domain.Ustadz, error)
	FindByName(name string) (*domain.Ustadz, error)
	FindAll(page, limit int) ([]domain.Ustadz, int64, error)
	Update(ustadz *domain.Ustadz) error
	Delete(id uint) error
}

type LikeRepository interface {
	Create(like *domain.Like) error
	FindByUUID(uuid string) (*domain.Like, error)
	FindByUserAndAudio(userID, audioID uint) (*domain.Like, error)
	FindByUserID(userID uint, page, limit int) ([]domain.Like, int64, error)
	Delete(id uint) error
	IncrementAudioLikeCount(audioID uint) error
	DecrementAudioLikeCount(audioID uint) error
}

type HadistRepository interface {
	Create(hadist *domain.Hadist) error
	FindByID(id uint) (*domain.Hadist, error)
	FindByUUID(uuid string) (*domain.Hadist, error)
	FindAll(page, limit int) ([]domain.Hadist, int64, error)
	FindByCategory(category string, page, limit int) ([]domain.Hadist, int64, error)
	FindByKitab(kitab string, page, limit int) ([]domain.Hadist, int64, error)
	FindRandom() (*domain.Hadist, error)
	Update(hadist *domain.Hadist) error
	Delete(id uint) error
	IncrementListeningCount(id uint) error
	// Like
	CreateLike(like *domain.HadistLike) error
	FindLikeByUserAndHadist(userID, hadistID uint) (*domain.HadistLike, error)
	DeleteLike(id uint) error
	IncrementLikeCount(hadistID uint) error
	DecrementLikeCount(hadistID uint) error
	// Bookmark
	CreateBookmark(bookmark *domain.HadistBookmark) error
	FindBookmarkByUserAndHadist(userID, hadistID uint) (*domain.HadistBookmark, error)
	DeleteBookmark(id uint) error
	IncrementBookmarkCount(hadistID uint) error
	DecrementBookmarkCount(hadistID uint) error
	GetUserBookmarks(userID uint, page, limit int) ([]domain.HadistBookmark, int64, error)
}

type DoaRepository interface {
	Create(doa *domain.Doa) error
	FindByID(id uint) (*domain.Doa, error)
	FindByUUID(uuid string) (*domain.Doa, error)
	FindAll(page, limit int) ([]domain.Doa, int64, error)
	FindByCategory(category string, page, limit int) ([]domain.Doa, int64, error)
	FindByHadistID(hadistID uint, page, limit int) ([]domain.Doa, int64, error)
	FindRandom() (*domain.Doa, error)
	Update(doa *domain.Doa) error
	Delete(id uint) error
	IncrementListeningCount(id uint) error
	// Like
	CreateLike(like *domain.DoaLike) error
	FindLikeByUserAndDoa(userID, doaID uint) (*domain.DoaLike, error)
	DeleteLike(id uint) error
	IncrementLikeCount(doaID uint) error
	DecrementLikeCount(doaID uint) error
	// Bookmark
	CreateBookmark(bookmark *domain.DoaBookmark) error
	FindBookmarkByUserAndDoa(userID, doaID uint) (*domain.DoaBookmark, error)
	DeleteBookmark(id uint) error
	IncrementBookmarkCount(doaID uint) error
	DecrementBookmarkCount(doaID uint) error
	GetUserBookmarks(userID uint, page, limit int) ([]domain.DoaBookmark, int64, error)
}
