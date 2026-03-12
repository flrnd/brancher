package cli

import "github.com/spf13/cobra"

func NewStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start <task-id>",
		Short: "Create a branch from a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
