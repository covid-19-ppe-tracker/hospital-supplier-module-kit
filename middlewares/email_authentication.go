package middlewares

import (
	"context"
	"net/http"

	"github.com/covid-19-ppe-tracker/hospital-module-kit/models"
)

func (p pattern) WithEmailAuthentication() pattern {
	//email authentication
	decorate := authenticateWithEmail()
	ro := routes[string(p)]
	h := decorate(ro)
	routes[string(p)] = h
	return p
}

func authenticateUser(r *http.Request) (models.User, bool, error) {
	i := r.Context().Value(DeserializerContextKey)
	if i == nil {
		return models.User{}, false, nil
	}
	user, ok := i.(models.User)
	if !ok {
		return models.User{}, false, nil
	}
	user, ok, err := user.AuthenticateEmail()
	if err != nil || !ok {
		return user, false, err
	}
	return user, true, nil
}

func authenticateWithEmail() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, isAuthenticated, err := authenticateUser(r)
			if err == models.ErrNoMongoConnection {
				user = models.User{}
				user.ErrorMessage = "Under maintenance. We'll be back shortly"
				SendJSONError(w, r, http.StatusUnauthorized, user)
				return
			}

			if !isAuthenticated {
				user = models.User{}
				user.ErrorMessage = "Invalid email or password"
				SendJSONError(w, r, http.StatusUnauthorized, user)
				return
			}

			//user is not a pointer so can't ne nil
			//noinspection GoNilness
			user.Password = ""

			//set user in request context
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
