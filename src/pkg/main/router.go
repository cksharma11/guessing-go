package main

import (
	dh "github.com/cksharma11/guessing/src/pkg/db_handler"
	h "github.com/cksharma11/guessing/src/pkg/handler"
	"github.com/gorilla/mux"
)

func router(dbClient *dh.DBHandler) *mux.Router {
	handlers := h.NewHandlerContext(dbClient)

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/signup/{username}", handlers.SignUp).Methods("POST")
	router.HandleFunc("/guess/{guess}", handlers.WrapAuth(handlers.Guess)).Methods("POST")
	router.HandleFunc("/current-level", handlers.WrapAuth(handlers.CurrentLevel)).Methods("GET")

	adminRoutes := router.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(h.WrapAdminAuth)
	adminRoutes.HandleFunc("/increment-level", handlers.IncrementLevel).Methods("POST")
	adminRoutes.HandleFunc("/current-level", handlers.CurrentLevel).Methods("GET")

	return router
}
