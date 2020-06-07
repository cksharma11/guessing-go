package handler

import (
	"fmt"
	"net/http"
)

func HelloAPI(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf(w, "Hello API")
}