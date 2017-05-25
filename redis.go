package main

import (
	"github.com/go-redis/redis"
	"os"
)

func initRedis() *redis.Client {
	redisUrl := os.Getenv("REDIS_URL")

	if redisUrl == "" {
		redisUrl = "127.0.0.1:6379"
	}

	return redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: "",
		DB:       0,
	})
}
