package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	repo "github.com/tatiananeda/todo/repository"
	"github.com/tatiananeda/todo/utils"
)

type Input struct {
	Title       string
	Description string
	Due_date    string
	IsComplete  bool
}

var Create = utils.WithErrorHandling(create)

func create(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	req, err := io.ReadAll(r.Body)

	if err != nil {
		return fmt.Errorf("Error parsing body %w;", err)
	}

	var input Input
	e := json.Unmarshal(req, &input)

	if e != nil {
		return utils.InvalidJSON(e)
	}

	if input.Title == "" {
		return utils.InvalidField("title is required")
	}

	if input.Due_date == "" {
		return utils.InvalidField("due_date is required")
	}

	uuid, err := uuid.NewRandom()

	if err != nil {
		return fmt.Errorf("Error generating uuid %w;", err)
	}

	t := repo.Task{
		Title:       input.Title,
		Description: input.Description,
		Due:         input.Due_date,
		Id:          uuid.String(),
		IsComplete:  input.IsComplete,
	}

	repo.Tasks = append(repo.Tasks, &t)

	if err := utils.WriteJSON(w, http.StatusOK, t); err != nil {
		return err
	}

	return nil
}
