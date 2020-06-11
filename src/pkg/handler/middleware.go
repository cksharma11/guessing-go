package handler

import (
	"fmt"
	"net/http"
)

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

func WrapAdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("auth")
		if token == "test" {
			next.ServeHTTP(w, r)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprint(w, "Invalid Auth")
	})
}
