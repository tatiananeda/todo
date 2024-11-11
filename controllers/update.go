package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	repo "github.com/tatiananeda/todo/repository"
	"io"
	"net/http"
)

func Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	req, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	var input Input
	e := json.Unmarshal(req, &input)

	if e != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	for _, task := range repo.Tasks {
		if task.Id == id {

			if input.Description != "" {
				task.Description = input.Description
			}

			if input.Title != "" {
				task.Title = input.Title
			}

			task.IsComplete = input.IsComplete

			json.NewEncoder(w).Encode(task)
			return
		}
	}

	http.Error(w, "Task "+id+" not found", 404)
}
