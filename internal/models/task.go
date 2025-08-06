package models

import "time"

type Status int

const (
	Pending Status = iota
	Ongoing
	Done
)

func (s Status) String() string {
	switch s {
	case Pending:
		return "Pending"
	case Ongoing:
		return "Ongoing"
	case Done:
		return "Done"
	default:
		return "Unknown"
	}
}

type Task struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	Status    Status     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	DueDate   *time.Time `json:"due_date"`
}

type TaskOutput struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	DueDate   string `json:"due_date"`
	TimeLeft  string `json:"time_left"`
}
