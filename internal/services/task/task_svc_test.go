package task_test

import (
	"errors"
	"testing"
	"time"

	"golang/tutorial/todo/internal/apperr"
	"golang/tutorial/todo/internal/models"
	"golang/tutorial/todo/internal/services/task"

	"github.com/google/uuid"
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
	t.Run("AddTask", func(t *testing.T) {
		type args struct {
			name    string
			dueDate string
		}
		tests := map[string]struct {
			args     args
			wantErr  bool
			wantCode apperr.Code
		}{
			"ok":               {args{"name", ""}, false, 0},
			"ok with due":      {args{"name", "2025-10-31"}, false, 0},
			"empty name":       {args{"", ""}, true, apperr.CodeInvalid},
			"invalid due date": {args{"name", "invalid date"}, true, apperr.CodeInvalid},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				svc := task.NewTaskService(&storageClientStub{})

				err := svc.AddTask(tc.args.name, tc.args.dueDate)

				if !tc.wantErr {
					assert.NoError(t, err)
					return
				}

				assert.Error(t, err)
				var ae *apperr.Error
				assert.ErrorAs(t, err, &ae)
				assert.Equal(t, tc.wantCode, ae.Code) // expected, actual
			})
		}
	})
	t.Run("GetTask", func(t *testing.T) {
		idstr := uuid.NewString()
		taskID, _ := models.ParseTaskID(idstr)
		tests := map[string]struct {
			taskID   string
			loadFunc func() (models.Tasks, error)
			wantErr  bool
			wantCode apperr.Code
		}{
			"ok": {
				taskID: idstr,
				loadFunc: func() (models.Tasks, error) {
					return models.Tasks{
						{ID: taskID, Name: "first"},
					}, nil
				},
				wantErr: false,
			},
			"invalid id": {
				taskID:   "abc",
				loadFunc: func() (models.Tasks, error) { return nil, nil },
				wantErr:  true,
				wantCode: apperr.CodeInvalid,
			},
			"not found": {
				taskID: uuid.NewString(),
				loadFunc: func() (models.Tasks, error) {
					return models.Tasks{{ID: taskID, Name: "first"}}, nil
				},
				wantErr:  true,
				wantCode: apperr.CodeNotFound,
			},
			"load error bubbles": {
				taskID: idstr,
				loadFunc: func() (models.Tasks, error) {
					return nil, errors.New("disk")
				},
				wantErr: true,
			},
		}

		for name, tc := range tests {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				st := &storageClientStub{LoadFunc: tc.loadFunc}
				svc := task.NewTaskService(st)

				got, err := svc.GetTask(tc.taskID)

				if !tc.wantErr {
					assert.NoError(t, err)
					assert.NotZero(t, got)
					return
				}
				assert.Error(t, err)
				if tc.wantCode != 0 {
					var ae *apperr.Error
					assert.ErrorAs(t, err, &ae)
					assert.Equal(t, tc.wantCode, ae.Code)
				}
			})
		}
	})
	t.Run("ListTasks", func(t *testing.T) {
		tests := map[string]struct {
			loadFunc func() (models.Tasks, error)
			wantLen  int
			wantErr  bool
			wantCode apperr.Code
		}{
			"ok: two tasks": {
				loadFunc: func() (models.Tasks, error) {
					id1, _ := models.ParseTaskID("1")
					id2, _ := models.ParseTaskID("2")
					return models.Tasks{
						{ID: id1, Name: "a"},
						{ID: id2, Name: "b"},
					}, nil
				},
				wantLen: 2,
			},
			"load error -> CodeUnknown": {
				loadFunc: func() (models.Tasks, error) { return nil, errors.New("io") },
				wantErr:  true,
				wantCode: apperr.CodeUnknown,
			},
		}

		for name, tc := range tests {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				st := &storageClientStub{LoadFunc: tc.loadFunc}
				svc := task.NewTaskService(st)

				got, err := svc.ListTasks()
				if !tc.wantErr {
					assert.NoError(t, err)
					assert.Len(t, got, tc.wantLen)
					return
				}
				assert.Error(t, err)
				var ae *apperr.Error
				assert.ErrorAs(t, err, &ae)
				assert.Equal(t, tc.wantCode, ae.Code)
			})
		}
	})
	t.Run("UpdateTask", func(t *testing.T) {
		idstr := uuid.NewString()
		taskID, _ := models.ParseTaskID(idstr)

		newName := "renamed"
		newDue, _ := time.Parse("2006-01-02", "2025-10-31")
		newDueStr := newDue.Format("2006-01-02")

		tests := map[string]struct {
			taskID   string
			updates  models.TaskUpdate
			loadFunc func() (models.Tasks, error)
			saveFunc func([]models.Task) error
			wantErr  bool
			wantCode apperr.Code
		}{
			"ok: update name only": {
				taskID:  idstr,
				updates: models.TaskUpdate{Name: &newName},
				loadFunc: func() (models.Tasks, error) {
					return models.Tasks{{ID: taskID, Name: "old"}}, nil
				},
				saveFunc: func(ts []models.Task) error {
					assert.Equal(t, newName, ts[0].Name)
					return nil
				},
			},
			"invalid id": {
				taskID:   "abc",
				updates:  models.TaskUpdate{Name: &newName},
				loadFunc: func() (models.Tasks, error) { return nil, nil },
				wantErr:  true,
				wantCode: apperr.CodeInvalid,
			},
			"not found": {
				taskID:  uuid.NewString(),
				updates: models.TaskUpdate{Name: &newName},
				loadFunc: func() (models.Tasks, error) {
					return models.Tasks{{ID: taskID, Name: "old"}}, nil
				},
				wantErr:  true,
				wantCode: apperr.CodeNotFound,
			},
			"load fails -> CodeInvalid (as written)": {
				taskID:   idstr,
				updates:  models.TaskUpdate{Name: &newName},
				loadFunc: func() (models.Tasks, error) { return nil, errors.New("io") },
				wantErr:  true,
				wantCode: apperr.CodeInvalid,
			},
			"save fails -> CodeUnknown": {
				taskID:  idstr,
				updates: models.TaskUpdate{Name: &newName},
				loadFunc: func() (models.Tasks, error) {
					return models.Tasks{{ID: taskID, Name: "old"}}, nil
				},
				saveFunc: func([]models.Task) error { return errors.New("disk") },
				wantErr:  true,
				wantCode: apperr.CodeUnknown,
			},
			"ok: update due only (parses date)": {
				taskID:  idstr,
				updates: models.TaskUpdate{Due: &newDueStr},
				loadFunc: func() (models.Tasks, error) {
					return models.Tasks{{ID: taskID, Name: "old"}}, nil
				},
				saveFunc: func(ts []models.Task) error {
					assert.NotNil(t, ts[0].DueDate)
					assert.Equal(t, newDue, ts[0].DueDate.UTC())
					return nil
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				st := &storageClientStub{LoadFunc: tc.loadFunc, SaveFunc: tc.saveFunc}
				svc := task.NewTaskService(st)

				err := svc.UpdateTask(tc.taskID, tc.updates)

				if !tc.wantErr {
					assert.NoError(t, err)
					return
				}
				assert.Error(t, err)
				var ae *apperr.Error
				assert.ErrorAs(t, err, &ae)
				assert.Equal(t, tc.wantCode, ae.Code)
			})
		}
	})
}
