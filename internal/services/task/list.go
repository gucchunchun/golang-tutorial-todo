package task

import (
	"golang/tutorial/todo/internal/apperr"
	"golang/tutorial/todo/internal/models"
)

func (s *TaskService) ListTasks() ([]models.TaskOutput, error) {
	tasks, err := s.storage.LoadTasks()
	if err != nil {
		return nil, apperr.E(apperr.CodeUnknown, "Failed to load tasks", err)
	}

	output := make([]models.TaskOutput, 0, len(tasks))
	for _, task := range tasks {
		output = append(output, task.TaskOutput())
	}

	return output, nil
}
