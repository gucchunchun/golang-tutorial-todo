package task

import (
	"golang/tutorial/todo/internal/models"
)

func (s *Service) ListTasks() ([]models.Task, error) {
	tasks, err := s.storage.LoadTasks()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
