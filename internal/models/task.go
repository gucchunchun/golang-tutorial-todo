package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TaskID uuid.UUID

func NewTaskID() TaskID { u := uuid.New(); return TaskID(u) }
func ParseTaskID(s string) (TaskID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return TaskID{}, fmt.Errorf("invalid UUID: %w", err)
	}
	return TaskID(u), nil
}
func (id TaskID) String() string { return uuid.UUID(id).String() }
func (id *TaskID) UnmarshalText(b []byte) error {
	u, err := uuid.ParseBytes(b)
	if err != nil {
		return err
	}
	*id = TaskID(u)
	return nil
}
func (id TaskID) MarshalJSON() ([]byte, error) { return json.Marshal(id.String()) }
func (id *TaskID) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return id.UnmarshalText([]byte(s))
}

type Status int

const (
	Pending Status = iota
	Ongoing
	Done
)

// NOTE: "[...]" はコンパイラに配列のサイズを推測させるために使用。
// ここでは、Statusの値に対応する文字列（３）を定義しています。
var statusNames = [...]string{"Pending", "Ongoing", "Done"}

func (s Status) String() string {
	if s < Pending || s > Done {
		return "Unknown"
	}
	return statusNames[s]
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
	ID        TaskID     `json:"id"`
	Name      string     `json:"name"`
	Status    Status     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	DueDate   *time.Time `json:"due_date"`
}

func NewTask(name string, due *time.Time, now time.Time) (Task, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Task{}, fmt.Errorf("name is required")
	}
	if due != nil {
		d := due.UTC()
		due = &d
	}
	return Task{
		ID:        NewTaskID(),
		Name:      name,
		Status:    Pending,
		CreatedAt: now,
		DueDate:   due,
	}, nil
}

type TaskOutput struct {
	ID        TaskID         `json:"id"`
	Name      string         `json:"name"`
	Status    Status         `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	DueDate   *time.Time     `json:"due_date"`
	TimeLeft  *time.Duration `json:"time_left"`
}

func (t *Task) TaskOutput() TaskOutput {
	var timeLeft *time.Duration = nil
	if t.DueDate != nil && !t.DueDate.IsZero() {
		localTime := time.Now()
		if t.DueDate.Before(localTime) {
			zero := time.Duration(0)
			timeLeft = &zero
		} else {
			d := t.DueDate.Sub(localTime)
			timeLeft = &d
		}
	}
	return TaskOutput{
		ID:        t.ID,
		Name:      t.Name,
		Status:    t.Status,
		CreatedAt: t.CreatedAt,
		DueDate:   t.DueDate,
		TimeLeft:  timeLeft,
	}
}

type TaskUpdate struct {
	Status *string
	Due    *string
	Name   *string
}

type Tasks []Task

func (t Tasks) FindByID(id TaskID) (Task, bool) {
	for _, task := range t {
		if task.ID == id {
			return task, true
		}
	}
	return Task{}, false
}
