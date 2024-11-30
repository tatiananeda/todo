package web

type TaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	IsComplete  bool   `json:"complete"`
}
