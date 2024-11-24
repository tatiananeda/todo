package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	c "github.com/tatiananeda/todo/controllers"
	r "github.com/tatiananeda/todo/repository"
	u "github.com/tatiananeda/todo/utils"
)

func TestDelete(t *testing.T) {
	handler := http.HandlerFunc(c.Delete)
	task := r.Task{
		Title:       title,
		Description: description,
		Due:         due,
		Id:          "4de53488-debd-490f-a65d-a15bf0473900",
		IsComplete:  false,
	}

	r.Tasks = append(r.Tasks, &task)

	fakeId := "d7833d46-3057-43da-b774-210a8747e884"

	tests := map[string]struct {
		id     string
		status int
		result u.APIError
	}{
		"Happy Path": {
			id:     task.Id,
			status: http.StatusOK,
		},
		"Not found": {
			id:     fakeId,
			status: http.StatusNotFound,
			result: *u.NotFound(fakeId),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, path+"/"+test.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": test.id})

			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			u.Check(t, rr.Code, test.status)

			if test.status != http.StatusOK {
				var r u.APIError
				response, err := u.ParseResponse(rr.Body, r)

				if err != nil {
					t.Fatal(err)
				}

				u.Check(t, test.result, response)
			}
		})
	}
}
