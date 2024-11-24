package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tatiananeda/todo/controllers"
	"github.com/tatiananeda/todo/middleware"
)

func main() {
	router := mux.NewRouter()
	router.Use(middleware.RecoverPanicMiddleware)

	router.HandleFunc("/tasks/all", controllers.GetMany).Methods("GET")
	router.HandleFunc("/tasks/{id}", controllers.Update).Methods("PUT")
	router.HandleFunc("/tasks/{id}", controllers.MarkComplete).Methods("PATCH")
	router.HandleFunc("/tasks/{id}", controllers.Delete).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", controllers.GetOne).Methods("GET")
	router.HandleFunc("/tasks", controllers.Create).Methods("POST")

	http.ListenAndServe(":3032", router)
}
