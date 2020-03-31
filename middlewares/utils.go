package middlewares

import (
	"fmt"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404 dude")
	}
	if status == http.StatusUnauthorized {
		fmt.Fprint(w, "401 dude")
	}
}
