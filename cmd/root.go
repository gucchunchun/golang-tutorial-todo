package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"

	"golang/tutorial/todo/internal/config"
	"golang/tutorial/todo/internal/db"
	"golang/tutorial/todo/internal/quotes"
	"golang/tutorial/todo/internal/services/quotesvc"
	"golang/tutorial/todo/internal/services/tasksvc"
	"golang/tutorial/todo/internal/storage/mysql"
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

func setupCommands(quoteSvc quotesvc.QuoteService, taskService tasksvc.TaskService) {
	rootCmd.AddCommand(newAddCmd(quoteSvc, taskService))
	rootCmd.AddCommand(newListCmd(taskService))
	rootCmd.AddCommand(newUpdateCmd(taskService))
	rootCmd.AddCommand(newAPICmd(quoteSvc, taskService))
}

func Execute() {
	// DBの初期化
	cfg := config.NewDBConfig()
	dbc := &db.Client{}
	if err := dbc.Connect(cfg); err != nil {
		log.Fatalf("DB connect failed: %v", err)
	}
	defer dbc.Close()

	quoteClient := quotes.NewHTTPClient(os.Getenv("QUOTES_BASE_URL"), 10*time.Second)
	quoteSvc := quotesvc.NewQuoteService(quoteClient)
	taskRepo := mysql.NewTaskRepo(dbc.DB())

	// サービスの初期化
	taskService := tasksvc.NewTaskService(taskRepo)

	setupCommands(*quoteSvc, *taskService)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
