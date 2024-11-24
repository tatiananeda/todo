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

func TestMarkComplete(t *testing.T) {
	handler := http.HandlerFunc(c.MarkComplete)
	task := r.Task{
		Title:       title,
		Description: description,
		Due:         due,
		Id:          "443387dd-47a1-42c0-bb73-50da2f1f7ce7",
		IsComplete:  false,
	}

	r.Tasks = append(r.Tasks, &task)

	fakeId := "d7833d46-3057-43da-b774-210a8747e884"

	t.Run("Happy Path", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPatch, path+"/"+task.Id, nil)
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

		u.Check(t, response, task)
		u.Check(t, response.IsComplete, true)
	})

	t.Run("Not found", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPatch, path+"/"+fakeId, nil)
		req = mux.SetURLVars(req, map[string]string{"id": fakeId})

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		u.Check(t, rr.Code, http.StatusNotFound)

		var r u.APIError
		response, err := u.ParseResponse(rr.Body, r)
		if err != nil {
			t.Fatal(err)
		}

		u.Check(t, response, *u.NotFound(fakeId))
	})
}
