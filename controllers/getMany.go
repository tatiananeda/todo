package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	repo "github.com/tatiananeda/todo/repository"
)

func GetMany(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	l := queryParams.Get("limit")
	p := queryParams.Get("page")
	c := queryParams.Get("completed")

	tasks := repo.Tasks

	if c != "" {
		completed, err := strconv.ParseBool(c)

		if err != nil {
			http.Error(w, "completed must be a boolean", 400)
			return
		}

		tasks = make([]*repo.Task, 0)

		for _, t := range repo.Tasks {
			if t.IsComplete == completed {
				tasks = append(tasks, t)
			}
		}
	}

	if l == "" && p == "" {
		json.NewEncoder(w).Encode(tasks)
		return
	}

	limit, err := strconv.ParseInt(l, 10, 64)
	if err != nil || limit < 1 {
		http.Error(w, "Limit must be a positive integer", 400)
		return
	}
	page, err := strconv.ParseInt(p, 10, 64)

	if err != nil || page < 1 {
		http.Error(w, "Page must be a positive integer", 400)
		return
	}

	startIdx := page - 1
	if startIdx > 0 {
		startIdx = startIdx * limit
	}

	if len(tasks) < int(startIdx) {
		json.NewEncoder(w).Encode([]*repo.Task{})
		return
	}

	json.NewEncoder(w).Encode(tasks[startIdx : limit*page])
}
