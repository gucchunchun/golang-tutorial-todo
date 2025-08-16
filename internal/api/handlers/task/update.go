package task

import (
	"errors"
	"net/http"
	"strings"

	"golang/tutorial/todo/internal/api/handlers"
	"golang/tutorial/todo/internal/models"
)

func (t *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: taskIDのバリデーションはserviceに移動
	taskID, ok := tailID(strings.TrimPrefix(r.URL.Path, "/tasks/"))
	if !ok {
		handlers.WriteError(w, errors.New("invalid taskID"))
		return
	}

	// リクエストボディの読み込み
	var input struct {
		Name    *string `json:"name"`
		DueDate *string `json:"due_date"`
		Status  *string `json:"status"`
	}
	if err := handlers.ReadJSON(r, &input); err != nil {
		// TODO: エラー返却方法を策定
		handlers.WriteError(w, errors.New("invalid request body"))
		return
	}

	taskUpdate := models.TaskUpdate{
		Status: input.Status,
		Due:    input.DueDate,
		Name:   input.Name,
	}
	if err := t.TaskService.UpdateTask(taskID, taskUpdate); err != nil {
		handlers.WriteError(w, err)
		return
	}
	handlers.WriteJSON(w, http.StatusOK, "Task updated successfully")
}
