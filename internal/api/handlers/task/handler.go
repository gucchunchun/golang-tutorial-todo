package task

import (
	"net/http"
	"strings"

	"golang/tutorial/todo/internal/models"
)

type TaskService interface {
	AddTask(taskName string, dueDate string) error
	GetTask(taskID models.TaskID) (models.TaskOutput, error)
	ListTasks() ([]models.TaskOutput, error)
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
	mux.HandleFunc("GET /tasks/{id}", h.get)
	mux.HandleFunc("GET /tasks", h.GetList)
	mux.HandleFunc("PATCH /tasks/{id}", h.Update)
}

func tailID(raw string) (models.TaskID, bool) {
	id, err := models.ParseTaskID(strings.TrimSpace(raw))
	return id, err == nil
}
