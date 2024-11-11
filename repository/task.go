package repository

type Task struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Due         string `json:"due_date"`
	IsComplete  bool   `json:"complete"`
	Id          string `json:"id"`
}

var Tasks []*Task = make([]*Task, 0)
