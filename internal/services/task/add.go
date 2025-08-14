package task

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"golang/tutorial/todo/internal/models"
	"golang/tutorial/todo/internal/storage"
	"golang/tutorial/todo/internal/utils"
)

func AddTask(filename, taskName string, setDueDate bool, dueDate string) error {
	if !utils.IsValidTaskName(taskName) {
		return fmt.Errorf("invalid task name")
	}

	tasks, err := storage.LoadTasks(filename)
	if err != nil {
		return err
	}

	var dueDateTime time.Time
	if setDueDate {
		parsedDueDate, err := utils.ParseDate(dueDate)
		if err != nil {
			return fmt.Errorf("failed to parse due date: %v", err)
		}
		dueDateTime = parsedDueDate
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("failed to generate unique ID: %v", err)
	}

	newTask := models.Task{
		ID:        id,
		Name:      strings.TrimSpace(taskName),
		Status:    models.Pending,
		CreatedAt: utils.Now(),
		DueDate:   &dueDateTime,
	}

	tasks = append(tasks, newTask)
	return storage.SaveTasks(filename, tasks)
}
