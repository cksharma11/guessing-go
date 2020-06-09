package dbhandler

import "github.com/go-redis/redis/v8"

type history struct {
	Winner  string           `json:"winner"`
	Level   string           `json:"level"`
	Guesses map[int][]string `json:"guesses"`
}

type IncrementLevelResponse struct {
	Result interface{} `json:"result"`
	Pl     string `json:"pl"`
	Cl     string `json:"cl"`
}

type DBHandler struct {
	redisClient *redis.Client
}
