package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func WithErrorHandling(handler APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			var e *APIError
			if errors.As(err, &e) {
				WriteJSON(w, e.StatusCode, err)
			} else {
				WriteJSON(w, http.StatusInternalServerError, InternalServerError)
			}
			log.Println("HTTP API error:", err.Error(), "method:", r.Method, "path:", r.URL.Path)
		}
	}
}
