package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/covid-19-ppe-tracker/hospital-module-kit/models"
)

func TestDeserializeJSONReturnsOK(t *testing.T) {
	u := models.User{}
	p := Handle("/users", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(DeserializerContextKey).(models.User)
		if u.Email == "blah@dumdum.com" {
			t.Errorf("chances of race condition")
		}
		if user.Email != "blah@dumdum.com" {
			t.Errorf("json deserializer returned wrong email: got %v want %v",
				user.Email, "blah@dumdum.com")
		}
		if user.Password != "dumdum" {
			t.Errorf("json deserializer returned wrong password: got %v want %v",
				user.Password, "dumdum")
		}
	}))
	rr := httptest.NewRecorder()
	p = p.WithJSONFromBodyIn(u)
	h := routes[string(p)]
	jsonBody := strings.NewReader("{\"email\":\"blah@dumdum.com\",\"password\":\"dumdum\"}")
	req, err := http.NewRequest("POST", "/users", jsonBody)
	if err != nil {
		t.Fatal(err)
	}
	h.ServeHTTP(rr, req)
}
