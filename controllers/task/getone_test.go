package tasks_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	tc "github.com/tatiananeda/todo/controllers/task"
	"github.com/tatiananeda/todo/entities/web"
	r "github.com/tatiananeda/todo/repository"
	"github.com/tatiananeda/todo/services"
	u "github.com/tatiananeda/todo/utils/testutils"
)

func TestGetOne(t *testing.T) {
	tRepo := r.NewRepository()
	task := r.Task{
		Title:       title,
		Description: description,
		Due:         due,
		Id:          "4e21a827-c643-47f2-a66d-90c8c5152064",
		IsComplete:  false,
	}
	tRepo.AddTask(&task)
	fakeId := "d7833d46-3057-43da-b774-210a8747e884"

	tests := map[string]struct {
		id    string
		check func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		"Happy Path": {
			id: task.Id,
			check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				u.Check(t, rr.Code, http.StatusOK)
				var r r.Task
				response, err := u.ParseResponse(rr.Body, r)
				if err != nil {
					t.Fatal(err)
				}

				u.Check(t, response, task)
			},
		},
		"Not found": {
			id: fakeId,
			check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				u.Check(t, rr.Code, http.StatusNotFound)
				var r web.APIError
				response, err := u.ParseResponse(rr.Body, r)
				if err != nil {
					t.Fatal(err)
				}
				u.Check(t, response, *web.NotFound(fakeId))
			},
		},
	}

	taskService := services.NewTaskService(tRepo)
	httpResponseService := services.NewHttpResponseService()
	getOne := tc.NewGetOneController(taskService, httpResponseService)
	handler := http.HandlerFunc(getOne.Handler)

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, path+"/"+test.id, nil)
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
