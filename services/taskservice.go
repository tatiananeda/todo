package services

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/tatiananeda/todo/entities/web"
	"github.com/tatiananeda/todo/repository"
)

type TaskService struct {
	repository *repository.Repository
}

func NewTaskService(r *repository.Repository) *TaskService {
	t := TaskService{repository: r}
	return &t
}

func (s TaskService) Create(input *web.TaskInput) (*repository.Task, error) {
	uuid, err := uuid.NewRandom()

	if err != nil {
		return nil, fmt.Errorf("Error generating uuid %w;", err)
	}

	t := repository.Task{
		Title:       input.Title,
		Description: input.Description,
		Due:         input.DueDate,
		Id:          uuid.String(),
		IsComplete:  input.IsComplete,
	}

	s.repository.AddTask(&t)

	return &t, nil
}

func (s TaskService) Update(id string, input *web.TaskInput) (*repository.Task, error) {
	task := s.repository.GetTaskById(id)
	if task == nil {
		return nil, web.NotFound(id)
	}

	if input.Description != "" {
		task.Description = input.Description
	}

	if input.Title != "" {
		task.Title = input.Title
	}

	task.IsComplete = input.IsComplete

	return task, nil
}

func (s TaskService) Delete(id string) error {
	res := s.repository.DeleteById(id)
	if res != true {
		return web.NotFound(id)
	}
	return nil
}

func (s TaskService) GetById(id string) (*repository.Task, error) {
	t := s.repository.GetTaskById(id)
	if t == nil {
		return nil, web.NotFound(id)
	}
	return t, nil
}

func (s TaskService) GetFilteredByCompleted(completed bool) []*repository.Task {
	return s.repository.GetByCompleted(completed)
}

func (s TaskService) GetAll() []*repository.Task {
	return s.repository.GetAll()
}

func (s TaskService) GetPage(limit, page string, slice []*repository.Task) ([]*repository.Task, error) {
	var tasks []*repository.Task = slice
	if tasks == nil {
		tasks = s.repository.GetAll()
	}

	l, err := strconv.ParseInt(limit, 10, 64)
	if err != nil || l < 1 {
		return nil, web.InvalidField("Limit must be a positive integer")
	}

	p, err := strconv.ParseInt(page, 10, 64)
	if err != nil || p < 1 {
		return nil, web.InvalidField("Page must be a positive integer")
	}

	lastIdx := len(tasks) - 1
	startIdx := p - 1
	if startIdx > 0 {
		startIdx = startIdx * l
	}

	endIdx := startIdx + l
	if len(tasks) < int(endIdx) {
		endIdx = int64(lastIdx)
	}

	if lastIdx < int(startIdx) {
		empty := make([]*repository.Task, 0)
		return empty, nil
	}

	result := (tasks)[startIdx:endIdx]
	return result, nil
}
