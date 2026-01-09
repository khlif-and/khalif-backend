package search

import (
	"net/http"
	"strconv"

	searchService "khalif-backend/internal/core/services/search"

	"github.com/gin-gonic/gin"
)

// SearchHandler handles search requests
type SearchHandler struct {
	service *searchService.SearchService
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(service *searchService.SearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

// SearchAll handles unified search across all indices
// @Summary Search all
// @Description Search across audios, ustadzs, and mood categories
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Result limit per category" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/search [get]
func (h *SearchHandler) SearchAll(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	result, err := h.service.SearchAll(query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query":   query,
		"results": result,
	})
}

// SearchAudios handles audio-specific search
// @Summary Search audios
// @Description Search only in audios index
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Result limit" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/search/audio [get]
func (h *SearchHandler) SearchAudios(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	results, err := h.service.SearchAudios(query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query":   query,
		"results": results,
		"total":   len(results),
	})
}

// SearchUstadzs handles ustadz-specific search
// @Summary Search ustadzs
// @Description Search only in ustadzs index
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Result limit" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/search/ustadz [get]
func (h *SearchHandler) SearchUstadzs(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	results, err := h.service.SearchUstadzs(query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query":   query,
		"results": results,
		"total":   len(results),
	})
}

// SearchMoodCategories handles mood category-specific search
// @Summary Search mood categories
// @Description Search only in mood categories index
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Result limit" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/search/mood [get]
func (h *SearchHandler) SearchMoodCategories(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	results, err := h.service.SearchMoodCategories(query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query":   query,
		"results": results,
		"total":   len(results),
	})
}

// SearchPlaylists handles playlist-specific search
// @Summary Search playlists
// @Description Search only in playlists index
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Result limit" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/search/playlist [get]
func (h *SearchHandler) SearchPlaylists(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	results, err := h.service.SearchPlaylists(query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query":   query,
		"results": results,
		"total":   len(results),
	})
}
