package domain

import (
	"time"

	"gorm.io/gorm"
)

// User represents the regular user in the system
type User struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	Username            string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email               string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Phone               string         `gorm:"uniqueIndex;size:20;not null" json:"phone"`
	PasswordHash        string         `gorm:"size:255;not null" json:"-"` // Never return password
	GoogleID            *string        `gorm:"uniqueIndex;size:100" json:"google_id,omitempty"`
	ProfilePicture      string         `gorm:"type:text" json:"profile_picture"`
	IsActivated         bool           `gorm:"default:false" json:"is_activated"`
	FailedLoginAttempts int            `gorm:"default:0" json:"-"`
	LockedUntil         *time.Time     `json:"locked_until,omitempty"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}

// UserAuditLog represents the immutable audit log for user actions
type UserAuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	ActionType string    `gorm:"size:50;not null" json:"action_type"` // UPDATE, DELETE
	OldData    []byte    `gorm:"type:jsonb" json:"old_data"`          // Changes stored as JSONB
	NewData    []byte    `gorm:"type:jsonb" json:"new_data"`
	ChangedAt  time.Time `gorm:"autoCreateTime" json:"changed_at"`
}

// UserRefreshToken represents the session token for users
type UserRefreshToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	TokenHash string    `gorm:"uniqueIndex;not null" json:"-"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	IsRevoked bool      `gorm:"default:false" json:"is_revoked"`
	UserAgent string    `gorm:"size:255" json:"user_agent"`
	IPAddress string    `gorm:"size:50" json:"ip_address"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// OTPToken represents OTP for email verification
type OTPToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	OTPCode   string    `gorm:"size:6;not null" json:"-"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	IsUsed    bool      `gorm:"default:false" json:"is_used"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// PasswordResetToken represents token for password reset
type PasswordResetToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	TokenHash string    `gorm:"uniqueIndex;not null" json:"-"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	IsUsed    bool      `gorm:"default:false" json:"is_used"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
