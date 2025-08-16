package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"golang/tutorial/todo/internal/services/quotesvc"
	"golang/tutorial/todo/internal/services/task"
)

var setDueDate bool
var sayQuote bool

func newAddCmd(quoteSvc quotesvc.QuoteService, svc task.TaskService) *cobra.Command {
	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add new task",
		Long:  "You can add a new task to your todo list with this command.",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			var dueDate string
			if len(args) == 1 {
				dueDate = ""
			} else {
				dueDate = args[1]
			}
			err := svc.AddTask(args[0], dueDate)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Println("task added successfully")

			if sayQuote {
				quote, err := quoteSvc.GetRandomQuote()
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
