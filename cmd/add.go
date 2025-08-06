package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"golang/tutorial/todo/internal/services/task"
)

var setDueDate bool

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new task",
	Long:  "You can add a new task to your todo list with this command.",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			fmt.Println("Error: Task name is required")
			return
		case 1:
			fmt.Printf("%s\n", task.AddTask("tasks.json", args[0], setDueDate, ""))
		case 2:
			fmt.Printf("%s\n", task.AddTask("tasks.json", args[0], setDueDate, args[1]))
		default:
			fmt.Println("Error: Too many arguments")
		}
	},
}

func init() {
	addCmd.Flags().BoolVarP(&setDueDate, "due", "d", false, "Set a due date for the task")
	rootCmd.AddCommand(addCmd)
}
