package http

import (
	"khalif-backend/internal/adapters/handlers"
	alquranHandler "khalif-backend/internal/adapters/handlers/alquran"
	audioHandler "khalif-backend/internal/adapters/handlers/audio"
	adminAuthHandler "khalif-backend/internal/adapters/handlers/auth/admin"
	userAuthHandler "khalif-backend/internal/adapters/handlers/auth/user"
	doaHandler "khalif-backend/internal/adapters/handlers/doa"
	engagementHandler "khalif-backend/internal/adapters/handlers/engagement"
	hadistHandler "khalif-backend/internal/adapters/handlers/hadist"
	moodCategoryHandler "khalif-backend/internal/adapters/handlers/mood_category"
	playlistHandler "khalif-backend/internal/adapters/handlers/playlist"
	prayerHandler "khalif-backend/internal/adapters/handlers/prayer"
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
	likeHdlr *engagementHandler.LikeHandler,
	bookmarkHdlr *engagementHandler.BookmarkHandler,
	searchHdlr *searchHandler.SearchHandler,
	playlistHdlr *playlistHandler.PlaylistHandler,
	hadistHdlr *hadistHandler.HadistHandler,
	doaHdlr *doaHandler.DoaHandler,
	prayerHdlr *prayerHandler.PrayerHandler,
	alquranHdlr *alquranHandler.AlquranHandler,
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

			admin.POST("/playlist", playlistHdlr.CreateAdmin)
			admin.PUT("/playlist/:id", playlistHdlr.UpdateAdmin)
			admin.DELETE("/playlist/:id", playlistHdlr.DeleteAdmin)

			admin.POST("/hadist", hadistHdlr.Create)
			admin.PUT("/hadist/:id", hadistHdlr.Update)
			admin.DELETE("/hadist/:id", hadistHdlr.Delete)

			admin.POST("/doa", doaHdlr.Create)
			admin.PUT("/doa/:id", doaHdlr.Update)
			admin.DELETE("/doa/:id", doaHdlr.Delete)
		}

		audio := api.Group("/audio")
		{
			audio.GET("", audioHdlr.GetAll)
			audio.GET("/:id", audioHdlr.GetByID)
		}

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

		search := api.Group("/search")
		{
			search.GET("", searchHdlr.SearchAll)
			search.GET("/audio", searchHdlr.SearchAudios)
			search.GET("/ustadz", searchHdlr.SearchUstadzs)
			search.GET("/mood", searchHdlr.SearchMoodCategories)
			search.GET("/playlist", searchHdlr.SearchPlaylists)
			search.GET("/doa", searchHdlr.SearchDoas)
		}

		playlist := api.Group("/playlist")
		{
			playlist.GET("", playlistHdlr.GetAll)
			playlist.GET("/:id", playlistHdlr.GetByID)
			playlist.POST("/:id/listen", playlistHdlr.IncrementListeningCount)
		}

		hadist := api.Group("/hadist")
		{
			hadist.GET("", hadistHdlr.GetAll)
			hadist.GET("/random", hadistHdlr.GetRandom)
			hadist.GET("/category", hadistHdlr.GetByCategory)
			hadist.GET("/kitab", hadistHdlr.GetByKitab)
			hadist.GET("/:id", hadistHdlr.GetByID)
			hadist.POST("/:id/listen", hadistHdlr.IncrementListeningCount)
		}

		doa := api.Group("/doa")
		{
			doa.GET("", doaHdlr.GetAll)
			doa.GET("/random", doaHdlr.GetRandom)
			doa.GET("/category", doaHdlr.GetByCategory)
			doa.GET("/hadist", doaHdlr.GetByHadist)
			doa.GET("/:id", doaHdlr.GetByID)
			doa.POST("/:id/listen", doaHdlr.IncrementListeningCount)
		}


		api.POST("/alquran", alquranHdlr.CreateEndpoint)
		api.GET("/alquran", alquranHdlr.GetAllEndpont)
		api.GET("/alquran/:id", alquranHdlr.GetByIDEndpoint)

		api.GET("/prayer-times", prayerHdlr.GetPrayerTimes)
		api.GET("/prayer-times/daily", prayerHdlr.GetDailyPrayerTimes)
	}

	users := api.Group("/users")
	usersAuth := users.Group("/auth")
	usersAuth.Use(middleware.RateLimitMiddleware(cfg))
	{
		usersAuth.POST("/register", userAuthHdlr.Register)
		usersAuth.POST("/login", userAuthHdlr.Login)
		usersAuth.POST("/refresh-token", userAuthHdlr.RefreshToken)
		usersAuth.POST("/forgot-password", userAuthHdlr.ForgotPassword)
		usersAuth.POST("/reset-password", userAuthHdlr.ResetPassword)
		usersAuth.POST("/google-login", userAuthHdlr.GoogleLogin)
		usersAuth.POST("/verify-otp", userAuthHdlr.VerifyOTP)
		usersAuth.POST("/resend-otp", userAuthHdlr.ResendOTP)

		// Protected auth routes (require valid JWT)
		usersAuthProtected := usersAuth.Group("/")
		usersAuthProtected.Use(middleware.UserAuthMiddleware(cfg))
		{
			usersAuthProtected.POST("/logout", userAuthHdlr.Logout)
			usersAuthProtected.GET("/me", userAuthHdlr.Me)
		}
	}

	usersProtected := users.Group("/")
	usersProtected.Use(middleware.UserAuthMiddleware(cfg))
	{
		usersProtected.PUT("/update", userHandler.UpdateProfile)
		usersProtected.POST("/audio/:id/listen", audioHdlr.IncrementListeningCount)
		usersProtected.GET("/listening-history", audioHdlr.GetListeningHistory)

		// Unified Like routes - :entity = audio|hadist|doa|playlist
		usersProtected.POST("/:entity/:id/like", likeHdlr.Like)
		usersProtected.DELETE("/:entity/:id/like", likeHdlr.Unlike)
		usersProtected.GET("/:entity/:id/is-liked", likeHdlr.IsLiked)
		usersProtected.GET("/likes", likeHdlr.GetUserLikes)

		// Unified Bookmark routes - :entity = hadist|doa
		usersProtected.POST("/:entity/:id/bookmark", bookmarkHdlr.Bookmark)
		usersProtected.DELETE("/:entity/:id/bookmark", bookmarkHdlr.Unbookmark)
		usersProtected.GET("/:entity/:id/is-bookmarked", bookmarkHdlr.IsBookmarked)

		// Playlist user routes
		usersProtected.POST("/playlist", playlistHdlr.CreateUser)
		usersProtected.GET("/playlist", playlistHdlr.GetMyPlaylists)
		usersProtected.PUT("/playlist/:id", playlistHdlr.UpdateUser)
		usersProtected.DELETE("/playlist/:id", playlistHdlr.DeleteUser)
		usersProtected.POST("/playlist/:id/audio/:audio_id", playlistHdlr.AddAudioToPlaylist)
		usersProtected.DELETE("/playlist/:id/audio/:audio_id", playlistHdlr.RemoveAudioFromPlaylist)
	}

	return r
}
