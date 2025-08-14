package task

import (
	"golang/tutorial/todo/internal/utils"
)

func (s *Service) ListTasks() ([]string, error) {
	tasks, err := s.storage.LoadTasks()
	if err != nil {
		return nil, err
	}

	formatted := make([]string, len(tasks))
	for i, t := range tasks {
		formatted[i] = utils.FormatTaskOutput(
			t.ID,
			t.Name,
			t.Status.String(),
			t.CreatedAt,
			*t.DueDate,
		)
	}
	return formatted, nil
}
