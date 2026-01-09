package mood_category

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

type moodCategoryService struct {
	moodRepo  ports.MoodCategoryRepository
	audioRepo ports.AudioRepository
}

func NewMoodCategoryService(moodRepo ports.MoodCategoryRepository, audioRepo ports.AudioRepository) ports.MoodCategoryService {
	return &moodCategoryService{
		moodRepo:  moodRepo,
		audioRepo: audioRepo,
	}
}

func (s *moodCategoryService) Create(req *domain.CreateMoodCategoryRequest) (*domain.MoodCategory, error) {
	if req.Name == "" {
		return nil, errors.New(messages.ErrMoodNameRequired)
	}

	existing, _ := s.moodRepo.FindByName(req.Name)
	if existing != nil {
		return nil, errors.New(messages.ErrMoodNameExists)
	}

	mood := &domain.MoodCategory{
		UUID:  uuid.New().String(),
		Name:  req.Name,
		Icon:  req.Icon,
		Color: req.Color,
	}

	if err := s.moodRepo.Create(mood); err != nil {
		logger.Log.Error("Failed to create mood category", zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	return mood, nil
}

func (s *moodCategoryService) GetByUUID(uuid string) (*domain.MoodCategory, error) {
	mood, err := s.moodRepo.FindByUUID(uuid)
	if err != nil {
		logger.Log.Error("Failed to find mood category", zap.String("uuid", uuid), zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}
	if mood == nil {
		return nil, errors.New(messages.ErrMoodNotFound)
	}
	return mood, nil
}

func (s *moodCategoryService) GetAll(page, limit int) (*domain.MoodCategoryListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	moods, total, err := s.moodRepo.FindAll(page, limit)
	if err != nil {
		logger.Log.Error("Failed to fetch mood categories", zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &domain.MoodCategoryListResponse{
		MoodCategories: moods,
		Total:          total,
		Page:           page,
		Limit:          limit,
		TotalPages:     totalPages,
	}, nil
}

func (s *moodCategoryService) GetAudiosByMoodUUID(moodUUID string, page, limit int) (*domain.AudioListResponse, error) {
	mood, err := s.moodRepo.FindByUUID(moodUUID)
	if err != nil || mood == nil {
		return nil, errors.New(messages.ErrMoodNotFound)
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	audios, total, err := s.audioRepo.FindByMoodCategoryID(mood.ID, page, limit)
	if err != nil {
		logger.Log.Error("Failed to fetch audios by mood", zap.String("mood_uuid", moodUUID), zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &domain.AudioListResponse{
		Audios:     audios,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *moodCategoryService) Update(uuid string, req *domain.UpdateMoodCategoryRequest) (*domain.MoodCategory, error) {
	mood, err := s.moodRepo.FindByUUID(uuid)
	if err != nil {
		logger.Log.Error("Failed to find mood category for update", zap.String("uuid", uuid), zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}
	if mood == nil {
		return nil, errors.New(messages.ErrMoodNotFound)
	}

	if req.Name != "" {
		existing, _ := s.moodRepo.FindByName(req.Name)
		if existing != nil && existing.ID != mood.ID {
			return nil, errors.New(messages.ErrMoodNameExists)
		}
		mood.Name = req.Name
	}
	if req.Icon != "" {
		mood.Icon = req.Icon
	}
	if req.Color != "" {
		mood.Color = req.Color
	}

	if err := s.moodRepo.Update(mood); err != nil {
		logger.Log.Error("Failed to update mood category", zap.String("uuid", uuid), zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	return mood, nil
}

func (s *moodCategoryService) Delete(uuid string) error {
	mood, err := s.moodRepo.FindByUUID(uuid)
	if err != nil {
		logger.Log.Error("Failed to find mood category for delete", zap.String("uuid", uuid), zap.Error(err))
		return errors.New(messages.ErrInternalServer)
	}
	if mood == nil {
		return errors.New(messages.ErrMoodNotFound)
	}

	if err := s.moodRepo.Delete(mood.ID); err != nil {
		logger.Log.Error("Failed to delete mood category", zap.String("uuid", uuid), zap.Error(err))
		return errors.New(messages.ErrInternalServer)
	}

	return nil
}
