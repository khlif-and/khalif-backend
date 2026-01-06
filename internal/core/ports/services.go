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
	Register(req *domain.RegisterRequest, userAgent, ipAddress string) (*domain.LoginResponse, error)
	Login(req *domain.LoginRequest, userAgent, ipAddress string) (*domain.LoginResponse, error)
	RefreshSession(refreshToken, userAgent, ipAddress string) (*domain.LoginResponse, error)
	GetMe(userID uint) (*domain.User, error)
	Logout(userID uint) error
}

type UserService interface {
	UpdateProfile(userID uint, req *domain.UpdateProfileRequest) (*domain.User, error)
}
