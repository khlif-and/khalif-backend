package engagement

import (
	"khalif-backend/internal/core/ports"
	"khalif-backend/pkg/messages"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LikeHandler struct {
	audioLikeSvc   ports.LikeService
	hadistSvc      ports.HadistService
	doaSvc         ports.DoaService
	playlistSvc    ports.PlaylistService
}

func NewLikeHandler(
	audioLikeSvc ports.LikeService,
	hadistSvc ports.HadistService,
	doaSvc ports.DoaService,
	playlistSvc ports.PlaylistService,
) *LikeHandler {
	return &LikeHandler{
		audioLikeSvc:   audioLikeSvc,
		hadistSvc:      hadistSvc,
		doaSvc:         doaSvc,
		playlistSvc:    playlistSvc,
	}
}

func (h *LikeHandler) Like(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": messages.ErrUnauthorized})
		return
	}

	entity := c.Param("entity")
	uuid := c.Param("id")

	var err error
	var message string

	switch entity {
	case "audio":
		_, err = h.audioLikeSvc.LikeAudio(userID.(uint), uuid)
		message = messages.MsgLikeCreated
	case "hadist":
		err = h.hadistSvc.LikeHadist(userID.(uint), uuid)
		message = messages.MsgHadistLiked
	case "doa":
		err = h.doaSvc.LikeDoa(userID.(uint), uuid)
		message = messages.MsgDoaLiked
	case "playlist":
		err = h.playlistSvc.LikePlaylist(uuid, userID.(uint))
		message = messages.MsgPlaylistLiked
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity type"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func (h *LikeHandler) Unlike(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": messages.ErrUnauthorized})
		return
	}

	entity := c.Param("entity")
	uuid := c.Param("id")

	var err error
	var message string

	switch entity {
	case "audio":
		err = h.audioLikeSvc.UnlikeAudio(userID.(uint), uuid)
		message = messages.MsgLikeDeleted
	case "hadist":
		err = h.hadistSvc.UnlikeHadist(userID.(uint), uuid)
		message = messages.MsgHadistUnliked
	case "doa":
		err = h.doaSvc.UnlikeDoa(userID.(uint), uuid)
		message = messages.MsgDoaUnliked
	case "playlist":
		err = h.playlistSvc.UnlikePlaylist(uuid, userID.(uint))
		message = messages.MsgPlaylistUnliked
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity type"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func (h *LikeHandler) IsLiked(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": messages.ErrUnauthorized})
		return
	}

	entity := c.Param("entity")
	uuid := c.Param("id")

	var isLiked bool
	var err error

	switch entity {
	case "audio":
		isLiked, err = h.audioLikeSvc.IsLiked(userID.(uint), uuid)
	case "hadist":
		isLiked, err = h.hadistSvc.IsLiked(userID.(uint), uuid)
	case "doa":
		isLiked, err = h.doaSvc.IsLiked(userID.(uint), uuid)
	case "playlist":
		isLiked, err = h.playlistSvc.IsLiked(uuid, userID.(uint))
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity type"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrInternalServer})
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_liked": isLiked})
}

func (h *LikeHandler) GetUserLikes(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": messages.ErrUnauthorized})
		return
	}

	entity := c.Query("entity")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	switch entity {
	case "audio", "":
		response, err := h.audioLikeSvc.GetUserLikes(userID.(uint), page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrInternalServer})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": response})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "entity likes list not supported yet"})
	}
}
