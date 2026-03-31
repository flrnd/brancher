// Package cli
package cli

import (
	"github.com/flrnd/brancher/internal/git"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "brancher",
		Short: "Create git branches from project tasks",
	}

	cmd.AddCommand(
		NewInitCommand(git.NewRemoteReader),
		NewTasksCommand(),
		NewStartCommand(),
		NewTasksCommand(),
	)

	return cmd
}
