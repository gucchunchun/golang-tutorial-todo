package task

import (
	"net/http"

	"golang/tutorial/todo/internal/api/handlers"
)

func (h *TaskHandler) GetList(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.TaskService.ListTasks()
	if err != nil {
		handlers.WriteError(w, err)
		return
	}

	handlers.WriteJSON(w, http.StatusOK, tasks)
}
