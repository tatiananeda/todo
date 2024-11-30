package main

import (
	"net/http"

	"github.com/gorilla/mux"
	tc "github.com/tatiananeda/todo/controllers/task"
	"github.com/tatiananeda/todo/middleware"
	"github.com/tatiananeda/todo/repository"
	"github.com/tatiananeda/todo/services"
)

func main() {
	tRepo := repository.NewRepository()
	taskService := services.NewTaskService(tRepo)
	httpResponseService := services.NewHttpResponseService()
	getOne := tc.NewGetOneController(taskService, httpResponseService)
	create := tc.NewCreateController(taskService, httpResponseService)
	getMany := tc.NewGetManyController(taskService, httpResponseService)
	update := tc.NewGetManyController(taskService, httpResponseService)
	delete := tc.NewGetManyController(taskService, httpResponseService)
	patch := tc.NewGetManyController(taskService, httpResponseService)
	recoverPanicMiddleware := middleware.NewRecoverPanicMiddleware(httpResponseService)

	router := mux.NewRouter()
	router.Use(recoverPanicMiddleware)

	router.HandleFunc("/tasks/{id}", getOne.Handler).Methods(http.MethodGet)
	router.HandleFunc("/tasks", getMany.Handler).Methods(http.MethodGet)
	router.HandleFunc("/tasks/{id}", update.Handler).Methods(http.MethodPut)
	router.HandleFunc("/tasks/{id}", patch.Handler).Methods(http.MethodPatch)
	router.HandleFunc("/tasks/{id}", delete.Handler).Methods(http.MethodDelete)
	router.HandleFunc("/tasks", create.Handler).Methods(http.MethodPost)

	http.ListenAndServe(":3032", router)
}
