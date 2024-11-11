package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tatiananeda/todo/controllers"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/tasks/all", controllers.GetMany).Methods("GET")
	router.HandleFunc("/tasks/{id}", controllers.Update).Methods("PUT")
	router.HandleFunc("/tasks/{id}", controllers.MarkComplete).Methods("PATCH")
	router.HandleFunc("/tasks/{id}", controllers.Delete).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", controllers.GetOne).Methods("GET")
	router.HandleFunc("/tasks", controllers.Create).Methods("POST")

	http.ListenAndServe(":3032", router)
}
