package storage

import (
	"encoding/json"
	"os"

	"golang/tutorial/todo/internal/models"
)

func LoadTasks(filename string) ([]models.Task, error) {
	// ファイルの存在チェック
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return []models.Task{}, nil
	}

	// ファイルの読み込み
	data, err := os.ReadFile(filename)
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

func SaveTasks(filename string, tasks []models.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
