package tasksvc

import (
	"context"
	"golang/tutorial/todo/internal/apperr"
	"golang/tutorial/todo/internal/models"
	"golang/tutorial/todo/internal/validation"
	"log"
)

type TaskService struct {
	repo Repository
}

func NewTaskService(repo Repository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) AddTask(c models.TaskCreate) (models.Task, error) {
	log.Println("--- AddTask ---")
	// バリデーション
	if err := validation.ValidateCreateTaskInput(validation.CreateTaskInput{
		Name:    c.Name,
		DueDate: c.DueAt,
	}); err != nil {
		return models.Task{}, apperr.E(apperr.CodeInvalid, "Validation error", err)
	}

	ctx := context.Background()

	task, err := s.repo.Create(ctx, c.Name, c.DueAt)
	if err != nil {
		return models.Task{}, apperr.E(apperr.CodeUnknown, "Failed to create task", err)
	}
	return task, nil
}

func (s *TaskService) GetTask(taskID string) (models.TaskOutput, error) {
	parsedID, err := models.ParseTaskID(taskID)
	if err != nil {
		return models.TaskOutput{}, apperr.E(apperr.CodeInvalid, "Validation error", err)
	}
	ctx := context.Background()
	task, err := s.repo.GetByID(ctx, parsedID)
	if err != nil {
		return models.TaskOutput{}, err
	}

	return task.TaskOutput(), nil
}

func (s *TaskService) ListTasks() ([]models.TaskOutput, error) {
	ctx := context.Background()
	tasks, err := s.repo.List(ctx, 0, 0)
	if err != nil {
		return nil, apperr.E(apperr.CodeUnknown, "Failed to load tasks", err)
	}

	return tasks.TaskOutputs(), nil
}

func (s *TaskService) UpdateTask(taskID string, updates models.TaskUpdate) (models.Task, error) {
	// バリデーション
	parsedID, err := models.ParseTaskID(taskID)
	if err != nil {
		return models.Task{}, apperr.E(apperr.CodeInvalid, "Validation error", err)
	}
	if err := validation.ValidateUpdateTaskInput(validation.UpdateTaskInput{
		Name:    updates.Name,
		DueDate: updates.DueAt,
		Status:  updates.Status,
	}); err != nil {
		return models.Task{}, apperr.E(apperr.CodeInvalid, "Validation error", err)
	}

	ctx := context.Background()
	task, err := s.repo.Update(ctx, parsedID, updates)
	if err != nil {
		return models.Task{}, apperr.E(apperr.CodeInvalid, "Validation error", err)
	}

	return task, nil
}
