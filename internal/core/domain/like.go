package domain

import (
	"time"

	"gorm.io/gorm"
)

type Like struct {
	ID        uint           `gorm:"primaryKey" json:"-"`
	UUID      string         `gorm:"uniqueIndex;size:36;not null" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"-"`
	User      *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	AudioID   uint           `gorm:"index;not null" json:"-"`
	Audio     *Audio         `gorm:"foreignKey:AudioID" json:"audio,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type LikeResponse struct {
	Like *Like `json:"like"`
}

type LikeListResponse struct {
	Likes      []Like `json:"likes"`
	Total      int64  `json:"total"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalPages int    `json:"total_pages"`
}
