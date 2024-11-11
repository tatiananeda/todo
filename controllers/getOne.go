package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	repo "github.com/tatiananeda/todo/repository"
	"net/http"
)

func GetOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for _, task := range repo.Tasks {
		if task.Id == id {
			json.NewEncoder(w).Encode(task)
			return
		}
	}

	http.Error(w, "Task "+id+" not found", 404)
}
