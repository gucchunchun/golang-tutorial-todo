package task

import (
	"fmt"
	"strings"
	"time"

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

	newTask := models.Task{
		ID:        uint(len(tasks) + 1),
		Name:      strings.TrimSpace(taskName),
		Status:    0,
		CreatedAt: time.Now().Format(time.RFC3339),
		DueDate:   dueDate,
	}

	tasks = append(tasks, newTask)
	return storage.SaveTasks(filename, tasks)
}
