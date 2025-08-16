package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"golang/tutorial/todo/internal/models"
	"golang/tutorial/todo/internal/services/task"
	"golang/tutorial/todo/internal/utils"
)

func newListCmd(svc task.TaskService) *cobra.Command {
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		Long:  "You can list all tasks in your todo list with this command.",
		Run: func(cmd *cobra.Command, args []string) {
			tasks, err := svc.ListTasks()
			if err != nil {
				fmt.Printf("Error listing tasks: %v\n", err)
				return
			}

			showTaskOutput(tasks)
		},
	}

	return listCmd
}

func showTaskOutput(tasks []models.TaskOutput) {
	// ヘッダーを表示
	fmt.Println("ID                                   | Name                 | Status   | CreatedAt           | DueDate     | TimeLeft")
	fmt.Println("-----------------------------------------------------------------------------------------")

	// タスクを表示
	for _, t := range tasks {
		fmt.Println(formatTask(t))
	}
}

func formatTask(task models.TaskOutput) string {
	formatCreateAt, err := utils.FormatDatetime(task.CreatedAt)
	if err != nil {
		formatCreateAt = "No date"
	}

	formatDueDate := "-"
	if task.DueDate != nil && !task.DueDate.IsZero() {
		formatDueDate, err = utils.FormatDate(*task.DueDate)
		if err != nil {
			formatDueDate = "-"
		}
	}

	timeLeft := "-"
	if task.DueDate != nil && task.TimeLeft != nil {
		timeLeft = utils.FormatDurationToDays(*task.TimeLeft)
		if timeLeft == "" {
			timeLeft = "Over due"
		}
	}

	return fmt.Sprintf("%-16s | %-20s | %-8s | %-19s | %-11s | %-15s",
		task.ID.String(), task.Name, task.Status, formatCreateAt, formatDueDate, timeLeft)
}
