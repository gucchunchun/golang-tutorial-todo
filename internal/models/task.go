package models

type Status int

const (
	Pending Status = iota
	Ongoing
	Done
)

type Task struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Status    Status `json:"status"`
	CreatedAt string `json:"created_at"`
	DueDate   string `json:"due_date"`
}
