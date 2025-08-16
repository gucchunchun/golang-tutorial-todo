package task

import (
	"fmt"
	"golang/tutorial/todo/internal/apperr"
	"golang/tutorial/todo/internal/models"
)

func (s *TaskService) GetTask(taskID models.TaskID) (models.TaskOutput, error) {
	tasks, err := s.storage.LoadTasks()
	if err != nil {
		return models.TaskOutput{}, err
	}

	target, ok := tasks.FindByID(taskID)
	if !ok {
		return models.TaskOutput{}, apperr.E(apperr.CodeNotFound, fmt.Sprintf("Task with the ID: %s not found", taskID), ErrNotFound)
	}
	return target.TaskOutput(), nil
}
