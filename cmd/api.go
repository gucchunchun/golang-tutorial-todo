package cmd

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"golang/tutorial/todo/internal/api"
	"golang/tutorial/todo/internal/services/quotesvc"
	"golang/tutorial/todo/internal/services/task"
)

func runAPI(ctx context.Context, addr string, svc task.TaskService) error {
	api := api.New(&svc)
	handler := api.Routes()

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	log.Printf("listening on %s", addr)
	go func() {
		// Ctrl-C / cancel()が呼ばれた時にサーバーをシャットダウンする
		<-ctx.Done()
		// NOTE: 5秒のタイムアウト(既存接続の完了のため）を設定してシャットダウン
		//       もし5秒以内にdeferが呼ばれなかった場合、強制的にserver.Shutdownが終了される
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown error: %v", err)
		}
	}()
	return http.ListenAndServe(addr, handler)
}

func newAPICmd(quoteSvc quotesvc.QuoteService, svc task.TaskService) *cobra.Command {
	var apiCmd = &cobra.Command{
		Use:   "api",
		Short: "Run the HTTP API server",
		Long:  "This command starts the HTTP API server for the todo application.",
		Run: func(cmd *cobra.Command, args []string) {
			addr, _ := cmd.Flags().GetString("addr")
			log.Printf("Starting API server on %s", addr)
			if err := runAPI(cmd.Context(), addr, svc); err != nil {
				log.Fatalf("Failed to start API server: %v", err)
			}
		},
	}

	apiCmd.Flags().StringP("addr", "a", ":8080", "Address to listen on")

	return apiCmd
}
