package taskhdl

import (
	"encoding/csv"
	"net/http"

	"golang/tutorial/todo/internal/api/handlers"
)

func (t *TaskHandler) downloadCSV(w http.ResponseWriter, r *http.Request) {
	t.Logger.Debug().Msg("---downloadCSV called---")
	// CSVダウンロード処理
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=tasks.csv")

	// Taskの読み込み
	t.Logger.Debug().Msg("fetching task data")
	tasks, err := t.TaskService.ListTasks()
	if err != nil {
		t.Logger.Error().Err(err).Msg("failed to fetch tasks")
		handlers.WriteError(w, err)
		return
	}
	t.Logger.Debug().Msg("fetching task data: ok")

	// CSVの書き込み
	t.Logger.Debug().Msg("csv.Writer creation")
	cw := csv.NewWriter(w)
	t.Logger.Debug().Msg("csv.Writer creation: ok")
	// 	バッファにあるデータを書き込み
	defer cw.Flush()

	// ヘッダーの書き込み
	t.Logger.Debug().Msg("write header")
	if err := cw.Write([]string{"ID", "Name", "Status", "Due Date", "TimeLeft"}); err != nil {
		t.Logger.Error().Err(err).Msg("failed to write CSV header")
		handlers.WriteError(w, err)
		return
	}
	t.Logger.Debug().Msg("write header: ok")

	// 各タスクの書き込み
	t.Logger.Debug().Msg("write tasks")
	for _, task := range tasks {
		// TODO: nilチェックが必要ないようにTask, TaskOutputを修正する
		dueDateAt := ""
		if task.DueAt != nil {
			dueDateAt = task.DueAt.Format()
		}
		timeLeft := ""
		if task.TimeLeft != nil {
			timeLeft = task.TimeLeft.String()
		}
		if err := cw.Write([]string{
			task.ID.String(),
			task.Name,
			task.Status.String(),
			dueDateAt,
			timeLeft,
		}); err != nil {
			t.Logger.Error().Err(err).Msg("failed to write task to CSV")
			handlers.WriteError(w, err)
			return
		}
	}
	t.Logger.Debug().Msg("write tasks: ok")

	if cw.Error() != nil {
		t.Logger.Error().Err(cw.Error()).Msg("failed to write CSV")
		handlers.WriteError(w, cw.Error())
		return
	}
}
