package handler

import (
	"encoding/json"
	dbHandler "github.com/cksharma11/guessing/pkg/db_handler"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

type Context struct {
	redisClient *dbHandler.DBHandler
}

type response struct {
	Message string `json:"message"`
	Err     bool   `json:"err"`
	Data    interface{} `json:"data"`
}

func sendResponse(w http.ResponseWriter, res response, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
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
		res := response{
			Message: "Username already registered",
			Err:     true,
			Data:    "",
		}
		sendResponse(w, res, http.StatusConflict)
		return
	}

	token := uuid.New().String()
	redisClient.redisClient.AssociateToken(username, token)

	res := response{
		Message: "User created",
		Err:     false,
		Data:    token,
	}
	sendResponse(w, res, http.StatusCreated)
}

func (redisClient *Context) Guess(w http.ResponseWriter, r *http.Request) {
	guess := mux.Vars(r)["guess"]
	username := r.Header.Get("username")
	redisClient.redisClient.RegisterGuess(username, guess)
	res := response{
		Message: "Guess has been registered",
		Err:     false,
		Data:    guess,
	}
	sendResponse(w, res, http.StatusCreated)
}

func (redisClient *Context) CurrentLevel(w http.ResponseWriter, r *http.Request) {
	level := redisClient.redisClient.CurrentLevel()
	res := response{
		Message: "Current level",
		Err:     false,
		Data:    level,
	}
	sendResponse(w, res, http.StatusOK)
}

func (redisClient *Context) IncrementLevel(w http.ResponseWriter, r *http.Request) {
	marshal := redisClient.redisClient.IncrementLevel()

	res := response{
		Message: "Level incremented",
		Err:     false,
		Data:    marshal,
	}

	sendResponse(w, res, http.StatusAccepted)
}
