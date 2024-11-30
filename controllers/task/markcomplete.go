package tasks

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tatiananeda/todo/services"
)

type MarkCompleteController struct {
	taskService         *services.TaskService
	httpResponseService *services.HttpResponseService
}

func NewMarkCompleteController(ts *services.TaskService, hs *services.HttpResponseService) *MarkCompleteController {
	c := MarkCompleteController{
		taskService:         ts,
		httpResponseService: hs,
	}
	return &c
}

func (c MarkCompleteController) Handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	t, err := c.taskService.GetById(id)
	if err != nil {
		c.httpResponseService.HandleErrorResponse(w, r, err)
		return
	}
	t.IsComplete = !t.IsComplete
	c.httpResponseService.WriteJSON(w, http.StatusOK, t)
}
