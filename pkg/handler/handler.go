package handler

import (
	"fmt"
	dbHandler "github.com/cksharma11/guessing/pkg/db_handler"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

type Context struct {
	redisClient *dbHandler.DBHandler
}

func NewHandlerContext(redisClient *dbHandler.DBHandler) *Context {
	if redisClient == nil {
		panic("nil redisClient!")
	}
	return &Context{redisClient: redisClient}
}

func (redisClient *Context) WrapAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("auth")
		isValidToken := redisClient.redisClient.ValidateToken(token)
		if isValidToken {
			username := redisClient.redisClient.GetUser(token)
			r.Header.Set("username", username)
			next.ServeHTTP(w, r)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprint(w, "Invalid Auth")
	}
}

func (redisClient *Context) SignUp(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	present := redisClient.redisClient.IsUserPresent(username)
	if present {
		_, _ = fmt.Fprint(w, "Username already registered")
		return
	}
	token := uuid.New().String()
	redisClient.redisClient.AssociateToken(token, username)
	w.WriteHeader(http.StatusCreated)
	_, _ = fmt.Fprint(w, token)
}

func (redisClient *Context) Guess(w http.ResponseWriter, r *http.Request) {
	guess := mux.Vars(r)["guess"]
	username := r.Header.Get("username")
	redisClient.redisClient.RegisterGuess(username, guess)
	_, _ = fmt.Fprint(w, "Guess has been registered")
}
