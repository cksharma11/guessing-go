package handler

import dbHandler "github.com/cksharma11/guessing/src/pkg/db_handler"

type Context struct {
	redisClient *dbHandler.DBHandler
}

type response struct {
	Message string      `json:"message"`
	Err     bool        `json:"err"`
	Data    interface{} `json:"data"`
}
