package repositories

import (
	"context"

	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"

	"gorm.io/gorm"
)

type alquranRepository struct {
	db *gorm.DB
}

func NewAlquranRepository(db *gorm.DB) ports.AlquranRepository {
	return &alquranRepository{
		db: db,
	}
}

func (r *alquranRepository) CreateSurah(ctx context.Context, surah *domain.Surah) error {
	return r.db.WithContext(ctx).Create(surah).Error
}

func (r *alquranRepository) CreateAyat(ctx context.Context, ayat *domain.Ayat) error {
	return r.db.WithContext(ctx).Create(ayat).Error
}

func (r *alquranRepository) FindAllSurah(ctx context.Context) ([]domain.Surah, error) {
	var surahs []domain.Surah
	// Retrieving all surahs, ordered by NoSurah
	// Omit ayats for list view to keep it light
	err := r.db.WithContext(ctx).Order("no_surah asc").Find(&surahs).Error
	return surahs, err
}

func (r *alquranRepository) FindSurahByUUID(ctx context.Context, uuid string) (*domain.Surah, error) {
	var surah domain.Surah
	// Preload Ayats ordered by number_of_ayat
	err := r.db.WithContext(ctx).
		Preload("Ayats", func(db *gorm.DB) *gorm.DB {
			return db.Order("number_of_ayat asc")
		}).
		Where("surah_id = ?", uuid).
		First(&surah).Error
	if err != nil {
		return nil, err
	}
	return &surah, nil
}
