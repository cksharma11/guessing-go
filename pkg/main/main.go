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
	router.HandleFunc("/guess/{guess}", handlers.WrapAuth(handlers.Guess)).Methods("POST")
	router.HandleFunc("/current-level", handlers.WrapAuth(handlers.CurrentLevel)).Methods("GET")

	adminRoutes := router.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(h.WrapAdminAuth)
	adminRoutes.HandleFunc("/increment-level", handlers.IncrementLevel).Methods("POST")
	adminRoutes.HandleFunc("/current-level", handlers.CurrentLevel).Methods("GET")

	server := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
