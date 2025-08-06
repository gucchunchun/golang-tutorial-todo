package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"golang/tutorial/todo/internal/services/task"
)

func newListCmd() *cobra.Command {
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		Long:  "You can list all tasks in your todo list with this command.",
		Run: func(cmd *cobra.Command, args []string) {
			rows, err := task.ListTasks("tasks.json")
			if err != nil {
				fmt.Printf("Error listing tasks: %v\n", err)
				return
			}

			// ヘッダーを表示
			fmt.Println("ID   | Name                 | Status   | CreatedAt           | DueDate     | TimeLeft")
			fmt.Println("-----------------------------------------------------------------------------------------")
			for _, row := range rows {
				fmt.Println(row)
			}
		},
	}

	return listCmd
}
