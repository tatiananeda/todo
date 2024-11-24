package utils

import (
	"bytes"
	"encoding/json"
	r "github.com/tatiananeda/todo/repository"
	"testing"
)

type ResponseType interface {
	r.Task | APIError | []r.Task
}

type AssertionType interface {
	APIError | string | int | bool | r.Task
}

func ParseResponse[T ResponseType](b *bytes.Buffer, r T) (T, error) {
	var result T
	err := json.NewDecoder(b).Decode(&result)
	return result, err
}

func Check[T AssertionType](t *testing.T, got, want T) {
	if got != want {
		t.Errorf("Expected %v; got %v", want, got)
	}
}
