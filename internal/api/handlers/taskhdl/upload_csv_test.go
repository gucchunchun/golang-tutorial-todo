package taskhdl

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"golang/tutorial/todo/internal/models"
)

// stubTaskService テスト用のTaskService
type stubTaskService struct {
	added  []models.TaskCreate
	nextID uint64
	fail   bool
}

func (s *stubTaskService) AddTask(ctx context.Context, c models.TaskCreate) (models.Task, error) {
	if s.fail {
		return models.Task{}, os.ErrInvalid
	}
	s.added = append(s.added, c)
	s.nextID++
	now := models.Date(time.Now())
	return models.NewTask(models.TaskID(s.nextID), c.Name, now, now, c.DueAt), nil
}
func (s *stubTaskService) GetTask(string) (models.TaskOutput, error) { return models.TaskOutput{}, nil }
func (s *stubTaskService) ListTasks() ([]models.TaskOutput, error)   { return nil, nil }
func (s *stubTaskService) UpdateTask(string, models.TaskUpdate) (models.Task, error) {
	return models.Task{}, nil
}

func TestBulkUploadCSV_OK(t *testing.T) {
	t.Parallel()

	// logger
	logger := zerolog.New(os.Stdout)

	// HTTPリクエストのmultipartボディを準備
	buf := &bytes.Buffer{}
	mw := multipart.NewWriter(buf)
	fw, err := mw.CreateFormFile("file", "tasks_valid.csv")
	require.NoError(t, err)
	fp := filepath.Join("testdata", "tasks_valid.csv")
	content, err := os.ReadFile(fp)
	require.NoError(t, err)
	_, err = fw.Write(content)
	require.NoError(t, err)
	mw.Close()

	svc := &stubTaskService{}
	h := NewTaskHandler(&logger, svc)

	r := httptest.NewRequest("POST", "/tasks/csv", buf)
	r.Header.Set("Content-Type", mw.FormDataContentType())
	w := httptest.NewRecorder()

	h.bulkUploadCSV(w, r)

	require.Equal(t, http.StatusCreated, w.Code)
	require.Len(t, svc.added, 3)
	require.Equal(t, "Buy milk", svc.added[0].Name)
}
