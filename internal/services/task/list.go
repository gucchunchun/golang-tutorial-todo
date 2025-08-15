package task

import (
	"golang/tutorial/todo/internal/models"
)

func (s *Service) ListTasks() ([]models.TaskOutput, error) {
	tasks, err := s.storage.LoadTasks()
	if err != nil {
		return nil, err
	}

	output := make([]models.TaskOutput, 0, len(tasks))
	for _, task := range tasks {
		output = append(output, task.TaskOutput())
	}

	return output, nil
}
