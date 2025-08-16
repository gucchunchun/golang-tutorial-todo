package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"golang/tutorial/todo/internal/models"
	"golang/tutorial/todo/internal/services/task"
)

func newUpdateCmd(svc task.TaskService) *cobra.Command {
	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update an existing task",
		Long:  "You can update an existing task in your todo list with this command.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			taskID, err := models.ParseTaskID(args[0])
			if err != nil {
				fmt.Printf("Error parsing task ID: %v\n", err)
				return
			}

			updates := models.TaskUpdate{}
			if status, _ := cmd.Flags().GetString("status"); status != "" {
				updates.Status = &status
			}
			if due, _ := cmd.Flags().GetString("due"); due != "" {
				updates.Due = &due
			}
			if name, _ := cmd.Flags().GetString("name"); name != "" {
				updates.Name = &name
			}
			err = svc.UpdateTask(taskID, updates)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Println("task updated successfully")
		},
	}

	updateCmd.Flags().StringP("status", "s", "", "Update the status of the task")
	updateCmd.Flags().StringP("due", "d", "", "Update the due date of the task")
	updateCmd.Flags().StringP("name", "n", "", "Update the name of the task")

	return updateCmd
}
