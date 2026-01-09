package hadist

import (
	"errors"
	"math"

	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/platform/logger"
	"khalif-backend/pkg/messages"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type HadistRepository interface {
	Create(hadist *domain.Hadist) error
	FindByID(id uint) (*domain.Hadist, error)
	FindByUUID(uuid string) (*domain.Hadist, error)
	FindAll(page, limit int) ([]domain.Hadist, int64, error)
	FindByCategory(category string, page, limit int) ([]domain.Hadist, int64, error)
	FindByKitab(kitab string, page, limit int) ([]domain.Hadist, int64, error)
	FindRandom() (*domain.Hadist, error)
	Update(hadist *domain.Hadist) error
	Delete(id uint) error
	IncrementListeningCount(id uint) error
	// Like
	CreateLike(like *domain.HadistLike) error
	FindLikeByUserAndHadist(userID, hadistID uint) (*domain.HadistLike, error)
	DeleteLike(id uint) error
	IncrementLikeCount(hadistID uint) error
	DecrementLikeCount(hadistID uint) error
	// Bookmark
	CreateBookmark(bookmark *domain.HadistBookmark) error
	FindBookmarkByUserAndHadist(userID, hadistID uint) (*domain.HadistBookmark, error)
	DeleteBookmark(id uint) error
	IncrementBookmarkCount(hadistID uint) error
	DecrementBookmarkCount(hadistID uint) error
	GetUserBookmarks(userID uint, page, limit int) ([]domain.HadistBookmark, int64, error)
}

type HadistService struct {
	repo HadistRepository
}

func NewHadistService(repo HadistRepository) *HadistService {
	return &HadistService{repo: repo}
}

func (s *HadistService) Create(req *domain.CreateHadistRequest) (*domain.Hadist, error) {
	if req.NamaHadist == "" {
		return nil, errors.New("nama hadist is required")
	}

	if req.ShahihStatus != "" && !req.ShahihStatus.IsValid() {
		return nil, errors.New("invalid shahih status")
	}

	hadist := req.ToEntity(uuid.New().String())

	if err := s.repo.Create(hadist); err != nil {
		logger.Log.Error("Failed to create hadist", zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	return hadist, nil
}

func (s *HadistService) GetByUUID(uuid string) (*domain.Hadist, error) {
	hadist, err := s.repo.FindByUUID(uuid)
	if err != nil {
		logger.Log.Error("Failed to find hadist", zap.String("uuid", uuid), zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}
	if hadist == nil {
		return nil, errors.New("hadist not found")
	}
	return hadist, nil
}

func (s *HadistService) GetAll(page, limit int) (*domain.HadistListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	hadists, total, err := s.repo.FindAll(page, limit)
	if err != nil {
		logger.Log.Error("Failed to fetch hadists", zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &domain.HadistListResponse{
		Hadists:    hadists,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *HadistService) GetByCategory(category string, page, limit int) (*domain.HadistListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	hadists, total, err := s.repo.FindByCategory(category, page, limit)
	if err != nil {
		logger.Log.Error("Failed to fetch hadists by category", zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &domain.HadistListResponse{
		Hadists:    hadists,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *HadistService) GetByKitab(kitab string, page, limit int) (*domain.HadistListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	hadists, total, err := s.repo.FindByKitab(kitab, page, limit)
	if err != nil {
		logger.Log.Error("Failed to fetch hadists by kitab", zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &domain.HadistListResponse{
		Hadists:    hadists,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *HadistService) GetRandom() (*domain.Hadist, error) {
	hadist, err := s.repo.FindRandom()
	if err != nil {
		logger.Log.Error("Failed to get random hadist", zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}
	if hadist == nil {
		return nil, errors.New("no hadist found")
	}
	return hadist, nil
}

func (s *HadistService) Update(uuid string, req *domain.UpdateHadistRequest) (*domain.Hadist, error) {
	hadist, err := s.repo.FindByUUID(uuid)
	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}
	if hadist == nil {
		return nil, errors.New("hadist not found")
	}

	if req.ShahihStatus != nil && !req.ShahihStatus.IsValid() {
		return nil, errors.New("invalid shahih status")
	}

	req.ApplyUpdates(hadist)

	if err := s.repo.Update(hadist); err != nil {
		logger.Log.Error("Failed to update hadist", zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	return hadist, nil
}

func (s *HadistService) Delete(uuid string) error {
	hadist, err := s.repo.FindByUUID(uuid)
	if err != nil {
		return errors.New(messages.ErrInternalServer)
	}
	if hadist == nil {
		return errors.New("hadist not found")
	}

	if err := s.repo.Delete(hadist.ID); err != nil {
		logger.Log.Error("Failed to delete hadist", zap.Error(err))
		return errors.New(messages.ErrInternalServer)
	}

	return nil
}

func (s *HadistService) IncrementListeningCount(uuid string) error {
	hadist, err := s.repo.FindByUUID(uuid)
	if err != nil {
		return errors.New(messages.ErrInternalServer)
	}
	if hadist == nil {
		return errors.New("hadist not found")
	}

	return s.repo.IncrementListeningCount(hadist.ID)
}

// Like operations
func (s *HadistService) LikeHadist(userID uint, hadistUUID string) error {
	hadist, err := s.repo.FindByUUID(hadistUUID)
	if err != nil || hadist == nil {
		return errors.New("hadist not found")
	}

	existing, _ := s.repo.FindLikeByUserAndHadist(userID, hadist.ID)
	if existing != nil {
		return errors.New("already liked")
	}

	like := &domain.HadistLike{
		UserID:   userID,
		HadistID: hadist.ID,
	}

	if err := s.repo.CreateLike(like); err != nil {
		return errors.New(messages.ErrInternalServer)
	}

	s.repo.IncrementLikeCount(hadist.ID)
	return nil
}

func (s *HadistService) UnlikeHadist(userID uint, hadistUUID string) error {
	hadist, err := s.repo.FindByUUID(hadistUUID)
	if err != nil || hadist == nil {
		return errors.New("hadist not found")
	}

	like, _ := s.repo.FindLikeByUserAndHadist(userID, hadist.ID)
	if like == nil {
		return errors.New("not liked")
	}

	if err := s.repo.DeleteLike(like.ID); err != nil {
		return errors.New(messages.ErrInternalServer)
	}

	s.repo.DecrementLikeCount(hadist.ID)
	return nil
}

func (s *HadistService) IsLiked(userID uint, hadistUUID string) (bool, error) {
	hadist, err := s.repo.FindByUUID(hadistUUID)
	if err != nil || hadist == nil {
		return false, errors.New("hadist not found")
	}

	like, _ := s.repo.FindLikeByUserAndHadist(userID, hadist.ID)
	return like != nil, nil
}

// Bookmark operations
func (s *HadistService) BookmarkHadist(userID uint, hadistUUID string) error {
	hadist, err := s.repo.FindByUUID(hadistUUID)
	if err != nil || hadist == nil {
		return errors.New("hadist not found")
	}

	existing, _ := s.repo.FindBookmarkByUserAndHadist(userID, hadist.ID)
	if existing != nil {
		return errors.New("already bookmarked")
	}

	bookmark := &domain.HadistBookmark{
		UserID:   userID,
		HadistID: hadist.ID,
	}

	if err := s.repo.CreateBookmark(bookmark); err != nil {
		return errors.New(messages.ErrInternalServer)
	}

	s.repo.IncrementBookmarkCount(hadist.ID)
	return nil
}

func (s *HadistService) UnbookmarkHadist(userID uint, hadistUUID string) error {
	hadist, err := s.repo.FindByUUID(hadistUUID)
	if err != nil || hadist == nil {
		return errors.New("hadist not found")
	}

	bookmark, _ := s.repo.FindBookmarkByUserAndHadist(userID, hadist.ID)
	if bookmark == nil {
		return errors.New("not bookmarked")
	}

	if err := s.repo.DeleteBookmark(bookmark.ID); err != nil {
		return errors.New(messages.ErrInternalServer)
	}

	s.repo.DecrementBookmarkCount(hadist.ID)
	return nil
}

func (s *HadistService) IsBookmarked(userID uint, hadistUUID string) (bool, error) {
	hadist, err := s.repo.FindByUUID(hadistUUID)
	if err != nil || hadist == nil {
		return false, errors.New("hadist not found")
	}

	bookmark, _ := s.repo.FindBookmarkByUserAndHadist(userID, hadist.ID)
	return bookmark != nil, nil
}
