package audio

import (
	"errors"

	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"

	"gorm.io/gorm"
)

type AudioRepo struct {
	db *gorm.DB
}

func NewAudioRepo(db *gorm.DB) ports.AudioRepository {
	return &AudioRepo{db: db}
}

func (r *AudioRepo) Create(audio *domain.Audio) error {
	return r.db.Create(audio).Error
}

func (r *AudioRepo) FindByID(id uint) (*domain.Audio, error) {
	var audio domain.Audio
	if err := r.db.Preload("MoodCategory").Preload("Ustadz").First(&audio, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &audio, nil
}

func (r *AudioRepo) FindByUUID(uuid string) (*domain.Audio, error) {
	var audio domain.Audio
	if err := r.db.Preload("MoodCategory").Preload("Ustadz").Where("uuid = ?", uuid).First(&audio).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &audio, nil
}

func (r *AudioRepo) FindAll(page, limit int) ([]domain.Audio, int64, error) {
	var audios []domain.Audio
	var total int64

	if err := r.db.Model(&domain.Audio{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	if err := r.db.Preload("MoodCategory").Preload("Ustadz").Order("created_at DESC").Offset(offset).Limit(limit).Find(&audios).Error; err != nil {
		return nil, 0, err
	}

	return audios, total, nil
}

func (r *AudioRepo) FindByMoodCategoryID(moodCategoryID uint, page, limit int) ([]domain.Audio, int64, error) {
	var audios []domain.Audio
	var total int64

	query := r.db.Model(&domain.Audio{}).Where("mood_category_id = ?", moodCategoryID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	if err := query.Preload("MoodCategory").Preload("Ustadz").Order("created_at DESC").Offset(offset).Limit(limit).Find(&audios).Error; err != nil {
		return nil, 0, err
	}

	return audios, total, nil
}

func (r *AudioRepo) Update(audio *domain.Audio) error {
	return r.db.Save(audio).Error
}

func (r *AudioRepo) Delete(id uint) error {
	return r.db.Delete(&domain.Audio{}, id).Error
}

func (r *AudioRepo) IncrementListeningCount(id uint) error {
	return r.db.Model(&domain.Audio{}).Where("id = ?", id).UpdateColumn("listening_count", gorm.Expr("listening_count + ?", 1)).Error
}
