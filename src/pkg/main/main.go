package main

import (
	dh "github.com/cksharma11/guessing/src/pkg/db_handler"
	"log"
	"net/http"
	"time"
)

func main() {
	redisClient := dh.GetDBHandler(GetRedisClient())

	server := &http.Server{
		Handler:      router(&redisClient),
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	print("Started! ->\n")
	log.Fatal(server.ListenAndServe())
}
