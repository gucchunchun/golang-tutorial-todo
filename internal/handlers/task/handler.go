package task

import (
	"net/http"
	"strings"

	"golang/tutorial/todo/internal/models"
)

type TaskService interface {
	AddTask(taskName string, dueDate string) error
	ListTasks() ([]models.Task, error)
	UpdateTask(taskID models.TaskID, updates models.TaskUpdate) error
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
	mux.HandleFunc("GET /tasks", h.List)
	mux.HandleFunc("PATCH /tasks/", h.Update)
}

func tailID(raw string) (models.TaskID, bool) {
	id, err := models.ParseTaskID(strings.TrimSpace(raw))
	return id, err == nil
}
