package task_test

import (
	"testing"

	"golang/tutorial/todo/internal/models"
	"golang/tutorial/todo/internal/services/task"

	"github.com/stretchr/testify/assert"
)

// stub
type storageClientStub struct {
	LoadFunc func() (models.Tasks, error)
	SaveFunc func(tasks []models.Task) error
}

func (s *storageClientStub) LoadTasks() (models.Tasks, error) {
	if s.LoadFunc == nil {
		return models.Tasks{}, nil
	}
	return s.LoadFunc()
}
func (s *storageClientStub) SaveTasks(tasks []models.Task) error {
	if s.SaveFunc == nil {
		return nil
	}
	return s.SaveFunc(tasks)
}

func TestTaskService(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		type args struct {
			name    string
			dueDate string
		}
		tests := map[string]struct {
			args    args
			wantErr error
		}{
			"ok":               {args{"name", ""}, nil},
			"ok with due":      {args{"name", "2025-10-31"}, nil},
			"empty name":       {args{"", ""}, task.ErrValidation},
			"invalid due date": {args{"name", "invalid date"}, task.ErrValidation},
		}

		for name, tc := range tests {
			tc := tc // capture
			t.Run(name, func(t *testing.T) {
				t.Parallel() // safe if svc/state is not shared
				svc := task.NewTaskService(&storageClientStub{})

				err := svc.AddTask(tc.args.name, tc.args.dueDate)

				if tc.wantErr != nil {
					assert.NoError(t, err)
					return
				}
				assert.Error(t, err)
				assert.ErrorAs(t, tc.wantErr, err)
			})
		}
	})
}
