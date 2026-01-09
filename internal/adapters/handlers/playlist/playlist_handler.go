package playlist

import (
	"net/http"
	"strconv"
	"strings"

	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/services/playlist"
	"khalif-backend/pkg/messages"
	"khalif-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type PlaylistHandler struct {
	service playlist.PlaylistService
}

func NewPlaylistHandler(service playlist.PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{service: service}
}

// Admin handlers

// CreateAdmin creates a playlist as admin (multipart/form-data)
func (h *PlaylistHandler) CreateAdmin(c *gin.Context) {
	adminID, exists := c.Get("adminID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": messages.ErrUnauthorized})
		return
	}

	title := strings.TrimSpace(c.PostForm("title"))
	description := strings.TrimSpace(c.PostForm("description"))
	isPublicStr := c.PostForm("is_public")

	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	isPublic := true
	if isPublicStr == "false" {
		isPublic = false
	}

	// Handle thumbnail file upload
	var thumbnailURL string
	thumbnailFile, err := c.FormFile("thumbnail_file")
	if err == nil && thumbnailFile != nil {
		result, err := utils.SaveUploadedFile(thumbnailFile, utils.ThumbnailDir)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		thumbnailURL = result.URL
	}

	req := &domain.CreatePlaylistRequest{
		Title:         title,
		Description:   description,
		ThumbnailFile: thumbnailURL,
		IsPublic:      isPublic,
	}

	playlist, err := h.service.Create(req, domain.AuthorTypeAdmin, adminID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": messages.MsgPlaylistCreated, "playlist": playlist})
}

// UpdateAdmin updates a playlist as admin (multipart/form-data)
func (h *PlaylistHandler) UpdateAdmin(c *gin.Context) {
	uuid := c.Param("id")

	adminID, exists := c.Get("adminID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": messages.ErrUnauthorized})
		return
	}

	title := strings.TrimSpace(c.PostForm("title"))
	description := strings.TrimSpace(c.PostForm("description"))
	isPublicStr := c.PostForm("is_public")

	// Handle thumbnail file upload
	var thumbnailURL string
	thumbnailFile, err := c.FormFile("thumbnail_file")
	if err == nil && thumbnailFile != nil {
		result, err := utils.SaveUploadedFile(thumbnailFile, utils.ThumbnailDir)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		thumbnailURL = result.URL
	}

	var isPublic *bool
	if isPublicStr != "" {
		val := isPublicStr == "true"
		isPublic = &val
	}

	req := &domain.UpdatePlaylistRequest{
		Title:         title,
		Description:   description,
		ThumbnailFile: thumbnailURL,
		IsPublic:      isPublic,
	}

	playlist, err := h.service.Update(uuid, req, domain.AuthorTypeAdmin, adminID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": messages.MsgPlaylistUpdated, "playlist": playlist})
}

// DeleteAdmin deletes a playlist as admin
func (h *PlaylistHandler) DeleteAdmin(c *gin.Context) {
	uuid := c.Param("id")

	adminID, exists := c.Get("adminID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": messages.ErrUnauthorized})
		return
	}

	if err := h.service.Delete(uuid, domain.AuthorTypeAdmin, adminID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": messages.MsgPlaylistDeleted})
}

// User handlers

// CreateUser creates a playlist as user (multipart/form-data)
func (h *PlaylistHandler) CreateUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": messages.ErrUnauthorized})
		return
	}

	title := strings.TrimSpace(c.PostForm("title"))
	description := strings.TrimSpace(c.PostForm("description"))
	isPublicStr := c.PostForm("is_public")

	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	isPublic := true
	if isPublicStr == "false" {
		isPublic = false
	}

	// Handle thumbnail file upload
	var thumbnailURL string
	thumbnailFile, err := c.FormFile("thumbnail_file")
	if err == nil && thumbnailFile != nil {
		result, err := utils.SaveUploadedFile(thumbnailFile, utils.ThumbnailDir)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		thumbnailURL = result.URL
	}

	req := &domain.CreatePlaylistRequest{
		Title:         title,
		Description:   description,
		ThumbnailFile: thumbnailURL,
		IsPublic:      isPublic,
	}

	playlist, err := h.service.Create(req, domain.AuthorTypeUser, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": messages.MsgPlaylistCreated, "playlist": playlist})
}

