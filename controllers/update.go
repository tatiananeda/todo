package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	repo "github.com/tatiananeda/todo/repository"
	"github.com/tatiananeda/todo/utils"
)

var Update = utils.WithErrorHandling(update)

func update(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	req, err := io.ReadAll(r.Body)

	if err != nil {
		return fmt.Errorf("Error parsing body %w;", err)
	}

	var input Input
	e := json.Unmarshal(req, &input)

	if e != nil {
		return utils.InvalidJSON(e)
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

			if err := utils.WriteJSON(w, http.StatusOK, task); err != nil {
				return err
			}
			return nil
		}
	}

	return utils.NotFound(id)
}
