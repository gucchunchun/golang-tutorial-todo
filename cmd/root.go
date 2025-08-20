package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"golang/tutorial/todo/internal/config"
	"golang/tutorial/todo/internal/db"

	"golang/tutorial/todo/internal/quote"
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

func setupCommands(log *zerolog.Logger, quoteSvc quotesvc.Service, taskSvc tasksvc.TaskService) {
	rootCmd.AddCommand(newAddCmd(quoteSvc, taskSvc))
	rootCmd.AddCommand(newListCmd(taskSvc))
	rootCmd.AddCommand(newUpdateCmd(taskSvc))
	rootCmd.AddCommand(newAPICmd(log, quoteSvc, taskSvc))
}

func Execute(ctx context.Context) error {
	// コンテキストからloggerを取得
	lg := zerolog.Ctx(ctx)

	// DBの初期化
	cfg := config.NewDBConfig()
	dbc := &db.Client{}
	if err := dbc.Connect(cfg); err != nil {
		return fmt.Errorf("DB connect failed: %v", err)
	}
	lg.Info().Msg("DB connected")
	defer dbc.Close()

	quoteClient := quote.New(os.Getenv("QUOTES_BASE_URL"), 10*time.Second)
	quoteSvc := quotesvc.New(quoteClient)
	taskRepo := mysql.NewTaskRepo(dbc.DB())

	// サービスの初期化
	taskService := tasksvc.NewTaskService(taskRepo)

	setupCommands(lg, *quoteSvc, *taskService)

	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("cobra command failed: %v", err)
	}

	return nil
}
