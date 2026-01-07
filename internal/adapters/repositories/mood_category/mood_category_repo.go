package mood_category

import (
	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"

	"gorm.io/gorm"
)

type moodCategoryRepo struct {
	db *gorm.DB
}

func NewMoodCategoryRepo(db *gorm.DB) ports.MoodCategoryRepository {
	return &moodCategoryRepo{db: db}
}

func (r *moodCategoryRepo) Create(mood *domain.MoodCategory) error {
	return r.db.Create(mood).Error
}

func (r *moodCategoryRepo) FindByID(id uint) (*domain.MoodCategory, error) {
	var mood domain.MoodCategory
	err := r.db.First(&mood, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &mood, err
}

func (r *moodCategoryRepo) FindByUUID(uuid string) (*domain.MoodCategory, error) {
	var mood domain.MoodCategory
	err := r.db.Where("uuid = ?", uuid).First(&mood).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &mood, err
}

func (r *moodCategoryRepo) FindByName(name string) (*domain.MoodCategory, error) {
	var mood domain.MoodCategory
	err := r.db.Where("name = ?", name).First(&mood).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &mood, err
}

func (r *moodCategoryRepo) FindAll() ([]domain.MoodCategory, int64, error) {
	var moods []domain.MoodCategory
	var total int64

	r.db.Model(&domain.MoodCategory{}).Count(&total)
	err := r.db.Find(&moods).Error

	return moods, total, err
}

func (r *moodCategoryRepo) Update(mood *domain.MoodCategory) error {
	return r.db.Save(mood).Error
}

func (r *moodCategoryRepo) Delete(id uint) error {
	return r.db.Delete(&domain.MoodCategory{}, id).Error
}
