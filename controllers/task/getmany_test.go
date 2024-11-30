package tasks_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	tc "github.com/tatiananeda/todo/controllers/task"
	"github.com/tatiananeda/todo/entities/web"
	r "github.com/tatiananeda/todo/repository"
	"github.com/tatiananeda/todo/services"
	u "github.com/tatiananeda/todo/utils/testutils"
)

func TestGetMany(t *testing.T) {
	tRepo := r.NewRepository()
	task1 := r.Task{
		Title:       title,
		Description: description,
		Due:         due,
		Id:          "4e21a827-c643-47f2-a66d-90c8c5152064",
		IsComplete:  false,
	}
	task2 := r.Task{
		Title:       title,
		Description: description,
		Due:         due,
		Id:          "5f5a0e15-be00-43cc-94ef-8366d2d0830a",
		IsComplete:  false,
	}
	task3 := r.Task{
		Title:       title,
		Description: description,
		Due:         due,
		Id:          "d700161f-9289-4d85-8942-d747a6b30718",
		IsComplete:  true,
	}
	task4 := r.Task{
		Title:       title,
		Description: description,
		Due:         due,
		Id:          "c9849c7d-1ac4-47d8-8a92-bc15d5adffee",
		IsComplete:  true,
	}
	tRepo.AddTask(&task1)
	tRepo.AddTask(&task2)
	tRepo.AddTask(&task3)
	tRepo.AddTask(&task4)

	tests := map[string]struct {
		params map[string]string
		check  func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		"With page and limit": {
			params: map[string]string{
				"limit": "2",
				"page":  "1",
			},
			check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				u.Check(t, rr.Code, http.StatusOK)
				var r []r.Task
				response, err := u.ParseResponse(rr.Body, r)
				if err != nil {
					t.Fatal(err)
				}
				u.Check(t, len(response), 2)
				for i, want := range tRepo.Tasks[0:2] {
					u.Check(t, response[i], *want)
				}
			},
		},
		"With completed": {
			params: map[string]string{
				"completed": "true",
			},
			check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				u.Check(t, rr.Code, http.StatusOK)
				var r []r.Task
				response, err := u.ParseResponse(rr.Body, r)
				if err != nil {
					t.Fatal(err)
				}
				u.Check(t, len(response), 2)
				for i, want := range tRepo.Tasks[2:4] {
					u.Check(t, response[i], *want)
				}
			},
		},
		"With completed, page and limit": {
			params: map[string]string{
				"completed": "false",
				"limit":     "2",
				"page":      "1",
			},
			check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				u.Check(t, rr.Code, http.StatusOK)
				var r []r.Task
				response, err := u.ParseResponse(rr.Body, r)
				if err != nil {
					t.Fatal(err)
				}
				u.Check(t, len(response), 2)
				for i, want := range tRepo.Tasks[0:2] {
					u.Check(t, response[i], *want)
				}
			},
		},
		"With no params": {
			params: map[string]string{},
			check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				u.Check(t, rr.Code, http.StatusOK)
				var r []r.Task
				response, err := u.ParseResponse(rr.Body, r)
				if err != nil {
					t.Fatal(err)
				}
				u.Check(t, len(response), 4)
				for i, want := range tRepo.Tasks {
					u.Check(t, response[i], *want)
				}
			},
		},
		"With completed not a boolean": {
			params: map[string]string{
				"completed": "wrong",
			},
			check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				u.Check(t, rr.Code, http.StatusBadRequest)

				var r web.APIError
				response, err := u.ParseResponse(rr.Body, r)
				if err != nil {
					t.Fatal(err)
				}
				u.Check(t, response, *web.InvalidField("completed must be a boolean"))
			},
		},
		"With limit not an int": {
			params: map[string]string{
				"limit": "wrong",
				"page":  "1",
			},
			check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				u.Check(t, rr.Code, http.StatusBadRequest)

				var r web.APIError
				response, err := u.ParseResponse(rr.Body, r)
				if err != nil {
					t.Fatal(err)
				}
				u.Check(t, response, *web.InvalidField("Limit must be a positive integer"))
			},
		},
		"With page not an int": {
			params: map[string]string{
				"page":  "wrong",
				"limit": "2",
			},
			check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				u.Check(t, rr.Code, http.StatusBadRequest)

				var r web.APIError
				response, err := u.ParseResponse(rr.Body, r)
				if err != nil {
					t.Fatal(err)
				}
				u.Check(t, response, *web.InvalidField("Page must be a positive integer"))
			},
		},
		"With page < 1": {
			params: map[string]string{
				"page":  "-3",
				"limit": "2",
			},
			check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				u.Check(t, rr.Code, http.StatusBadRequest)

				var r web.APIError
				response, err := u.ParseResponse(rr.Body, r)
				if err != nil {
					t.Fatal(err)
				}
				u.Check(t, response, *web.InvalidField("Page must be a positive integer"))
			},
		},
		"With limit < 1": {
			params: map[string]string{
				"limit": "0",
				"page":  "2",
			},
			check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				u.Check(t, rr.Code, http.StatusBadRequest)

				var r web.APIError
				response, err := u.ParseResponse(rr.Body, r)
				if err != nil {
					t.Fatal(err)
				}
				u.Check(t, response, *web.InvalidField("Limit must be a positive integer"))
			},
		},
	}

	taskService := services.NewTaskService(tRepo)
	httpResponseService := services.NewHttpResponseService()
	getMany := tc.NewGetManyController(taskService, httpResponseService)
	handler := http.HandlerFunc(getMany.Handler)

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, path, nil)
			if err != nil {
				t.Fatal(err)
			}
			q := req.URL.Query()
			for k, v := range test.params {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			test.check(t, rr)
		})
	}
}

func BenchmarkTestGetMany(b *testing.B) {
	tRepo := r.NewRepository()
	taskService := services.NewTaskService(tRepo)
	httpResponseService := services.NewHttpResponseService()
	getMany := tc.NewGetManyController(taskService, httpResponseService)
	handler := http.HandlerFunc(getMany.Handler)

	b.ReportAllocs()
	req, _ := http.NewRequest("GET", path+"/all", nil)
	rr := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(rr, req)
	}
}
