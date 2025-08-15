package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"golang/tutorial/todo/internal/models"
	"golang/tutorial/todo/internal/services/task"
	"golang/tutorial/todo/internal/utils"
)

func newUpdateCmd(svc task.Service) *cobra.Command {
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
				newStatus, err := models.ParseStatus(status)
				if err != nil {
					fmt.Printf("Error parsing status: %v\n", err)
					return
				}
				updates.Status = &newStatus
			}
			if due, _ := cmd.Flags().GetString("due"); due != "" {
				newDue, err := utils.ParseDate(due)
				if err != nil {
					fmt.Printf("Error parsing due date: %v\n", err)
					return
				}
				updates.Due = &newDue
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
