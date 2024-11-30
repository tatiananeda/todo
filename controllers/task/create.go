package tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tatiananeda/todo/entities/web"
	"github.com/tatiananeda/todo/services"
)

type CreateController struct {
	taskService         *services.TaskService
	httpResponseService *services.HttpResponseService
}

func NewCreateController(ts *services.TaskService, hs *services.HttpResponseService) *CreateController {
	c := CreateController{
		taskService:         ts,
		httpResponseService: hs,
	}
	return &c
}

func (c CreateController) Handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

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

	if input.Title == "" {
		c.httpResponseService.HandleErrorResponse(w, r, web.InvalidField("title is required"))
		return
	}

	if input.DueDate == "" {
		c.httpResponseService.HandleErrorResponse(w, r, web.InvalidField("due_date is required"))
		return
	}

	t, err := c.taskService.Create(&input)
	if err != nil {
		c.httpResponseService.HandleErrorResponse(w, r, err)
	}

	c.httpResponseService.WriteJSON(w, http.StatusOK, t)
}
