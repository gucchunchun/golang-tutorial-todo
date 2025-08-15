package cmd

import (
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"golang/tutorial/todo/internal/httpapi"
)

func RunAPI(addr string, svc TaskService) error {
	srv := httpapi.New(svc)
	log.Printf("listening on %s", addr)
	return http.ListenAndServe(addr, srv.Routes())
}

func newApiCmd(svc TaskService) *cobra.Command {
	var apiCmd = &cobra.Command{
		Use:   "api",
		Short: "Run the HTTP API server",
		Long:  "This command starts the HTTP API server for the todo application.",
		Run: func(cmd *cobra.Command, args []string) {
			addr, _ := cmd.Flags().GetString("addr")
			log.Printf("Starting API server on %s", addr)
			if err := RunAPI(addr, svc); err != nil {
				log.Fatalf("Failed to start API server: %v", err)
			}
		},
	}

	apiCmd.Flags().StringP("addr", "a", ":8080", "Address to listen on")

	return apiCmd
}
