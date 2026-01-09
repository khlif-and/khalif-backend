package doa

import (
	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"
	"khalif-backend/pkg/messages"
	"khalif-backend/pkg/utils"
	"net/http"
	"strconv"

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
		"message": messages.MsgDoaCreated,
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
		"message": messages.MsgDoaUpdated,
		"data":    doa,
	})
}

func (h *DoaHandler) Delete(c *gin.Context) {
	uuid := c.Param("id")

	if err := h.service.Delete(uuid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": messages.MsgDoaDeleted})
}

// Public Handlers

func (h *DoaHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	response, err := h.service.GetAll(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrInternalServer})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *DoaHandler) GetByID(c *gin.Context) {
	uuid := c.Param("id")
	doa, err := h.service.GetByUUID(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": messages.ErrDoaNotFound})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": doa})
}

func (h *DoaHandler) GetByCategory(c *gin.Context) {
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

func (h *DoaHandler) GetByHadist(c *gin.Context) {
	hadistID := c.Query("hadist_id")
	if hadistID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrMissingFields})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	response, err := h.service.GetByHadist(hadistID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrInternalServer})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *DoaHandler) GetRandom(c *gin.Context) {
	doa, err := h.service.GetRandom()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": messages.ErrDoaNotFound})
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
	c.JSON(http.StatusOK, gin.H{"message": messages.MsgListeningCountIncremented})
}
