package services

import (
	"context"
	"errors"

	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"
	"khalif-backend/internal/platform/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type alquranService struct {
	repo ports.AlquranRepository
}

func NewAlquranService(repo ports.AlquranRepository) ports.AlquranService {
	return &alquranService{
		repo: repo,
	}
}

func (s *alquranService) CreateAlquran(ctx context.Context, req *domain.CreateAlquranRequest) error {
	// 2. Generate UUIDs if missing
	if req.SurahID == "" {
		req.SurahID = uuid.New().String()
	}
	
	// Ensure Ayats have UUIDs
	for i := range req.Ayats {
		if req.Ayats[i].AyatID == "" {
			req.Ayats[i].AyatID = uuid.New().String()
		}
	}

	// 3. Map to Entity
	surah := req.ToEntity()

	// 4. Save Surah (Parent)
	// Note: GORM can handle deep creation if configured, or we save one by one.
	// We'll try saving the Surah with Ayats because we defined the relationship in domain.
	// But let's look at `repo.CreateSurah`. It calls `db.Create(surah)`.
	// If `surah.Ayats` is populated, GORM default behavior is to create associations too.
	
	if err := s.repo.CreateSurah(ctx, &surah); err != nil {
		logger.Log.Error("Failed to create surah", zap.Error(err))
		return errors.New("failed to create alquran data")
	}

	return nil
}

func (s *alquranService) GetAllSurah(ctx context.Context) ([]domain.Surah, error) {
	return s.repo.FindAllSurah(ctx)
}

func (s *alquranService) GetSurahByUUID(ctx context.Context, uuid string) (*domain.Surah, error) {
	surah, err := s.repo.FindSurahByUUID(ctx, uuid)
	if err != nil {
		logger.Log.Error("Failed to find surah", zap.String("uuid", uuid), zap.Error(err))
		return nil, errors.New("surah not found")
	}
	return surah, nil
}
