package http

import (
	"khalif-backend/internal/adapters/handlers"
	adminAuthHandler "khalif-backend/internal/adapters/handlers/auth/admin"
	userAuthHandler "khalif-backend/internal/adapters/handlers/auth/user"
	"khalif-backend/internal/platform/config"
	"khalif-backend/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	cfg *config.Config,
	authHandler *adminAuthHandler.AuthHandler,
	adminHandler *adminAuthHandler.AdminHandler,
	userAuthHdlr *userAuthHandler.AuthHandler,
	userHandler *userAuthHandler.UserHandler,
) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Global Middleware
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RequestIDMiddleware()) // Add RequestID before Logger
	r.Use(middleware.ZapLoggerMiddleware())

	r.Static("/uploads", "./uploads")

	healthHandler := handlers.NewHealthHandler()
	r.GET("/health", healthHandler.Health)
	r.GET("/ready", healthHandler.Ready)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		auth.Use(middleware.RateLimitMiddleware(cfg))
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)

			protected := auth.Group("/")
			protected.Use(middleware.AdminAuthMiddleware(cfg))
			{
				protected.GET("/me", authHandler.Me)
				protected.POST("/logout", authHandler.Logout)
			}
		}

		admin := api.Group("/admin")
		admin.Use(middleware.AdminAuthMiddleware(cfg))
		{
			admin.PUT("/update", adminHandler.UpdateProfile)
		}
	}

	users := api.Group("/users")
	usersAuth := users.Group("/auth")
	usersAuth.Use(middleware.RateLimitMiddleware(cfg))
	{
		usersAuth.POST("/register", userAuthHdlr.Register)
		usersAuth.POST("/login", userAuthHdlr.Login)
		usersAuth.POST("/refresh", userAuthHdlr.RefreshToken)
		
		usersProtected := usersAuth.Group("/")
		usersProtected.Use(middleware.UserAuthMiddleware(cfg))
		{
			usersProtected.GET("/me", userAuthHdlr.Me)
			usersProtected.POST("/logout", userAuthHdlr.Logout)
		}
	}
	
	usersProtected := users.Group("/")
	usersProtected.Use(middleware.UserAuthMiddleware(cfg))
	{
		usersProtected.PUT("/update", userHandler.UpdateProfile)
	}

	return r
}
