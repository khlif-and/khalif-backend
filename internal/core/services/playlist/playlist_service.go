package playlist

import (
	"errors"
	"math"

	"khalif-backend/internal/adapters/repositories/playlist"
	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"
	"khalif-backend/pkg/messages"

	"github.com/google/uuid"
)

type PlaylistService interface {
	// CRUD
	Create(req *domain.CreatePlaylistRequest, authorType domain.AuthorType, authorID uint) (*domain.Playlist, error)
	GetByUUID(uuid string) (*domain.PlaylistResponse, error)
	GetAll(page, limit int) (*domain.PlaylistListResponse, error)
	GetMyPlaylists(authorType domain.AuthorType, authorID uint, page, limit int) (*domain.PlaylistListResponse, error)
	Update(uuid string, req *domain.UpdatePlaylistRequest, authorType domain.AuthorType, authorID uint) (*domain.Playlist, error)
	Delete(uuid string, authorType domain.AuthorType, authorID uint) error

	// Audio management
	AddAudio(playlistUUID, audioUUID string, position int, authorType domain.AuthorType, authorID uint) error
	RemoveAudio(playlistUUID, audioUUID string, authorType domain.AuthorType, authorID uint) error

	// Likes
	LikePlaylist(playlistUUID string, userID uint) error
	UnlikePlaylist(playlistUUID string, userID uint) error
	IsLiked(playlistUUID string, userID uint) (bool, error)

	// Counters
	IncrementListeningCount(playlistUUID string) error
}

type playlistService struct {
	repo      playlist.PlaylistRepository
	audioRepo ports.AudioRepository
}

func NewPlaylistService(repo playlist.PlaylistRepository, audioRepo ports.AudioRepository) PlaylistService {
	return &playlistService{repo: repo, audioRepo: audioRepo}
}

func (s *playlistService) Create(req *domain.CreatePlaylistRequest, authorType domain.AuthorType, authorID uint) (*domain.Playlist, error) {
	playlist := &domain.Playlist{
		UUID:          uuid.New().String(),
		Title:         req.Title,
		Description:   req.Description,
		AuthorType:    authorType,
		AuthorID:      authorID,
		ThumbnailFile: req.ThumbnailFile,
		IsPublic:      req.IsPublic,
	}

	if err := s.repo.Create(playlist); err != nil {
		return nil, err
	}

	return playlist, nil
}

func (s *playlistService) GetByUUID(uuid string) (*domain.PlaylistResponse, error) {
	playlist, err := s.repo.FindByUUID(uuid)
	if err != nil {
		return nil, errors.New(messages.MsgPlaylistNotFound)
	}

	// Get playlist audios
	playlistAudios, err := s.repo.GetPlaylistAudios(playlist.ID)
	if err != nil {
		return nil, err
	}

	// Calculate totals and extract audios
	var audios []domain.Audio
	totalDuration := 0
	for _, pa := range playlistAudios {
		if pa.Audio != nil {
			audios = append(audios, *pa.Audio)
			totalDuration += pa.Audio.DurationAudio
		}
	}

	playlist.TotalAudio = len(audios)
	playlist.TotalPlayingAudio = totalDuration

	return &domain.PlaylistResponse{
		Playlist: playlist,
		Audios:   audios,
	}, nil
}

