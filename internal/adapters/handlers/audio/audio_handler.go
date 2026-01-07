package audio

import (
	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"
	"khalif-backend/pkg/messages"
	"khalif-backend/pkg/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AudioHandler struct {
	service ports.AudioService
}

func NewAudioHandler(service ports.AudioService) *AudioHandler {
	return &AudioHandler{service: service}
}

func (h *AudioHandler) Create(c *gin.Context) {
	title := strings.TrimSpace(c.PostForm("title"))
	ustadzUUID := strings.TrimSpace(c.PostForm("ustadz_id"))
	moodCategoryUUID := strings.TrimSpace(c.PostForm("mood_category_id"))

	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrAudioTitleRequired})
		return
	}

	var audioURL string
	var duration int
	audioFile, err := c.FormFile("audio_file")
	if err == nil && audioFile != nil {
		result, err := utils.SaveAudioFile(audioFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		audioURL = result.URL

		extractedDuration, err := utils.GetAudioDuration(result.Path)
		if err == nil {
			duration = extractedDuration
		}
	}

	if audioURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrAudioFileRequired})
		return
	}

	var thumbnailURL string
	var finalColor string
	thumbnailFile, err := c.FormFile("thumbnail_file")
	if err == nil && thumbnailFile != nil {
		result, err := utils.SaveUploadedFile(thumbnailFile, utils.ThumbnailDir)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		thumbnailURL = result.URL
		finalColor = utils.GetThumbnailColorPath(result.Path)
	}

	req := &domain.CreateAudioRequest{
		Title:               title,
		AudioFile:           audioURL,
		ThumbnailFile:       thumbnailURL,
		ColorThumbnailAudio: finalColor,
		UstadzUUID:          ustadzUUID,
		DurationAudio:       duration,
		MoodCategoryUUID:    moodCategoryUUID,
	}

	audio, err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": messages.MsgAudioCreated,
		"data":    audio,
	})
}

func (h *AudioHandler) GetAll(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	response, err := h.service.GetAll(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrInternalServer})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (h *AudioHandler) GetByID(c *gin.Context) {
	uuid := c.Param("id")

	audio, err := h.service.GetByUUID(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": messages.ErrAudioNotFound})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": audio,
	})
}

func (h *AudioHandler) Update(c *gin.Context) {
	uuid := c.Param("id")

	title := strings.TrimSpace(c.PostForm("title"))
	ustadzUUID := strings.TrimSpace(c.PostForm("ustadz_id"))
	moodCategoryUUID := strings.TrimSpace(c.PostForm("mood_category_id"))

	var audioURL string
	var duration int
	audioFile, err := c.FormFile("audio_file")
	if err == nil && audioFile != nil {
		result, err := utils.SaveAudioFile(audioFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		audioURL = result.URL

		extractedDuration, err := utils.GetAudioDuration(result.Path)
		if err == nil {
			duration = extractedDuration
		}
	}

	var thumbnailURL string
	var finalColor string
	thumbnailFile, err := c.FormFile("thumbnail_file")
	if err == nil && thumbnailFile != nil {
		result, err := utils.SaveUploadedFile(thumbnailFile, utils.ThumbnailDir)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		thumbnailURL = result.URL
		finalColor = utils.GetThumbnailColorPath(result.Path)
	}

	req := &domain.UpdateAudioRequest{
		Title:               title,
		AudioFile:           audioURL,
		ThumbnailFile:       thumbnailURL,
		ColorThumbnailAudio: finalColor,
		UstadzUUID:          ustadzUUID,
		DurationAudio:       duration,
		MoodCategoryUUID:    moodCategoryUUID,
	}

	audio, err := h.service.Update(uuid, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": messages.MsgAudioUpdated,
		"data":    audio,
	})
}

func (h *AudioHandler) Delete(c *gin.Context) {
	uuid := c.Param("id")

	if err := h.service.Delete(uuid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": messages.MsgAudioDeleted,
	})
}

func (h *AudioHandler) IncrementListeningCount(c *gin.Context) {
	uuid := c.Param("id")

	if err := h.service.IncrementListeningCount(uuid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": messages.MsgListeningCountIncremented,
	})
}
