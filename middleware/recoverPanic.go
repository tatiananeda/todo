package middleware

import (
	"github.com/gorilla/mux"
	"github.com/tatiananeda/todo/entities/web"
	"github.com/tatiananeda/todo/services"
	"log"
	"net/http"
)

func NewRecoverPanicMiddleware(httpResponseService *services.HttpResponseService) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.Header().Set("Connection", "close")
					log.Println("Panic:", err, "method:", r.Method, "path:", r.URL.Path)

					httpResponseService.WriteJSON(w, http.StatusInternalServerError, web.InternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
