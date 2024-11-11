package controllers

import (
	"encoding/json"
	"github.com/google/uuid"
	repo "github.com/tatiananeda/todo/repository"
	"io"
	"net/http"
)

type Input struct {
	Title       string
	Description string
	Due_date    string
	IsComplete  bool
}

func Create(w http.ResponseWriter, r *http.Request) {
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

	if input.Title == "" {
		http.Error(w, "title is required", 400)
		return
	}

	if input.Due_date == "" {
		http.Error(w, "due_date is required", 400)
		return
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	t := repo.Task{
		Title:       input.Title,
		Description: input.Description,
		Due:         input.Due_date,
		Id:          uuid.String(),
		IsComplete:  input.IsComplete,
	}

	repo.Tasks = append(repo.Tasks, &t)

	json.NewEncoder(w).Encode(t)
}
