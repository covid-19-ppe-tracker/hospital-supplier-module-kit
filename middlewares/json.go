package middlewares

import (
	"context"
	"net/http"

	"github.com/covid-19-ppe-tracker/hospital-module-kit/models"
)

type ContextJSONValue struct {
	D      models.JSONSerializer
	Status int
}

func (p pattern) SerializeJSON() pattern {
	decorate := handleSerializeJSON()
	ro := routes[string(p)]
	h := decorate(ro)
	routes[string(p)] = h
	return p
}

func handleSerializeJSON() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			h.ServeHTTP(w, r)
			jsonContextValue := r.Context().Value(SerializerContextKey).(ContextJSONValue)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(jsonContextValue.Status)
			err := jsonContextValue.D.SerializeToJSON(w)
			if err != nil {
				//profile this error. eat the returned error
				//_ = getdone_log.Profile("Error in serializing response: " + err.Error())
			}
		})
	}
}

func (p pattern) WithJSONFromBodyIn(d models.JSONBodyDeserializer) pattern {
	decorate := handleDeserializeJSON(d)
	ro := routes[string(p)]
	h := decorate(ro)
	routes[string(p)] = h
	return p
}

func handleDeserializeJSON(d models.JSONBodyDeserializer) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := d.DeserializeFromJSONInBody(r)
			ctx := context.WithValue(r.Context(), DeserializerContextKey, data)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func sendJSONResponse(w http.ResponseWriter, r *http.Request, status int, data models.JSONSerializer) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := data.SerializeToJSON(w)
	if err != nil {
		//profile this error. eat the returned error
		//_ = getdone_log.Profile("Error in serializing response: " + err.Error())
	}
}

func SendJSONError(w http.ResponseWriter, r *http.Request, status int, data models.JSONSerializer) {
	sendJSONResponse(w, r, status, data)
}

func SendJSONSuccess(w http.ResponseWriter, r *http.Request, data models.JSONSerializer) {
	sendJSONResponse(w, r, http.StatusOK, data)
}
