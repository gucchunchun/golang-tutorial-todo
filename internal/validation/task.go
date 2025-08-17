package validation

import (
	"golang/tutorial/todo/internal/models"
)

type CreateTaskInput struct {
	Name    string       `json:"name" validate:"required,min=1,max=255"`
	DueDate *models.Date `json:"due_date" validate:"omitempty,datetime=2006-01-02"`
}

func ValidateCreateTaskInput(input CreateTaskInput) Errors {
	return Validate(input)
}

type UpdateTaskInput struct {
	Name     *string        `json:"name" validate:"omitempty,min=1,max=255"`
	Status   *models.Status `json:"status" validate:"omitempty,oneof=pending ongoing done"`
	DueDate  *models.Date   `json:"due_date" validate:"omitempty,datetime=2006-01-02"`
	ClearDue bool           `json:"clear_due" validate:"omitempty"`
}

func ValidateUpdateTaskInput(input UpdateTaskInput) Errors {
	return Validate(input)
}
