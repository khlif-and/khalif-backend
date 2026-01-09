package search

import (
	"encoding/json"

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
	IndexPlaylists      = "playlists"
	IndexDoas           = "doas"
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

	// Create playlists index
	_, err = m.client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        IndexPlaylists,
		PrimaryKey: "id",
	})
	if err != nil {
		logger.Log.Debug("Playlists index may already exist", zap.Error(err))
	}

	// Create doas index
	_, err = m.client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        IndexDoas,
		PrimaryKey: "id",
	})
	if err != nil {
		logger.Log.Debug("Doas index may already exist", zap.Error(err))
	}

	// // Configure filterable attributes for playlists
	// m.client.Index(IndexPlaylists).UpdateFilterableAttributes(&[]string{
	// 	"id",
	// })
	
	// Configure searchable attributes for doas
	m.client.Index(IndexDoas).UpdateSearchableAttributes(&[]string{
		"judul_doa",
		"arabic_doa",
		"latin_doa",
		"translate_doa",
		"description_doa",
		"tags",
		"category_doa",
	})

	// // Configure filterable attributes for doas
	// m.client.Index(IndexDoas).UpdateFilterableAttributes(&[]string{
	// 	"id",
	// 	"category_doa",
	// 	"tags",
	// })

	// Configure searchable attributes for playlists
	m.client.Index(IndexPlaylists).UpdateSearchableAttributes(&[]string{
		"title",
		"description",
		"author_name",
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

// PlaylistDocument represents a playlist document in Meilisearch
type PlaylistDocument struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	AuthorName     string `json:"author_name"`
	ThumbnailFile  string `json:"thumbnail_file"`
	LikeCount      int64  `json:"like_count"`
	ListeningCount int64  `json:"listening_count"`
	TotalAudio     int    `json:"total_audio"`
}

// IndexAudio adds or updates an audio document
func (m *MeilisearchClient) IndexAudio(doc AudioDocument) error {
	_, err := m.client.Index(IndexAudios).AddDocuments([]AudioDocument{doc}, nil)
	return err
}

// IndexUstadz adds or updates an ustadz document
func (m *MeilisearchClient) IndexUstadz(doc UstadzDocument) error {
	_, err := m.client.Index(IndexUstadzs).AddDocuments([]UstadzDocument{doc}, nil)
	return err
}

// IndexMoodCategory adds or updates a mood category document
func (m *MeilisearchClient) IndexMoodCategory(doc MoodCategoryDocument) error {
	_, err := m.client.Index(IndexMoodCategories).AddDocuments([]MoodCategoryDocument{doc}, nil)
	return err
}

// DeleteAudio removes an audio document
func (m *MeilisearchClient) DeleteAudio(id string) error {
	_, err := m.client.Index(IndexAudios).DeleteDocument(id, nil)
	return err
}

// DeleteUstadz removes an ustadz document
func (m *MeilisearchClient) DeleteUstadz(id string) error {
	_, err := m.client.Index(IndexUstadzs).DeleteDocument(id, nil)
	return err
}

// DeleteMoodCategory removes a mood category document
func (m *MeilisearchClient) DeleteMoodCategory(id string) error {
	_, err := m.client.Index(IndexMoodCategories).DeleteDocument(id, nil)
	return err
}

// IndexPlaylist adds or updates a playlist document
func (m *MeilisearchClient) IndexPlaylist(doc PlaylistDocument) error {
	_, err := m.client.Index(IndexPlaylists).AddDocuments([]PlaylistDocument{doc}, nil)
	return err
}

// DeletePlaylist removes a playlist document
func (m *MeilisearchClient) DeletePlaylist(id string) error {
	_, err := m.client.Index(IndexPlaylists).DeleteDocument(id, nil)
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
		results = append(results, mapToAudioDocument(hit))
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
		results = append(results, mapToUstadzDocument(hit))
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
		results = append(results, mapToMoodCategoryDocument(hit))
	}
	return results, nil
}

