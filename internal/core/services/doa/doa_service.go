package doa

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

type doaService struct {
	repo       ports.DoaRepository
	hadistRepo ports.HadistRepository
}

func NewDoaService(repo ports.DoaRepository, hRepo ports.HadistRepository) ports.DoaService {
	return &doaService{
		repo:       repo,
		hadistRepo: hRepo,
	}
}

func (s *doaService) Create(req *domain.CreateDoaRequest) (*domain.Doa, error) {
	if req.JudulDoa == "" {
		return nil, errors.New("judul doa is required")
	}

	var hadistID *uint
	if req.HadistID != "" {
		hadist, err := s.hadistRepo.FindByUUID(req.HadistID)
		if err != nil {
			logger.Log.Warn("Failed to find hadist for linking", zap.String("uuid", req.HadistID))
			// We can choose to error or ignore. Let's error if ID provided but invalid.
			if hadist == nil {
				return nil, errors.New("referenced hadist not found") 
			}
		}
		if hadist != nil {
			hadistID = &hadist.ID
		}
	}

	doa := req.ToEntity(uuid.New().String(), hadistID)

	if err := s.repo.Create(doa); err != nil {
		logger.Log.Error("Failed to create doa", zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	return doa, nil
}

func (s *doaService) Update(uuid string, req *domain.UpdateDoaRequest) (*domain.Doa, error) {
	doa, err := s.repo.FindByUUID(uuid)
	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}
	if doa == nil {
		return nil, errors.New("doa not found")
	}
	
	var hadistID *uint
	if req.HadistID != "" {
		hadist, err := s.hadistRepo.FindByUUID(req.HadistID)
		if err != nil || hadist == nil {
			return nil, errors.New("referenced hadist not found")
		}
		hadistID = &hadist.ID
	}

	req.ApplyUpdates(doa, hadistID)

	if err := s.repo.Update(doa); err != nil {
		logger.Log.Error("Failed to update doa", zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	return doa, nil
}

func (s *doaService) Delete(uuid string) error {
	doa, err := s.repo.FindByUUID(uuid)
	if err != nil || doa == nil {
		return errors.New("doa not found")
	}
	return s.repo.Delete(doa.ID)
}

func (s *doaService) GetAll(page, limit int) (*domain.DoaListResponse, error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 20 }

	doas, total, err := s.repo.FindAll(page, limit)
	if err != nil {
		logger.Log.Error("Failed to fetch doas", zap.Error(err))
		return nil, errors.New(messages.ErrInternalServer)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &domain.DoaListResponse{
		Doas:       doas,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *doaService) GetByUUID(uuid string) (*domain.Doa, error) {
	doa, err := s.repo.FindByUUID(uuid)
	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}
	if doa == nil {
		return nil, errors.New("doa not found")
	}
	return doa, nil
}

func (s *doaService) GetByCategory(category string, page, limit int) (*domain.DoaListResponse, error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 20 }

	doas, total, err := s.repo.FindByCategory(category, page, limit)
	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}
	
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return &domain.DoaListResponse{
		Doas:       doas,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *doaService) GetByHadist(hadistUUID string, page, limit int) (*domain.DoaListResponse, error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 20 }
	
	hadist, err := s.hadistRepo.FindByUUID(hadistUUID)
	if err != nil || hadist == nil {
		return nil, errors.New("hadist not found")
	}

	doas, total, err := s.repo.FindByHadistID(hadist.ID, page, limit)
	if err != nil {
		return nil, errors.New(messages.ErrInternalServer)
	}
	
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return &domain.DoaListResponse{
		Doas:       doas,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}


func (s *doaService) GetRandom() (*domain.Doa, error) {
	doa, err := s.repo.FindRandom()
	if err != nil || doa == nil {
		return nil, errors.New("no doa found")
	}
	return doa, nil
}

func (s *doaService) IncrementListeningCount(uuid string) error {
	doa, err := s.repo.FindByUUID(uuid)
	if err != nil || doa == nil {
		return errors.New("doa not found")
	}
	return s.repo.IncrementListeningCount(doa.ID)
}

// Like Logic
func (s *doaService) LikeDoa(userID uint, doaUUID string) error {
	doa, err := s.repo.FindByUUID(doaUUID)
	if err != nil || doa == nil {
		return errors.New("doa not found")
	}
	
	existing, _ := s.repo.FindLikeByUserAndDoa(userID, doa.ID)
	if existing != nil {
		return errors.New("already liked")
	}

	like := &domain.DoaLike{UserID: userID, DoaID: doa.ID}
	if err := s.repo.CreateLike(like); err != nil {
		return errors.New(messages.ErrInternalServer)
	}
	s.repo.IncrementLikeCount(doa.ID)
	return nil
}

func (s *doaService) UnlikeDoa(userID uint, doaUUID string) error {
	doa, err := s.repo.FindByUUID(doaUUID)
	if err != nil || doa == nil {
		return errors.New("doa not found")
	}
	
	like, _ := s.repo.FindLikeByUserAndDoa(userID, doa.ID)
	if like == nil {
		return errors.New("not liked")
	}
	
	if err := s.repo.DeleteLike(like.ID); err != nil {
		return errors.New(messages.ErrInternalServer)
	}
	s.repo.DecrementLikeCount(doa.ID)
	return nil
}

func (s *doaService) IsLiked(userID uint, doaUUID string) (bool, error) {
	doa, err := s.repo.FindByUUID(doaUUID)
	if err != nil || doa == nil {
		return false, errors.New("doa not found")
	}
	like, _ := s.repo.FindLikeByUserAndDoa(userID, doa.ID)
	return like != nil, nil
}

// Bookmark Logic
func (s *doaService) BookmarkDoa(userID uint, doaUUID string) error {
	doa, err := s.repo.FindByUUID(doaUUID)
	if err != nil || doa == nil {
		return errors.New("doa not found")
	}
	
	existing, _ := s.repo.FindBookmarkByUserAndDoa(userID, doa.ID)
	if existing != nil {
		return errors.New("already bookmarked")
	}

	bookmark := &domain.DoaBookmark{UserID: userID, DoaID: doa.ID}
	if err := s.repo.CreateBookmark(bookmark); err != nil {
		return errors.New(messages.ErrInternalServer)
	}
	s.repo.IncrementBookmarkCount(doa.ID)
	return nil
}

func (s *doaService) UnbookmarkDoa(userID uint, doaUUID string) error {
	doa, err := s.repo.FindByUUID(doaUUID)
	if err != nil || doa == nil {
		return errors.New("doa not found")
	}
	
	bookmark, _ := s.repo.FindBookmarkByUserAndDoa(userID, doa.ID)
	if bookmark == nil {
		return errors.New("not bookmarked")
	}
	
	if err := s.repo.DeleteBookmark(bookmark.ID); err != nil {
		return errors.New(messages.ErrInternalServer)
	}
	s.repo.DecrementBookmarkCount(doa.ID)
	return nil
}

func (s *doaService) IsBookmarked(userID uint, doaUUID string) (bool, error) {
	doa, err := s.repo.FindByUUID(doaUUID)
	if err != nil || doa == nil {
		return false, errors.New("doa not found")
	}
	bookmark, _ := s.repo.FindBookmarkByUserAndDoa(userID, doa.ID)
	return bookmark != nil, nil
}
