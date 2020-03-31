package middlewares

import (
	"net/http"
	"time"
)

type TimeoutCallback func()

type responseWriterWithStatusCode struct {
	http.ResponseWriter
	statusCode int
}

func (rws *responseWriterWithStatusCode) WriteHeader(code int) {
	rws.statusCode = code
	rws.ResponseWriter.WriteHeader(code)
}

func (p pattern) WithTimeout(dt time.Duration, msg string, f func()) pattern {
	decorate := handleTimeout(dt, msg, f)
	ro := routes[string(p)]
	h := decorate(ro)
	routes[string(p)] = h
	return p
}

func handleTimeout(dt time.Duration, msg string, f TimeoutCallback) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wS := &responseWriterWithStatusCode{}
			wS.ResponseWriter = w
			http.TimeoutHandler(h, dt, msg).ServeHTTP(wS, r)
			if wS.statusCode == http.StatusServiceUnavailable {
				f()
			}
		})
	}
}
