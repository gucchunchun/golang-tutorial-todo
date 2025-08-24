package taskhdl

import (
	"encoding/csv"
	"golang/tutorial/todo/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownloadCSV_Success(t *testing.T) {
	d := models.Date(time.Date(2025, 8, 30, 0, 0, 0, 0, time.UTC))
	task := models.Task{
		ID:     1,
		Name:   "Task 1",
		Status: models.StatusPending,
		DueAt:  &d,
	}
	out := (models.Tasks{task}).TaskOutputs()

	svc := &stubTaskService{
		listTasksFunc: func(s *stubTaskService) ([]models.TaskOutput, error) {
			return out, nil
		},
	}
	h := &TaskHandler{Logger: getTestLogger(), TaskService: svc}

	req := httptest.NewRequest(http.MethodGet, "/tasks/csv", nil)
	rec := httptest.NewRecorder()
	h.downloadCSV(rec, req)
	res := rec.Result()

	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "text/csv", res.Header.Get("Content-Type"))
	assert.Contains(t, res.Header.Get("Content-Disposition"), "tasks.csv")

	r := csv.NewReader(rec.Body)
	records, err := r.ReadAll()
	require.NoError(t, err)
	require.Len(t, records, 2, "expect header + single data row")

	expectedHeader := []string{"ID", "タスク名", "ステータス", "作成日", "更新日", "期限", "残り時間"}
	require.Equal(t, expectedHeader, records[0], "CSV header mismatch")

	row := records[1]
	require.GreaterOrEqual(t, len(row), len(expectedHeader))

	assert.Equal(t, task.ID.String(), row[0])
	assert.Equal(t, task.Name, row[1])
	assert.NotEmpty(t, row[2], "status should not be empty")
	assert.NotEmpty(t, row[5], "due date should not be empty")
	assert.NotEmpty(t, row[6], "remaining time should not be empty (future due date)")
}
