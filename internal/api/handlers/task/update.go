package task

import (
	"net/http"
	"time"

	"golang/tutorial/todo/internal/adapters/convert"
	"golang/tutorial/todo/internal/api/handlers"
)

func (t *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	// IDの取得
	taskID, err := handlers.ParseID(r, "id")
	if err != nil {
		handlers.WriteError(w, err)
		return
	}

	// リクエストボディの読み込み
	var in updateRequest
	if err := handlers.ReadJSON(r, &in); err != nil {
		// TODO: apperr
		handlers.WriteError(w, err)
		return
	}

	// TODO: configで設定する
	loc, _ := time.LoadLocation("Asia/Tokyo")

	params := convert.UpdateParams{
		Name:     in.Name,
		DueAt:    in.DueAt,
		Status:   in.Status,
		ClearDue: in.ClearDue != nil && *in.ClearDue,
	}
	upd, err := convert.FromUpdateInput(params, loc)
	if err != nil {
		// TODO: apperr
		handlers.WriteError(w, err)
		return
	}

	task, err := t.TaskService.UpdateTask(taskID, upd)
	if err != nil {
		handlers.WriteError(w, err)
		return
	}
	handlers.WriteJSON(w, http.StatusOK, task)
}
