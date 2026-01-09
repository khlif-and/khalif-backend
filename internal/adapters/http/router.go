package http

import (
	"khalif-backend/internal/adapters/handlers"
	audioHandler "khalif-backend/internal/adapters/handlers/audio"
	adminAuthHandler "khalif-backend/internal/adapters/handlers/auth/admin"
	userAuthHandler "khalif-backend/internal/adapters/handlers/auth/user"
	likeHandler "khalif-backend/internal/adapters/handlers/like"
	moodCategoryHandler "khalif-backend/internal/adapters/handlers/mood_category"
	playlistHandler "khalif-backend/internal/adapters/handlers/playlist"
	searchHandler "khalif-backend/internal/adapters/handlers/search"
	ustadzHandler "khalif-backend/internal/adapters/handlers/ustadz"
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
	audioHdlr *audioHandler.AudioHandler,
	moodHdlr *moodCategoryHandler.MoodCategoryHandler,
	ustadzHdlr *ustadzHandler.UstadzHandler,
	likeHdlr *likeHandler.LikeHandler,
	searchHdlr *searchHandler.SearchHandler,
	playlistHdlr *playlistHandler.PlaylistHandler,
) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RequestIDMiddleware())
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

			admin.POST("/audio", audioHdlr.Create)
			admin.PUT("/audio/:id", audioHdlr.Update)
			admin.DELETE("/audio/:id", audioHdlr.Delete)

			admin.POST("/mood-categories", moodHdlr.Create)
			admin.PUT("/mood-categories/:id", moodHdlr.Update)
			admin.DELETE("/mood-categories/:id", moodHdlr.Delete)

			admin.POST("/ustadz", ustadzHdlr.Create)
			admin.PUT("/ustadz/:id", ustadzHdlr.Update)
			admin.DELETE("/ustadz/:id", ustadzHdlr.Delete)

			// Playlist admin routes
			admin.POST("/playlist", playlistHdlr.CreateAdmin)
			admin.PUT("/playlist/:id", playlistHdlr.UpdateAdmin)
			admin.DELETE("/playlist/:id", playlistHdlr.DeleteAdmin)
		}

		audio := api.Group("/audio")
		{
			audio.GET("", audioHdlr.GetAll)
			audio.GET("/:id", audioHdlr.GetByID)
		}

		// Radio (public)
		api.GET("/radio/:id", audioHdlr.GetRadio)

		moods := api.Group("/mood-categories")
		{
			moods.GET("", moodHdlr.GetAll)
			moods.GET("/:id", moodHdlr.GetByID)
			moods.GET("/:id/audios", moodHdlr.GetAudiosByMoodID)
		}

		ustadzGroup := api.Group("/ustadz")
		{
			ustadzGroup.GET("", ustadzHdlr.GetAll)
			ustadzGroup.GET("/:id", ustadzHdlr.GetByID)
		}

		// Search routes (public)
		search := api.Group("/search")
		{
			search.GET("", searchHdlr.SearchAll)
			search.GET("/audio", searchHdlr.SearchAudios)
			search.GET("/ustadz", searchHdlr.SearchUstadzs)
			search.GET("/mood", searchHdlr.SearchMoodCategories)
			search.GET("/playlist", searchHdlr.SearchPlaylists)
		}

		// Playlist public routes
		playlist := api.Group("/playlist")
		{
			playlist.GET("", playlistHdlr.GetAll)
			playlist.GET("/:id", playlistHdlr.GetByID)
			playlist.POST("/:id/listen", playlistHdlr.IncrementListeningCount)
		}
	}

	users := api.Group("/users")
	usersAuth := users.Group("/auth")
	usersAuth.Use(middleware.RateLimitMiddleware(cfg))
	{
	usersAuth.POST("/register", userAuthHdlr.Register)
		usersAuth.POST("/login", userAuthHdlr.Login)
		usersAuth.POST("/refresh", userAuthHdlr.RefreshToken)
		usersAuth.POST("/verify-otp", userAuthHdlr.VerifyOTP)
		usersAuth.POST("/resend-otp", userAuthHdlr.ResendOTP)
		usersAuth.POST("/forgot-password", userAuthHdlr.ForgotPassword)
		usersAuth.POST("/reset-password", userAuthHdlr.ResetPassword)
		
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
		usersProtected.POST("/audio/:id/listen", audioHdlr.IncrementListeningCount)
		usersProtected.GET("/listening-history", audioHdlr.GetListeningHistory)

		usersProtected.POST("/audio/:id/like", likeHdlr.LikeAudio)
		usersProtected.DELETE("/audio/:id/like", likeHdlr.UnlikeAudio)
		usersProtected.GET("/audio/:id/is-liked", likeHdlr.IsLiked)
		usersProtected.GET("/likes", likeHdlr.GetUserLikes)

		// Playlist user routes
		usersProtected.POST("/playlist", playlistHdlr.CreateUser)
		usersProtected.GET("/playlist", playlistHdlr.GetMyPlaylists)
		usersProtected.PUT("/playlist/:id", playlistHdlr.UpdateUser)
		usersProtected.DELETE("/playlist/:id", playlistHdlr.DeleteUser)
		usersProtected.POST("/playlist/:id/audio/:audio_id", playlistHdlr.AddAudioToPlaylist)
		usersProtected.DELETE("/playlist/:id/audio/:audio_id", playlistHdlr.RemoveAudioFromPlaylist)
		usersProtected.POST("/playlist/:id/like", playlistHdlr.LikePlaylist)
		usersProtected.DELETE("/playlist/:id/like", playlistHdlr.UnlikePlaylist)
		usersProtected.GET("/playlist/:id/is-liked", playlistHdlr.IsLiked)
	}

	return r
}
