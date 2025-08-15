package task

import (
	"net/http"

	"golang/tutorial/todo/internal/api/handlers"
)

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.TaskService.ListTasks()
	if err != nil {
		handlers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
	}
	if len(tasks) == 0 {
		handlers.WriteJSONError(w, http.StatusNotFound, "No tasks found")
	}

	handlers.WriteJSON(w, http.StatusOK, tasks)
}
