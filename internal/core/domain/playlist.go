package domain

import (
	"time"

	"gorm.io/gorm"
)

// AuthorType represents the type of playlist author
type AuthorType string

const (
	AuthorTypeAdmin AuthorType = "admin"
	AuthorTypeUser  AuthorType = "user"
)

// Playlist represents a collection of audio tracks
type Playlist struct {
	ID             uint           `gorm:"primaryKey" json:"-"`
	UUID           string         `gorm:"uniqueIndex;size:36;not null" json:"id"`
	Title          string         `gorm:"size:255;not null" json:"title"`
	Description    string         `gorm:"type:text" json:"description"`
	AuthorType     AuthorType     `gorm:"size:10;not null" json:"author_type"`
	AuthorID       uint           `gorm:"not null;index" json:"-"`
	AuthorName     string         `gorm:"-" json:"author_name,omitempty"`
	ThumbnailFile  string         `gorm:"type:text" json:"thumbnail_file"`
	LikeCount      int64          `gorm:"default:0" json:"like_count"`
	ListeningCount int64          `gorm:"default:0" json:"listening_count"`
	IsPublic       bool           `gorm:"default:true" json:"is_public"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	PlaylistAudios []PlaylistAudio `gorm:"foreignKey:PlaylistID" json:"-"`

	// Computed fields (not stored in DB)
	TotalAudio        int `gorm:"-" json:"total_audio"`
	TotalPlayingAudio int `gorm:"-" json:"total_playing_audio"` // in seconds
}

// PlaylistAudio represents the junction table between Playlist and Audio
type PlaylistAudio struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PlaylistID uint      `gorm:"not null;index" json:"playlist_id"`
	AudioID    uint      `gorm:"not null;index" json:"audio_id"`
	Audio      *Audio    `gorm:"foreignKey:AudioID" json:"audio,omitempty"`
	Position   int       `gorm:"default:0" json:"position"`
	AddedAt    time.Time `gorm:"autoCreateTime" json:"added_at"`
}

// PlaylistLike represents a user's like on a playlist
type PlaylistLike struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	PlaylistID uint      `gorm:"not null;index" json:"playlist_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// Request DTOs
type CreatePlaylistRequest struct {
	Title         string `json:"title" binding:"required"`
	Description   string `json:"description"`
	ThumbnailFile string `json:"thumbnail_file"`
	IsPublic      bool   `json:"is_public"`
}

type UpdatePlaylistRequest struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	ThumbnailFile string `json:"thumbnail_file"`
	IsPublic      *bool  `json:"is_public"`
}

type AddAudioToPlaylistRequest struct {
	AudioUUID string `json:"audio_id" binding:"required"`
	Position  int    `json:"position"`
}

// Response DTOs
type PlaylistResponse struct {
	Playlist *Playlist `json:"playlist"`
	Audios   []Audio   `json:"audios,omitempty"`
}

type PlaylistListResponse struct {
	Playlists  []Playlist `json:"playlists"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
	TotalPages int        `json:"total_pages"`
}
