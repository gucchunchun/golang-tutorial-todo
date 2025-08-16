package task

import (
	"net/http"

	"golang/tutorial/todo/internal/models"
)

type TaskService interface {
	AddTask(taskName string, dueDate string) error
	GetTask(taskID string) (models.TaskOutput, error)
	ListTasks() ([]models.TaskOutput, error)
	UpdateTask(taskID string, updates models.TaskUpdate) error
}

type TaskHandler struct {
	TaskService TaskService
}

func NewTaskHandler(taskService TaskService) *TaskHandler {
	return &TaskHandler{
		TaskService: taskService,
	}
}

func (h *TaskHandler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /tasks", h.Add)
	mux.HandleFunc("GET /tasks/{id}", h.get)
	mux.HandleFunc("GET /tasks", h.GetList)
	mux.HandleFunc("PATCH /tasks/{id}", h.Update)
}
