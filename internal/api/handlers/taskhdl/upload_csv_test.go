package taskhdl

import (
	"bytes"
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

func TestBulkUploadCSV_Success(t *testing.T) {

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

	addFunc := func(s *stubTaskService, c models.TaskCreate) (models.Task, error) {
		if s.Fail {
			return models.Task{}, os.ErrInvalid
		}
		s.Tasks = append(s.Tasks, models.NewTask(models.TaskID(s.NextID), c.Name, models.Date(time.Now()), models.Date(time.Now()), c.DueAt))
		s.NextID++
		return s.Tasks[len(s.Tasks)-1], nil
	}

	svc := &stubTaskService{
		addTaskFunc: addFunc,
	}
	h := NewTaskHandler(&logger, svc)

	r := httptest.NewRequest("POST", "/tasks/csv", buf)
	r.Header.Set("Content-Type", mw.FormDataContentType())
	w := httptest.NewRecorder()

	h.bulkUploadCSV(w, r)

	require.Equal(t, http.StatusCreated, w.Code)
	require.Len(t, svc.Tasks, 3)
	require.Equal(t, "Buy milk", svc.Tasks[0].Name)
}
