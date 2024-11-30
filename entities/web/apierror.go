package web

import (
	"errors"
	"net/http"
)

type APIError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func (e APIError) Error() string {
	return e.Message
}

func NewAPIError(statusCode int, err error) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    err.Error(),
	}
}

func InvalidJSON(e error) *APIError {
	return NewAPIError(http.StatusBadRequest, e)
}

func NotFound(id string) *APIError {
	return NewAPIError(http.StatusNotFound, errors.New("Not found task with id "+id))
}

func InvalidField(msg string) *APIError {
	return NewAPIError(http.StatusBadRequest, errors.New(msg))
}
