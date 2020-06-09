package dbhandler

import (
	"context"
	"encoding/json"
	"github.com/cksharma11/guessing/src/pkg/dbkeys"
	"github.com/go-redis/redis/v8"
)

func GetDBHandler(redisClient *redis.Client) DBHandler {
	return DBHandler{redisClient}
}

func (client *DBHandler) IsUserPresent(username string) bool {
	return client.redisClient.SIsMember(context.Background(), dbkeys.UsersKey, username).Val()
}

func (client *DBHandler) AssociateToken(username string, token string) {
	pipe := client.redisClient.TxPipeline()
	pipe.SAdd(context.Background(), dbkeys.UsersKey, username)
	pipe.HSet(context.Background(), dbkeys.Tokens, token, username)

	_, err := pipe.Exec(context.Background())
	if err != nil {
		panic("Transaction failed")
	}
}

func (client *DBHandler) ValidateToken(token string) bool {
	return client.redisClient.HExists(context.Background(), dbkeys.Tokens, token).Val()
}

func (client *DBHandler) GetUser(token string) string {
	return client.redisClient.HGet(context.Background(), dbkeys.Tokens, token).Val()
}

func (client *DBHandler) RegisterGuess(username string, guess string) {
	client.redisClient.HSet(context.Background(), dbkeys.GuessesKey, username, guess)
}

func (client *DBHandler) CurrentLevel() string {
	currentLevel := client.redisClient.Get(context.Background(), dbkeys.CurrentLevel).Val()
	if currentLevel == "" {
		client.redisClient.GetSet(context.Background(), dbkeys.CurrentLevel, 1)
		return client.redisClient.Get(context.Background(), dbkeys.CurrentLevel).Val()
	}
	return currentLevel
}

func (client *DBHandler) IncrementLevel() IncrementLevelResponse {
	guesses, _ := client.redisClient.HGetAll(context.Background(), dbkeys.GuessesKey).Result()
	winner := getWinner(guesses)

	h := history{
		Winner:  winner,
		Level:   client.CurrentLevel(),
		Guesses: getGuessMap(guesses),
	}
	historyJson, _ := json.Marshal(h)
	cl := client.CurrentLevel()

	pipe := client.redisClient.TxPipeline()
	pipe.Del(context.Background(), dbkeys.GuessesKey)
	pipe.RPush(context.Background(), dbkeys.History, historyJson)
	pipe.Incr(context.Background(), dbkeys.CurrentLevel)

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
