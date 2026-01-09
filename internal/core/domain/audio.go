package domain

import (
	"time"

	"gorm.io/gorm"
)

type Audio struct {
	ID                  uint           `gorm:"primaryKey" json:"-"`
	UUID                string         `gorm:"uniqueIndex;size:36;not null" json:"id"`
	Title               string         `gorm:"size:255;not null" json:"title"`
	AudioFile           string         `gorm:"type:text;not null" json:"audio_file"`
	ThumbnailFile       string         `gorm:"type:text" json:"thumbnail_file"`
	ColorThumbnailAudio string         `gorm:"size:20" json:"color_thumbnail_audio"`
	UstadzID            *uint          `gorm:"index" json:"-"`
	Ustadz              *Ustadz        `gorm:"foreignKey:UstadzID" json:"ustadz,omitempty"`
	MoodCategoryID      *uint          `gorm:"index" json:"-"`
	MoodCategory        *MoodCategory  `gorm:"foreignKey:MoodCategoryID" json:"mood_category,omitempty"`
	ListeningCount      int64          `gorm:"default:0" json:"listening_count"`
	LikeCount           int64          `gorm:"default:0" json:"like_count"`
	DurationAudio       int            `gorm:"default:0" json:"duration_audio"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateAudioRequest struct {
	Title               string `json:"title"`
	AudioFile           string `json:"audio_file"`
	ThumbnailFile       string `json:"thumbnail_file"`
	ColorThumbnailAudio string `json:"color_thumbnail_audio"`
	UstadzUUID          string `json:"ustadz_id"`
	DurationAudio       int    `json:"duration_audio"`
	MoodCategoryUUID    string `json:"mood_category_id"`
}

type UpdateAudioRequest struct {
	Title               string `json:"title"`
	AudioFile           string `json:"audio_file"`
	ThumbnailFile       string `json:"thumbnail_file"`
	ColorThumbnailAudio string `json:"color_thumbnail_audio"`
	UstadzUUID          string `json:"ustadz_id"`
	DurationAudio       int    `json:"duration_audio"`
	MoodCategoryUUID    string `json:"mood_category_id"`
}

type AudioResponse struct {
	Audio *Audio `json:"audio"`
}

type AudioListResponse struct {
	Audios     []Audio `json:"audios"`
	Total      int64   `json:"total"`
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
	TotalPages int     `json:"total_pages"`
}

// ListeningHistory tracks which users listened to which audios
type ListeningHistory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	AudioID    uint      `gorm:"not null;index" json:"audio_id"`
	Audio      *Audio    `gorm:"foreignKey:AudioID" json:"audio,omitempty"`
	ListenedAt time.Time `gorm:"autoCreateTime" json:"listened_at"`
}

type ListeningHistoryResponse struct {
	History    []ListeningHistory `json:"history"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}

// RadioResponse contains the seed audio and queue of similar audios
type RadioResponse struct {
	SeedAudio *Audio  `json:"seed_audio"`
	Queue     []Audio `json:"queue"`
	Total     int     `json:"total"`
}
