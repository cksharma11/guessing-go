package main

import (
	"log"
	"net/http"
	"time"

	handler "github.com/cksharma11/guessing/pkg/handler"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", handler.HelloAPI)
	router.HandleFunc("/signup/{username}", handler.Signup).Methods("POST")

	server := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
