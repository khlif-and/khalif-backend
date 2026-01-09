package domain

import (
	"time"

	"khalif-backend/internal/core/domain/enums"

	"gorm.io/gorm"
)

// Hadist represents a hadist entity
type Hadist struct {
	ID                uint               `gorm:"primaryKey" json:"-"`
	UUID              string             `gorm:"uniqueIndex;size:36;not null" json:"id"`
	NamaHadist        string             `gorm:"size:255;not null" json:"nama_hadist"`
	PerawiHadist      string             `gorm:"size:255" json:"perawi_hadist"`
	NomorHadist       int                `gorm:"default:0" json:"nomor_hadist"`
	ShahihStatus      enums.ShahihStatus `gorm:"size:20;not null;default:'shahih'" json:"shahih_status"`
	KitabHadist       string             `gorm:"size:255" json:"kitab_hadist"`

	// Content
	ArabicHadist      string `gorm:"type:text" json:"arabic_hadist"`
	LatinHadist       string `gorm:"type:text" json:"latin_hadist"`
	TranslateHadist   string `gorm:"type:text" json:"translate_hadist"`
	DescriptionHadist string `gorm:"type:text" json:"description_hadist"`
	AudioHadist       string `gorm:"type:text" json:"audio_hadist"`

	// Category (simple string)
	CategoryHadist string `gorm:"size:100" json:"category_hadist"`

	// Engagement
	LikeCount      int64 `gorm:"default:0" json:"like_count"`
	BookmarkCount  int64 `gorm:"default:0" json:"bookmark_count"`
	ListeningCount int64 `gorm:"default:0" json:"listening_count"`

	// Source
	SourceLink string `gorm:"type:text" json:"source_link"`
	Tags       string `gorm:"size:500" json:"tags"` // Comma separated

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// HadistLike tracks user likes on hadists
type HadistLike struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	HadistID  uint      `gorm:"not null;index" json:"hadist_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// HadistBookmark tracks user bookmarks on hadists
type HadistBookmark struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	HadistID  uint      `gorm:"not null;index" json:"hadist_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// Request/Response DTOs

type CreateHadistRequest struct {
	NamaHadist        string             `json:"nama_hadist" form:"nama_hadist" binding:"required"`
	PerawiHadist      string             `json:"perawi_hadist" form:"perawi_hadist"`
	NomorHadist       int                `json:"nomor_hadist" form:"nomor_hadist"`
	ShahihStatus      enums.ShahihStatus `json:"shahih_status" form:"shahih_status"`
	KitabHadist       string             `json:"kitab_hadist" form:"kitab_hadist"`
	ArabicHadist      string             `json:"arabic_hadist" form:"arabic_hadist"`
	LatinHadist       string             `json:"latin_hadist" form:"latin_hadist"`
	TranslateHadist   string             `json:"translate_hadist" form:"translate_hadist"`
	DescriptionHadist string             `json:"description_hadist" form:"description_hadist"`
	AudioHadist       string             `json:"audio_hadist" form:"audio_hadist"` // From file upload
	CategoryHadist    string             `json:"category_hadist" form:"category_hadist"`
	SourceLink        string             `json:"source_link" form:"source_link"`
	Tags              string             `json:"tags" form:"tags"`
}

type UpdateHadistRequest struct {
	NamaHadist        string              `json:"nama_hadist" form:"nama_hadist"`
	PerawiHadist      string              `json:"perawi_hadist" form:"perawi_hadist"`
	NomorHadist       *int                `json:"nomor_hadist" form:"nomor_hadist"`
	ShahihStatus      *enums.ShahihStatus `json:"shahih_status" form:"shahih_status"`
	KitabHadist       string              `json:"kitab_hadist" form:"kitab_hadist"`
	ArabicHadist      string              `json:"arabic_hadist" form:"arabic_hadist"`
	LatinHadist       string              `json:"latin_hadist" form:"latin_hadist"`
	TranslateHadist   string              `json:"translate_hadist" form:"translate_hadist"`
	DescriptionHadist string              `json:"description_hadist" form:"description_hadist"`
	AudioHadist       string              `json:"audio_hadist" form:"audio_hadist"`
	CategoryHadist    string              `json:"category_hadist" form:"category_hadist"`
	SourceLink        string              `json:"source_link" form:"source_link"`
	Tags              string              `json:"tags" form:"tags"`
}

func (r *CreateHadistRequest) ToEntity(uuid string) *Hadist {
	status := r.ShahihStatus
	if !status.IsValid() {
		status = enums.ShahihStatusShahih
	}

	return &Hadist{
		UUID:              uuid,
		NamaHadist:        r.NamaHadist,
		PerawiHadist:      r.PerawiHadist,
		NomorHadist:       r.NomorHadist,
		ShahihStatus:      status,
		KitabHadist:       r.KitabHadist,
		ArabicHadist:      r.ArabicHadist,
		LatinHadist:       r.LatinHadist,
		TranslateHadist:   r.TranslateHadist,
		DescriptionHadist: r.DescriptionHadist,
		AudioHadist:       r.AudioHadist,
		CategoryHadist:    r.CategoryHadist,
		SourceLink:        r.SourceLink,
		Tags:              r.Tags,
	}
}

func (r *UpdateHadistRequest) ApplyUpdates(h *Hadist) {
	if r.NamaHadist != "" {
		h.NamaHadist = r.NamaHadist
	}
	if r.PerawiHadist != "" {
		h.PerawiHadist = r.PerawiHadist
	}
	if r.NomorHadist != nil {
		h.NomorHadist = *r.NomorHadist
	}
	if r.ShahihStatus != nil && r.ShahihStatus.IsValid() {
		h.ShahihStatus = *r.ShahihStatus
	}
	if r.KitabHadist != "" {
		h.KitabHadist = r.KitabHadist
	}
	if r.ArabicHadist != "" {
		h.ArabicHadist = r.ArabicHadist
	}
	if r.LatinHadist != "" {
		h.LatinHadist = r.LatinHadist
	}
	if r.TranslateHadist != "" {
		h.TranslateHadist = r.TranslateHadist
	}
	if r.DescriptionHadist != "" {
		h.DescriptionHadist = r.DescriptionHadist
	}
	if r.AudioHadist != "" {
		h.AudioHadist = r.AudioHadist
	}
	if r.CategoryHadist != "" {
		h.CategoryHadist = r.CategoryHadist
	}
	if r.SourceLink != "" {
		h.SourceLink = r.SourceLink
	}
	if r.Tags != "" {
		h.Tags = r.Tags
	}
}

type HadistListResponse struct {
	Hadists    []Hadist `json:"hadists"`
	Total      int64    `json:"total"`
	Page       int      `json:"page"`
	Limit      int      `json:"limit"`
	TotalPages int      `json:"total_pages"`
}
