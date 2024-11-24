package middleware

import (
	"github.com/tatiananeda/todo/utils"
	"log"
	"net/http"
)

func RecoverPanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				log.Println("Panic:", err, "method:", r.Method, "path:", r.URL.Path)

				utils.WriteJSON(w, http.StatusInternalServerError, utils.InternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
