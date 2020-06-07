package dbhandler

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type DBHandler struct {
	redisClient *redis.Client
}

func GetDBHandler(redisClient *redis.Client) DBHandler {
	return DBHandler{redisClient}
}

func (client *DBHandler) IsUserPresent(username string) bool {
	return client.redisClient.SIsMember(context.Background(), "users", username).Val()
}

func (client *DBHandler) AssociateToken(username string, token string) {
	client.redisClient.SAdd(context.Background(), "users", username)
	client.redisClient.HSet(context.Background(), "tokens", username, token)
}

func (client *DBHandler) ValidateToken(token string) bool {
	return client.redisClient.HExists(context.Background(), "tokens", token).Val()
}

func (client *DBHandler) GetUser(token string) string {
	return client.redisClient.HGet(context.Background(), "tokens", token).Val()
}

func (client *DBHandler) RegisterGuess(username string, guess string) {
	client.redisClient.HSet(context.Background(), "guesses", username, guess)
}

func (client *DBHandler) CurrentLevel() string {
	currentLevel := client.redisClient.Get(context.Background(), "current-level").Val()
	if currentLevel == "" {
		client.redisClient.GetSet(context.Background(), "current-level", 1)
		return client.redisClient.Get(context.Background(), "current-level").Val()
	}
	return currentLevel
}

func (client *DBHandler) IncrementLevel() {
	client.redisClient.Incr(context.Background(), "current-level")
}
