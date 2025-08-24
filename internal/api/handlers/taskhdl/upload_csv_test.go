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

	"github.com/stretchr/testify/require"

	"golang/tutorial/todo/internal/models"
)

func TestBulkUploadCSV_Success(t *testing.T) {
	// ハンドラーの準備
	addFunc := func(s *stubTaskService, c models.TaskCreate) (models.Task, error) {
		if s.Fail {
			return models.Task{}, os.ErrInvalid
		}
		s.Tasks = append(s.Tasks, models.NewTask(models.TaskID(s.NextID), c.Name, models.Date(time.Now()), models.Date(time.Now()), c.DueAt))
		s.NextID++
		return s.Tasks[len(s.Tasks)-1], nil
	}

	type testTaskInput struct {
		Name  string
		DueAt *models.Date
	}

	type testCase struct {
		Name    string
		CSVFile string
		Tasks   []testTaskInput
	}

	loc := time.Local
	tdue1, err := time.ParseInLocation("2006-01-02", "2025-09-01", loc)
	require.NoError(t, err)
	due1 := models.Date(tdue1)
	tdue3, err := time.ParseInLocation("2006-01-02", "2025-12-31", loc)
	require.NoError(t, err)
	due3 := models.Date(tdue3)
	tests := []testCase{
		{
			"valid",
			"tasks_valid.csv",
			[]testTaskInput{
				{Name: "Buy milk", DueAt: &due1},
				{Name: "Pay bills", DueAt: nil},
				{Name: "Read book", DueAt: &due3},
			},
		},
		{
			"valid with BOM",
			"tasks_valid_with_bom.csv",
			[]testTaskInput{
				{Name: "Buy milk", DueAt: &due1},
				{Name: "Pay bills", DueAt: nil},
				{Name: "Read book", DueAt: &due3},
			},
		},
		{
			"valid with comment out line",
			"tasks_valid_with_comment.csv",
			[]testTaskInput{
				{Name: "Buy milk", DueAt: &due1},
				{Name: "Pay bills", DueAt: nil},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// テスト用のサービスとハンドラーを準備
			svc := &stubTaskService{
				addTaskFunc: addFunc,
			}
			h := NewTaskHandler(getTestLogger(), svc)

			// HTTPリクエストのmultipartボディを準備
			buf := &bytes.Buffer{}
			mw := multipart.NewWriter(buf)
			fw, err := mw.CreateFormFile("file", test.CSVFile)
			require.NoError(t, err)
			fp := filepath.Join("testdata", test.CSVFile)
			content, err := os.ReadFile(fp)
			require.NoError(t, err)
			_, err = fw.Write(content)
			require.NoError(t, err)
			mw.Close()

			// HTTPリクエストを作成
			r := httptest.NewRequest("POST", "/tasks/csv", buf)
			r.Header.Set("Content-Type", mw.FormDataContentType())
			w := httptest.NewRecorder()

			// ハンドラーを呼び出す
			h.bulkUploadCSV(w, r)

			// 結果を検証
			require.Equal(t, http.StatusCreated, w.Code)
			require.Len(t, svc.Tasks, len(test.Tasks))
			for i, task := range test.Tasks {
				require.Equal(t, task.Name, svc.Tasks[i].Name)
				if task.DueAt != nil {
					require.Equal(t, task.DueAt, svc.Tasks[i].DueAt)
				} else {
					require.Nil(t, svc.Tasks[i].DueAt)
				}
			}
		})
	}
}
