package tasks

import (
	"net/http"
	"strconv"

	"github.com/tatiananeda/todo/entities/web"
	"github.com/tatiananeda/todo/services"
)

type GetManyController struct {
	taskService         *services.TaskService
	httpResponseService *services.HttpResponseService
}

func NewGetManyController(ts *services.TaskService, hs *services.HttpResponseService) *GetManyController {
	c := GetManyController{
		taskService:         ts,
		httpResponseService: hs,
	}
	return &c
}

func (c GetManyController) Handler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	l := queryParams.Get("limit")
	p := queryParams.Get("page")
	cFilter := queryParams.Get("completed")

	tasks := c.taskService.GetAll()

	if cFilter != "" {
		completed, err := strconv.ParseBool(cFilter)

		if err != nil {
			c.httpResponseService.HandleErrorResponse(w, r, web.InvalidField("completed must be a boolean"))
			return
		}

		tasks = c.taskService.GetFilteredByCompleted(completed)
	}

	if l == "" && p == "" {
		c.httpResponseService.WriteJSON(w, http.StatusOK, tasks)
		return
	}

	tasks, err := c.taskService.GetPage(l, p, tasks)
	if err != nil {
		c.httpResponseService.HandleErrorResponse(w, r, err)
		return
	}
	c.httpResponseService.WriteJSON(w, http.StatusOK, tasks)
}
