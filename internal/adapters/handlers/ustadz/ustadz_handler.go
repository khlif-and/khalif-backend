package ustadz

import (
	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"
	"khalif-backend/pkg/messages"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UstadzHandler struct {
	service ports.UstadzService
}

func NewUstadzHandler(service ports.UstadzService) *UstadzHandler {
	return &UstadzHandler{service: service}
}

func (h *UstadzHandler) Create(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	description := strings.TrimSpace(c.PostForm("description"))
	wikipediaLink := strings.TrimSpace(c.PostForm("wikipedia_link"))

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrUstadzNameRequired})
		return
	}

	req := &domain.CreateUstadzRequest{
		Name:          name,
		Description:   description,
		WikipediaLink: wikipediaLink,
	}

	ustadz, err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": messages.MsgUstadzCreated,
		"data":    ustadz,
	})
}

func (h *UstadzHandler) GetAll(c *gin.Context) {
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

func (h *UstadzHandler) GetByID(c *gin.Context) {
	uuid := c.Param("id")

	ustadz, err := h.service.GetByUUID(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": messages.ErrUstadzNotFound})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": ustadz,
	})
}

func (h *UstadzHandler) Update(c *gin.Context) {
	uuid := c.Param("id")

	name := strings.TrimSpace(c.PostForm("name"))
	description := strings.TrimSpace(c.PostForm("description"))
	wikipediaLink := strings.TrimSpace(c.PostForm("wikipedia_link"))

	req := &domain.UpdateUstadzRequest{
		Name:          name,
		Description:   description,
		WikipediaLink: wikipediaLink,
	}

	ustadz, err := h.service.Update(uuid, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": messages.MsgUstadzUpdated,
		"data":    ustadz,
	})
}

func (h *UstadzHandler) Delete(c *gin.Context) {
	uuid := c.Param("id")

	if err := h.service.Delete(uuid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": messages.MsgUstadzDeleted,
	})
}
