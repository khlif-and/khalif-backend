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
	ErrInvalidPhone        = "invalid phone format (must be 7-15 digits, optional + prefix)"
	ErrUsernameRequired    = "username is required"
	ErrEmailRequired       = "email is required"
	ErrPhoneRequired       = "phone is required"
	ErrPasswordRequired    = "password is required"
	ErrPasswordTooShort    = "password must be at least 6 characters"

	// Audio errors
	ErrAudioNotFound       = "audio not found"
	ErrAudioTitleRequired  = "audio title is required"
	ErrAudioFileRequired   = "audio file is required"
	ErrUstadzNameRequired  = "ustadz name is required"
	ErrInvalidAudioFile    = "invalid audio file format"
	ErrInvalidDuration     = "invalid duration format"

	ErrMoodNotFound     = "mood category not found"
	ErrMoodNameRequired = "mood category name is required"
	ErrMoodNameExists   = "mood category name already exists"

	ErrUstadzNotFound   = "ustadz not found"
	ErrUstadzNameExists = "ustadz name already exists"

	ErrLikeNotFound = "like not found"
	ErrAlreadyLiked = "already liked this audio"

	// OTP errors
	ErrAccountNotActivated     = "account not activated, please verify your email"
	ErrInvalidOTP              = "invalid or expired OTP"
	ErrOTPAlreadyUsed          = "OTP has already been used"
	ErrAccountAlreadyActivated = "account is already activated"

	// Google Auth errors
	ErrInvalidGoogleToken     = "invalid google token"
	ErrGoogleEmailNotVerified = "google email not verified"
)

// Success Messages
const (
	MsgRegisterSuccess     = "registered successfully"
	MsgLoginSuccess        = "login successful"
	MsgUpdateSuccess       = "updated successfully"
	MsgRefreshTokenSuccess = "token refreshed successfully"
	MsgLogoutSuccess       = "logout successful"

	// Audio success
	MsgAudioCreated = "audio created successfully"
	MsgAudioUpdated = "audio updated successfully"
	MsgAudioDeleted = "audio deleted successfully"
	MsgListeningCountIncremented = "listening count incremented successfully"

	MsgMoodCreated = "mood category created successfully"
	MsgMoodUpdated = "mood category updated successfully"
	MsgMoodDeleted = "mood category deleted successfully"

	MsgUstadzCreated = "ustadz created successfully"
	MsgUstadzUpdated = "ustadz updated successfully"
	MsgUstadzDeleted = "ustadz deleted successfully"

	MsgLikeCreated = "audio liked successfully"
	MsgLikeDeleted = "audio unliked successfully"

	// OTP success
	MsgOTPSent          = "registration successful, please check your email for OTP"
	MsgOTPResent        = "OTP has been resent to your email"
	MsgAccountActivated = "account activated successfully"

	// Password Reset success
	MsgPasswordResetSent    = "if the email exists, a password reset link has been sent"
	MsgPasswordResetSuccess = "password has been reset successfully"

	// Email Subjects
	EmailSubjectOTP           = "Kode Masuk - Khalif App"
	EmailSubjectPasswordReset = "Reset Password - Khalif App"
	EmailSubjectWelcome       = "Selamat Datang di Khalif App 🎉"

	// Email Template Names
	EmailTemplateOTP           = "otp.html"
	EmailTemplatePasswordReset = "password_reset.html"
	EmailTemplateWelcome       = "welcome.html"

	// Playlist
	MsgPlaylistNotFound     = "playlist not found"
	MsgPlaylistCreated      = "playlist created successfully"
	MsgPlaylistUpdated      = "playlist updated successfully"
	MsgPlaylistDeleted      = "playlist deleted successfully"
	MsgAudioAddedToPlaylist = "audio added to playlist successfully"
	MsgAudioRemovedFromPlaylist = "audio removed from playlist successfully"
)
