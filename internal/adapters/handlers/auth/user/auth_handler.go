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

	err = h.service.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": messages.MsgOTPSent,
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
		if err.Error() == messages.ErrAccountNotActivated {
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

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req domain.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrBadRequest})
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	resp, err := h.service.VerifyOTP(&req, userAgent, ipAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": messages.MsgAccountActivated,
		"data":    resp,
	})
}

func (h *AuthHandler) ResendOTP(c *gin.Context) {
	var req domain.ResendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrBadRequest})
		return
	}

	err := h.service.ResendOTP(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": messages.MsgOTPResent,
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

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req domain.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrBadRequest})
		return
	}

	if err := h.service.ForgotPassword(req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Always return success to prevent email enumeration
	c.JSON(http.StatusOK, gin.H{
		"message": messages.MsgPasswordResetSent,
	})
}

// ResetPassword handles password reset
// @Summary Reset password
// @Description Reset password using token
// @Tags auth-user
// @Accept json
// @Produce json
// @Param request body domain.ResetPasswordRequest true "Request body"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/users/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req domain.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ResetPassword(req.Token, req.NewPassword); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == messages.ErrInvalidToken || err.Error() == messages.ErrUserNotFound {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// GoogleLogin handles Google Sign-In
// @Summary Login with Google
// @Description Login or Register using Google ID Token
// @Tags auth-user
// @Accept json
// @Produce json
// @Param request body map[string]string true "Request body with id_token"
// @Success 200 {object} domain.LoginResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/users/auth/google-login [post]
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req struct {
		IDToken string `json:"id_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userAgent := c.Request.UserAgent()
	ipAddress := c.ClientIP()

	resp, err := h.service.LoginWithGoogle(req.IDToken, userAgent, ipAddress)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "invalid google token" {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
