package cmd

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"golang/tutorial/todo/internal/api"
	"golang/tutorial/todo/internal/api/handlers/quotehdl"
	"golang/tutorial/todo/internal/api/handlers/taskhdl"
)

func runAPI(ctx context.Context, addr string, log *zerolog.Logger, quoteSvc quotehdl.QuoteService, taskSvc taskhdl.TaskService) error {
	api := api.New(log, quoteSvc, taskSvc)
	handler := api.Routes()

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	log.Info().Msgf("listening on %s", addr)

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

func newAPICmd(log *zerolog.Logger, quoteSvc quotehdl.QuoteService, taskSvc taskhdl.TaskService) *cobra.Command {
	var apiCmd = &cobra.Command{
		Use:   "api",
		Short: "Run the HTTP API server",
		Long:  "This command starts the HTTP API server for the todo application.",
		Run: func(cmd *cobra.Command, args []string) {
			addr, _ := cmd.Flags().GetString("addr")
			log.Info().Msgf("Starting API server on %s", addr)
			if err := runAPI(cmd.Context(), addr, log, quoteSvc, taskSvc); err != nil {
				log.Fatal().Msgf("Failed to start API server: %v", err)
			}
		},
	}

	apiCmd.Flags().StringP("addr", "a", ":8080", "Address to listen on")

	return apiCmd
}
