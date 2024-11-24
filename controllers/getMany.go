package controllers

import (
	repo "github.com/tatiananeda/todo/repository"
	"github.com/tatiananeda/todo/utils"
	"net/http"
	"strconv"
)

var GetMany = utils.WithErrorHandling(getMany)

func getMany(w http.ResponseWriter, r *http.Request) error {
	queryParams := r.URL.Query()
	l := queryParams.Get("limit")
	p := queryParams.Get("page")
	c := queryParams.Get("completed")

	tasks := repo.Tasks

	if c != "" {
		completed, err := strconv.ParseBool(c)

		if err != nil {
			return utils.InvalidField("completed must be a boolean")
		}

		tasks = make([]*repo.Task, 0)

		for _, t := range repo.Tasks {
			if t.IsComplete == completed {
				tasks = append(tasks, t)
			}
		}
	}

	if l == "" && p == "" {
		if err := utils.WriteJSON(w, http.StatusOK, tasks); err != nil {
			return err
		}
		return nil
	}

	limit, err := strconv.ParseInt(l, 10, 64)
	if err != nil || limit < 1 {
		return utils.InvalidField("Limit must be a positive integer")
	}

	page, err := strconv.ParseInt(p, 10, 64)
	if err != nil || page < 1 {
		return utils.InvalidField("Page must be a positive integer")
	}

	startIdx := page - 1
	if startIdx > 0 {
		startIdx = startIdx * limit
	}

	if len(tasks) < int(startIdx) {
		if err := utils.WriteJSON(w, http.StatusOK, []*repo.Task{}); err != nil {
			return err
		}
		return nil
	}

	if err := utils.WriteJSON(w, http.StatusOK, tasks[startIdx:limit*page]); err != nil {
		return err
	}
	return nil
}
