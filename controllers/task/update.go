package tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tatiananeda/todo/entities/web"
	"github.com/tatiananeda/todo/services"
)

type UpdateController struct {
	taskService         *services.TaskService
	httpResponseService *services.HttpResponseService
}

func NewUpdateController(ts *services.TaskService, hs *services.HttpResponseService) *UpdateController {
	c := UpdateController{
		taskService:         ts,
		httpResponseService: hs,
	}
	return &c
}

func (c UpdateController) Handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	req, err := io.ReadAll(r.Body)

	if err != nil {
		c.httpResponseService.HandleErrorResponse(w, r, fmt.Errorf("Error parsing body %w;", err))
		return
	}

	var input web.TaskInput
	e := json.Unmarshal(req, &input)

	if e != nil {
		c.httpResponseService.HandleErrorResponse(w, r, web.InvalidJSON(e))
		return
	}

	t, err := c.taskService.GetById(id)
	if err != nil {
		c.httpResponseService.HandleErrorResponse(w, r, err)
		return
	}

	t, err = c.taskService.Update(id, &input)
	if err != nil {
		c.httpResponseService.HandleErrorResponse(w, r, err)
		return
	}
	c.httpResponseService.WriteJSON(w, http.StatusOK, t)
}
