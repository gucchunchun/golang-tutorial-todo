package task

import (
	"golang/tutorial/todo/internal/storage"
	"golang/tutorial/todo/internal/utils"
)

func ListTasks(filename string) ([]string, error) {
	tasks, err := storage.LoadTasks(filename)
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
