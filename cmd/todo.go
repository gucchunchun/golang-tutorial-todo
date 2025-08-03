package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

type Task struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`
	DueDate   string `json:"due_date"`
}

func Add(input string) (result string) {
	// Task27: Regex
	if match, err := regexp.MatchString("[a-z]+", input); err != nil || !match {
		fmt.Println("Error: Task must be a valid string")
		return
	}

	task := strings.ReplaceAll(input, `"`, ``)
	if task == "" {
		fmt.Println("Error: Task cannot be empty")
		return
	}

	if err := AddTask("tasks.json", task); err != nil {
		fmt.Println("Error: Saving tasks:", err)
		return
	}
	return fmt.Sprintf("Task added: %s\n", task)

}
func AddTask(filename, input string) error {
	// Task27: Regex
	if match, err := regexp.MatchString("[a-zA-Z]+", input); err != nil || !match {
		return fmt.Errorf("invalid task name")
	}

	// Load existing tasks
	tasks, err := LoadTasks(filename)
	if err != nil {
		return err
	}

	// Create new task
	newTask := Task{
		ID:        uint(len(tasks) + 1), // Simple ID generation
		Name:      strings.TrimSpace(input),
		Status:    0,
		CreatedAt: time.Now().Format(time.RFC3339),
		DueDate:   "",
	}

	// Append and save
	tasks = append(tasks, newTask)
	return SaveTasks(filename, tasks)
}

func LoadTasks(filename string) ([]Task, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return []Task{}, nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Check if the file is empty
	if len(data) == 0 {
		return []Task{}, nil
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func SaveTasks(filename string, tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
