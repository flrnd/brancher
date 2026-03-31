package cli

import (
	"fmt"

	"github.com/flrnd/brancher/internal/cli/output"
	"github.com/spf13/cobra"
)

func NewTasksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tasks",
		Short: "List available tasks from the provider",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return fmt.Errorf("repository not initialized. Run 'brancher init'")
			}

			provider, err := newProvider(cfg)
			if err != nil {
				return err
			}

			tasks, err := provider.ListTasks(cmd.Context())
			if err != nil {
				return err
			}

			for _, t := range tasks {
				output.Task(cmd, t.ID, t.Title)
			}

			return nil
		},
	}

	return cmd
}
