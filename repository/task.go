package repository

type Task struct {
	Title       string
	Description string
	Due         string
	IsComplete  bool
	Id          string
}
type Repository struct {
	Tasks []*Task
}

func NewRepository() *Repository {
	r := Repository{Tasks: make([]*Task, 0)}
	return &r
}

func (r *Repository) AddTask(t *Task) {
	r.Tasks = append(r.Tasks, t)
}

func (r *Repository) GetTaskById(id string) *Task {
	for _, task := range r.Tasks {
		if task.Id == id {
			return task
		}
	}
	return nil
}

func (r *Repository) DeleteById(id string) bool {
	for idx, task := range r.Tasks {
		if task.Id == id {
			r.Tasks = append(r.Tasks[:idx], r.Tasks[idx+1:]...)
			return true
		}
	}
	return false
}

func (r *Repository) GetByCompleted(completed bool) []*Task {
	tasks := make([]*Task, 0)
	for _, t := range r.Tasks {
		if t.IsComplete == completed {
			tasks = append(tasks, t)
		}
	}
	return tasks
}

func (r *Repository) GetAll() []*Task {
	return r.Tasks
}
