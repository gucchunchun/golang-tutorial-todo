package taskhdl

import (
	"encoding/csv"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"golang/tutorial/todo/internal/models"
)

func TestDownloadCSV_Success(t *testing.T) {
	due := models.Date(time.Date(2025, 8, 30, 0, 0, 0, 0, time.UTC))
	task := models.Task{ID: 1, Name: "Task 1", Status: models.StatusPending, DueAt: &due}
	outputs := models.Tasks{task}.TaskOutputs()

	svc := &stubTaskService{listTasksFunc: func(s *stubTaskService) ([]models.TaskOutput, error) { return outputs, nil }}
	h := &TaskHandler{Logger: getTestLogger(), TaskService: svc}

	req := httptest.NewRequest(http.MethodGet, "/tasks/csv", nil)
	rec := httptest.NewRecorder()
	h.downloadCSV(rec, req)

	res := rec.Result()
	if ct := res.Header.Get("Content-Type"); ct != "text/csv" {
		t.Fatalf("expected Content-Type text/csv, got %q", ct)
	}
	if cd := res.Header.Get("Content-Disposition"); !strings.Contains(cd, "tasks.csv") {
		t.Fatalf("expected filename in Content-Disposition, got %q", cd)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	reader := csv.NewReader(rec.Body)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed reading CSV: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(records))
	}
	header := records[0]
	expectedHeader := []string{"ID", "Name", "Status", "Due Date", "TimeLeft"}
	if len(header) != len(expectedHeader) {
		t.Fatalf("unexpected header length: %d", len(header))
	}
	for i, col := range expectedHeader {
		if header[i] != col {
			t.Fatalf("header[%d]=%q want %q", i, header[i], col)
		}
	}
	row := records[1]
	if row[0] != task.ID.String() {
		t.Errorf("ID=%q want %q", row[0], task.ID.String())
	}
	if row[1] != task.Name {
		t.Errorf("Name=%q want %q", row[1], task.Name)
	}
	if row[2] == "" {
		t.Errorf("Status empty")
	}
	if row[3] == "" {
		t.Errorf("Due Date empty")
	}
	if row[4] == "" {
		t.Errorf("TimeLeft empty; expected a duration since due date is in future")
	}
}
