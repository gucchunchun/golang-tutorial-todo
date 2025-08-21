package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"golang/tutorial/todo/internal/adapters/convert"
	"golang/tutorial/todo/internal/services/quotesvc"
	"golang/tutorial/todo/internal/services/tasksvc"
)

var setDueDate bool
var sayQuote bool

func newAddCmd(quoteSvc quotesvc.Service, svc tasksvc.TaskService) *cobra.Command {
	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add new task",
		Long:  "You can add a new task to your todo list with this command.",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: configで設定する
			loc, _ := time.LoadLocation("Asia/Tokyo")

			dueDate := ""
			if len(args) != 1 {
				dueDate = args[1]
			}

			params := convert.CreateParams{
				Name:  args[0],
				DueAt: &dueDate,
			}

			c, err := convert.FromCreateInput(params, loc)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			_, err = svc.AddTask(c)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Println("task added successfully")

			if sayQuote {
				quote, err := quoteSvc.GetRandomQuote(context.Background())
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
