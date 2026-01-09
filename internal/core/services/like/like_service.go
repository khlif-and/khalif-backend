package like

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

type likeService struct {
	likeRepo  ports.LikeRepository
	audioRepo ports.AudioRepository
}

func NewLikeService(likeRepo ports.LikeRepository, audioRepo ports.AudioRepository) ports.LikeService {
	return &likeService{
		likeRepo:  likeRepo,
		audioRepo: audioRepo,
	}
}

func (s *likeService) LikeAudio(userID uint, audioUUID string) (*domain.Like, error) {
	audio, err := s.audioRepo.FindByUUID(audioUUID)
	if err != nil {
		logger.Log.Error("Failed to find audio for like", zap.String("audio_uuid", audioUUID), zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}
	if audio == nil {
		return nil, errors.New(messages.ErrAudioNotFound)
	}

	existing, _ := s.likeRepo.FindByUserAndAudio(userID, audio.ID)
	if existing != nil {
		return nil, errors.New(messages.ErrAlreadyLiked)
	}

	like := &domain.Like{
		UUID:    uuid.New().String(),
		UserID:  userID,
		AudioID: audio.ID,
	}

	if err := s.likeRepo.Create(like); err != nil {
		logger.Log.Error("Failed to create like", zap.Uint("user_id", userID), zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	if err := s.likeRepo.IncrementAudioLikeCount(audio.ID); err != nil {
		logger.Log.Error("Failed to increment like count", zap.Uint("audio_id", audio.ID), zap.Error(err))
	}

	return like, nil
}

func (s *likeService) UnlikeAudio(userID uint, audioUUID string) error {
	audio, err := s.audioRepo.FindByUUID(audioUUID)
	if err != nil {
		logger.Log.Error("Failed to find audio for unlike", zap.String("audio_uuid", audioUUID), zap.Error(err))
		return errors.New(messages.ErrInternalServer)
	}
	if audio == nil {
		return errors.New(messages.ErrAudioNotFound)
	}

	like, err := s.likeRepo.FindByUserAndAudio(userID, audio.ID)
	if err != nil {
		logger.Log.Error("Failed to find like", zap.Uint("user_id", userID), zap.Error(err))
		return errors.New(messages.ErrInternalServer)
	}
	if like == nil {
		return errors.New(messages.ErrLikeNotFound)
	}

	if err := s.likeRepo.Delete(like.ID); err != nil {
		logger.Log.Error("Failed to delete like", zap.Uint("like_id", like.ID), zap.Error(err))
		return errors.New(messages.ErrInternalServer)
	}

	if err := s.likeRepo.DecrementAudioLikeCount(audio.ID); err != nil {
		logger.Log.Error("Failed to decrement like count", zap.Uint("audio_id", audio.ID), zap.Error(err))
	}

	return nil
}

func (s *likeService) GetUserLikes(userID uint, page, limit int) (*domain.LikeListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	likes, total, err := s.likeRepo.FindByUserID(userID, page, limit)
	if err != nil {
		logger.Log.Error("Failed to fetch user likes", zap.Uint("user_id", userID), zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &domain.LikeListResponse{
		Likes:      likes,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *likeService) IsLiked(userID uint, audioUUID string) (bool, error) {
	audio, err := s.audioRepo.FindByUUID(audioUUID)
	if err != nil || audio == nil {
		return false, nil
	}

	like, err := s.likeRepo.FindByUserAndAudio(userID, audio.ID)
	if err != nil {
		return false, errors.New(messages.ErrInternalServer)
	}

	return like != nil, nil
}
