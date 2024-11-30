package test_utils

import (
	"bytes"
	"encoding/json"
	"github.com/tatiananeda/todo/entities/web"
	r "github.com/tatiananeda/todo/repository"
	"testing"
)

type ResponseType interface {
	web.TaskInput | web.APIError | []r.Task | r.Task
}

type AssertionType interface {
	web.APIError | string | int | bool | r.Task
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
