package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"khalif-backend/internal/platform/config"
	"khalif-backend/internal/platform/database"
)

func RateLimitMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip rate limiting in test environment
		if cfg.AppEnv == "test" {
			c.Next()
			return
		}

		if database.RedisClient == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "redis unavailable"})
			return
		}

		ip := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s", ip)
		limit := 5
		window := 1 * time.Second

		ctx := context.Background()

		pipe := database.RedisClient.Pipeline()
		incr := pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, window)
		_, err := pipe.Exec(ctx)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "rate limit error"})
			return
		}

		if incr.Val() > int64(limit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			return
		}

		c.Next()
	}
}
