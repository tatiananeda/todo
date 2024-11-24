package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	c "github.com/tatiananeda/todo/controllers"
	r "github.com/tatiananeda/todo/repository"
	u "github.com/tatiananeda/todo/utils"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const path, title, description, due = "/tasks", "Do Homework", "Test coverage", "25.11.2024"

func getCreatePayload(t *testing.T, title string, description string, due string) io.Reader {
	data := r.Task{
		Title:       title,
		Description: description,
		Due:         due,
	}

	p, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	return bytes.NewBuffer(p)
}

func TestCreate(t *testing.T) {
	handler := http.HandlerFunc(c.Create)

	tests := map[string]struct {
		input  io.Reader
		result u.APIError
	}{
		"Invalid JSON": {
			input:  strings.NewReader(`{"test": wrong}`),
			result: *u.InvalidJSON(fmt.Errorf("invalid character 'w' looking for beginning of value")),
		},
		"No title": {
			input:  getCreatePayload(t, "", description, due),
			result: *u.InvalidField("title is required"),
		},
		"No due_date": {
			input:  getCreatePayload(t, title, description, ""),
			result: *u.InvalidField("due_date is required"),
		},
	}

	t.Run("Happy path", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, path, getCreatePayload(t, title, description, due))

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		u.Check(t, rr.Code, http.StatusOK)

		var r r.Task
		response, err := u.ParseResponse(rr.Body, r)

		if err != nil {
			t.Fatal(err)
		}

		u.Check(t, response.Description, description)
		u.Check(t, response.Due, due)
		u.Check(t, response.Title, title)
		u.Check(t, response.IsComplete, false)

		if response.Id == "" {
			t.Errorf("Create: expected task created with id")
		}
	})

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, path, test.input)

			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			u.Check(t, rr.Code, http.StatusBadRequest)

			var r u.APIError
			response, err := u.ParseResponse(rr.Body, r)

			if err != nil {
				t.Fatal(err)
			}

			u.Check(t, test.result, response)
		})
	}
}
