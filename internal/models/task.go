package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ID
type TaskID uint64

func ParseTaskID(s string) (TaskID, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return TaskID(0), fmt.Errorf("invalid TaskID: %w", err)
	}
	return TaskID(i), nil
}
func ParseTaskIDInt(i int) TaskID { return TaskID(i) }
func (id TaskID) String() string  { return strconv.FormatUint(uint64(id), 10) }
func (id *TaskID) UnmarshalText(b []byte) error {
	i, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		return err
	}
	*id = TaskID(i)
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

// Status
type Status int

const (
	StatusPending Status = iota
	StatusOngoing
	StatusDone
)

// NOTE: "[...]" はコンパイラに配列のサイズを推測させるために使用
// ここでは、Statusの値に対応する文字列（３）を定義している
var statusNames = [...]string{"Pending", "Ongoing", "Done"}

func (s Status) String() string {
	if s < StatusPending || s > StatusDone {
		return "Unknown"
	}
	return statusNames[s]
}
func ParseStatus(input string) (Status, error) {
	switch strings.ToLower(input) {
	case "pending":
		return StatusPending, nil
	case "ongoing":
		return StatusOngoing, nil
	case "done":
		return StatusDone, nil
	default:
		return -1, fmt.Errorf("invalid status: %s", input)
	}
}

// Date
type Date = time.Time

// TimeLeft
type TimeLeft = time.Duration

// Task
type Task struct {
	ID        TaskID `json:"id"`
	Name      string `json:"name"`
	Status    Status `json:"status"`
	CreatedAt Date   `json:"created_at"`
	UpdatedAt Date   `json:"updated_at"`
	DueAt     *Date  `json:"due_date"`
}

type TaskOutput struct {
	ID        TaskID    `json:"id"`
	Name      string    `json:"name"`
	Status    Status    `json:"status"`
	CreatedAt Date      `json:"created_at"`
	UpdatedAt Date      `json:"updated_at"`
	DueAt     *Date     `json:"due_date"`
	TimeLeft  *TimeLeft `json:"time_left"`
}

func (t *Task) TaskOutput() TaskOutput {
	var timeLeft *time.Duration = nil
	if t.DueAt != nil && !t.DueAt.IsZero() {
		localTime := time.Now()
		if t.DueAt.Before(localTime) {
			zero := time.Duration(0)
			timeLeft = &zero
		} else {
			d := t.DueAt.Sub(localTime)
			timeLeft = &d
		}
	}
	return TaskOutput{
		ID:        t.ID,
		Name:      t.Name,
		Status:    t.Status,
		CreatedAt: t.CreatedAt,
		DueAt:     t.DueAt,
		TimeLeft:  timeLeft,
	}
}

type TaskCreate struct {
	Name  string
	DueAt *Date
}

func NewTask(id TaskID, name string, createdAt, updatedAt Date, dueAt *Date) Task {
	return Task{
		ID:        id,
		Name:      name,
		Status:    StatusPending,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		DueAt:     dueAt,
	}
}

type TaskUpdate struct {
	Name     *string
	Status   *Status
	DueAt    *Date
	ClearDue bool
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
func (t Tasks) TaskOutputs() []TaskOutput {
	outputs := make([]TaskOutput, 0, len(t))
	for _, task := range t {
		outputs = append(outputs, task.TaskOutput())
	}
	return outputs
}
