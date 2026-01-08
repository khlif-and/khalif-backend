package audio

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

type audioService struct {
	audioRepo    ports.AudioRepository
	moodRepo     ports.MoodCategoryRepository
	ustadzRepo   ports.UstadzRepository
}

func NewAudioService(audioRepo ports.AudioRepository, moodRepo ports.MoodCategoryRepository, ustadzRepo ports.UstadzRepository) ports.AudioService {
	return &audioService{
		audioRepo:  audioRepo,
		moodRepo:   moodRepo,
		ustadzRepo: ustadzRepo,
	}
}

func (s *audioService) Create(req *domain.CreateAudioRequest) (*domain.Audio, error) {
	if req.Title == "" {
		return nil, errors.New(messages.ErrAudioTitleRequired)
	}
	if req.AudioFile == "" {
		return nil, errors.New(messages.ErrAudioFileRequired)
	}

	var ustadzID *uint
	if req.UstadzUUID != "" {
		ustadz, err := s.ustadzRepo.FindByUUID(req.UstadzUUID)
		if err != nil || ustadz == nil {
			return nil, errors.New(messages.ErrUstadzNotFound)
		}
		ustadzID = &ustadz.ID
	}

	var moodCategoryID *uint
	if req.MoodCategoryUUID != "" {
		mood, err := s.moodRepo.FindByUUID(req.MoodCategoryUUID)
		if err != nil || mood == nil {
			return nil, errors.New(messages.ErrMoodNotFound)
		}
		moodCategoryID = &mood.ID
	}

	audio := &domain.Audio{
		UUID:                uuid.New().String(),
		Title:               req.Title,
		AudioFile:           req.AudioFile,
		ThumbnailFile:       req.ThumbnailFile,
		ColorThumbnailAudio: req.ColorThumbnailAudio,
		UstadzID:            ustadzID,
		DurationAudio:       req.DurationAudio,
		MoodCategoryID:      moodCategoryID,
		ListeningCount:      0,
		LikeCount:           0,
	}

	if err := s.audioRepo.Create(audio); err != nil {
		logger.Log.Error("Failed to create audio", zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	return audio, nil
}

func (s *audioService) GetByUUID(uuid string) (*domain.Audio, error) {
	audio, err := s.audioRepo.FindByUUID(uuid)
	if err != nil {
		logger.Log.Error("Failed to find audio", zap.String("uuid", uuid), zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}
	if audio == nil {
		return nil, errors.New(messages.ErrAudioNotFound)
	}
	return audio, nil
}

func (s *audioService) GetAll(page, limit int) (*domain.AudioListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	audios, total, err := s.audioRepo.FindAll(page, limit)
	if err != nil {
		logger.Log.Error("Failed to fetch audios", zap.Error(err))
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

func (s *audioService) Update(uuid string, req *domain.UpdateAudioRequest) (*domain.Audio, error) {
	audio, err := s.audioRepo.FindByUUID(uuid)
	if err != nil {
		logger.Log.Error("Failed to find audio for update", zap.String("uuid", uuid), zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}
	if audio == nil {
		return nil, errors.New(messages.ErrAudioNotFound)
	}

	if req.Title != "" {
		audio.Title = req.Title
	}
	if req.UstadzUUID != "" {
		ustadz, err := s.ustadzRepo.FindByUUID(req.UstadzUUID)
		if err != nil || ustadz == nil {
			return nil, errors.New(messages.ErrUstadzNotFound)
		}
		audio.UstadzID = &ustadz.ID
	}
	if req.DurationAudio > 0 {
		audio.DurationAudio = req.DurationAudio
	}
	if req.ColorThumbnailAudio != "" {
		audio.ColorThumbnailAudio = req.ColorThumbnailAudio
	}
	if req.AudioFile != "" {
		audio.AudioFile = req.AudioFile
	}
	if req.ThumbnailFile != "" {
		audio.ThumbnailFile = req.ThumbnailFile
	}
	if req.MoodCategoryUUID != "" {
		mood, err := s.moodRepo.FindByUUID(req.MoodCategoryUUID)
		if err != nil || mood == nil {
			return nil, errors.New(messages.ErrMoodNotFound)
		}
		audio.MoodCategoryID = &mood.ID
	}

	if err := s.audioRepo.Update(audio); err != nil {
		logger.Log.Error("Failed to update audio", zap.String("uuid", uuid), zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	return audio, nil
}

func (s *audioService) Delete(uuid string) error {
	audio, err := s.audioRepo.FindByUUID(uuid)
	if err != nil {
		logger.Log.Error("Failed to find audio for delete", zap.String("uuid", uuid), zap.Error(err))
		return errors.New(messages.ErrInternalServer)
	}
	if audio == nil {
		return errors.New(messages.ErrAudioNotFound)
	}

	if err := s.audioRepo.Delete(audio.ID); err != nil {
		logger.Log.Error("Failed to delete audio", zap.String("uuid", uuid), zap.Error(err))
		return errors.New(messages.ErrInternalServer)
	}

	return nil
}

func (s *audioService) IncrementListeningCount(uuid string) error {
	audio, err := s.audioRepo.FindByUUID(uuid)
	if err != nil {
		logger.Log.Error("Failed to find audio for increment", zap.String("uuid", uuid), zap.Error(err))
		return errors.New(messages.ErrInternalServer)
	}
	if audio == nil {
		return errors.New(messages.ErrAudioNotFound)
	}

	if err := s.audioRepo.IncrementListeningCount(audio.ID); err != nil {
		logger.Log.Error("Failed to increment listening count", zap.String("uuid", uuid), zap.Error(err))
		return errors.New(messages.ErrInternalServer)
	}

	return nil
}

// RecordListening records user listening via SP (prevents spam - 1 per user per audio per day)
func (s *audioService) RecordListening(userID uint, audioUUID string) (alreadyListened bool, newCount int64, err error) {
	audio, err := s.audioRepo.FindByUUID(audioUUID)
	if err != nil {
		logger.Log.Error("Failed to find audio for recording", zap.String("uuid", audioUUID), zap.Error(err))
		return false, 0, errors.New(messages.ErrInternalServer)
	}
	if audio == nil {
		return false, 0, errors.New(messages.ErrAudioNotFound)
	}

	already, count, err := s.audioRepo.RecordListening(userID, audio.ID)
	if err != nil {
		logger.Log.Error("Failed to record listening", zap.Uint("userID", userID), zap.String("audioUUID", audioUUID), zap.Error(err))
		return false, 0, errors.New(messages.ErrInternalServer)
	}

	return already, count, nil
}

// GetUserListeningHistory returns paginated listening history for a user
func (s *audioService) GetUserListeningHistory(userID uint, page, limit int) (*domain.ListeningHistoryResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	history, total, err := s.audioRepo.GetUserListeningHistory(userID, page, limit)
	if err != nil {
		logger.Log.Error("Failed to fetch listening history", zap.Uint("userID", userID), zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &domain.ListeningHistoryResponse{
		History:    history,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
