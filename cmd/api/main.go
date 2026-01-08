package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	audioHandler "khalif-backend/internal/adapters/handlers/audio"
	adminAuthHandler "khalif-backend/internal/adapters/handlers/auth/admin"
	userAuthHandler "khalif-backend/internal/adapters/handlers/auth/user"
	likeHandler "khalif-backend/internal/adapters/handlers/like"
	moodCategoryHandler "khalif-backend/internal/adapters/handlers/mood_category"
	ustadzHandler "khalif-backend/internal/adapters/handlers/ustadz"
	appRouter "khalif-backend/internal/adapters/http"
	audioRepo "khalif-backend/internal/adapters/repositories/audio"
	adminAuthRepo "khalif-backend/internal/adapters/repositories/auth/admin"
	userAuthRepo "khalif-backend/internal/adapters/repositories/auth/user"
	likeRepo "khalif-backend/internal/adapters/repositories/like"
	moodCategoryRepo "khalif-backend/internal/adapters/repositories/mood_category"
	ustadzRepo "khalif-backend/internal/adapters/repositories/ustadz"
	audioService "khalif-backend/internal/core/services/audio"
	adminAuthService "khalif-backend/internal/core/services/auth/admin"
	userAuthService "khalif-backend/internal/core/services/auth/user"
	likeService "khalif-backend/internal/core/services/like"
	moodCategoryService "khalif-backend/internal/core/services/mood_category"
	ustadzService "khalif-backend/internal/core/services/ustadz"
	"khalif-backend/internal/infrastructure/email"
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

	adminRepo := adminAuthRepo.NewAdminRepo(db)
	userRepo := userAuthRepo.NewUserRepo(db)
	authRepo := adminAuthRepo.NewAuthRepo(db)
	userAuthRepoInstance := userAuthRepo.NewAuthRepo(db)
	audioRepoInstance := audioRepo.NewAudioRepo(db)
	moodCategoryRepoInstance := moodCategoryRepo.NewMoodCategoryRepo(db)
	ustadzRepoInstance := ustadzRepo.NewUstadzRepo(db)
	likeRepoInstance := likeRepo.NewLikeRepo(db)

	authService := adminAuthService.NewAuthService(adminRepo, authRepo, cfg)
	adminService := adminAuthService.NewAdminService(adminRepo)
	emailService := email.NewBrevoEmailService(cfg)
	userAuthSvc := userAuthService.NewAuthService(userRepo, userAuthRepoInstance, emailService, cfg)
	userService := userAuthService.NewUserService(userRepo)
	audioSvc := audioService.NewAudioService(audioRepoInstance, moodCategoryRepoInstance, ustadzRepoInstance)
	moodCategorySvc := moodCategoryService.NewMoodCategoryService(moodCategoryRepoInstance, audioRepoInstance)
	ustadzSvc := ustadzService.NewUstadzService(ustadzRepoInstance)
	likeSvc := likeService.NewLikeService(likeRepoInstance, audioRepoInstance)

	authHandler := adminAuthHandler.NewAuthHandler(authService)
	adminHandler := adminAuthHandler.NewAdminHandler(adminService)
	userAuthHdlr := userAuthHandler.NewAuthHandler(userAuthSvc)
	userHandler := userAuthHandler.NewUserHandler(userService)
	audioHdlr := audioHandler.NewAudioHandler(audioSvc)
	moodCategoryHdlr := moodCategoryHandler.NewMoodCategoryHandler(moodCategorySvc)
	ustadzHdlr := ustadzHandler.NewUstadzHandler(ustadzSvc)
	likeHdlr := likeHandler.NewLikeHandler(likeSvc)

	router := appRouter.NewRouter(cfg, authHandler, adminHandler, userAuthHdlr, userHandler, audioHdlr, moodCategoryHdlr, ustadzHdlr, likeHdlr)

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
