package controllers_test

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
	c "github.com/tatiananeda/todo/controllers"
	r "github.com/tatiananeda/todo/repository"
	u "github.com/tatiananeda/todo/utils"
)

func getUpdatePayload(t *testing.T, title string, description string) io.Reader {
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

func TestUpdate(t *testing.T) {
	handler := http.HandlerFunc(c.Update)
	task := r.Task{
		Title:       "Another Title",
		Description: "Some optional description",
		Due:         "12.01.2024",
		Id:          "321bc166-1279-4b29-8b20-64994dba752f",
		IsComplete:  false,
	}

	r.Tasks = append(r.Tasks, &task)

	fakeId := "d7833d46-3057-43da-b774-210a8747e884"

	okCases := map[string]struct {
		input io.Reader
	}{
		"With just Title": {
			input: getUpdatePayload(t, title, ""),
		},
		"With just Description": {
			input: getUpdatePayload(t, "", description),
		},
		"With both Title & Description": {
			input: getUpdatePayload(t, title, description),
		},
	}

	errorCases := map[string]struct {
		id     string
		input  io.Reader
		result u.APIError
		status int
	}{
		"Invalid JSON": {
			id:     task.Id,
			input:  strings.NewReader(`{"test": wrong}`),
			result: *u.InvalidJSON(fmt.Errorf("invalid character 'w' looking for beginning of value")),
			status: http.StatusBadRequest,
		},
		"Not Found": {
			id:     fakeId,
			input:  getUpdatePayload(t, title, description),
			result: *u.NotFound(fakeId),
			status: http.StatusNotFound,
		},
	}

	for name, test := range okCases {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPut, path+"/"+task.Id, test.input)
			req = mux.SetURLVars(req, map[string]string{"id": task.Id})

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

			u.Check(t, task, response)
		})
	}

	for name, test := range errorCases {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPut, path+"/"+test.id, test.input)
			req = mux.SetURLVars(req, map[string]string{"id": test.id})

			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			u.Check(t, rr.Code, test.status)

			var r u.APIError
			response, err := u.ParseResponse(rr.Body, r)

			if err != nil {
				t.Fatal(err)
			}

			u.Check(t, test.result, response)
		})
	}

}
