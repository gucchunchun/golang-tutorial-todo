package task

import (
	"net/http"
	"strings"
	"time"

	"golang/tutorial/todo/internal/handlers"
	"golang/tutorial/todo/internal/models"
	"golang/tutorial/todo/internal/utils"
	"golang/tutorial/todo/internal/validation"
)

func (t *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	taskID, ok := tailID(strings.TrimPrefix(r.URL.Path, "/tasks/"))
	if !ok {
		handlers.WriteJSONError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// リクエストボディの読み込み
	var input struct {
		Name    *string `json:"name"`
		DueDate *string `json:"due_date"`
		Status  *string `json:"status"`
	}
	if err := handlers.ReadJSON(r, &input); err != nil {
		handlers.WriteRequestError(w)
		return
	}

	// バリデーション
	if vErr := validation.ValidateUpdateTaskInput(validation.UpdateTaskInput{
		Name:    input.Name,
		DueDate: input.DueDate,
		Status:  input.Status,
	}); vErr != nil {
		handlers.WriteValidationError(w, vErr)
		return
	}

	var status *models.Status
	if input.Status == nil {
		status = nil
	} else {
		parsedStatus, err := models.ParseStatus(*input.Status)
		if err != nil {
			handlers.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		status = &parsedStatus
	}

	var dueDate *time.Time
	if input.DueDate == nil {
		dueDate = nil
	} else {
		parsedDate, err := utils.ParseDate(*input.DueDate)
		if err != nil {
			handlers.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		dueDate = &parsedDate
	}

	taskUpdate := models.TaskUpdate{
		Status: status,
		Due:    dueDate,
		Name:   input.Name,
	}
	if err := t.TaskService.UpdateTask(taskID, taskUpdate); err != nil {
		handlers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	handlers.WriteJSON(w, http.StatusOK, "Task updated successfully")
}
