package utils

import "net/http"

var InternalServerError = map[string]any{
	"statusCode": http.StatusInternalServerError,
	"message":    http.StatusText(http.StatusInternalServerError),
}
