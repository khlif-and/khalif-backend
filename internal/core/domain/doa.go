package domain

import (
	"time"

	"gorm.io/gorm"
)

// Doa represents a supplication entity
type Doa struct {
	ID             uint           `gorm:"primaryKey" json:"-"`
	UUID           string         `gorm:"uniqueIndex;size:36;not null" json:"id"`
	JudulDoa       string         `gorm:"size:255;not null" json:"judul_doa"`
	ArabicDoa      string         `gorm:"type:text" json:"arabic_doa"`
	LatinDoa       string         `gorm:"type:text" json:"latin_doa"`
	TranslateDoa   string         `gorm:"type:text" json:"translate_doa"`
	DescriptionDoa string         `gorm:"type:text" json:"description_doa"` // Tentang Doa
	AudioDoa       string         `gorm:"text" json:"audio_doa"`
	CategoryDoa    string         `gorm:"size:100;index" json:"category_doa"`
	SourceLink     string         `gorm:"size:255" json:"source_link"`
	Tags           string         `gorm:"type:text" json:"tags"`
	
	// Engagement
	LikeCount      int `gorm:"default:0" json:"like_count"`
	BookmarkCount  int `gorm:"default:0" json:"bookmark_count"`
	ListeningCount int `gorm:"default:0" json:"listening_count"`

	// Relationship (Optional link to a Hadist)
	HadistID *uint   `gorm:"index" json:"-"`
	Hadist   *Hadist `gorm:"foreignKey:HadistID" json:"hadist,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// User Engagement Models for Doa
type DoaLike struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index:idx_user_doa_like,unique" json:"user_id"`
	DoaID     uint      `gorm:"not null;index:idx_user_doa_like,unique" json:"doa_id"`
	CreatedAt time.Time `json:"created_at"`
}

type DoaBookmark struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index:idx_user_doa_bookmark,unique" json:"user_id"`
	DoaID     uint      `gorm:"not null;index:idx_user_doa_bookmark,unique" json:"doa_id"`
	CreatedAt time.Time `json:"created_at"`
	Doa       Doa       `gorm:"foreignKey:DoaID" json:"doa"`
}

// DTOs

type CreateDoaRequest struct {
	JudulDoa       string `form:"judul_doa" json:"judul_doa"`
	ArabicDoa      string `form:"arabic_doa" json:"arabic_doa"`
	LatinDoa       string `form:"latin_doa" json:"latin_doa"`
	TranslateDoa   string `form:"translate_doa" json:"translate_doa"`
	DescriptionDoa string `form:"description_doa" json:"description_doa"`
	AudioDoa       string `form:"audio_doa" json:"audio_doa"` // URL from upload
	CategoryDoa    string `form:"category_doa" json:"category_doa"`
	SourceLink     string `form:"source_link" json:"source_link"`
	Tags           string `form:"tags" json:"tags"`
	HadistID       string `form:"hadist_id" json:"hadist_id"` // Receive UUID of hadist
}

type UpdateDoaRequest struct {
	JudulDoa       string `form:"judul_doa" json:"judul_doa"`
	ArabicDoa      string `form:"arabic_doa" json:"arabic_doa"`
	LatinDoa       string `form:"latin_doa" json:"latin_doa"`
	TranslateDoa   string `form:"translate_doa" json:"translate_doa"`
	DescriptionDoa string `form:"description_doa" json:"description_doa"`
	AudioDoa       string `form:"audio_doa" json:"audio_doa"`
	CategoryDoa    string `form:"category_doa" json:"category_doa"`
	SourceLink     string `form:"source_link" json:"source_link"`
	Tags           string `form:"tags" json:"tags"`
	HadistID       string `form:"hadist_id" json:"hadist_id"` // Receive UUID of hadist
}

type DoaListResponse struct {
	Doas       []Doa `json:"doas"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
}

// Helper Methods

func (r *CreateDoaRequest) ToEntity(uuid string, hadistID *uint) *Doa {
	return &Doa{
		UUID:           uuid,
		JudulDoa:       r.JudulDoa,
		ArabicDoa:      r.ArabicDoa,
		LatinDoa:       r.LatinDoa,
		TranslateDoa:   r.TranslateDoa,
		DescriptionDoa: r.DescriptionDoa,
		AudioDoa:       r.AudioDoa,
		CategoryDoa:    r.CategoryDoa,
		SourceLink:     r.SourceLink,
		Tags:           r.Tags,
		HadistID:       hadistID,
	}
}

func (r *UpdateDoaRequest) ApplyUpdates(d *Doa, hadistID *uint) {
	if r.JudulDoa != "" {
		d.JudulDoa = r.JudulDoa
	}
	if r.ArabicDoa != "" {
		d.ArabicDoa = r.ArabicDoa
	}
	if r.LatinDoa != "" {
		d.LatinDoa = r.LatinDoa
	}
	if r.TranslateDoa != "" {
		d.TranslateDoa = r.TranslateDoa
	}
	if r.DescriptionDoa != "" {
		d.DescriptionDoa = r.DescriptionDoa
	}
	if r.AudioDoa != "" {
		d.AudioDoa = r.AudioDoa
	}
	if r.CategoryDoa != "" {
		d.CategoryDoa = r.CategoryDoa
	}
	if r.SourceLink != "" {
		d.SourceLink = r.SourceLink
	}
	if r.Tags != "" {
		d.Tags = r.Tags
	}
	
	// Only update HadistID if it's explicitly provided (handled by service logic usually, 
	// but here we accept the pointer passed from service)
	if hadistID != nil {
		d.HadistID = hadistID
	}
}
