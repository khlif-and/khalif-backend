package playlist

import (
	"khalif-backend/internal/core/domain"

	"gorm.io/gorm"
)

type PlaylistRepository interface {
	Create(playlist *domain.Playlist) error
	FindByID(id uint) (*domain.Playlist, error)
	FindByUUID(uuid string) (*domain.Playlist, error)
	FindAll(page, limit int, isPublic *bool) ([]domain.Playlist, int64, error)
	FindByAuthor(authorType domain.AuthorType, authorID uint, page, limit int) ([]domain.Playlist, int64, error)
	Update(playlist *domain.Playlist) error
	Delete(id uint) error

	// Audio management
	AddAudio(playlistID, audioID uint, position int) error
	RemoveAudio(playlistID, audioID uint) error
	GetPlaylistAudios(playlistID uint) ([]domain.PlaylistAudio, error)
	UpdateAudioPosition(playlistID, audioID uint, position int) error

	// Likes
	AddLike(userID, playlistID uint) error
	RemoveLike(userID, playlistID uint) error
	IsLiked(userID, playlistID uint) (bool, error)

	// Counters
	IncrementListeningCount(playlistID uint) error
}

type playlistRepo struct {
	db *gorm.DB
}

func NewPlaylistRepo(db *gorm.DB) PlaylistRepository {
	return &playlistRepo{db: db}
}

func (r *playlistRepo) Create(playlist *domain.Playlist) error {
	return r.db.Create(playlist).Error
}

func (r *playlistRepo) FindByID(id uint) (*domain.Playlist, error) {
	var playlist domain.Playlist
	err := r.db.First(&playlist, id).Error
	if err != nil {
		return nil, err
	}
	return &playlist, nil
}

func (r *playlistRepo) FindByUUID(uuid string) (*domain.Playlist, error) {
	var playlist domain.Playlist
	err := r.db.Where("uuid = ?", uuid).First(&playlist).Error
	if err != nil {
		return nil, err
	}
	return &playlist, nil
}

func (r *playlistRepo) FindAll(page, limit int, isPublic *bool) ([]domain.Playlist, int64, error) {
	var playlists []domain.Playlist
	var total int64

	query := r.db.Model(&domain.Playlist{})
	if isPublic != nil {
		query = query.Where("is_public = ?", *isPublic)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&playlists).Error
	if err != nil {
		return nil, 0, err
	}

	return playlists, total, nil
}

func (r *playlistRepo) FindByAuthor(authorType domain.AuthorType, authorID uint, page, limit int) ([]domain.Playlist, int64, error) {
	var playlists []domain.Playlist
	var total int64

	query := r.db.Model(&domain.Playlist{}).Where("author_type = ? AND author_id = ?", authorType, authorID)
	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&playlists).Error
	if err != nil {
		return nil, 0, err
	}

	return playlists, total, nil
}

func (r *playlistRepo) Update(playlist *domain.Playlist) error {
	return r.db.Save(playlist).Error
}

func (r *playlistRepo) Delete(id uint) error {
	return r.db.Delete(&domain.Playlist{}, id).Error
}

// Audio management
func (r *playlistRepo) AddAudio(playlistID, audioID uint, position int) error {
	pa := domain.PlaylistAudio{
		PlaylistID: playlistID,
		AudioID:    audioID,
		Position:   position,
	}
	return r.db.Create(&pa).Error
}

func (r *playlistRepo) RemoveAudio(playlistID, audioID uint) error {
	return r.db.Where("playlist_id = ? AND audio_id = ?", playlistID, audioID).Delete(&domain.PlaylistAudio{}).Error
}

func (r *playlistRepo) GetPlaylistAudios(playlistID uint) ([]domain.PlaylistAudio, error) {
	var audios []domain.PlaylistAudio
	err := r.db.Where("playlist_id = ?", playlistID).
		Preload("Audio").
		Preload("Audio.Ustadz").
		Preload("Audio.MoodCategory").
		Order("position ASC").
		Find(&audios).Error
	if err != nil {
		return nil, err
	}
	return audios, nil
}

func (r *playlistRepo) UpdateAudioPosition(playlistID, audioID uint, position int) error {
	return r.db.Model(&domain.PlaylistAudio{}).
		Where("playlist_id = ? AND audio_id = ?", playlistID, audioID).
		Update("position", position).Error
}

// Likes
func (r *playlistRepo) AddLike(userID, playlistID uint) error {
	like := domain.PlaylistLike{
		UserID:     userID,
		PlaylistID: playlistID,
	}
	err := r.db.Create(&like).Error
	if err != nil {
		return err
	}
	// Increment like count
	return r.db.Model(&domain.Playlist{}).Where("id = ?", playlistID).
		UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}

func (r *playlistRepo) RemoveLike(userID, playlistID uint) error {
	err := r.db.Where("user_id = ? AND playlist_id = ?", userID, playlistID).Delete(&domain.PlaylistLike{}).Error
	if err != nil {
		return err
	}
	// Decrement like count
	return r.db.Model(&domain.Playlist{}).Where("id = ? AND like_count > 0", playlistID).
		UpdateColumn("like_count", gorm.Expr("like_count - 1")).Error
}

func (r *playlistRepo) IsLiked(userID, playlistID uint) (bool, error) {
	var count int64
	err := r.db.Model(&domain.PlaylistLike{}).Where("user_id = ? AND playlist_id = ?", userID, playlistID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Counters
func (r *playlistRepo) IncrementListeningCount(playlistID uint) error {
	return r.db.Model(&domain.Playlist{}).Where("id = ?", playlistID).
		UpdateColumn("listening_count", gorm.Expr("listening_count + 1")).Error
}
