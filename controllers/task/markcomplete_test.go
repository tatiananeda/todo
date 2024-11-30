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

func TestMarkComplete(t *testing.T) {
	tRepo := r.NewRepository()
	task := r.Task{
		Title:       title,
		Description: description,
		Due:         due,
		Id:          "443387dd-47a1-42c0-bb73-50da2f1f7ce7",
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
				u.Check(t, response.IsComplete, true)
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
	mark := tc.NewMarkCompleteController(taskService, httpResponseService)
	handler := http.HandlerFunc(mark.Handler)

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPatch, path+"/"+test.id, nil)
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
