package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	m "github.com/tatiananeda/todo/middleware"
	"github.com/tatiananeda/todo/services"
	u "github.com/tatiananeda/todo/utils/testutils"
)

func TestRecoverPanicMiddleware(t *testing.T) {
	next := func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	}
	req, err := http.NewRequest(http.MethodGet, "/tasks", nil)

	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	httpResponseService := services.NewHttpResponseService()
	mid := m.NewRecoverPanicMiddleware(httpResponseService)
	h := mid.Middleware(http.HandlerFunc(next))
	h.ServeHTTP(rr, req)
	u.Check(t, rr.Code, http.StatusInternalServerError)
}
