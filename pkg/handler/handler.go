package handler

import (
	"encoding/json"
	dbHandler "github.com/cksharma11/guessing/pkg/db_handler"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Context struct {
	redisClient *dbHandler.DBHandler
}

type response struct {
	Message string `json:"message"`
	Err     bool   `json:"err"`
	Data    string `json:"data"`
}

func sendResponse(w http.ResponseWriter, res response) {
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
}

func NewHandlerContext(redisClient *dbHandler.DBHandler) *Context {
	if redisClient == nil {
		panic("nil redisClient!")
	}
	return &Context{redisClient: redisClient}
}

func (redisClient *Context) SignUp(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	present := redisClient.redisClient.IsUserPresent(username)
	if present {
		w.WriteHeader(http.StatusConflict)
		res := response{
			Message: "Username already registered",
			Err:     true,
			Data:    "",
		}
		sendResponse(w, res)
		return
	}

	token := uuid.New().String()
	redisClient.redisClient.AssociateToken(username, token)
	w.WriteHeader(http.StatusCreated)
	res := response{
		Message: "User created",
		Err:     false,
		Data:    token,
	}

	sendResponse(w, res)
}

func (redisClient *Context) Guess(w http.ResponseWriter, r *http.Request) {
	guess := mux.Vars(r)["guess"]
	username := r.Header.Get("username")
	redisClient.redisClient.RegisterGuess(username, guess)
	w.WriteHeader(http.StatusCreated)
	res := response{
		Message: "Guess has been registered",
		Err:     false,
		Data:    guess,
	}
	sendResponse(w, res)
}

func (redisClient *Context) CurrentLevel(w http.ResponseWriter, r *http.Request) {
	level := redisClient.redisClient.CurrentLevel()
	res := response{
		Message: "Current level",
		Err:     false,
		Data:    level,
	}
	sendResponse(w, res)
}

func (redisClient *Context) IncrementLevel(w http.ResponseWriter, r *http.Request) {
	marshal, err := json.Marshal(redisClient.redisClient.IncrementLevel())
	if err != nil {
		log.Fatal("Error in marshaling increment response")
	}

	responseData := string(marshal)
	res := response{
		Message: "Level incremented",
		Err:     false,
		Data:    responseData,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusAccepted)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
}
