package user

import (
	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"
	"khalif-backend/pkg/messages"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service ports.UserService
}

func NewUserHandler(service ports.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": messages.ErrUnauthorized})
		return
	}

	var req domain.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrBadRequest})
		return
	}

	updatedUser, err := h.service.UpdateProfile(userID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": messages.MsgUpdateSuccess,
		"data":    updatedUser,
	})
}
