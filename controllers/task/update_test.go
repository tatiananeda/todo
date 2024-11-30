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

	"github.com/gorilla/mux"
	tc "github.com/tatiananeda/todo/controllers/task"
	"github.com/tatiananeda/todo/entities/web"
	r "github.com/tatiananeda/todo/repository"
	"github.com/tatiananeda/todo/services"
	u "github.com/tatiananeda/todo/utils/testutils"
)

func TestUpdate(t *testing.T) {
	tRepo := r.NewRepository()
	task := r.Task{
		Title:       "Another Title",
		Description: "Some optional description",
		Due:         "12.01.2024",
		Id:          "321bc166-1279-4b29-8b20-64994dba752f",
		IsComplete:  false,
	}
	tRepo.AddTask(&task)
	fakeId := "d7833d46-3057-43da-b774-210a8747e884"

	tests := map[string]struct {
		id    string
		input io.Reader
		check func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		"Allows updating title": {
			id:    task.Id,
			input: getUpdatePayload(t, title, ""),
			check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				u.Check(t, rr.Code, http.StatusOK)
				var r r.Task
				response, err := u.ParseResponse(rr.Body, r)
				if err != nil {
					t.Fatal(err)
				}
				u.Check(t, task, response)
			},
		},
		"Allows updating description": {
			id:    task.Id,
			input: getUpdatePayload(t, "", description),
			check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				u.Check(t, rr.Code, http.StatusOK)
				var r r.Task
				response, err := u.ParseResponse(rr.Body, r)
				if err != nil {
					t.Fatal(err)
				}
				u.Check(t, task, response)
			},
		},
		"Allows updating both title & description": {
			id:    task.Id,
			input: getUpdatePayload(t, title, description),
			check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				u.Check(t, rr.Code, http.StatusOK)
				var r r.Task
				response, err := u.ParseResponse(rr.Body, r)
				if err != nil {
					t.Fatal(err)
				}
				u.Check(t, task, response)
			},
		},
		"Invalid JSON": {
			id:    task.Id,
			input: strings.NewReader(`{"test": wrong}`),
			check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				u.Check(t, rr.Code, http.StatusBadRequest)
				var r web.APIError
				response, err := u.ParseResponse(rr.Body, r)
				if err != nil {
					t.Fatal(err)
				}
				u.Check(
					t,
					*web.InvalidJSON(fmt.Errorf("invalid character 'w' looking for beginning of value")),
					response,
				)
			},
		},
		"Not Found": {
			id:    fakeId,
			input: getUpdatePayload(t, title, description),
			check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				u.Check(t, rr.Code, http.StatusNotFound)
				var r web.APIError
				response, err := u.ParseResponse(rr.Body, r)
				if err != nil {
					t.Fatal(err)
				}
				u.Check(
					t,
					*web.NotFound(fakeId),
					response,
				)
			},
		},
	}

	taskService := services.NewTaskService(tRepo)
	httpResponseService := services.NewHttpResponseService()
	update := tc.NewUpdateController(taskService, httpResponseService)
	handler := http.HandlerFunc(update.Handler)

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPut, path+"/"+test.id, test.input)
			req = mux.SetURLVars(req, map[string]string{"id": test.id})
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			test.check(t, rr)
		})
	}
}

func getUpdatePayload(t *testing.T, title string, description string) io.Reader {
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
