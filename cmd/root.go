package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"golang/tutorial/todo/internal/models"
)

type TaskService interface {
	AddTask(title string, setDueDate bool, dueDate string) error
	GetRandomQuote() (string, error)
	ListTasks() ([]models.Task, error)
	UpdateTask(id models.TaskID, updates models.TaskUpdate) error
}

var rootCmd = &cobra.Command{
	// Use: holds the text used to invoke usage of the command.
	Use: "todo",
	// Short: represents a short description of the command. This is shown in the CLI help output.
	Short: "short description of the command",
	// Long: similar to Short but holds a longer description of the command.
	Long: "long description of the command",
	// Run: holds the function to be executed on the invocation of the command.
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func setupCommands(taskService TaskService) {
	rootCmd.AddCommand(newAddCmd(taskService))
	rootCmd.AddCommand(newListCmd(taskService))
	rootCmd.AddCommand(newUpdateCmd(taskService))
	rootCmd.AddCommand(newApiCmd(taskService))
}

func Execute(taskService TaskService) {
	setupCommands(taskService)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
