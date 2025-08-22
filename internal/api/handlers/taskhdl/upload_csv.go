package taskhdl

import (
	"encoding/csv"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang/tutorial/todo/internal/adapters/convert"
	"golang/tutorial/todo/internal/api/handlers"
)

func (h *TaskHandler) bulkUploadCSV(w http.ResponseWriter, r *http.Request) {
	if h.TaskService == nil {
		http.Error(w, "task service not initialized", http.StatusInternalServerError)
		return
	}

	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data") {
		http.Error(w, "content-type must be multipart/form-data", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file field 'file' is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	/*
		Reference: O'REILLY「実用GO言語」8.2 p.180
		基本的にCSVに関しては、encoding/json パッケージを使用してJSON形式に変換することができる
	*/
	reader := csv.NewReader(file)

	var createdIDs []string
	line := 0
	const maxRows = 5000

	for {
		rec, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			http.Error(w, "invalid CSV", http.StatusBadRequest)
			return
		}
		line++

		// 空白行はスキップ
		if len(rec) == 0 || (len(rec) == 1 && strings.TrimSpace(rec[0]) == "") {
			continue
		}

		if line == 1 {
			first := strings.ToLower(strings.TrimSpace(rec[0]))
			if first == "name" {
				continue
			}
		}

		name := strings.TrimSpace(rec[0])
		if name == "" {
			http.Error(w, "line "+strconv.Itoa(line)+": name required", http.StatusBadRequest)
			return
		}
		h.Logger.Debug().Msgf("name :%v", name)

		var duePtr *string
		if len(rec) > 1 {
			ds := strings.TrimSpace(rec[1])
			if ds != "" {
				duePtr = &ds
			}
		}
		h.Logger.Debug().Msgf("duePtr :%v", duePtr)

		cparams := convert.CreateParams{Name: name, DueAt: duePtr}
		h.Logger.Debug().Msgf("cparams :%v", cparams)
		taskCreate, err := convert.FromCreateInput(cparams, time.Local)
		if err != nil {
			http.Error(w, "line "+strconv.Itoa(line)+": invalid due date (expected YYYY-MM-DD)", http.StatusBadRequest)
			return
		}

		task, err := h.TaskService.AddTask(r.Context(), taskCreate)
		if err != nil {
			if h.Logger != nil {
				h.Logger.Error().
					Int("line", line).
					Str("name", name).
					Err(err).
					Msg("bulk csv add task failed")
			}
			w.Header().Set("X-Error-Line", strconv.Itoa(line))
			handlers.WriteError(w, err)
			return
		}
		h.Logger.Debug().Msgf("task :%v", task)
		createdIDs = append(createdIDs, task.ID.String())
		h.Logger.Debug().Msgf("createdIDs :%v", createdIDs)

		if len(createdIDs) > maxRows {
			http.Error(w, "row limit exceeded ("+strconv.Itoa(maxRows)+")", http.StatusBadRequest)
			return
		}
	}

	handlers.WriteJSON(w, http.StatusCreated, map[string]any{
		"created": createdIDs,
		"count":   len(createdIDs),
	})
}
