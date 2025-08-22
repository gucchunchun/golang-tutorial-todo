package taskhdl

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog"

	"golang/tutorial/todo/internal/adapters/convert"
	"golang/tutorial/todo/internal/api/handlers"
	"golang/tutorial/todo/internal/models"
)

type TaskService interface {
	AddTask(ctx context.Context, c models.TaskCreate) (models.Task, error)
	GetTask(taskID string) (models.TaskOutput, error)
	ListTasks() ([]models.TaskOutput, error)
	UpdateTask(taskID string, updates models.TaskUpdate) (models.Task, error)
}

type TaskHandler struct {
	Logger      *zerolog.Logger
	TaskService TaskService
}

func NewTaskHandler(logger *zerolog.Logger, taskService TaskService) *TaskHandler {
	return &TaskHandler{
		Logger:      logger,
		TaskService: taskService,
	}
}

func (h *TaskHandler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /tasks", h.add)
	mux.HandleFunc("POST /tasks/csv", h.bulkUploadCSV)
	mux.HandleFunc("GET /tasks/{id}", h.get)
	mux.HandleFunc("GET /tasks", h.getList)
	mux.HandleFunc("GET /tasks/csv", h.downloadCSV)
	mux.HandleFunc("PATCH /tasks/{id}", h.update)
}

func (h *TaskHandler) add(w http.ResponseWriter, r *http.Request) {
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

	task, err := h.TaskService.AddTask(r.Context(), c)
	if err != nil {
		handlers.WriteError(w, err)
		return
	}

	handlers.WriteJSON(w, http.StatusCreated, task)
}

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

func (h *TaskHandler) getList(w http.ResponseWriter, r *http.Request) {
	/*
		Reference: O'REILLY「実用GO言語」8.1 p.173
		スライスをエンコードする場合, nilスライスはnullとして扱われる。
		そのためから配列としてエンコードするには、空のスライスを使用する必要がある。
	*/
	tasks, err := h.TaskService.ListTasks()
	if err != nil {
		handlers.WriteError(w, err)
		return
	}

	handlers.WriteJSON(w, http.StatusOK, tasks)
}

func (t *TaskHandler) update(w http.ResponseWriter, r *http.Request) {
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
