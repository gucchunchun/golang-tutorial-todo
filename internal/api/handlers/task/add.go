package task

import (
	"net/http"

	"golang/tutorial/todo/internal/adapters/convert"
	"golang/tutorial/todo/internal/api/handlers"
)

func (h *TaskHandler) Add(w http.ResponseWriter, r *http.Request) {
	// リクエストボディの読み込み
	var in createRequest
	if err := handlers.ReadJSON(r, &in); err != nil {
		handlers.WriteError(w, err)
		return
	}

	param := convert.CreateParams{
		Name:  in.Name,
		DueAt: in.DueAt,
	}
	c, err := convert.FromCreateInput(param, nil)
	if err != nil {
		handlers.WriteError(w, err)
		return

	}

	task, err := h.TaskService.AddTask(c)
	if err != nil {
		handlers.WriteError(w, err)
		return
	}

	handlers.WriteJSON(w, http.StatusCreated, task)
}
