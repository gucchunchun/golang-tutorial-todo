package storage_test

import (
	"os"
	"testing"
	"time"

	"golang/tutorial/todo/internal/models"
	"golang/tutorial/todo/internal/storage"
)

func TestStorage_LoadEmptyTasks(t *testing.T) {
	file, err := os.CreateTemp("", "test_tasks.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	storage := storage.New(file.Name())

	tasks, err := storage.LoadTasks()
	if err != nil {
		t.Fatalf("LoadTasks returned error: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("LoadTasks returned %d tasks, want 0", len(tasks))
	}
}

func TestStorage_StoreTask(t *testing.T) {
	file, err := os.CreateTemp("", "test_tasks.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	storage := storage.New(file.Name())

	task1, err := models.NewTask("Buy milk", nil, time.Now())
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	tasks := []models.Task{task1}

	err = storage.SaveTasks(tasks)
	if err != nil {
		t.Fatalf("SaveTasks returned error: %v", err)
	}
}
