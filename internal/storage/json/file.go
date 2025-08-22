package json

import (
	"context"
	"encoding/json"
	"os"

	"golang/tutorial/todo/internal/apperr"
	"golang/tutorial/todo/internal/models"
	"golang/tutorial/todo/internal/utils"
)

type Storage struct {
	fileName string
}

func NewStorage(fileName string) *Storage {
	return &Storage{
		fileName: fileName,
	}
}

func (s *Storage) loadTasks() (models.Tasks, error) {
	// ファイルの存在チェック
	if _, err := os.Stat(s.fileName); os.IsNotExist(err) {
		return []models.Task{}, nil
	}

	// ファイルの読み込み
	data, err := os.ReadFile(s.fileName)
	if err != nil {
		return nil, err
	}

	// ファイルが空の場合
	if len(data) == 0 {
		return []models.Task{}, nil
	}

	// ファイルのデコード
	var tasks []models.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Storage) saveTasks(tasks []models.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.fileName, data, 0644)
}

func (s *Storage) Create(ctx context.Context, name string, dueAt *models.Date) (models.Task, error) {
	tasks, err := s.loadTasks()
	if err != nil {
		return models.Task{}, apperr.E(apperr.CodeUnknown, "Failed to load task", err)
	}

	newTask := models.NewTask(models.ParseTaskIDInt(len(tasks)+1), name, models.Date(utils.SystemTime()), models.Date(utils.SystemTime()), dueAt)
	if err != nil {
		return models.Task{}, apperr.E(apperr.CodeInvalid, "Failed to create task", err)
	}

	return newTask, nil
}

func (s *Storage) GetByID(ctx context.Context, id models.TaskID) (models.Task, error) {
	tasks, err := s.loadTasks()
	if err != nil {
		return models.Task{}, apperr.E(apperr.CodeUnknown, "Failed to load task", err)
	}

	task, ok := tasks.FindByID(id)
	if !ok {
		return models.Task{}, apperr.E(apperr.CodeNotFound, "Task not found", nil)
	}
	return task, nil
}

func (s *Storage) List(ctx context.Context, limit, offset int) (models.Tasks, error) {
	return s.loadTasks()
}

func (s *Storage) Update(ctx context.Context, id models.TaskID, upd models.TaskUpdate) (models.Task, error) {
	tasks, err := s.loadTasks()
	if err != nil {
		return models.Task{}, apperr.E(apperr.CodeUnknown, "Failed to load task", err)
	}

	found := false
	for i, task := range tasks {
		if task.ID == id {
			found = true
			if upd.Name != nil {
				tasks[i].Name = *upd.Name
			}
			if upd.Status != nil {
				tasks[i].Status = *upd.Status
			}
			if upd.DueAt != nil {
				tasks[i].DueAt = upd.DueAt
			}
			break
		}
	}

	if !found {
		return models.Task{}, apperr.E(apperr.CodeNotFound, "Task not found", nil)
	}

	err = s.saveTasks(tasks)
	if err != nil {
		return models.Task{}, apperr.E(apperr.CodeUnknown, "Failed to update task", err)
	}

	return s.GetByID(ctx, id)
}

func (s *Storage) Delete(ctx context.Context, id models.TaskID) error {
	tasks, err := s.loadTasks()
	if err != nil {
		return apperr.E(apperr.CodeUnknown, "Failed to load task", err)
	}

	found := false
	newTasks := make([]models.Task, 0, len(tasks))
	for _, task := range tasks {
		if task.ID == id {
			found = true
			continue
		}
		newTasks = append(newTasks, task)
	}

	if !found {
		return apperr.E(apperr.CodeNotFound, "Task not found", nil)
	}
	return nil
}
