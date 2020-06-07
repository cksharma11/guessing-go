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

func (redisClient *Context) SignUp(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	present := redisClient.redisClient.IsUserPresent(username)
	if present {
		_, _ = fmt.Fprint(w, "Username already registered")
		return
	}
	token := uuid.New().String()
	redisClient.redisClient.AssociateToken(username, token)
	_, _ = fmt.Fprint(w, token)
}
