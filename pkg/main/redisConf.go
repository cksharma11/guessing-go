package main

import (
	redis "github.com/go-redis/redis/v8"
)

// GetRedisClient : Creates and return redis client
func GetRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       2,
	})
	return client
}
