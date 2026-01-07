package user

import (
	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"
	"khalif-backend/pkg/messages"
	"khalif-backend/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service ports.UserAuthService
}

func NewAuthHandler(service ports.UserAuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(c *gin.Context) {
	username := strings.TrimSpace(c.PostForm("username"))
	email := strings.TrimSpace(c.PostForm("email"))
	phone := strings.TrimSpace(c.PostForm("phone"))
	password := c.PostForm("password")

	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrUsernameRequired})
		return
	}
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrEmailRequired})
		return
	}
	if !utils.IsValidEmail(email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrInvalidEmail})
		return
	}
	if phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrPhoneRequired})
		return
	}
	if password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrPasswordRequired})
		return
	}
	if len(password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrPasswordTooShort})
		return
	}
	if !utils.IsValidPhone(phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrInvalidPhone})
		return
	}

	var profilePicURL string

	file, err := c.FormFile("profile_picture")
	if err == nil && file != nil {
		result, err := utils.SaveUploadedFile(file, utils.ProfilePicDir)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		profilePicURL = result.URL
	}

	req := &domain.RegisterRequest{
		Username:       username,
		Email:          email,
		Phone:          phone,
		Password:       password,
		ProfilePicture: profilePicURL,
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	resp, err := h.service.Register(req, userAgent, ipAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": messages.MsgRegisterSuccess,
		"data":    resp,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	email := strings.TrimSpace(c.PostForm("email"))
	password := c.PostForm("password")

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrEmailRequired})
		return
	}
	if !utils.IsValidEmail(email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrInvalidEmail})
		return
	}
	if password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrPasswordRequired})
		return
	}

	req := &domain.LoginRequest{
		Email:    email,
		Password: password,
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	resp, err := h.service.Login(req, userAgent, ipAddress)
	if err != nil {
		status := http.StatusUnauthorized
		if err.Error() == messages.ErrAccountLocked {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": messages.MsgLoginSuccess,
		"data":    resp,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req domain.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrBadRequest})
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	resp, err := h.service.RefreshSession(req.RefreshToken, userAgent, ipAddress)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": messages.MsgRefreshTokenSuccess,
		"data":    resp,
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": messages.ErrUnauthorized})
		return
	}

	user, err := h.service.GetMe(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrInternalServer})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": messages.ErrUnauthorized})
		return
	}

	if err := h.service.Logout(userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrInternalServer})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": messages.MsgLogoutSuccess})
}
