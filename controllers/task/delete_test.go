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

func TestDelete(t *testing.T) {
	tRepo := r.NewRepository()
	task := r.Task{
		Title:       title,
		Description: description,
		Due:         due,
		Id:          "4de53488-debd-490f-a65d-a15bf0473900",
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
				u.Check(t, *web.NotFound(fakeId), response)
			},
		},
	}

	taskService := services.NewTaskService(tRepo)
	httpResponseService := services.NewHttpResponseService()
	delete := tc.NewDeleteController(taskService, httpResponseService)
	handler := http.HandlerFunc(delete.Handler)

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, path+"/"+test.id, nil)
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
