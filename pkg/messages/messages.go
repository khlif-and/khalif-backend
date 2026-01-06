package messages

// Error Messages
const (
	ErrAdminLimitReached   = "admin limit reached (max 2)"
	ErrInvalidCredentials  = "invalid email or password"
	ErrInternalServer      = "internal server error"
	ErrBadRequest          = "bad request"
	ErrEmailAlreadyExists  = "email already exists"
	ErrPhoneAlreadyExists  = "phone number already exists"
	ErrUsernameExists      = "username already exists"
	ErrUnauthorized        = "unauthorized access"
	ErrAccountLocked       = "account is locked due to too many failed attempts"
	ErrUserNotFound        = "user not found"
	ErrInvalidToken        = "invalid or expired token"
	ErrMissingFields       = "missing required fields"
	ErrInvalidEmail        = "invalid email format"
)

// Success Messages
const (
	MsgRegisterSuccess     = "registered successfully"
	MsgLoginSuccess        = "login successful"
	MsgUpdateSuccess       = "updated successfully"
	MsgRefreshTokenSuccess = "token refreshed successfully"
	MsgLogoutSuccess       = "logout successful"
)
