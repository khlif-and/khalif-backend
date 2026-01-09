package ustadz

import (
	"errors"
	"math"

	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"
	"khalif-backend/internal/platform/logger"
	"khalif-backend/pkg/messages"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ustadzService struct {
	ustadzRepo ports.UstadzRepository
}

func NewUstadzService(ustadzRepo ports.UstadzRepository) ports.UstadzService {
	return &ustadzService{
		ustadzRepo: ustadzRepo,
	}
}

func (s *ustadzService) Create(req *domain.CreateUstadzRequest) (*domain.Ustadz, error) {
	if req.Name == "" {
		return nil, errors.New(messages.ErrUstadzNameRequired)
	}

	existing, _ := s.ustadzRepo.FindByName(req.Name)
	if existing != nil {
		return nil, errors.New(messages.ErrUstadzNameExists)
	}

	ustadz := &domain.Ustadz{
		UUID:          uuid.New().String(),
		Name:          req.Name,
		Description:   req.Description,
		WikipediaLink: req.WikipediaLink,
	}

	if err := s.ustadzRepo.Create(ustadz); err != nil {
		logger.Log.Error("Failed to create ustadz", zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	return ustadz, nil
}

func (s *ustadzService) GetByUUID(uuid string) (*domain.Ustadz, error) {
	ustadz, err := s.ustadzRepo.FindByUUID(uuid)
	if err != nil {
		logger.Log.Error("Failed to find ustadz", zap.String("uuid", uuid), zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}
	if ustadz == nil {
		return nil, errors.New(messages.ErrUstadzNotFound)
	}
	return ustadz, nil
}

func (s *ustadzService) GetAll(page, limit int) (*domain.UstadzListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	ustadzList, total, err := s.ustadzRepo.FindAll(page, limit)
	if err != nil {
		logger.Log.Error("Failed to fetch ustadz list", zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &domain.UstadzListResponse{
		UstadzList: ustadzList,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *ustadzService) Update(uuid string, req *domain.UpdateUstadzRequest) (*domain.Ustadz, error) {
	ustadz, err := s.ustadzRepo.FindByUUID(uuid)
	if err != nil {
		logger.Log.Error("Failed to find ustadz for update", zap.String("uuid", uuid), zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}
	if ustadz == nil {
		return nil, errors.New(messages.ErrUstadzNotFound)
	}

	if req.Name != "" {
		existing, _ := s.ustadzRepo.FindByName(req.Name)
		if existing != nil && existing.ID != ustadz.ID {
			return nil, errors.New(messages.ErrUstadzNameExists)
		}
		ustadz.Name = req.Name
	}
	if req.Description != "" {
		ustadz.Description = req.Description
	}
	if req.WikipediaLink != "" {
		ustadz.WikipediaLink = req.WikipediaLink
	}

	if err := s.ustadzRepo.Update(ustadz); err != nil {
		logger.Log.Error("Failed to update ustadz", zap.String("uuid", uuid), zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	return ustadz, nil
}

func (s *ustadzService) Delete(uuid string) error {
	ustadz, err := s.ustadzRepo.FindByUUID(uuid)
	if err != nil {
		logger.Log.Error("Failed to find ustadz for delete", zap.String("uuid", uuid), zap.Error(err))
		return errors.New(messages.ErrInternalServer)
	}
	if ustadz == nil {
		return errors.New(messages.ErrUstadzNotFound)
	}

	if err := s.ustadzRepo.Delete(ustadz.ID); err != nil {
		logger.Log.Error("Failed to delete ustadz", zap.String("uuid", uuid), zap.Error(err))
		return errors.New(messages.ErrInternalServer)
	}

	return nil
}
