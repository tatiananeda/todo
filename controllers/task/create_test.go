package tasks_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	tc "github.com/tatiananeda/todo/controllers/task"
	"github.com/tatiananeda/todo/entities/web"
	r "github.com/tatiananeda/todo/repository"
	"github.com/tatiananeda/todo/services"
	u "github.com/tatiananeda/todo/utils/testutils"
)

const path, title, description, due = "/tasks", "Do Homework", "Test coverage", "25.11.2024"

func TestCreate(t *testing.T) {
	tests := map[string]struct {
		input io.Reader
		check func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		"Happy Path": {
			input: getCreatePayload(t, title, description, due),
			check: func(t *testing.T, rr *httptest.ResponseRecorder) {
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
			},
		},
		"Invalid JSON": {
			input: strings.NewReader(`{"test": wrong}`),
			check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				u.Check(t, rr.Code, http.StatusBadRequest)
				var r web.APIError
				response, err := u.ParseResponse(rr.Body, r)
				if err != nil {
					t.Fatal(err)
				}
				u.Check(t, *web.InvalidJSON(fmt.Errorf("invalid character 'w' looking for beginning of value")), response)
			},
		},
		"No title": {
			input: getCreatePayload(t, "", description, due),
			check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				u.Check(t, rr.Code, http.StatusBadRequest)
				var r web.APIError
				response, err := u.ParseResponse(rr.Body, r)
				if err != nil {
					t.Fatal(err)
				}
				u.Check(t, *web.InvalidField("title is required"), response)
			},
		},
		"No due_date": {
			input: getCreatePayload(t, title, description, ""),
			check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				u.Check(t, rr.Code, http.StatusBadRequest)
				var r web.APIError
				response, err := u.ParseResponse(rr.Body, r)
				if err != nil {
					t.Fatal(err)
				}
				u.Check(t, *web.InvalidField("due_date is required"), response)
			},
		},
	}

	tRepo := r.NewRepository()
	taskService := services.NewTaskService(tRepo)
	httpResponseService := services.NewHttpResponseService()
	create := tc.NewCreateController(taskService, httpResponseService)
	handler := http.HandlerFunc(create.Handler)

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, path, test.input)

			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			test.check(t, rr)
		})
	}
}

func getCreatePayload(t *testing.T, title string, description string, due string) io.Reader {
	data := web.TaskInput{
		Title:       title,
		Description: description,
		DueDate:     due,
	}

	p, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	return bytes.NewBuffer(p)
}
