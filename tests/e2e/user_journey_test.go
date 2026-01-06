package e2e_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"khalif-backend/internal/adapters/handlers/auth/admin"
	userAuthHandler "khalif-backend/internal/adapters/handlers/auth/user"
	appRouter "khalif-backend/internal/adapters/http"
	adminAuthRepo "khalif-backend/internal/adapters/repositories/auth/admin"
	userAuthRepo "khalif-backend/internal/adapters/repositories/auth/user"
	adminAuthService "khalif-backend/internal/core/services/auth/admin"
	userAuthService "khalif-backend/internal/core/services/auth/user"
	"khalif-backend/internal/platform/config"
	"khalif-backend/internal/platform/database"
	"khalif-backend/internal/platform/logger"
	"khalif-backend/pkg/utils"
	"os"
)

// Setup Helper for E2E (Similar to Integration Setup)
func setupE2EServer(t *testing.T) *httptest.Server {
	// Change to project root so migrations can be found
	_ = os.Chdir("../../")
	
	// Load .env from project root
	// Ignore error if specific env file doesn't exist (relies on system env)
	_ = godotenv.Load(".env")

	cfg := config.LoadConfig()
	cfg.AppEnv = "test" // Force test env
	
	logger.InitLogger(cfg.AppEnv)

	db := database.InitDB(cfg)
	database.InitRedis(cfg)
	utils.InitUploadDirs()

	// Clean State
	db.Exec("TRUNCATE TABLE users, user_refresh_tokens RESTART IDENTITY CASCADE")

	// Dependencies
	adminRepo := adminAuthRepo.NewAdminRepo(db)
	userRepo := userAuthRepo.NewUserRepo(db)
	authRepo := adminAuthRepo.NewAuthRepo(db)
	userAuthRepoInstance := userAuthRepo.NewAuthRepo(db)

	authService := adminAuthService.NewAuthService(adminRepo, authRepo, cfg)
	adminService := adminAuthService.NewAdminService(adminRepo)
	userAuthSvc := userAuthService.NewAuthService(userRepo, userAuthRepoInstance, cfg)
	userService := userAuthService.NewUserService(userRepo)

	authHandler := admin.NewAuthHandler(authService)
	adminHandler := admin.NewAdminHandler(adminService)
	userAuthHdlr := userAuthHandler.NewAuthHandler(userAuthSvc)
	userHandler := userAuthHandler.NewUserHandler(userService)

	router := appRouter.NewRouter(cfg, authHandler, adminHandler, userAuthHdlr, userHandler)

	return httptest.NewServer(router)
}

func createMultipartRequest(fields map[string]string) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range fields {
		_ = writer.WriteField(key, val)
	}
	writer.Close()
	return body, writer.FormDataContentType()
}

func TestE2E_UserJourney(t *testing.T) {
	server := setupE2EServer(t)
	defer server.Close()

	client := &http.Client{}
	var accessToken, refreshToken string

	// Step 1: Register User
	t.Run("1_Register", func(t *testing.T) {
		fields := map[string]string{
			"username": "journeyuser",
			"email":    "journey@example.com",
			"phone":    "08888888888",
			"password": "password123",
		}
		body, contentType := createMultipartRequest(fields)
		resp, err := http.Post(server.URL+"/api/v1/users/auth/register", contentType, body)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	// Step 2: Login to get Tokens
	t.Run("2_Login", func(t *testing.T) {
		fields := map[string]string{
			"email":    "journey@example.com",
			"password": "password123",
		}
		body, contentType := createMultipartRequest(fields)
		resp, err := http.Post(server.URL+"/api/v1/users/auth/login", contentType, body)
		assert.NoError(t, err)
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		data := result["data"].(map[string]interface{})
		accessToken = data["token"].(string)
		refreshToken = data["refresh_token"].(string)
	})

	// Step 3: Get Profile (Me)
	t.Run("3_GetMe", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL+"/api/v1/users/auth/me", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		data := result["data"].(map[string]interface{})
		assert.Equal(t, "journeyuser", data["username"])
	})

	// Step 4: Update Profile
	t.Run("4_UpdateProfile", func(t *testing.T) {
		updateData := map[string]interface{}{
			"username": "journey_updated",
			"email": "journey@example.com", // Keeping email same to avoid unique constraint if not handling change
			"phone": "08888888888", 
		}
		jsonBody, _ := json.Marshal(updateData)
		req, _ := http.NewRequest("PUT", server.URL+"/api/v1/users/update", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Content-Type", "application/json") // Handlers likely expect JSON for update
		
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		
		// Note: Check if UpdateProfile uses JSON or Multipart.
		// Usually updates are JSON. If it fails, I'll check Handler.
		assert.Equal(t, http.StatusOK, resp.StatusCode) 
	})

	// Step 5: Refresh Token
	t.Run("5_RefreshToken", func(t *testing.T) {
		reqBody := map[string]string{"refresh_token": refreshToken}
		jsonBody, _ := json.Marshal(reqBody)
		
		resp, err := http.Post(server.URL+"/api/v1/users/auth/refresh", "application/json", bytes.NewBuffer(jsonBody))
		assert.NoError(t, err)
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		data := result["data"].(map[string]interface{})
		accessToken = data["token"].(string) // Update Access Token
	})

	// Step 6: Logout
	t.Run("6_Logout", func(t *testing.T) {
		req, _ := http.NewRequest("POST", server.URL+"/api/v1/users/auth/logout", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// Step 7: Verify Token Revoked (Refresh should fail)
	t.Run("7_VerifyLogout", func(t *testing.T) {
		reqBody := map[string]string{"refresh_token": refreshToken}
		jsonBody, _ := json.Marshal(reqBody)
		
		resp, err := http.Post(server.URL+"/api/v1/users/auth/refresh", "application/json", bytes.NewBuffer(jsonBody))
		assert.NoError(t, err)
		defer resp.Body.Close()
		
		// Should be 401 Unauthorized because refresh token was revoked
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
