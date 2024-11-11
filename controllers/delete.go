package controllers

import (
	"github.com/gorilla/mux"
	repo "github.com/tatiananeda/todo/repository"
	"net/http"
)

func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	for idx, task := range repo.Tasks {
		if task.Id == id {
			repo.Tasks = append(repo.Tasks[:idx], repo.Tasks[idx+1:]...)
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(w, "Task "+id+" not found", 404)
}
