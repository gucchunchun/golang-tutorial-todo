package storage

import (
	"encoding/json"
	"os"

	"golang/tutorial/todo/internal/models"
)

type Storage struct {
	fileName string
}

func NewStorage(fileName string) *Storage {
	return &Storage{
		fileName: fileName,
	}
}

func (s *Storage) LoadTasks() (models.Tasks, error) {
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

func (s *Storage) SaveTasks(tasks []models.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.fileName, data, 0644)
}
