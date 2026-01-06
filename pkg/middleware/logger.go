package middleware

import (
	"khalif-backend/internal/platform/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ZapLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.Log.Error(e)
			}
		} else {
			logger.Log.Info("request",
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("request_id", c.GetString(RequestIDKey)),
				zap.Duration("latency", latency),
			)
		}
	}
}
