package task

import (
	"net/http"

	"golang/tutorial/todo/internal/handlers"
	"golang/tutorial/todo/internal/validation"
)

func (h *TaskHandler) Add(w http.ResponseWriter, r *http.Request) {
	// リクエストボディの読み込み
	var input struct {
		Name    string  `json:"name"`
		DueDate *string `json:"due_date"`
	}
	if err := handlers.ReadJSON(r, &input); err != nil {
		handlers.WriteRequestError(w)
		return
	}

	// バリデーション
	if vErr := validation.ValidateCreateTaskInput(validation.CreateTaskInput{
		Name:    input.Name,
		DueDate: input.DueDate,
	}); vErr != nil {
		handlers.WriteValidationError(w, vErr)
		return
	}

	var dueDate string
	if input.DueDate == nil {
		dueDate = ""
	} else {
		dueDate = *input.DueDate
	}

	if err := h.TaskService.AddTask(input.Name, dueDate); err != nil {
		handlers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handlers.WriteJSON(w, http.StatusCreated)
}
