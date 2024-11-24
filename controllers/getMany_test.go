package controllers_test

import (
	c "github.com/tatiananeda/todo/controllers"
	r "github.com/tatiananeda/todo/repository"
	u "github.com/tatiananeda/todo/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMany(t *testing.T) {
	handler := http.HandlerFunc(c.GetMany)

	r.Tasks = make([]*r.Task, 0)

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

	r.Tasks = append(r.Tasks, &task1, &task2, &task3, &task4)

	// fakeId := "d7833d46-3057-43da-b774-210a8747e884"

	t.Run("Happy path", func(t *testing.T) {
		okCases := map[string]struct {
			params       map[string]string
			statusCode   int
			expectLength int
			matchTasks   []*r.Task
		}{
			"With page and limit": {
				statusCode: http.StatusOK,
				params: map[string]string{
					"limit": "2",
					"page":  "1",
				},
				expectLength: 2,
				matchTasks:   r.Tasks[0:2],
			},
			"With completed": {
				statusCode: http.StatusOK,
				params: map[string]string{
					"completed": "true",
				},
				expectLength: 2,
				matchTasks:   r.Tasks[2:4],
			},
			"With completed, page and limit": {
				statusCode: http.StatusOK,
				params: map[string]string{
					"completed": "false",
					"limit":     "2",
					"page":      "1",
				},
				expectLength: 2,
				matchTasks:   r.Tasks[0:2],
			},
			"With no params": {
				statusCode:   http.StatusOK,
				params:       map[string]string{},
				expectLength: 4,
				matchTasks:   r.Tasks,
			},
		}

		for name, test := range okCases {
			t.Run(name, func(t *testing.T) {
				req, err := http.NewRequest(http.MethodGet, path+"/all", nil)

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

				u.Check(t, rr.Code, http.StatusOK)

				var r []r.Task
				response, err := u.ParseResponse(rr.Body, r)
				if err != nil {
					t.Fatal(err)
				}
				u.Check(t, len(response), test.expectLength)

				for i, want := range test.matchTasks {
					u.Check(t, response[i], *want)
				}
			})
		}
	})

	t.Run("Error path", func(t *testing.T) {
		errorCases := map[string]struct {
			params     map[string]string
			statusCode int
			result     u.APIError
		}{
			"With completed not a boolean": {
				statusCode: http.StatusBadRequest,
				params: map[string]string{
					"completed": "wrong",
				},
				result: *u.InvalidField("completed must be a boolean"),
			},
			"With limit not an int": {
				statusCode: http.StatusBadRequest,
				params: map[string]string{
					"limit": "wrong",
					"page":  "1",
				},
				result: *u.InvalidField("Limit must be a positive integer"),
			},
			"With page not an int": {
				statusCode: http.StatusBadRequest,
				params: map[string]string{
					"page":  "wrong",
					"limit": "2",
				},
				result: *u.InvalidField("Page must be a positive integer"),
			},
			"With page < 1": {
				statusCode: http.StatusBadRequest,
				params: map[string]string{
					"page":  "-3",
					"limit": "2",
				},
				result: *u.InvalidField("Page must be a positive integer"),
			},
			"With limit < 1": {
				statusCode: http.StatusBadRequest,
				params: map[string]string{
					"limit": "0",
					"page":  "2",
				},
				result: *u.InvalidField("Limit must be a positive integer"),
			},
		}

		for name, test := range errorCases {
			t.Run(name, func(t *testing.T) {
				req, err := http.NewRequest(http.MethodGet, path+"/all", nil)

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

				u.Check(t, rr.Code, test.statusCode)

				var r u.APIError
				response, err := u.ParseResponse(rr.Body, r)
				if err != nil {
					t.Fatal(err)
				}
				u.Check(t, response, test.result)
			})
		}
	})
}

func BenchmarkTestGetMany(b *testing.B) {
	b.ReportAllocs()

	handler := http.HandlerFunc(c.GetMany)
	req, _ := http.NewRequest("GET", path+"/all", nil)
	rr := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(rr, req)
	}
}
