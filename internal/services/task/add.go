package task

import (
	"fmt"
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

	var dueDateTime time.Time
	if setDueDate {
		parsedDueDate, err := utils.ParseDate(dueDate)
		if err != nil {
			return fmt.Errorf("failed to parse due date: %v", err)
		}
		dueDateTime = parsedDueDate
	}

	newTask, err := models.NewTask(taskName, &dueDateTime, utils.Now())
	if err != nil {
		return fmt.Errorf("failed to create new task: %v", err)
	}

	tasks = append(tasks, newTask)
	return storage.SaveTasks(filename, tasks)
}
