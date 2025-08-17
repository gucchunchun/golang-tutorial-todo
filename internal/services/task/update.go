package task

import (
	"fmt"
	"time"

	"golang/tutorial/todo/internal/apperr"
	"golang/tutorial/todo/internal/models"
	"golang/tutorial/todo/internal/utils"
	"golang/tutorial/todo/internal/validation"
)

func (s *TaskService) UpdateTask(taskID string, updates models.TaskUpdate) error {
	// バリデーション
	parsedID, err := models.ParseTaskID(taskID)
	if err != nil {
		return apperr.E(apperr.CodeInvalid, "Validation error", err)
	}
	if err := validation.ValidateUpdateTaskInput(validation.UpdateTaskInput{
		Name:    updates.Name,
		DueDate: updates.Due,
		Status:  updates.Status,
	}); err != nil {
		return apperr.E(apperr.CodeInvalid, "Validation error", err)
	}

	var status *models.Status
	if updates.Status == nil {
		status = nil
	} else {
		parsedStatus, err := models.ParseStatus(*updates.Status)
		if err != nil {
			return apperr.E(apperr.CodeUnknown, "Validation error", err)
		}
		status = &parsedStatus
	}

	var dueDate *time.Time
	if updates.Due == nil {
		dueDate = nil
	} else {
		parsedDate, err := utils.ParseDate(*updates.Due)
		if err != nil {
			return apperr.E(apperr.CodeUnknown, "Validation error", err)
		}
		dueDate = &parsedDate
	}

	tasks, err := s.storage.LoadTasks()
	if err != nil {
		return apperr.E(apperr.CodeInvalid, "Validation error", err)
	}

	found := false
	for i, task := range tasks {
		if task.ID == parsedID {
			found = true
			if updates.Status != nil {
				tasks[i].Status = *status
			}
			if updates.Due != nil {
				tasks[i].DueDate = dueDate
			}
			if updates.Name != nil {
				tasks[i].Name = *updates.Name
			}
			break
		}
	}

	if !found {
		return apperr.E(apperr.CodeNotFound, fmt.Sprintf("Task not found with ID: %s,", taskID), err)
	}

	err = s.storage.SaveTasks(tasks)
	if err != nil {
		return apperr.E(apperr.CodeUnknown, "Failed to update task", err)
	}
	return nil
}
