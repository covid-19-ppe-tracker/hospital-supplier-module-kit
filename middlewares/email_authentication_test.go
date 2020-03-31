package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/covid-19-ppe-tracker/hospital-module-kit/models"
)

func TestEmailAuthenticationReturnsOK(t *testing.T) {
	p := Handle("/users", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserContextKey).(models.User)
		if user.Email != "blah@dumdum.com" {
			t.Errorf("email authentication returned wrong email: got %v want %v",
				user.Email, "blah@dumdum.com")
		}
	}))
	u := &models.User{}
	rr := httptest.NewRecorder()
	p = p.WithEmailAuthentication().WithJSONFromBodyIn(u)
	h := routes[string(p)]
	jsonBody := strings.NewReader("{\"email\":\"blah@dumdum.com\",\"password\":\"dumdum\"}")
	req, err := http.NewRequest("POST", "/users", jsonBody)
	if err != nil {
		t.Fatal(err)
	}
	h.ServeHTTP(rr, req)
}
