package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"khalif-backend/internal/platform/config"
	appLogger "khalif-backend/internal/platform/logger"
	"go.uber.org/zap"
)

var RedisClient *redis.Client

func InitRedis(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		appLogger.Log.Fatal("failed to connect to redis", zap.Error(err))
	}

	appLogger.Log.Info("Redis connected successfully")
	RedisClient = rdb
	return rdb
}
