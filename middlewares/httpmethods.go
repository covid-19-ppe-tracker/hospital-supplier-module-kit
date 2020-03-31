package middlewares

import "net/http"

func (p pattern) WithMethods(m ...string) pattern {
	decorate := checkMethods(m)
	ro := routes[string(p)]
	h := decorate(ro)
	routes[string(p)] = h
	return p
}

func checkMethods(m []string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, method := range m {
				if method == r.Method {
					h.ServeHTTP(w, r)
					return
				}
			}
			ErrorHandler(w, r, http.StatusNotFound)
			return
		})
	}
}
