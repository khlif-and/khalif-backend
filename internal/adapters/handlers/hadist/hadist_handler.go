package hadist

import (
	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"
	"khalif-backend/pkg/messages"
	"khalif-backend/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HadistHandler struct {
	service ports.HadistService
}

func NewHadistHandler(service ports.HadistService) *HadistHandler {
	return &HadistHandler{service: service}
}

// Admin Handlers

func (h *HadistHandler) Create(c *gin.Context) {
	var req domain.CreateHadistRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	audioFile, err := c.FormFile("audio_file")
	if err == nil && audioFile != nil {
		result, err := utils.SaveAudioFile(audioFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req.AudioHadist = result.URL
	}

	hadist, err := h.service.Create(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": messages.MsgHadistCreated,
		"data":    hadist,
	})
}

func (h *HadistHandler) Update(c *gin.Context) {
	uuid := c.Param("id")
	var req domain.UpdateHadistRequest
	
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	audioFile, err := c.FormFile("audio_file")
	if err == nil && audioFile != nil {
		result, err := utils.SaveAudioFile(audioFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req.AudioHadist = result.URL
	}

	hadist, err := h.service.Update(uuid, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": messages.MsgHadistUpdated,
		"data":    hadist,
	})
}

func (h *HadistHandler) Delete(c *gin.Context) {
	uuid := c.Param("id")

	if err := h.service.Delete(uuid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": messages.MsgHadistDeleted})
}

// Public Handlers

func (h *HadistHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	response, err := h.service.GetAll(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrInternalServer})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *HadistHandler) GetByID(c *gin.Context) {
	uuid := c.Param("id")

	hadist, err := h.service.GetByUUID(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": messages.ErrHadistNotFound})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": hadist})
}

func (h *HadistHandler) GetByCategory(c *gin.Context) {
	category := c.Query("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrMissingFields})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	response, err := h.service.GetByCategory(category, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrInternalServer})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *HadistHandler) GetByKitab(c *gin.Context) {
	kitab := c.Query("kitab")
	if kitab == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrMissingFields})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	response, err := h.service.GetByKitab(kitab, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrInternalServer})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *HadistHandler) GetRandom(c *gin.Context) {
	hadist, err := h.service.GetRandom()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": messages.ErrHadistNotFound})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": hadist})
}

func (h *HadistHandler) IncrementListeningCount(c *gin.Context) {
	uuid := c.Param("id")

	if err := h.service.IncrementListeningCount(uuid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": messages.MsgListeningCountIncremented})
}
