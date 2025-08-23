package taskhdl

import (
	"context"
	"io"

	"github.com/rs/zerolog"

	"golang/tutorial/todo/internal/models"
)

type stubTaskService struct {
	Tasks          models.Tasks
	NextID         uint64
	Fail           bool
	addTaskFunc    func(s *stubTaskService, c models.TaskCreate) (models.Task, error)
	getTaskFunc    func(s *stubTaskService, taskID string) (models.TaskOutput, error)
	listTasksFunc  func(s *stubTaskService) ([]models.TaskOutput, error)
	updateTaskFunc func(s *stubTaskService, taskID string, updates models.TaskUpdate) (models.Task, error)
}

func (s *stubTaskService) AddTask(ctx context.Context, c models.TaskCreate) (models.Task, error) {
	if s.addTaskFunc == nil {
		panic("addTaskFunc is not set")
	}
	return s.addTaskFunc(s, c)
}

func (s *stubTaskService) GetTask(taskID string) (models.TaskOutput, error) {
	if s.getTaskFunc == nil {
		panic("getTaskFunc is not set")
	}
	return s.getTaskFunc(s, taskID)
}

func (s *stubTaskService) ListTasks() ([]models.TaskOutput, error) {
	if s.listTasksFunc == nil {
		panic("listTasksFunc is not set")
	}
	return s.listTasksFunc(s)
}

func (s *stubTaskService) UpdateTask(taskID string, updates models.TaskUpdate) (models.Task, error) {
	if s.updateTaskFunc == nil {
		panic("updateTaskFunc is not set")
	}
	return s.updateTaskFunc(s, taskID, updates)
}

func getTestLogger() *zerolog.Logger {
	// writerがnilの場合はpanicするため、io.Discardを指定
	l := zerolog.New(io.Discard).Level(zerolog.DebugLevel)
	return &l
}
