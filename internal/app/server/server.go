package server

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"

	"golang/tutorial/todo/internal/api"
	"golang/tutorial/todo/internal/app/bootstrap"
	"golang/tutorial/todo/internal/quote"
	"golang/tutorial/todo/internal/services/quotesvc"
	"golang/tutorial/todo/internal/services/tasksvc"
	"golang/tutorial/todo/internal/storage/mysql"
)

func RunHTTPServer(ctx context.Context) error {
	deps, err := bootstrap.Init(ctx, bootstrap.Options{
		Service:         "todo-server",
		LogFile:         "",
		Level:           zerolog.DebugLevel,
		ConsoleWarnOnly: false,
		UseGlobal:       true,
	})
	if err != nil {
		return err
	}
	defer deps.Cleanup()

	// サービスの初期化
	quoteClient := quote.New(os.Getenv("QUOTES_BASE_URL"), 10*time.Second)
	quoteSvc := quotesvc.New(quoteClient)
	taskRepo := mysql.NewTaskRepo(deps.DB.DB())
	taskService := tasksvc.NewTaskService(taskRepo)

	api := api.New(deps.Log, quoteSvc, taskService)
	handler := api.Routes()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	return srv.ListenAndServe()
}
