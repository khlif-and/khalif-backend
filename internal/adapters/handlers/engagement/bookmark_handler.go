package engagement

import (
	"khalif-backend/internal/core/ports"
	"khalif-backend/pkg/messages"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookmarkHandler struct {
	hadistSvc ports.HadistService
	doaSvc    ports.DoaService
}

func NewBookmarkHandler(hadistSvc ports.HadistService, doaSvc ports.DoaService) *BookmarkHandler {
	return &BookmarkHandler{
		hadistSvc: hadistSvc,
		doaSvc:    doaSvc,
	}
}

func (h *BookmarkHandler) Bookmark(c *gin.Context) {
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
	case "hadist":
		err = h.hadistSvc.BookmarkHadist(userID.(uint), uuid)
		message = messages.MsgHadistBookmarked
	case "doa":
		err = h.doaSvc.BookmarkDoa(userID.(uint), uuid)
		message = messages.MsgDoaBookmarked
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity type for bookmark"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func (h *BookmarkHandler) Unbookmark(c *gin.Context) {
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
	case "hadist":
		err = h.hadistSvc.UnbookmarkHadist(userID.(uint), uuid)
		message = messages.MsgHadistUnbookmarked
	case "doa":
		err = h.doaSvc.UnbookmarkDoa(userID.(uint), uuid)
		message = messages.MsgDoaUnbookmarked
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity type for bookmark"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func (h *BookmarkHandler) IsBookmarked(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": messages.ErrUnauthorized})
		return
	}

	entity := c.Param("entity")
	uuid := c.Param("id")

	var isBookmarked bool
	var err error

	switch entity {
	case "hadist":
		isBookmarked, err = h.hadistSvc.IsBookmarked(userID.(uint), uuid)
	case "doa":
		isBookmarked, err = h.doaSvc.IsBookmarked(userID.(uint), uuid)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity type for bookmark"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrInternalServer})
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_bookmarked": isBookmarked})
}