// SearchPlaylists performs search on playlists index
func (m *MeilisearchClient) SearchPlaylists(query string, limit int64) ([]PlaylistDocument, error) {
	searchRes, err := m.client.Index(IndexPlaylists).Search(query, &meilisearch.SearchRequest{
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}

	var results []PlaylistDocument
	hits, _ := json.Marshal(searchRes.Hits)
	json.Unmarshal(hits, &results)

	return results, nil
}

// SearchDoas performs search on doas index
func (m *MeilisearchClient) SearchDoas(query string, limit int64) ([]DoaDocument, error) {
	searchRes, err := m.client.Index(IndexDoas).Search(query, &meilisearch.SearchRequest{
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}

	var results []DoaDocument
	hits, _ := json.Marshal(searchRes.Hits)
	json.Unmarshal(hits, &results)

	return results, nil
}

// Helper functions to map search results to structs
func mapToAudioDocument(hit meilisearch.Hit) AudioDocument {
	var doc AudioDocument
	// Hit is map[string]json.RawMessage, need to convert
	data := make(map[string]interface{})
	for k, v := range hit {
		var val interface{}
		if err := json.Unmarshal(v, &val); err == nil {
			data[k] = val
		}
	}
	doc.ID = getStringFromMap(data, "id")
	doc.Title = getStringFromMap(data, "title")
	doc.UstadzName = getStringFromMap(data, "ustadz_name")
	doc.MoodCategoryName = getStringFromMap(data, "mood_category_name")
	doc.ThumbnailFile = getStringFromMap(data, "thumbnail_file")
	doc.AudioFile = getStringFromMap(data, "audio_file")
	doc.ListeningCount = getInt64FromMap(data, "listening_count")
	doc.LikeCount = getInt64FromMap(data, "like_count")
	doc.DurationAudio = int(getInt64FromMap(data, "duration_audio"))
	return doc
}

func mapToUstadzDocument(hit meilisearch.Hit) UstadzDocument {
	data := make(map[string]interface{})
	for k, v := range hit {
		var val interface{}
		if err := json.Unmarshal(v, &val); err == nil {
			data[k] = val
		}
	}
	return UstadzDocument{
		ID:            getStringFromMap(data, "id"),
		Name:          getStringFromMap(data, "name"),
		Description:   getStringFromMap(data, "description"),
		WikipediaLink: getStringFromMap(data, "wikipedia_link"),
	}
}

func mapToMoodCategoryDocument(hit meilisearch.Hit) MoodCategoryDocument {
	data := make(map[string]interface{})
	for k, v := range hit {
		var val interface{}
		if err := json.Unmarshal(v, &val); err == nil {
			data[k] = val
		}
	}
	return MoodCategoryDocument{
		ID:    getStringFromMap(data, "id"),
		Name:  getStringFromMap(data, "name"),
		Icon:  getStringFromMap(data, "icon"),
		Color: getStringFromMap(data, "color"),
	}
}

func mapToPlaylistDocument(hit meilisearch.Hit) PlaylistDocument {
	data := make(map[string]interface{})
	for k, v := range hit {
		var val interface{}
		if err := json.Unmarshal(v, &val); err == nil {
			data[k] = val
		}
	}
	return PlaylistDocument{
		ID:             getStringFromMap(data, "id"),
		Title:          getStringFromMap(data, "title"),
		Description:    getStringFromMap(data, "description"),
		AuthorName:     getStringFromMap(data, "author_name"),
		ThumbnailFile:  getStringFromMap(data, "thumbnail_file"),
		LikeCount:      getInt64FromMap(data, "like_count"),
		ListeningCount: getInt64FromMap(data, "listening_count"),
		TotalAudio:     int(getInt64FromMap(data, "total_audio")),
	}
}

func getStringFromMap(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

func getInt64FromMap(data map[string]interface{}, key string) int64 {
	if val, ok := data[key].(float64); ok {
		return int64(val)
	}
	return 0
}
