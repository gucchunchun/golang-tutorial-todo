package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"golang/tutorial/todo/internal/adapters/convert"
	"golang/tutorial/todo/internal/services/tasksvc"
)

func newUpdateCmd(svc tasksvc.TaskService) *cobra.Command {
	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update an existing task",
		Long:  "You can update an existing task in your todo list with this command.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: configで設定する
			loc, _ := time.LoadLocation("Asia/Tokyo")

			params := convert.UpdateParams{}
			if status, _ := cmd.Flags().GetString("status"); status != "" {
				params.Status = &status
			}
			if due, _ := cmd.Flags().GetString("due"); due != "" {
				params.DueAt = &due
			}
			if name, _ := cmd.Flags().GetString("name"); name != "" {
				params.Name = &name
			}
			upd, err := convert.FromUpdateInput(params, loc)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			_, err = svc.UpdateTask(args[0], upd)
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
