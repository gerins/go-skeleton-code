package redis

import (
	"context"

	"go-skeleton-code/pkg/log"

	"github.com/go-redis/redis/v8"

	"go-skeleton-code/config"
)

func Init(cfg config.Cache) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.Database,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal("failed connecting to redis server")
	}

	return redisClient
}
