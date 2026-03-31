package cli

import (
	"fmt"

	"github.com/flrnd/brancher/internal/branch"
	"github.com/flrnd/brancher/internal/cli/output"
	"github.com/flrnd/brancher/internal/config"
	"github.com/flrnd/brancher/internal/git"
	"github.com/flrnd/brancher/internal/task"
	"github.com/spf13/cobra"
)

var loadConfig = func() (*config.Config, error) {
	return config.Load()
}

var newProvider = func(cfg *config.Config) (task.Provider, error) {
	return task.NewProvider(cfg)
}

var newGitDriver = func() (git.Driver, error) {
	return git.NewDriver()
}

func NewStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start <task-id>",
		Short: "Create a branch from a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskRef := args[0]

			cfg, err := loadConfig()
			if err != nil {
				return fmt.Errorf("repository not initialized. Run 'brancher init'")
			}

			provider, err := newProvider(cfg)
			if err != nil {
				return err
			}

			t, err := provider.GetTask(cmd.Context(), taskRef)
			if err != nil {
				return err
			}

			strategy, err := branch.ResolveStrategy(branch.Strategy(cfg.Branch.Strategy))
			if err != nil {
				return err
			}

			generator := branch.NewGenerator(strategy)
			branchName := generator.Generate(t)

			driver, err := newGitDriver()
			if err != nil {
				return err
			}

			if err := driver.CreateBranch(branchName); err != nil {
				return err
			}

			output.Task(cmd, t.ID, t.Title)
			output.BranchCreated(cmd, branchName)

			return nil
		},
	}

	return cmd
}
