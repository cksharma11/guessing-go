package dbhandler

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"strconv"
)

type history struct {
	Winner  string           `json:"winner"`
	Level   string           `json:"level"`
	Guesses map[int][]string `json:"guesses"`
}

type incrementLevelResponse struct {
	Result string `json:"result"`
	Pl     string `json:"pl"`
	Cl     string `json:"cl"`
}

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

func (client *DBHandler) IncrementLevel() ([]byte, error) {
	guesses, _ := client.redisClient.HGetAll(context.Background(), "guesses").Result()
	winner := getWinner(guesses)

	h := history{
		Winner:  winner,
		Level:   client.CurrentLevel(),
		Guesses: getGuessMap(guesses),
	}
	historyJson, _ := json.Marshal(h)
	cl := client.CurrentLevel()

	client.redisClient.Del(context.Background(), "guesses")
	client.redisClient.RPush(context.Background(), "history", historyJson)
	client.redisClient.Incr(context.Background(), "current-level")

	return json.Marshal(incrementLevelResponse{
		Result: string(historyJson),
		Pl:     cl,
		Cl:     client.CurrentLevel(),
	})
}

func getWinner(guesses map[string]string) string {
	guessMap := getGuessMap(guesses)
	maxGuess := 0
	var maxGuessBy []string
	for guess, guessedBy := range guessMap {
		if maxGuess < guess {
			maxGuessBy = guessedBy
		}
	}
	return maxGuessBy[rand.Intn(len(maxGuessBy))]
}

func getGuessMap(guesses map[string]string) map[int][]string {
	guessMap := make(map[int][]string)
	for _, key := range guesses {
		guess, _ := strconv.Atoi(key)
		guessMap[guess] = []string{}
	}
	for value, key := range guesses {
		guess, _ := strconv.Atoi(key)
		guessMap[guess] = append(guessMap[guess], value)
	}
	return guessMap
}
