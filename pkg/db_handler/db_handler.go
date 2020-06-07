package dbhandler

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type DBHandler struct {
	redisClient *redis.Client
}

type User struct {
	username string
	token    string
}

func GetDBHandler(redisClient *redis.Client) DBHandler {
	return DBHandler{redisClient}
}

func (client *DBHandler) IsUserPresent(username string) bool {
	member := client.redisClient.SIsMember(context.Background(), "users", username).Val()
	return member
}

func (client *DBHandler) AssociateToken(username string, token string) {
	client.redisClient.SAdd(context.Background(), "users", username)
	client.redisClient.HSet(context.Background(), "token", username, token)
}
