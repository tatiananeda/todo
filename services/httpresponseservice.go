package services

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/tatiananeda/todo/entities/web"
)

type HttpResponseService struct{}

func NewHttpResponseService() *HttpResponseService {
	return &HttpResponseService{}
}

func (s HttpResponseService) WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func (s HttpResponseService) HandleErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Println("HTTP API error:", err.Error(), "method:", r.Method, "path:", r.URL.Path)
	var e *web.APIError
	if errors.As(err, &e) {
		s.WriteJSON(w, e.StatusCode, err)
	} else {
		s.WriteJSON(w, http.StatusInternalServerError, web.InternalServerError)
	}
}
