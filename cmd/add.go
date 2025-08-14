package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"golang/tutorial/todo/internal/quotes"
	"golang/tutorial/todo/internal/services/task"
	"golang/tutorial/todo/internal/storage"
)

var setDueDate bool
var sayQuote bool

func newAddCmd() *cobra.Command {
	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add new task",
		Long:  "You can add a new task to your todo list with this command.",
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			svc := task.NewService(quotes.NewHTTPClient(os.Getenv("QUOTES_BASE_URL"), 10*time.Second), storage.New("tasks.json"))
			var err error
			switch len(args) {
			case 0:
				fmt.Println("Error: Task name is required")
				return
			case 1:
				err = svc.AddTask(args[0], setDueDate, "")
			case 2:
				err = svc.AddTask(args[0], setDueDate, args[1])
			default:
				fmt.Println("too many arguments")
				return
			}
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Println("task added successfully")

			if sayQuote {
				quote, err := svc.GetRandomQuote()
				if err != nil {
					fmt.Printf("Error fetching quote: %v\n", err)
					return
				}
				fmt.Printf("Here's a quote for you: %s\n", quote)
			}
		},
	}

	addCmd.Flags().BoolVarP(&setDueDate, "due", "d", false, "Set a due date for the task")
	addCmd.Flags().BoolVarP(&sayQuote, "quote", "q", false, "Say a quote after adding the task")

	return addCmd
}
