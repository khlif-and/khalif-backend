package search

import (
	"khalif-backend/internal/infrastructure/search"
	"khalif-backend/internal/platform/logger"

	"go.uber.org/zap"
)

// SearchService handles search operations
type SearchService struct {
	meili *search.MeilisearchClient
}

// NewSearchService creates a new search service
func NewSearchService(meili *search.MeilisearchClient) *SearchService {
	return &SearchService{meili: meili}
}

// UnifiedSearchResult contains results from all indices
type UnifiedSearchResult struct {
	Audios         []search.AudioDocument        `json:"audios"`
	Ustadzs        []search.UstadzDocument       `json:"ustadzs"`
	MoodCategories []search.MoodCategoryDocument `json:"mood_categories"`
	Playlists      []search.PlaylistDocument     `json:"playlists"`
}

// SearchAll performs a unified search across all indices
func (s *SearchService) SearchAll(query string, limit int64) (*UnifiedSearchResult, error) {
	if limit <= 0 {
		limit = 10
	}

	result := &UnifiedSearchResult{}

	// Search audios
	audios, err := s.meili.SearchAudios(query, limit)
	if err != nil {
		logger.Log.Error("Failed to search audios", zap.Error(err))
	} else {
		result.Audios = audios
	}

	// Search ustadzs
	ustadzs, err := s.meili.SearchUstadzs(query, limit)
	if err != nil {
		logger.Log.Error("Failed to search ustadzs", zap.Error(err))
	} else {
		result.Ustadzs = ustadzs
	}

	// Search mood categories
	moods, err := s.meili.SearchMoodCategories(query, limit)
	if err != nil {
		logger.Log.Error("Failed to search mood categories", zap.Error(err))
	} else {
		result.MoodCategories = moods
	}

	// Search playlists
	playlists, err := s.meili.SearchPlaylists(query, limit)
	if err != nil {
		logger.Log.Error("Failed to search playlists", zap.Error(err))
	} else {
		result.Playlists = playlists
	}

	return result, nil
}

// SearchAudios searches only the audios index
func (s *SearchService) SearchAudios(query string, limit int64) ([]search.AudioDocument, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.meili.SearchAudios(query, limit)
}

// SearchUstadzs searches only the ustadzs index
func (s *SearchService) SearchUstadzs(query string, limit int64) ([]search.UstadzDocument, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.meili.SearchUstadzs(query, limit)
}

// SearchMoodCategories searches only the mood categories index
func (s *SearchService) SearchMoodCategories(query string, limit int64) ([]search.MoodCategoryDocument, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.meili.SearchMoodCategories(query, limit)
}

// SearchPlaylists searches only the playlists index
func (s *SearchService) SearchPlaylists(query string, limit int64) ([]search.PlaylistDocument, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.meili.SearchPlaylists(query, limit)
}
