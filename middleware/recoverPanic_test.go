package middleware_test

import (
	m "github.com/tatiananeda/todo/middleware"
	u "github.com/tatiananeda/todo/utils"
	"net/http"
	"net/http/httptest"
	"testing"
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

	h := m.RecoverPanicMiddleware(http.HandlerFunc(next))
	h.ServeHTTP(rr, req)
	u.Check(t, rr.Code, http.StatusInternalServerError)
}