func (s *playlistService) GetAll(page, limit int) (*domain.PlaylistListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	isPublic := true
	playlists, total, err := s.repo.FindAll(page, limit, &isPublic)
	if err != nil {
		return nil, err
	}

	// Calculate totals for each playlist
	for i := range playlists {
		audios, _ := s.repo.GetPlaylistAudios(playlists[i].ID)
		totalDuration := 0
		for _, pa := range audios {
			if pa.Audio != nil {
				totalDuration += pa.Audio.DurationAudio
			}
		}
		playlists[i].TotalAudio = len(audios)
		playlists[i].TotalPlayingAudio = totalDuration
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &domain.PlaylistListResponse{
		Playlists:  playlists,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *playlistService) GetMyPlaylists(authorType domain.AuthorType, authorID uint, page, limit int) (*domain.PlaylistListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	playlists, total, err := s.repo.FindByAuthor(authorType, authorID, page, limit)
	if err != nil {
		return nil, err
	}

	// Calculate totals for each playlist
	for i := range playlists {
		audios, _ := s.repo.GetPlaylistAudios(playlists[i].ID)
		totalDuration := 0
		for _, pa := range audios {
			if pa.Audio != nil {
				totalDuration += pa.Audio.DurationAudio
			}
		}
		playlists[i].TotalAudio = len(audios)
		playlists[i].TotalPlayingAudio = totalDuration
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &domain.PlaylistListResponse{
		Playlists:  playlists,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *playlistService) Update(uuid string, req *domain.UpdatePlaylistRequest, authorType domain.AuthorType, authorID uint) (*domain.Playlist, error) {
	playlist, err := s.repo.FindByUUID(uuid)
	if err != nil {
		return nil, errors.New(messages.MsgPlaylistNotFound)
	}

	// Check ownership
	if playlist.AuthorType != authorType || playlist.AuthorID != authorID {
		return nil, errors.New(messages.ErrUnauthorized)
	}

	if req.Title != "" {
		playlist.Title = req.Title
	}
	if req.Description != "" {
		playlist.Description = req.Description
	}
	if req.ThumbnailFile != "" {
		playlist.ThumbnailFile = req.ThumbnailFile
	}
	if req.IsPublic != nil {
		playlist.IsPublic = *req.IsPublic
	}

	if err := s.repo.Update(playlist); err != nil {
		return nil, err
	}

	return playlist, nil
}

func (s *playlistService) Delete(uuid string, authorType domain.AuthorType, authorID uint) error {
	playlist, err := s.repo.FindByUUID(uuid)
	if err != nil {
		return errors.New(messages.MsgPlaylistNotFound)
	}

	// Check ownership (admin can delete any)
	if authorType != domain.AuthorTypeAdmin {
		if playlist.AuthorType != authorType || playlist.AuthorID != authorID {
			return errors.New(messages.ErrUnauthorized)
		}
	}

	return s.repo.Delete(playlist.ID)
}

func (s *playlistService) AddAudio(playlistUUID, audioUUID string, position int, authorType domain.AuthorType, authorID uint) error {
	playlist, err := s.repo.FindByUUID(playlistUUID)
	if err != nil {
		return errors.New(messages.MsgPlaylistNotFound)
	}

	// Check ownership
	if playlist.AuthorType != authorType || playlist.AuthorID != authorID {
		return errors.New(messages.ErrUnauthorized)
	}

	audio, err := s.audioRepo.FindByUUID(audioUUID)
	if err != nil {
		return errors.New(messages.ErrAudioNotFound)
	}

	// Get current position if not specified
	if position == 0 {
		audios, _ := s.repo.GetPlaylistAudios(playlist.ID)
		position = len(audios) + 1
	}

	return s.repo.AddAudio(playlist.ID, audio.ID, position)
}

func (s *playlistService) RemoveAudio(playlistUUID, audioUUID string, authorType domain.AuthorType, authorID uint) error {
	playlist, err := s.repo.FindByUUID(playlistUUID)
	if err != nil {
		return errors.New(messages.MsgPlaylistNotFound)
	}

	// Check ownership
	if playlist.AuthorType != authorType || playlist.AuthorID != authorID {
		return errors.New(messages.ErrUnauthorized)
	}

	audio, err := s.audioRepo.FindByUUID(audioUUID)
	if err != nil {
		return errors.New(messages.ErrAudioNotFound)
	}

	return s.repo.RemoveAudio(playlist.ID, audio.ID)
}

func (s *playlistService) LikePlaylist(playlistUUID string, userID uint) error {
	playlist, err := s.repo.FindByUUID(playlistUUID)
	if err != nil {
		return errors.New(messages.MsgPlaylistNotFound)
	}

	// Check if already liked
	liked, _ := s.repo.IsLiked(userID, playlist.ID)
	if liked {
		return errors.New("playlist already liked")
	}

	return s.repo.AddLike(userID, playlist.ID)
}

func (s *playlistService) UnlikePlaylist(playlistUUID string, userID uint) error {
	playlist, err := s.repo.FindByUUID(playlistUUID)
	if err != nil {
		return errors.New(messages.MsgPlaylistNotFound)
	}

	return s.repo.RemoveLike(userID, playlist.ID)
}

func (s *playlistService) IsLiked(playlistUUID string, userID uint) (bool, error) {
	playlist, err := s.repo.FindByUUID(playlistUUID)
	if err != nil {
		return false, errors.New(messages.MsgPlaylistNotFound)
	}

	return s.repo.IsLiked(userID, playlist.ID)
}

func (s *playlistService) IncrementListeningCount(playlistUUID string) error {
	playlist, err := s.repo.FindByUUID(playlistUUID)
	if err != nil {
		return errors.New(messages.MsgPlaylistNotFound)
	}

	return s.repo.IncrementListeningCount(playlist.ID)
}
