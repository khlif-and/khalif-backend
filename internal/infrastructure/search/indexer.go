package search

import (
	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/platform/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Indexer syncs data from PostgreSQL to Meilisearch
type Indexer struct {
	db    *gorm.DB
	meili *MeilisearchClient
}

// NewIndexer creates a new indexer
func NewIndexer(db *gorm.DB, meili *MeilisearchClient) *Indexer {
	return &Indexer{db: db, meili: meili}
}

// SyncAll syncs all data from PostgreSQL to Meilisearch
func (i *Indexer) SyncAll() error {
	logger.Log.Info("Starting full index sync...")

	if err := i.SyncAudios(); err != nil {
		logger.Log.Error("Failed to sync audios", zap.Error(err))
	}

	if err := i.SyncUstadzs(); err != nil {
		logger.Log.Error("Failed to sync ustadzs", zap.Error(err))
	}

	if err := i.SyncMoodCategories(); err != nil {
		logger.Log.Error("Failed to sync mood categories", zap.Error(err))
	}

	if err := i.SyncPlaylists(); err != nil {
		logger.Log.Error("Failed to sync playlists", zap.Error(err))
	}

	logger.Log.Info("Full index sync completed")
	return nil
}

// SyncAudios syncs all audios to Meilisearch
func (i *Indexer) SyncAudios() error {
	var audios []domain.Audio
	if err := i.db.Preload("Ustadz").Preload("MoodCategory").Find(&audios).Error; err != nil {
		return err
	}

	var docs []AudioDocument
	for _, audio := range audios {
		doc := AudioDocument{
			ID:             audio.UUID,
			Title:          audio.Title,
			ThumbnailFile:  audio.ThumbnailFile,
			AudioFile:      audio.AudioFile,
			ListeningCount: audio.ListeningCount,
			LikeCount:      audio.LikeCount,
			DurationAudio:  audio.DurationAudio,
		}

		if audio.Ustadz != nil {
			doc.UstadzName = audio.Ustadz.Name
		}
		if audio.MoodCategory != nil {
			doc.MoodCategoryName = audio.MoodCategory.Name
		}

		docs = append(docs, doc)
	}

	if len(docs) > 0 {
		_, err := i.meili.GetClient().Index(IndexAudios).AddDocuments(docs, nil)
		if err != nil {
			return err
		}
	}

	logger.Log.Info("Synced audios to Meilisearch", zap.Int("count", len(docs)))
	return nil
}

// SyncUstadzs syncs all ustadzs to Meilisearch
func (i *Indexer) SyncUstadzs() error {
	var ustadzs []domain.Ustadz
	if err := i.db.Find(&ustadzs).Error; err != nil {
		return err
	}

	var docs []UstadzDocument
	for _, ustadz := range ustadzs {
		docs = append(docs, UstadzDocument{
			ID:            ustadz.UUID,
			Name:          ustadz.Name,
			Description:   ustadz.Description,
			WikipediaLink: ustadz.WikipediaLink,
		})
	}

	if len(docs) > 0 {
		_, err := i.meili.GetClient().Index(IndexUstadzs).AddDocuments(docs, nil)
		if err != nil {
			return err
		}
	}

	logger.Log.Info("Synced ustadzs to Meilisearch", zap.Int("count", len(docs)))
	return nil
}

// SyncMoodCategories syncs all mood categories to Meilisearch
func (i *Indexer) SyncMoodCategories() error {
	var moods []domain.MoodCategory
	if err := i.db.Find(&moods).Error; err != nil {
		return err
	}

	var docs []MoodCategoryDocument
	for _, mood := range moods {
		docs = append(docs, MoodCategoryDocument{
			ID:    mood.UUID,
			Name:  mood.Name,
			Icon:  mood.Icon,
			Color: mood.Color,
		})
	}

	if len(docs) > 0 {
		_, err := i.meili.GetClient().Index(IndexMoodCategories).AddDocuments(docs, nil)
		if err != nil {
			return err
		}
	}

	logger.Log.Info("Synced mood categories to Meilisearch", zap.Int("count", len(docs)))
	return nil
}

// IndexAudioFromDomain indexes a single audio from domain model
func (i *Indexer) IndexAudioFromDomain(audio *domain.Audio) error {
	doc := AudioDocument{
		ID:             audio.UUID,
		Title:          audio.Title,
		ThumbnailFile:  audio.ThumbnailFile,
		AudioFile:      audio.AudioFile,
		ListeningCount: audio.ListeningCount,
		LikeCount:      audio.LikeCount,
		DurationAudio:  audio.DurationAudio,
	}

	if audio.Ustadz != nil {
		doc.UstadzName = audio.Ustadz.Name
	}
	if audio.MoodCategory != nil {
		doc.MoodCategoryName = audio.MoodCategory.Name
	}

	return i.meili.IndexAudio(doc)
}

// IndexUstadzFromDomain indexes a single ustadz from domain model
func (i *Indexer) IndexUstadzFromDomain(ustadz *domain.Ustadz) error {
	doc := UstadzDocument{
		ID:            ustadz.UUID,
		Name:          ustadz.Name,
		Description:   ustadz.Description,
		WikipediaLink: ustadz.WikipediaLink,
	}
	return i.meili.IndexUstadz(doc)
}

// IndexMoodCategoryFromDomain indexes a single mood category from domain model
func (i *Indexer) IndexMoodCategoryFromDomain(mood *domain.MoodCategory) error {
	doc := MoodCategoryDocument{
		ID:    mood.UUID,
		Name:  mood.Name,
		Icon:  mood.Icon,
		Color: mood.Color,
	}
	return i.meili.IndexMoodCategory(doc)
}

// SyncPlaylists syncs all playlists to Meilisearch
func (i *Indexer) SyncPlaylists() error {
	var playlists []domain.Playlist
	if err := i.db.Where("is_public = ?", true).Find(&playlists).Error; err != nil {
		return err
	}

	var docs []PlaylistDocument
	for _, playlist := range playlists {
		// Get audio count
		var audioCount int64
		i.db.Model(&domain.PlaylistAudio{}).Where("playlist_id = ?", playlist.ID).Count(&audioCount)

		docs = append(docs, PlaylistDocument{
			ID:             playlist.UUID,
			Title:          playlist.Title,
			Description:    playlist.Description,
			AuthorName:     playlist.AuthorName,
			ThumbnailFile:  playlist.ThumbnailFile,
			LikeCount:      playlist.LikeCount,
			ListeningCount: playlist.ListeningCount,
			TotalAudio:     int(audioCount),
		})
	}

	if len(docs) > 0 {
		_, err := i.meili.GetClient().Index(IndexPlaylists).AddDocuments(docs, nil)
		if err != nil {
			return err
		}
	}

	logger.Log.Info("Synced playlists to Meilisearch", zap.Int("count", len(docs)))
	return nil
}

// IndexPlaylistFromDomain indexes a single playlist from domain model
func (i *Indexer) IndexPlaylistFromDomain(playlist *domain.Playlist) error {
	// Get audio count
	var audioCount int64
	i.db.Model(&domain.PlaylistAudio{}).Where("playlist_id = ?", playlist.ID).Count(&audioCount)

	doc := PlaylistDocument{
		ID:             playlist.UUID,
		Title:          playlist.Title,
		Description:    playlist.Description,
		AuthorName:     playlist.AuthorName,
		ThumbnailFile:  playlist.ThumbnailFile,
		LikeCount:      playlist.LikeCount,
		ListeningCount: playlist.ListeningCount,
		TotalAudio:     int(audioCount),
	}
	return i.meili.IndexPlaylist(doc)
}
