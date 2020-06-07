package main

import (
	redis "github.com/go-redis/redis/v8"
)

// RedisClient : Creates and return redis client
func RedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6367",
		Password: "",
		DB:       2,
	})
	return client
}
