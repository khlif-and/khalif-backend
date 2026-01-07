package mood_category

import (
	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"
	"khalif-backend/pkg/messages"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type MoodCategoryHandler struct {
	service ports.MoodCategoryService
}

func NewMoodCategoryHandler(service ports.MoodCategoryService) *MoodCategoryHandler {
	return &MoodCategoryHandler{service: service}
}

func (h *MoodCategoryHandler) Create(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	icon := strings.TrimSpace(c.PostForm("icon"))
	color := strings.TrimSpace(c.PostForm("color"))

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrMoodNameRequired})
		return
	}

	req := &domain.CreateMoodCategoryRequest{
		Name:  name,
		Icon:  icon,
		Color: color,
	}

	mood, err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": messages.MsgMoodCreated,
		"data":    mood,
	})
}

func (h *MoodCategoryHandler) GetAll(c *gin.Context) {
	response, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrInternalServer})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (h *MoodCategoryHandler) GetByID(c *gin.Context) {
	uuid := c.Param("id")

	mood, err := h.service.GetByUUID(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": messages.ErrMoodNotFound})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": mood,
	})
}

func (h *MoodCategoryHandler) GetAudiosByMoodID(c *gin.Context) {
	uuid := c.Param("id")

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	response, err := h.service.GetAudiosByMoodUUID(uuid, page, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (h *MoodCategoryHandler) Update(c *gin.Context) {
	uuid := c.Param("id")

	name := strings.TrimSpace(c.PostForm("name"))
	icon := strings.TrimSpace(c.PostForm("icon"))
	color := strings.TrimSpace(c.PostForm("color"))

	req := &domain.UpdateMoodCategoryRequest{
		Name:  name,
		Icon:  icon,
		Color: color,
	}

	mood, err := h.service.Update(uuid, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": messages.MsgMoodUpdated,
		"data":    mood,
	})
}

func (h *MoodCategoryHandler) Delete(c *gin.Context) {
	uuid := c.Param("id")

	if err := h.service.Delete(uuid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": messages.MsgMoodDeleted,
	})
}
