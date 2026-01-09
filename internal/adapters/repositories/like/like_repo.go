package like

import (
	"errors"

	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"

	"gorm.io/gorm"
)

type likeRepo struct {
	db *gorm.DB
}

func NewLikeRepo(db *gorm.DB) ports.LikeRepository {
	return &likeRepo{db: db}
}

func (r *likeRepo) Create(like *domain.Like) error {
	return r.db.Create(like).Error
}

func (r *likeRepo) FindByUUID(uuid string) (*domain.Like, error) {
	var like domain.Like
	err := r.db.Preload("Audio").Preload("Audio.Ustadz").Preload("Audio.MoodCategory").Where("uuid = ?", uuid).First(&like).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &like, err
}

func (r *likeRepo) FindByUserAndAudio(userID, audioID uint) (*domain.Like, error) {
	var like domain.Like
	err := r.db.Where("user_id = ? AND audio_id = ?", userID, audioID).First(&like).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &like, err
}

func (r *likeRepo) FindByUserID(userID uint, page, limit int) ([]domain.Like, int64, error) {
	var likes []domain.Like
	var total int64

	r.db.Model(&domain.Like{}).Where("user_id = ?", userID).Count(&total)

	offset := (page - 1) * limit
	err := r.db.Preload("Audio").Preload("Audio.Ustadz").Preload("Audio.MoodCategory").
		Where("user_id = ?", userID).Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&likes).Error

	return likes, total, err
}

func (r *likeRepo) Delete(id uint) error {
	return r.db.Delete(&domain.Like{}, id).Error
}

func (r *likeRepo) IncrementAudioLikeCount(audioID uint) error {
	return r.db.Model(&domain.Audio{}).Where("id = ?", audioID).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
}

func (r *likeRepo) DecrementAudioLikeCount(audioID uint) error {
	return r.db.Model(&domain.Audio{}).Where("id = ?", audioID).UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - 1, 0)")).Error
}
