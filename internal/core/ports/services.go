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
}

type MoodCategoryService interface {
	Create(req *domain.CreateMoodCategoryRequest) (*domain.MoodCategory, error)
	GetByUUID(uuid string) (*domain.MoodCategory, error)
	GetAll() (*domain.MoodCategoryListResponse, error)
	GetAudiosByMoodUUID(moodUUID string, page, limit int) (*domain.AudioListResponse, error)
	Update(uuid string, req *domain.UpdateMoodCategoryRequest) (*domain.MoodCategory, error)
	Delete(uuid string) error
}

type UstadzService interface {
	Create(req *domain.CreateUstadzRequest) (*domain.Ustadz, error)
	GetByUUID(uuid string) (*domain.Ustadz, error)
	GetAll() (*domain.UstadzListResponse, error)
	Update(uuid string, req *domain.UpdateUstadzRequest) (*domain.Ustadz, error)
	Delete(uuid string) error
}

type LikeService interface {
	LikeAudio(userID uint, audioUUID string) (*domain.Like, error)
	UnlikeAudio(userID uint, audioUUID string) error
	GetUserLikes(userID uint) (*domain.LikeListResponse, error)
	IsLiked(userID uint, audioUUID string) (bool, error)
}
