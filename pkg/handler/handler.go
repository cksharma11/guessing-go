package handler

import (
	"fmt"
	dbHandler "github.com/cksharma11/guessing/pkg/db_handler"
	"github.com/gorilla/mux"
	"net/http"
)

type Context struct {
	redisClient *dbHandler.DBHandler
}

func NewHandlerContext(redisClient *dbHandler.DBHandler) *Context {
	if redisClient == nil {
		panic("nil MongoDB session!")
	}
	return &Context{redisClient: redisClient}
}

func (redisClient *Context) HelloAPI(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello API")
}

func (redisClient *Context) SignUp(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	present := redisClient.redisClient.IsUserPresent(username)
	if present {
		_, _ = fmt.Fprint(w, "Username already registered")
		return
	}
	redisClient.redisClient.AssociateToken(username, "123")
	_, _ = fmt.Fprint(w, "123")
}
