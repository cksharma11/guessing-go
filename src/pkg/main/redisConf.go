package main

import (
	"github.com/go-redis/redis/v8"
)

// GetRedisClient : Creates and return redis client
func GetRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       2,
	})
	return client
}
