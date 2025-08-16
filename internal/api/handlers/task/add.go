package task

import (
	"net/http"

	"golang/tutorial/todo/internal/api/handlers"
)

func (h *TaskHandler) Add(w http.ResponseWriter, r *http.Request) {
	// リクエストボディの読み込み
	var input struct {
		Name    string  `json:"name"`
		DueDate *string `json:"due_date"`
	}
	if err := handlers.ReadJSON(r, &input); err != nil {
		handlers.WriteError(w, err)
		return
	}

	var dueDate string
	if input.DueDate == nil {
		dueDate = ""
	} else {
		dueDate = *input.DueDate
	}

	if err := h.TaskService.AddTask(input.Name, dueDate); err != nil {
		handlers.WriteError(w, err)
		return
	}

	handlers.WriteJSON(w, http.StatusCreated)
}
