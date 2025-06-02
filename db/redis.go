package db

import (
	"context"
	"go-postgres-gorm-gin-api/config"
	"log"

	"github.com/go-redis/redis/v8"
)

// RedisClient is the Redis client
var RedisClient *redis.Client

func ConnectRedis(cfg *config.Config) (*redis.Client, error) {
	// Initialize Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	var err error
	_, err = RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		return nil, err
	}

	log.Println("Connected to Redis")

	return RedisClient, nil
}
