package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

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

func ParseStatus(input string) (Status, error) {
	switch strings.ToLower(input) {
	case "pending":
		return Pending, nil
	case "ongoing":
		return Ongoing, nil
	case "done":
		return Done, nil
	default:
		return -1, fmt.Errorf("invalid status: %s", input)
	}
}

type Task struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Status    Status     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	DueDate   *time.Time `json:"due_date"`
}

type TaskOutput struct {
	ID        string
	Name      string
	Status    string
	CreatedAt string
	DueDate   string
	TimeLeft  string
}

type TaskUpdate struct {
	Status *Status
	Due    *time.Time
	Name   *string
}
