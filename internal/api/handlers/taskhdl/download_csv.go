package taskhdl

import (
	"net/http"

	"github.com/gocarina/gocsv"

	"golang/tutorial/todo/internal/api/handlers"
)

func (t *TaskHandler) downloadCSV(w http.ResponseWriter, r *http.Request) {
	t.Logger.Debug().Msg("--- downloadCSV ---")
	// CSVダウンロード処理
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=tasks.csv")

	// Taskの読み込み
	tasks, err := t.TaskService.ListTasks()
	if err != nil {
		t.Logger.Error().Err(err).Msg("failed to fetch tasks")
		handlers.WriteError(w, err)
		return
	}

	t.Logger.Debug().Msgf("tasks: %+v", tasks)

	/*
		Reference: O'REILLY「実用GO言語」8.2. p.183
		構造体にcsvタグを付与することで、gocsvパッケージを使用してCSVの読み書きが可能になる。
	*/
	if err := gocsv.Marshal(tasks, w); err != nil {
		t.Logger.Error().Err(err).Msg("failed to write CSV")
		handlers.WriteError(w, err)
		return
	}

	// encoding/csvを使った場合のコード例：
	// // CSVの書き込み
	// cw := csv.NewWriter(w)
	// // 	バッファにあるデータを書き込み
	// defer cw.Flush()

	// // ヘッダーの書き込み
	// if err := cw.Write([]string{"ID", "Name", "Status", "Due Date", "TimeLeft"}); err != nil {
	// 	t.Logger.Error().Err(err).Msg("failed to write CSV header")
	// 	handlers.WriteError(w, err)
	// 	return
	// }

	// // 各タスクの書き込み
	// for _, task := range tasks {
	// 	// TODO: nilチェックが必要ないようにTask, TaskOutputを修正する
	// 	dueDateAt := ""
	// 	if task.DueAt != nil {
	// 		dueDateAt = task.DueAt.Format()
	// 	}
	// 	timeLeft := ""
	// 	if task.TimeLeft != nil {
	// 		timeLeft = task.TimeLeft.String()
	// 	}
	// 	if err := cw.Write([]string{
	// 		task.ID.String(),
	// 		task.Name,
	// 		task.Status.String(),
	// 		dueDateAt,
	// 		timeLeft,
	// 	}); err != nil {
	// 		t.Logger.Error().Err(err).Msg("failed to write task to CSV")
	// 		handlers.WriteError(w, err)
	// 		return
	// 	}
	// }

	// if cw.Error() != nil {
	// 	t.Logger.Error().Err(cw.Error()).Msg("failed to write CSV")
	// 	handlers.WriteError(w, cw.Error())
	// 	return
	// }
}