// GetMyPlaylists gets user's own playlists
func (h *PlaylistHandler) GetMyPlaylists(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": messages.ErrUnauthorized})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	result, err := h.service.GetMyPlaylists(domain.AuthorTypeUser, userID.(uint), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateUser updates user's own playlist (multipart/form-data)
func (h *PlaylistHandler) UpdateUser(c *gin.Context) {
	uuid := c.Param("id")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": messages.ErrUnauthorized})
		return
	}

	title := strings.TrimSpace(c.PostForm("title"))
	description := strings.TrimSpace(c.PostForm("description"))
	isPublicStr := c.PostForm("is_public")

	// Handle thumbnail file upload
	var thumbnailURL string
	thumbnailFile, err := c.FormFile("thumbnail_file")
	if err == nil && thumbnailFile != nil {
		result, err := utils.SaveUploadedFile(thumbnailFile, utils.ThumbnailDir)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		thumbnailURL = result.URL
	}

	var isPublic *bool
	if isPublicStr != "" {
		val := isPublicStr == "true"
		isPublic = &val
	}

	req := &domain.UpdatePlaylistRequest{
		Title:         title,
		Description:   description,
		ThumbnailFile: thumbnailURL,
		IsPublic:      isPublic,
	}

	playlist, err := h.service.Update(uuid, req, domain.AuthorTypeUser, userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": messages.MsgPlaylistUpdated, "playlist": playlist})
}

// DeleteUser deletes user's own playlist
func (h *PlaylistHandler) DeleteUser(c *gin.Context) {
	uuid := c.Param("id")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": messages.ErrUnauthorized})
		return
	}

	if err := h.service.Delete(uuid, domain.AuthorTypeUser, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": messages.MsgPlaylistDeleted})
}

// AddAudioToPlaylist adds an audio to user's playlist
func (h *PlaylistHandler) AddAudioToPlaylist(c *gin.Context) {
	playlistUUID := c.Param("id")
	audioUUID := c.Param("audio_id")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": messages.ErrUnauthorized})
		return
	}

	position, _ := strconv.Atoi(c.DefaultQuery("position", "0"))

	if err := h.service.AddAudio(playlistUUID, audioUUID, position, domain.AuthorTypeUser, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": messages.MsgAudioAddedToPlaylist})
}

// RemoveAudioFromPlaylist removes an audio from user's playlist
func (h *PlaylistHandler) RemoveAudioFromPlaylist(c *gin.Context) {
	playlistUUID := c.Param("id")
	audioUUID := c.Param("audio_id")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": messages.ErrUnauthorized})
		return
	}

	if err := h.service.RemoveAudio(playlistUUID, audioUUID, domain.AuthorTypeUser, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": messages.MsgAudioRemovedFromPlaylist})
}

// Public handlers

// GetAll gets all public playlists
func (h *PlaylistHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	result, err := h.service.GetAll(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetByID gets a playlist by UUID with its audios
func (h *PlaylistHandler) GetByID(c *gin.Context) {
	uuid := c.Param("id")

	result, err := h.service.GetByUUID(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// LikePlaylist likes a playlist
func (h *PlaylistHandler) LikePlaylist(c *gin.Context) {
	uuid := c.Param("id")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": messages.ErrUnauthorized})
		return
	}

	if err := h.service.LikePlaylist(uuid, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "playlist liked"})
}

// UnlikePlaylist unlikes a playlist
func (h *PlaylistHandler) UnlikePlaylist(c *gin.Context) {
	uuid := c.Param("id")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": messages.ErrUnauthorized})
		return
	}

	if err := h.service.UnlikePlaylist(uuid, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "playlist unliked"})
}

// IsLiked checks if user has liked a playlist
func (h *PlaylistHandler) IsLiked(c *gin.Context) {
	uuid := c.Param("id")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": messages.ErrUnauthorized})
		return
	}

	liked, err := h.service.IsLiked(uuid, userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_liked": liked})
}

// IncrementListeningCount increments playlist listening count
func (h *PlaylistHandler) IncrementListeningCount(c *gin.Context) {
	uuid := c.Param("id")

	if err := h.service.IncrementListeningCount(uuid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "listening count incremented"})
}
