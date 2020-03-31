package middlewares

import "net/http"

func (p pattern) AllowCORS() pattern {
	decorate := handleCORS()
	ro := routes[string(p)]
	h := decorate(ro)
	routes[string(p)] = h
	return p
}

func handleCORS() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestOrigin := r.Header.Get("Origin")
			if r.Method == "OPTIONS" {
				requestMethod := r.Header.Get("Access-Control-Request-Method")
				requestHeaders := r.Header.Get("Access-Control-Request-Headers")
				if requestMethod != "" && requestHeaders != "" {
					w.Header().Set("Access-Control-Allow-Origin", requestOrigin)
					w.Header().Set("Access-Control-Allow-Headers", requestHeaders)
					w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
					w.Header().Set("Access-Control-Max-Age", "86400")
					w.Header().Set("Vary", "Accept-Encoding, Origin")
					w.Header().Set("Content-Encoding", "gzip")
					w.Header().Set("Content-Length", "0")
					w.Header().Set("Keep-Alive", "timeout=2, max=100")
					w.Header().Set("Connection", "Keep-Alive")
					w.Header().Set("Content-Type", "text/plain")
					return
				}
			}

			if requestOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", requestOrigin)
				w.Header().Set("Vary", "Accept-Encoding, Origin")
			}

			h.ServeHTTP(w, r)
		})
	}
}
