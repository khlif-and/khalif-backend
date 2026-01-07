package domain

import (
	"time"

	"gorm.io/gorm"
)

type MoodCategory struct {
	ID        uint           `gorm:"primaryKey" json:"-"`
	UUID      string         `gorm:"uniqueIndex;size:36;not null" json:"id"`
	Name      string         `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Icon      string         `gorm:"size:50" json:"icon"`
	Color     string         `gorm:"size:20" json:"color"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateMoodCategoryRequest struct {
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
}

type UpdateMoodCategoryRequest struct {
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
}

type MoodCategoryListResponse struct {
	MoodCategories []MoodCategory `json:"mood_categories"`
	Total          int64          `json:"total"`
}
