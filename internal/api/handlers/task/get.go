package task

import (
	"errors"
	"net/http"
	"strings"

	"golang/tutorial/todo/internal/api/handlers"
)

func (s *TaskHandler) get(w http.ResponseWriter, r *http.Request) {
	// TODO: taskIDのバリデーションはserviceに移動
	taskID, ok := tailID(strings.TrimPrefix(r.URL.Path, "/tasks/"))
	if !ok {
		handlers.WriteError(w, errors.New("invalid taskID"))
		return
	}

	task, err := s.TaskService.GetTask(taskID)
	if err != nil {
		handlers.WriteError(w, err)
		return
	}

	handlers.WriteJSON(w, http.StatusOK, task)
}
