package main

import (
	dh "github.com/cksharma11/guessing/pkg/db_handler"
	"log"
	"net/http"
	"time"

	h "github.com/cksharma11/guessing/pkg/handler"

	"github.com/gorilla/mux"
)

func main() {
	redisClient := dh.GetDBHandler(GetRedisClient())
	handlers := h.NewHandlerContext(&redisClient)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/signup/{username}", handlers.SignUp).Methods("POST")

	server := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
