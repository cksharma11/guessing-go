package handler

import (
	"fmt"
	"net/http"
)

// HelloAPI : Hello from guessing server
func HelloAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello API")
}

// Signup : Register user
func Signup(w http.ResponseWriter, r *http.Request) {
}
