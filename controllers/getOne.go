package controllers

import (
	"github.com/gorilla/mux"
	repo "github.com/tatiananeda/todo/repository"
	"github.com/tatiananeda/todo/utils"
	"net/http"
)

var GetOne = utils.WithErrorHandling(getOne)

func getOne(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	for _, task := range repo.Tasks {
		if task.Id == id {
			if err := utils.WriteJSON(w, http.StatusOK, task); err != nil {
				return err
			}
			return nil
		}
	}

	return utils.NotFound(id)
}
