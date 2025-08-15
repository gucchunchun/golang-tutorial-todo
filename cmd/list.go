package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"golang/tutorial/todo/internal/utils"
)

func newListCmd(svc TaskService) *cobra.Command {
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

			rows := make([]string, len(tasks))
			for i, t := range tasks {
				rows[i] = utils.FormatTaskOutput(
					t.ID,
					t.Name,
					t.Status.String(),
					t.CreatedAt,
					*t.DueDate,
				)
			}

			// ヘッダーを表示
			fmt.Println("ID                                   | Name                 | Status   | CreatedAt           | DueDate     | TimeLeft")
			fmt.Println("-----------------------------------------------------------------------------------------")
			for _, row := range rows {
				fmt.Println(row)
			}
		},
	}

	return listCmd
}
