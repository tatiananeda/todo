package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	repo "github.com/tatiananeda/todo/repository"
	"net/http"
)

func MarkComplete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for _, task := range repo.Tasks {
		if task.Id == id {
			task.IsComplete = !task.IsComplete

			json.NewEncoder(w).Encode(task)
			return
		}
	}

	http.Error(w, "Task "+id+" not found", 404)
}
