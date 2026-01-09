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

// RecordListening uses stored procedure to record listening and prevent spam
func (r *AudioRepo) RecordListening(userID, audioID uint) (alreadyListened bool, newCount int64, err error) {
	var already bool
	var count int64

	row := r.db.Raw("SELECT * FROM sp_record_listening($1, $2)", userID, audioID).Row()
	if err := row.Scan(&already, &count); err != nil {
		return false, 0, err
	}

	return already, count, nil
}

// GetUserListeningHistory returns paginated listening history for a user
func (r *AudioRepo) GetUserListeningHistory(userID uint, page, limit int) ([]domain.ListeningHistory, int64, error) {
	var history []domain.ListeningHistory
	var total int64

	if err := r.db.Model(&domain.ListeningHistory{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	if err := r.db.Preload("Audio").Preload("Audio.MoodCategory").Preload("Audio.Ustadz").
		Where("user_id = ?", userID).
		Order("listened_at DESC").
		Offset(offset).Limit(limit).
		Find(&history).Error; err != nil {
		return nil, 0, err
	}

	return history, total, nil
}

// GetRadioQueue returns a queue of similar audios based on seed audio
// Uses scoring: same ustadz (+3), same mood (+2), popular (+1)
func (r *AudioRepo) GetRadioQueue(seedAudio *domain.Audio, limit int) ([]domain.Audio, error) {
	var audios []domain.Audio

	// Build raw SQL for scoring
	// Priority: 1) Same ustadz, 2) Same mood, 3) Popular, then random
	orderClause := gorm.Expr(`
		CASE WHEN ustadz_id = ? THEN 3 ELSE 0 END +
		CASE WHEN mood_category_id = ? THEN 2 ELSE 0 END +
		CASE WHEN listening_count > 100 THEN 1 ELSE 0 END DESC,
		listening_count DESC,
		RANDOM()
	`, seedAudio.UstadzID, seedAudio.MoodCategoryID)

	query := r.db.Model(&domain.Audio{}).
		Preload("MoodCategory").
		Preload("Ustadz").
		Where("id != ?", seedAudio.ID). // Exclude seed audio
		Order(orderClause).
		Limit(limit)

	if err := query.Find(&audios).Error; err != nil {
		return nil, err
	}

	return audios, nil
}
