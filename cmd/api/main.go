package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfig()

	logger.InitLogger(cfg.AppEnv)
	logger.Log.Info("Starting Khalif Backend Auth Service...")

	db := database.InitDB(cfg)
	database.InitRedis(cfg)

	if err := utils.InitUploadDirs(); err != nil {
		logger.Log.Fatal(fmt.Sprintf("Failed to init upload dirs: %v", err))
	}

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

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Log.Info(fmt.Sprintf("Server running on port %s", cfg.ServerPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Error("Server forced to shutdown", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
		logger.Log.Info("Database connection closed")
	}

	if database.RedisClient != nil {
		database.RedisClient.Close()
		logger.Log.Info("Redis connection closed")
	}

	logger.Log.Info("Server exited gracefully")
}
