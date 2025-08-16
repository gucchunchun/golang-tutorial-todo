package task

import (
	"net/http"

	"golang/tutorial/todo/internal/api/handlers"
)

func (s *TaskHandler) get(w http.ResponseWriter, r *http.Request) {
	taskID, err := handlers.ParseID(r, "id")
	if err != nil {
		handlers.WriteError(w, err)
		return
	}

	task, err := s.TaskService.GetTask(taskID)
	if err != nil {
		handlers.WriteError(w, err)
		return
	}

	handlers.WriteJSON(w, http.StatusOK, task)
}
