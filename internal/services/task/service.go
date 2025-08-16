package task

import (
	"golang/tutorial/todo/internal/models"
)

type StorageClient interface {
	LoadTasks() (models.Tasks, error)
	SaveTasks(tasks []models.Task) error
}

type TaskService struct {
	storage StorageClient
}

func NewTaskService(storage StorageClient) *TaskService {
	return &TaskService{
		storage: storage,
	}
}
