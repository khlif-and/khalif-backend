package alquran

import (
	"net/http"

	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type AlquranHandler struct {
	service ports.AlquranService
}

func NewAlquranHandler(service ports.AlquranService) *AlquranHandler {
	return &AlquranHandler{service: service}
}

func (h *AlquranHandler) CreateEndpoint(c *gin.Context) {
	var req domain.CreateAlquranRequest

	// Bind JSON body to struct
	// This directly expects a JSON payload matching CreateAlquranRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle Audio File Upload is removed in favor of JSON-only payload with audio_url string.
	// User requested to prioritize JSON structure. 
	// If file upload is needed, it should be done via a separate upload endpoint
	// that returns the URL, which is then passed here.

	// Call Service
	if err := h.service.CreateAlquran(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Alquran data created successfully",
	})
}

func (h *AlquranHandler) GetAllEndpont(c *gin.Context) {
	surahs, err := h.service.GetAllSurah(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": surahs,
	})
}

func (h *AlquranHandler) GetByIDEndpoint(c *gin.Context) {
	id := c.Param("id")
	surah, err := h.service.GetSurahByUUID(c.Request.Context(), id)
	if err != nil {
		// Differentiate error types if needed, for generic "not found" use 404
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": surah,
	})
}
