package domain

import (
	"time"

	"gorm.io/gorm"
)

type Ustadz struct {
	ID            uint           `gorm:"primaryKey" json:"-"`
	UUID          string         `gorm:"uniqueIndex;size:36;not null" json:"id"`
	Name          string         `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Description   string         `gorm:"type:text" json:"description"`
	WikipediaLink string         `gorm:"size:500" json:"wikipedia_link"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateUstadzRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	WikipediaLink string `json:"wikipedia_link"`
}

type UpdateUstadzRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	WikipediaLink string `json:"wikipedia_link"`
}

type UstadzListResponse struct {
	UstadzList []Ustadz `json:"ustadz_list"`
	Total      int64    `json:"total"`
	Page       int      `json:"page"`
	Limit      int      `json:"limit"`
	TotalPages int      `json:"total_pages"`
}
