package task

import (
	"context"
	"fmt"
	"time"

	"golang/tutorial/todo/internal/models"
	"golang/tutorial/todo/internal/utils"
)

func (s *Service) AddTask(taskName string, dueDate string) error {
	if !utils.IsValidTaskName(taskName) {
		return fmt.Errorf("invalid task name")
	}

	tasks, err := s.storage.LoadTasks()
	if err != nil {
		return err
	}

	var dueDateTime time.Time
	if dueDate != "" {
		parsedDueDate, err := utils.ParseDate(dueDate)
		if err != nil {
			return fmt.Errorf("failed to parse due date: %v", err)
		}
		dueDateTime = parsedDueDate
	}

	newTask, err := models.NewTask(taskName, &dueDateTime, utils.SystemTime())
	if err != nil {
		return fmt.Errorf("failed to create new task: %v", err)
	}

	tasks = append(tasks, newTask)
	return s.storage.SaveTasks(tasks)
}

func (s *Service) GetRandomQuote() (string, error) {
	quote, err := s.quoteClient.RandomQuote(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get random quote: %v", err)
	}
	return fmt.Sprintf("%s - %s", quote.Author, quote.Text), nil
}
