package convert

import (
	"fmt"
	"golang/tutorial/todo/internal/models"
	"time"
)

type CreateParams struct {
	Name  string
	DueAt *string
}

type UpdateParams struct {
	Name     *string
	DueAt    *string
	Status   *string
	ClearDue bool
}

func ParseStatusPtr(s *string) (*models.Status, error) {
	if s == nil {
		return nil, nil
	}
	st, err := models.ParseStatus(*s)
	if err != nil {
		return nil, err
	}
	return &st, nil
}

func ParseDatePtr(s *string, loc *time.Location) (*models.Date, error) {
	if s == nil {
		return nil, nil
	}
	if loc == nil {
		loc = time.Local
	}
	if t, err := time.ParseInLocation("2006-01-02", *s, loc); err == nil {
		tt := t
		return &tt, nil
	}
	return nil, fmt.Errorf("invalid due_date format: %q (YYYY-MM-DD)", *s)
}

func FromCreateInput(in CreateParams, loc *time.Location) (models.TaskCreate, error) {
	var out models.TaskCreate
	out.Name = in.Name

	due, err := ParseDatePtr(in.DueAt, loc)
	if err != nil {
		return out, err
	}
	out.DueAt = due
	return out, nil
}

func FromUpdateInput(in UpdateParams, loc *time.Location) (models.TaskUpdate, error) {
	var out models.TaskUpdate
	out.Name = in.Name

	st, err := ParseStatusPtr(in.Status)
	if err != nil {
		return out, err
	}
	out.Status = st

	due, err := ParseDatePtr(in.DueAt, loc)
	if err != nil {
		return out, err
	}
	out.DueAt = due

	out.ClearDue = in.ClearDue
	return out, nil
}
