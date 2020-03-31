package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/covid-19-ppe-tracker/hospital-module-kit/models"
	"github.com/dgrijalva/jwt-go"
)

func (p pattern) WithJWTAuthentication(secret string) pattern {
	//jwt authentication
	decorate := authenticateWithJWT(secret)
	ro := routes[string(p)]
	h := decorate(ro)
	routes[string(p)] = h
	return p
}

func authenticateWithJWT(secret string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ah := r.Header.Get("authorization")
			if ah == "" {
				cookie, err := r.Cookie("session")
				if err == nil {
					ah = "BEARER " + cookie.Value
				}
			}
			if ah != "" {
				// Should be a bearer token
				if len(ah) > 6 && strings.ToUpper(ah[0:7]) == "BEARER " {
					parsedToken, err := jwt.Parse(ah[7:], func(token *jwt.Token) (interface{}, error) {
						return []byte(secret), nil
					})

					if err == nil && parsedToken.Valid {
						claims := parsedToken.Claims.(jwt.MapClaims)
						//user, isAuthenticated, err := models.AuthenticateWithJWT(claims["id"], r)

						//if err == models.ErrNoMongoConnection {
						//	user = &models.User{}
						//	user.ErrorMessage = "Under maintenance. We'll be back shortly"
						//	SendJSONError(w, r, http.StatusUnauthorized, user)
						//	return
						//}
						//if !isAuthenticated {
						//	user = &models.User{}
						//	user.ErrorMessage = "Not authorized"
						//	SendJSONError(w, r, http.StatusUnauthorized, user)
						//	return
						//}
						user := claims["user"].(models.User)
						//if user != nil {
						user.Password = ""
						//}
						ctx := context.WithValue(r.Context(), UserContextKey, user)
						h.ServeHTTP(w, r.WithContext(ctx))
						return
					}
				}
			}
			user := &models.User{}
			user.ErrorMessage = "Not authorized"
			SendJSONError(w, r, http.StatusUnauthorized, user)
		})
	}
}
