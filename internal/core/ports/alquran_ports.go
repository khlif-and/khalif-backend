package ports

import (
	"context"

	"khalif-backend/internal/core/domain"
)

type AlquranRepository interface {
	CreateSurah(ctx context.Context, surah *domain.Surah) error
	CreateAyat(ctx context.Context, ayat *domain.Ayat) error
	// Add other methods (Get, List, etc) as needed later
	FindAllSurah(ctx context.Context) ([]domain.Surah, error)
	FindSurahByUUID(ctx context.Context, uuid string) (*domain.Surah, error)
}

type AlquranService interface {
	CreateAlquran(ctx context.Context, req *domain.CreateAlquranRequest) error
	GetAllSurah(ctx context.Context) ([]domain.Surah, error)
	GetSurahByUUID(ctx context.Context, uuid string) (*domain.Surah, error)
}
