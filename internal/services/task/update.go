package task

import (
	"fmt"

	"golang/tutorial/todo/internal/models"
	"golang/tutorial/todo/internal/storage"
)

func UpdateTask(filePath string, taskID models.TaskID, updates models.TaskUpdate) error {
	// Load existing tasks from the file
	tasks, err := storage.LoadTasks(filePath)
	if err != nil {
		return err
	}

	found := false
	// Find the task by ID and update its status
	for i, task := range tasks {
		if task.ID == taskID {
			found = true
			if updates.Status != nil {
				tasks[i].Status = *updates.Status
			}
			if updates.Due != nil {
				*updates.Due = updates.Due.UTC()
				tasks[i].DueDate = updates.Due
			}
			if updates.Name != nil {
				tasks[i].Name = *updates.Name
			}
			break
		}
	}

	if !found {
		return fmt.Errorf("task with ID %d not found", taskID)
	}

	// Save the updated tasks back to the file
	return storage.SaveTasks(filePath, tasks)
}
