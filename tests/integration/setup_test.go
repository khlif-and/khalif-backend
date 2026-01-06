package integration_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	adminAuthHandler "khalif-backend/internal/adapters/handlers/auth/admin"
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
)

func setupTestServer(t *testing.T) *httptest.Server {
	// Change to project root so migrations can be found
	_ = os.Chdir("../../")
	
	// Load .env from project root
	_ = godotenv.Load(".env")
	
	cfg := config.LoadConfig()
	cfg.AppEnv = "test"
	
	// Initialize logger before database
	logger.InitLogger(cfg.AppEnv)

	defer func() {
		if r := recover(); r != nil {
			t.Skipf("Skipping Integration Test: %v", r)
		}
	}()

	db := database.InitDB(cfg)
	database.InitRedis(cfg)
	utils.InitUploadDirs()

	db.Exec("TRUNCATE TABLE admins, refresh_tokens RESTART IDENTITY CASCADE")

	// Repositories
	adminRepo := adminAuthRepo.NewAdminRepo(db)
	userRepo := userAuthRepo.NewUserRepo(db)
	authRepo := adminAuthRepo.NewAuthRepo(db)
	userAuthRepoInstance := userAuthRepo.NewAuthRepo(db)

	// Services
	authService := adminAuthService.NewAuthService(adminRepo, authRepo, cfg)
	adminService := adminAuthService.NewAdminService(adminRepo)
	userAuthSvc := userAuthService.NewAuthService(userRepo, userAuthRepoInstance, cfg)
	userService := userAuthService.NewUserService(userRepo)

	// Handlers
	authHandler := adminAuthHandler.NewAuthHandler(authService)
	adminHandler := adminAuthHandler.NewAdminHandler(adminService)
	userAuthHdlr := userAuthHandler.NewAuthHandler(userAuthSvc)
	userHandler := userAuthHandler.NewUserHandler(userService)

	router := appRouter.NewRouter(cfg, authHandler, adminHandler, userAuthHdlr, userHandler)

	return httptest.NewServer(router)
}

func createMultipartRequest(fields map[string]string, file []byte, fileName, fileFieldName string) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, val := range fields {
		_ = writer.WriteField(key, val)
	}

	if file != nil && fileName != "" {
		part, _ := writer.CreateFormFile(fileFieldName, fileName)
		part.Write(file)
	}

	writer.Close()
	return body, writer.FormDataContentType()
}

func postJSON(server *httptest.Server, path string, data map[string]string) (*http.Response, error) {
	jsonBody, _ := json.Marshal(data)
	return http.Post(server.URL+path, "application/json", bytes.NewBuffer(jsonBody))
}
