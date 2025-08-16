package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"golang/tutorial/todo/internal/quotes"
	"golang/tutorial/todo/internal/services/quotesvc"
	"golang/tutorial/todo/internal/services/task"
	"golang/tutorial/todo/internal/storage"
)

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

func setupCommands(quoteSvc quotesvc.QuoteService, taskService task.TaskService) {
	rootCmd.AddCommand(newAddCmd(quoteSvc, taskService))
	rootCmd.AddCommand(newListCmd(taskService))
	rootCmd.AddCommand(newUpdateCmd(taskService))
	rootCmd.AddCommand(newAPICmd(quoteSvc, taskService))
}

func Execute() {
	// サービスの初期化
	quoteClient := quotes.NewHTTPClient(os.Getenv("QUOTES_BASE_URL"), 10*time.Second)
	quoteSvc := quotesvc.NewQuoteService(quoteClient)
	storageClient := storage.NewStorage(os.Getenv("STORAGE_FILE_PATH"))
	taskService := task.NewTaskService(storageClient)

	setupCommands(*quoteSvc, *taskService)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
