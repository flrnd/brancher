// Package cli
package cli

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "brancher",
		Short: "Create git branches from project tasks",
	}

	cmd.AddCommand(
		NewInitCommand(),
		NewTasksCommand(),
		NewStartCommand(),
		NewTasksCommand(),
	)

	return cmd
}
