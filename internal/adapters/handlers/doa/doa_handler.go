package doa

import (
	"net/http"
	"strconv"

	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"
	"khalif-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type DoaHandler struct {
	service ports.DoaService
}

func NewDoaHandler(service ports.DoaService) *DoaHandler {
	return &DoaHandler{service: service}
}

// Admin Handlers

func (h *DoaHandler) Create(c *gin.Context) {
	var req domain.CreateDoaRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle audio upload
	audioFile, err := c.FormFile("audio_file")
	if err == nil && audioFile != nil {
		result, err := utils.SaveAudioFile(audioFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req.AudioDoa = result.URL
	}

	doa, err := h.service.Create(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Doa created successfully",
		"data":    doa,
	})
}

func (h *DoaHandler) Update(c *gin.Context) {
	uuid := c.Param("id")
	var req domain.UpdateDoaRequest
	
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle audio upload
	audioFile, err := c.FormFile("audio_file")
	if err == nil && audioFile != nil {
		result, err := utils.SaveAudioFile(audioFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req.AudioDoa = result.URL
	}

	doa, err := h.service.Update(uuid, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Doa updated successfully",
		"data":    doa,
	})
}

func (h *DoaHandler) Delete(c *gin.Context) {
	uuid := c.Param("id")

	if err := h.service.Delete(uuid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Doa deleted successfully"})
}

// Public Handlers

func (h *DoaHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	response, err := h.service.GetAll(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *DoaHandler) GetByID(c *gin.Context) {
	uuid := c.Param("id")
	doa, err := h.service.GetByUUID(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": doa})
}

func (h *DoaHandler) GetByCategory(c *gin.Context) {
	category := c.Query("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category is required"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	response, err := h.service.GetByCategory(category, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *DoaHandler) GetByHadist(c *gin.Context) {
	hadistID := c.Query("hadist_id") // UUID of Hadist
	if hadistID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hadist_id is required"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	response, err := h.service.GetByHadist(hadistID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *DoaHandler) GetRandom(c *gin.Context) {
	doa, err := h.service.GetRandom()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": doa})
}

func (h *DoaHandler) IncrementListeningCount(c *gin.Context) {
	uuid := c.Param("id")
	if err := h.service.IncrementListeningCount(uuid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Listening count incremented"})
}

// User Handlers

func (h *DoaHandler) LikeDoa(c *gin.Context) {
	userID := c.GetUint("userID")
	uuid := c.Param("id")
	if err := h.service.LikeDoa(userID, uuid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Doa liked"})
}

func (h *DoaHandler) UnlikeDoa(c *gin.Context) {
	userID := c.GetUint("userID")
	uuid := c.Param("id")
	if err := h.service.UnlikeDoa(userID, uuid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Doa unliked"})
}

func (h *DoaHandler) IsLiked(c *gin.Context) {
	userID := c.GetUint("userID")
	uuid := c.Param("id")
	isLiked, err := h.service.IsLiked(userID, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"is_liked": isLiked})
}

func (h *DoaHandler) BookmarkDoa(c *gin.Context) {
	userID := c.GetUint("userID")
	uuid := c.Param("id")
	if err := h.service.BookmarkDoa(userID, uuid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Doa bookmarked"})
}

func (h *DoaHandler) UnbookmarkDoa(c *gin.Context) {
	userID := c.GetUint("userID")
	uuid := c.Param("id")
	if err := h.service.UnbookmarkDoa(userID, uuid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Doa unbookmarked"})
}

func (h *DoaHandler) IsBookmarked(c *gin.Context) {
	userID := c.GetUint("userID")
	uuid := c.Param("id")
	isBookmarked, err := h.service.IsBookmarked(userID, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"is_bookmarked": isBookmarked})
}
