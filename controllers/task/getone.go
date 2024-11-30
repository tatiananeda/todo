package tasks

import (
	"github.com/gorilla/mux"
	"github.com/tatiananeda/todo/services"
	"net/http"
)

type GetOneController struct {
	taskService         *services.TaskService
	httpResponseService *services.HttpResponseService
}

func NewGetOneController(ts *services.TaskService, hs *services.HttpResponseService) *GetOneController {
	c := GetOneController{
		taskService:         ts,
		httpResponseService: hs,
	}
	return &c
}

func (c GetOneController) Handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	t, err := c.taskService.GetById(id)
	if err != nil {
		c.httpResponseService.HandleErrorResponse(w, r, err)
		return
	}
	c.httpResponseService.WriteJSON(w, http.StatusOK, t)
}
