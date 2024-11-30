package tasks

import (
	"github.com/gorilla/mux"
	"github.com/tatiananeda/todo/services"
	"net/http"
)

type DeleteController struct {
	taskService         *services.TaskService
	httpResponseService *services.HttpResponseService
}

func NewDeleteController(ts *services.TaskService, hs *services.HttpResponseService) *DeleteController {
	c := DeleteController{
		taskService:         ts,
		httpResponseService: hs,
	}
	return &c
}

func (c DeleteController) Handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := c.taskService.Delete(id)
	if err != nil {
		c.httpResponseService.HandleErrorResponse(w, r, err)
		return
	}
	c.httpResponseService.WriteJSON(w, http.StatusOK, nil)
}
