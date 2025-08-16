package task

import (
	"errors"
	"net/http"
	"strings"

	"golang/tutorial/todo/internal/api/handlers"
	"golang/tutorial/todo/internal/apperr"
)

func (s *TaskHandler) get(w http.ResponseWriter, r *http.Request) {
	taskID, ok := tailID(strings.TrimPrefix(r.URL.Path, "/tasks/"))
	if !ok {
		handlers.WriteJSONError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := s.TaskService.GetTask(taskID)
	if err != nil {
		var ae *apperr.Error
		if errors.As(err, &ae) {
			switch ae.Code {
			case apperr.CodeNotFound:
				handlers.WriteJSONError(w, http.StatusNotFound, ae.Error())
				return
			}
		}
		handlers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handlers.WriteJSON(w, http.StatusOK, task)
}
