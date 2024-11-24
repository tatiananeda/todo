package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	repo "github.com/tatiananeda/todo/repository"
	"github.com/tatiananeda/todo/utils"
)

var Delete = utils.WithErrorHandling(delete)

func delete(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	for idx, task := range repo.Tasks {
		if task.Id == id {
			repo.Tasks = append(repo.Tasks[:idx], repo.Tasks[idx+1:]...)
			w.WriteHeader(http.StatusOK)
			return nil
		}
	}

	return utils.NotFound(id)
}
