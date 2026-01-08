package search

import (
	"khalif-backend/internal/platform/config"
	"khalif-backend/internal/platform/logger"

	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
)

// Index names
const (
	IndexAudios         = "audios"
	IndexUstadzs        = "ustadzs"
	IndexMoodCategories = "mood_categories"
)

// MeilisearchClient wraps the Meilisearch client
type MeilisearchClient struct {
	client meilisearch.ServiceManager
}

// NewMeilisearchClient creates a new Meilisearch client
func NewMeilisearchClient(cfg *config.Config) *MeilisearchClient {
	client := meilisearch.New(cfg.MeilisearchHost, meilisearch.WithAPIKey(cfg.MeilisearchAPIKey))

	ms := &MeilisearchClient{client: client}

	// Initialize indices
	if err := ms.initializeIndices(); err != nil {
		logger.Log.Warn("Failed to initialize Meilisearch indices", zap.Error(err))
	} else {
		logger.Log.Info("Meilisearch indices initialized successfully")
	}

	return ms
}

// initializeIndices creates and configures all search indices
func (m *MeilisearchClient) initializeIndices() error {
	// Create audios index
	_, err := m.client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        IndexAudios,
		PrimaryKey: "id",
	})
	if err != nil {
		logger.Log.Debug("Audios index may already exist", zap.Error(err))
	}

	// Create ustadzs index
	_, err = m.client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        IndexUstadzs,
		PrimaryKey: "id",
	})
	if err != nil {
		logger.Log.Debug("Ustadzs index may already exist", zap.Error(err))
	}

	// Create mood_categories index
	_, err = m.client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        IndexMoodCategories,
		PrimaryKey: "id",
	})
	if err != nil {
		logger.Log.Debug("MoodCategories index may already exist", zap.Error(err))
	}

	// Configure searchable attributes for audios
	m.client.Index(IndexAudios).UpdateSearchableAttributes(&[]string{
		"title",
		"ustadz_name",
		"mood_category_name",
	})

	// Configure searchable attributes for ustadzs
	m.client.Index(IndexUstadzs).UpdateSearchableAttributes(&[]string{
		"name",
		"description",
	})

	// Configure searchable attributes for mood_categories
	m.client.Index(IndexMoodCategories).UpdateSearchableAttributes(&[]string{
		"name",
	})

	return nil
}

// GetClient returns the underlying Meilisearch client
func (m *MeilisearchClient) GetClient() meilisearch.ServiceManager {
	return m.client
}

// AudioDocument represents an audio document in Meilisearch
type AudioDocument struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	UstadzName       string `json:"ustadz_name"`
	MoodCategoryName string `json:"mood_category_name"`
	ThumbnailFile    string `json:"thumbnail_file"`
	AudioFile        string `json:"audio_file"`
	ListeningCount   int64  `json:"listening_count"`
	LikeCount        int64  `json:"like_count"`
	DurationAudio    int    `json:"duration_audio"`
}

// UstadzDocument represents an ustadz document in Meilisearch
type UstadzDocument struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	WikipediaLink string `json:"wikipedia_link"`
}

// MoodCategoryDocument represents a mood category document in Meilisearch
type MoodCategoryDocument struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
}

// IndexAudio adds or updates an audio document
func (m *MeilisearchClient) IndexAudio(doc AudioDocument) error {
	_, err := m.client.Index(IndexAudios).AddDocuments([]AudioDocument{doc}, "id")
	return err
}

// IndexUstadz adds or updates an ustadz document
func (m *MeilisearchClient) IndexUstadz(doc UstadzDocument) error {
	_, err := m.client.Index(IndexUstadzs).AddDocuments([]UstadzDocument{doc}, "id")
	return err
}

// IndexMoodCategory adds or updates a mood category document
func (m *MeilisearchClient) IndexMoodCategory(doc MoodCategoryDocument) error {
	_, err := m.client.Index(IndexMoodCategories).AddDocuments([]MoodCategoryDocument{doc}, "id")
	return err
}

// DeleteAudio removes an audio document
func (m *MeilisearchClient) DeleteAudio(id string) error {
	_, err := m.client.Index(IndexAudios).DeleteDocument(id)
	return err
}

// DeleteUstadz removes an ustadz document
func (m *MeilisearchClient) DeleteUstadz(id string) error {
	_, err := m.client.Index(IndexUstadzs).DeleteDocument(id)
	return err
}

// DeleteMoodCategory removes a mood category document
func (m *MeilisearchClient) DeleteMoodCategory(id string) error {
	_, err := m.client.Index(IndexMoodCategories).DeleteDocument(id)
	return err
}

// SearchAudios searches the audios index
func (m *MeilisearchClient) SearchAudios(query string, limit int64) ([]AudioDocument, error) {
	resp, err := m.client.Index(IndexAudios).Search(query, &meilisearch.SearchRequest{
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}

	var results []AudioDocument
	for _, hit := range resp.Hits {
		if doc, ok := hit.(map[string]interface{}); ok {
			results = append(results, mapToAudioDocument(doc))
		}
	}
	return results, nil
}

// SearchUstadzs searches the ustadzs index
func (m *MeilisearchClient) SearchUstadzs(query string, limit int64) ([]UstadzDocument, error) {
	resp, err := m.client.Index(IndexUstadzs).Search(query, &meilisearch.SearchRequest{
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}

	var results []UstadzDocument
	for _, hit := range resp.Hits {
		if doc, ok := hit.(map[string]interface{}); ok {
			results = append(results, mapToUstadzDocument(doc))
		}
	}
	return results, nil
}

// SearchMoodCategories searches the mood_categories index
func (m *MeilisearchClient) SearchMoodCategories(query string, limit int64) ([]MoodCategoryDocument, error) {
	resp, err := m.client.Index(IndexMoodCategories).Search(query, &meilisearch.SearchRequest{
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}

	var results []MoodCategoryDocument
	for _, hit := range resp.Hits {
		if doc, ok := hit.(map[string]interface{}); ok {
			results = append(results, mapToMoodCategoryDocument(doc))
		}
	}
	return results, nil
}

// Helper functions to map search results to structs
func mapToAudioDocument(doc map[string]interface{}) AudioDocument {
	return AudioDocument{
		ID:               getString(doc, "id"),
		Title:            getString(doc, "title"),
		UstadzName:       getString(doc, "ustadz_name"),
		MoodCategoryName: getString(doc, "mood_category_name"),
		ThumbnailFile:    getString(doc, "thumbnail_file"),
		AudioFile:        getString(doc, "audio_file"),
		ListeningCount:   getInt64(doc, "listening_count"),
		LikeCount:        getInt64(doc, "like_count"),
		DurationAudio:    int(getInt64(doc, "duration_audio")),
	}
}

func mapToUstadzDocument(doc map[string]interface{}) UstadzDocument {
	return UstadzDocument{
		ID:            getString(doc, "id"),
		Name:          getString(doc, "name"),
		Description:   getString(doc, "description"),
		WikipediaLink: getString(doc, "wikipedia_link"),
	}
}

func mapToMoodCategoryDocument(doc map[string]interface{}) MoodCategoryDocument {
	return MoodCategoryDocument{
		ID:    getString(doc, "id"),
		Name:  getString(doc, "name"),
		Icon:  getString(doc, "icon"),
		Color: getString(doc, "color"),
	}
}

func getString(doc map[string]interface{}, key string) string {
	if val, ok := doc[key].(string); ok {
		return val
	}
	return ""
}

func getInt64(doc map[string]interface{}, key string) int64 {
	if val, ok := doc[key].(float64); ok {
		return int64(val)
	}
	return 0
}
