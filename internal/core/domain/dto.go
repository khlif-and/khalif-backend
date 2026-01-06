package domain

// DTOs (Data Transfer Objects) for Service Layer

type RegisterRequest struct {
	Username       string `json:"username" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Phone          string `json:"phone" binding:"required"`
	Password       string `json:"password" binding:"required,min=6"`
	ProfilePicture string `json:"profile_picture"` // Optional
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	Admin        *Admin `json:"admin,omitempty"`
	User         *User  `json:"user,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UpdateProfileRequest struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	ProfilePicture string `json:"profile_picture"`
}
