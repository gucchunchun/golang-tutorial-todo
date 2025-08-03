package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new task",
	Long:  "You can add a new task to your todo list with this command.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", Add(args[0]))
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
