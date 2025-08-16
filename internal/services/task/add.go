package task

import (
	"fmt"
	"time"

	"golang/tutorial/todo/internal/apperr"
	"golang/tutorial/todo/internal/models"
	"golang/tutorial/todo/internal/utils"
	"golang/tutorial/todo/internal/validation"
)

func (s *TaskService) AddTask(taskName string, dueDate string) error {
	// バリデーション
	if err := validation.ValidateCreateTaskInput(validation.CreateTaskInput{
		Name:    taskName,
		DueDate: dueDate,
	}); err != nil {
		return apperr.E(apperr.CodeInvalid, "Validation error", ErrValidation)
	}

	var dueDateTime time.Time
	if dueDate != "" {
		parsedDueDate, err := utils.ParseDate(dueDate)
		if err != nil {
			// NOTE: この時点で日付フォーマットが間違っていることは想定外
			return apperr.E(apperr.CodeUnknown, fmt.Sprintf("Failed to parse due date: %s", s), nil)
		}
		dueDateTime = parsedDueDate
	}

	// NOTE: task29の実装
	// if !utils.IsValidTaskName(taskName) {
	// 	return apperr.E(apperr.CodeInvalid, "Invalid task name", ErrValidation)
	// }

	tasks, err := s.storage.LoadTasks()
	if err != nil {
		return apperr.E(apperr.CodeUnknown, "Failed to load tasks", ErrDatabase)
	}

	newTask, err := models.NewTask(taskName, &dueDateTime, utils.SystemTime())
	if err != nil {
		return apperr.E(apperr.CodeUnknown, "Failed to create new task", nil)
	}

	tasks = append(tasks, newTask)
	return s.storage.SaveTasks(tasks)
}
