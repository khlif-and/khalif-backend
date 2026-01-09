package ports

import "khalif-backend/internal/core/domain"

type AuthService interface {
	Register(req *domain.RegisterRequest, userAgent, ipAddress string) (*domain.LoginResponse, error)
	Login(req *domain.LoginRequest, userAgent, ipAddress string) (*domain.LoginResponse, error)
	RefreshSession(refreshToken, userAgent, ipAddress string) (*domain.LoginResponse, error)
	GetMe(adminID uint) (*domain.Admin, error)
	Logout(adminID uint) error
}

type AdminService interface {
	UpdateProfile(adminID uint, req *domain.UpdateProfileRequest) (*domain.Admin, error)
}

type UserAuthService interface {
	Register(req *domain.RegisterRequest) error // No auto-login, creates unactivated user
	VerifyOTP(req *domain.VerifyOTPRequest, userAgent, ipAddress string) (*domain.LoginResponse, error)
	ResendOTP(email string) error
	Login(req *domain.LoginRequest, userAgent, ipAddress string) (*domain.LoginResponse, error)
	RefreshSession(refreshToken, userAgent, ipAddress string) (*domain.LoginResponse, error)
	GetMe(userID uint) (*domain.User, error)
	Logout(userID uint) error
	// Password Reset
	ForgotPassword(email string) error
	ResetPassword(token, newPassword string) error
	// Social Login
	LoginWithGoogle(idToken, userAgent, ipAddress string) (*domain.LoginResponse, error)
}

type UserService interface {
	UpdateProfile(userID uint, req *domain.UpdateProfileRequest) (*domain.User, error)
}

type AudioService interface {
	Create(req *domain.CreateAudioRequest) (*domain.Audio, error)
	GetByUUID(uuid string) (*domain.Audio, error)
	GetAll(page, limit int) (*domain.AudioListResponse, error)
	Update(uuid string, req *domain.UpdateAudioRequest) (*domain.Audio, error)
	Delete(uuid string) error
	IncrementListeningCount(uuid string) error
	// Listening History with SP
	RecordListening(userID uint, audioUUID string) (alreadyListened bool, newCount int64, err error)
	GetUserListeningHistory(userID uint, page, limit int) (*domain.ListeningHistoryResponse, error)
	// Radio
	GenerateRadio(seedAudioUUID string, limit int) (*domain.RadioResponse, error)
}

type MoodCategoryService interface {
	Create(req *domain.CreateMoodCategoryRequest) (*domain.MoodCategory, error)
	GetByUUID(uuid string) (*domain.MoodCategory, error)
	GetAll(page, limit int) (*domain.MoodCategoryListResponse, error)
	GetAudiosByMoodUUID(moodUUID string, page, limit int) (*domain.AudioListResponse, error)
	Update(uuid string, req *domain.UpdateMoodCategoryRequest) (*domain.MoodCategory, error)
	Delete(uuid string) error
}

type UstadzService interface {
	Create(req *domain.CreateUstadzRequest) (*domain.Ustadz, error)
	GetByUUID(uuid string) (*domain.Ustadz, error)
	GetAll(page, limit int) (*domain.UstadzListResponse, error)
	Update(uuid string, req *domain.UpdateUstadzRequest) (*domain.Ustadz, error)
	Delete(uuid string) error
}

type LikeService interface {
	LikeAudio(userID uint, audioUUID string) (*domain.Like, error)
	UnlikeAudio(userID uint, audioUUID string) error
	GetUserLikes(userID uint, page, limit int) (*domain.LikeListResponse, error)
	IsLiked(userID uint, audioUUID string) (bool, error)
}

type HadistService interface {
	Create(req *domain.CreateHadistRequest) (*domain.Hadist, error)
	GetByUUID(uuid string) (*domain.Hadist, error)
	GetAll(page, limit int) (*domain.HadistListResponse, error)
	GetByCategory(category string, page, limit int) (*domain.HadistListResponse, error)
	GetByKitab(kitab string, page, limit int) (*domain.HadistListResponse, error)
	GetRandom() (*domain.Hadist, error)
	Update(uuid string, req *domain.UpdateHadistRequest) (*domain.Hadist, error)
	Delete(uuid string) error
	IncrementListeningCount(uuid string) error
	// User engagement
	LikeHadist(userID uint, uuid string) error
	UnlikeHadist(userID uint, uuid string) error
	IsLiked(userID uint, uuid string) (bool, error)
	BookmarkHadist(userID uint, uuid string) error
	UnbookmarkHadist(userID uint, uuid string) error
	IsBookmarked(userID uint, uuid string) (bool, error)
}

type DoaService interface {
	Create(req *domain.CreateDoaRequest) (*domain.Doa, error)
	GetByUUID(uuid string) (*domain.Doa, error)
	GetAll(page, limit int) (*domain.DoaListResponse, error)
	GetByCategory(category string, page, limit int) (*domain.DoaListResponse, error)
	GetByHadist(hadistUUID string, page, limit int) (*domain.DoaListResponse, error)
	GetRandom() (*domain.Doa, error)
	Update(uuid string, req *domain.UpdateDoaRequest) (*domain.Doa, error)
	Delete(uuid string) error
	IncrementListeningCount(uuid string) error
	// User engagement
	LikeDoa(userID uint, uuid string) error
	UnlikeDoa(userID uint, uuid string) error
	IsLiked(userID uint, uuid string) (bool, error)
	BookmarkDoa(userID uint, uuid string) error
	UnbookmarkDoa(userID uint, uuid string) error
	IsBookmarked(userID uint, uuid string) (bool, error)
}

type PrayerTimeService interface {
	GetPrayerTimes(req *domain.PrayerTimesRequest) (*domain.PrayerTimesResponse, error)
	GetDailyPrayerTimes(req *domain.PrayerTimesRequest) (*domain.PrayerTimesResponse, error)
}
