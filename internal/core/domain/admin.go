package domain

import (
	"time"

	"gorm.io/gorm"
)

// Admin represents the admin user in the system
type Admin struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	Username            string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email               string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Phone               string         `gorm:"uniqueIndex;size:20;not null" json:"phone"`
	PasswordHash        string         `gorm:"size:255;not null" json:"-"` // Never return password
	ProfilePicture      string         `gorm:"type:text" json:"profile_picture"`
	FailedLoginAttempts int            `gorm:"default:0" json:"-"`
	LockedUntil         *time.Time     `json:"locked_until,omitempty"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}

// AdminAuditLog represents the immutable audit log for admin actions
type AdminAuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	AdminID    uint      `gorm:"not null" json:"admin_id"`
	ActionType string    `gorm:"size:50;not null" json:"action_type"` // UPDATE, DELETE
	OldData    []byte    `gorm:"type:jsonb" json:"old_data"`          // Changes stored as JSONB
	NewData    []byte    `gorm:"type:jsonb" json:"new_data"`
	ChangedAt  time.Time `gorm:"autoCreateTime" json:"changed_at"`
}

// RefreshToken represents the session token in database
type RefreshToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AdminID   uint      `gorm:"not null;index" json:"admin_id"`
	TokenHash string    `gorm:"uniqueIndex;not null" json:"-"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	IsRevoked bool      `gorm:"default:false" json:"is_revoked"`
	UserAgent string    `gorm:"size:255" json:"user_agent"`
	IPAddress string    `gorm:"size:50" json:"ip_address"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
