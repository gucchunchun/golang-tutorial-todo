package task

import (
	"net/http"

	"golang/tutorial/todo/internal/api/handlers"
	"golang/tutorial/todo/internal/models"
)

func (t *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	taskID, err := handlers.ParseID(r, "id")
	if err != nil {
		handlers.WriteError(w, err)
		return
	}

	// リクエストボディの読み込み
	var input struct {
		Name    *string `json:"name"`
		DueDate *string `json:"due_date"`
		Status  *string `json:"status"`
	}
	if err := handlers.ReadJSON(r, &input); err != nil {
		handlers.WriteError(w, err)
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
