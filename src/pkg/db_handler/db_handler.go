package dbhandler

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
)

func GetDBHandler(redisClient *redis.Client) DBHandler {
	return DBHandler{redisClient}
}

func (client *DBHandler) IsUserPresent(username string) bool {
	return client.redisClient.SIsMember(context.Background(), "users", username).Val()
}

func (client *DBHandler) AssociateToken(username string, token string) {
	pipe := client.redisClient.TxPipeline()
	pipe.SAdd(context.Background(), "users", username)
	pipe.HSet(context.Background(), "tokens", token, username)

	_, err := pipe.Exec(context.Background())
	if err != nil {
		panic("Transaction failed")
	}
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

func (client *DBHandler) IncrementLevel() IncrementLevelResponse {
	guesses, _ := client.redisClient.HGetAll(context.Background(), "guesses").Result()
	winner := getWinner(guesses)

	h := history{
		Winner:  winner,
		Level:   client.CurrentLevel(),
		Guesses: getGuessMap(guesses),
	}
	historyJson, _ := json.Marshal(h)
	cl := client.CurrentLevel()

	pipe := client.redisClient.TxPipeline()
	pipe.Del(context.Background(), "guesses")
	pipe.RPush(context.Background(), "history", historyJson)
	pipe.Incr(context.Background(), "current-level")

	_, err := pipe.Exec(context.Background())

	if err != nil {
		panic("Transaction failed")
	}

	return IncrementLevelResponse{
		Result: h,
		Pl:     cl,
		Cl:     client.CurrentLevel(),
	}
}
