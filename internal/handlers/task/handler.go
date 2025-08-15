package task

import (
	"net/http"

	"golang/tutorial/todo/internal/models"
)

type TaskService interface {
	// Create(title string) error
	// Get(id string) (models.Task, error)
	ListTasks() ([]models.Task, error)
	// Update(id string, title string) error
	// Delete(id string) error
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
	// mux.HandleFunc("POST /tasks", h.Add)
	mux.HandleFunc("GET /tasks", h.List)
	// mux.HandleFunc("PATCH /tasks/", h.Update)
	// mux.HandleFunc("DELETE /tasks/", h.Delete)
}
