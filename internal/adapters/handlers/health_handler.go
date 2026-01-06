package handlers

import (
	"context"
	"net/http"

	"khalif-backend/internal/platform/database"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}

func (h *HealthHandler) Ready(c *gin.Context) {
	dbOK := checkDatabase()
	redisOK := checkRedis()

	if dbOK && redisOK {
		c.JSON(http.StatusOK, gin.H{
			"status":   "ready",
			"database": "connected",
			"redis":    "connected",
		})
		return
	}

	response := gin.H{
		"status":   "not ready",
		"database": boolToStatus(dbOK),
		"redis":    boolToStatus(redisOK),
	}

	c.JSON(http.StatusServiceUnavailable, response)
}

func checkDatabase() bool {
	if database.DB == nil {
		return false
	}
	sqlDB, err := database.DB.DB()
	if err != nil {
		return false
	}
	return sqlDB.Ping() == nil
}

func checkRedis() bool {
	if database.RedisClient == nil {
		return false
	}
	ctx := context.Background()
	return database.RedisClient.Ping(ctx).Err() == nil
}

func boolToStatus(ok bool) string {
	if ok {
		return "connected"
	}
	return "disconnected"
}
