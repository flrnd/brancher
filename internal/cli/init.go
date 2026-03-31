package cli

import (
	"fmt"

	"github.com/flrnd/brancher/internal/cli/input"
	"github.com/flrnd/brancher/internal/cli/output"
	"github.com/flrnd/brancher/internal/config"
	"github.com/flrnd/brancher/internal/git"
	"github.com/flrnd/brancher/internal/provider"
	"github.com/spf13/cobra"
)

func NewInitCommand(
	newDriver func() (git.RemoteReader, error),
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize brancher configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			output.Println(cmd, "Initializing brancher")

			driver, err := newDriver()
			if err != nil {
				return fmt.Errorf("not a git repository")
			}

			output.Println(cmd, "Git repository detected")

			if _, err := config.Load(); err == nil {
				return fmt.Errorf("brancher is already initialized")
			}

			output.Println(cmd, "Select provider:")
			output.Println(cmd, "1) github")

			choice := input.Ask(cmd, ">")

			var prov provider.Name

			switch choice {
			case "", "1":
				prov = provider.GitHub
			default:
				return fmt.Errorf("invalid provider selection")
			}

			var detectedOwner string
			var detectedRepo string

			if url, err := driver.OriginURL(); err == nil {
				if o, r, err := git.ParseRemote(url); err == nil {
					detectedOwner = o
					detectedRepo = r
				}
			}

			if detectedOwner == "" || detectedRepo == "" {
				return fmt.Errorf("owner and repo are required")
			}

			ownerPrompt := "Repository owner:"
			if detectedOwner != "" {
				ownerPrompt = fmt.Sprintf("Repository owner [%s]:", detectedOwner)
			}

			owner := input.Ask(cmd, ownerPrompt)
			if owner == "" {
				owner = detectedOwner
			}

			repoPrompt := "Repository:"
			if detectedRepo != "" {
				repoPrompt = fmt.Sprintf("Repository [%s]:", detectedRepo)
			}

			repo := input.Ask(cmd, repoPrompt)
			if repo == "" {
				repo = detectedRepo
			}

			cfg := &config.Config{
				Provider: prov,
				Project: config.ProjectConfig{
					Owner: owner,
					Repo:  repo,
				},
				Branch: config.BranchConfig{
					Strategy: "title",
				},
			}

			if err := cfg.Save(); err != nil {
				return err
			}

			output.Println(cmd, "Brancher initialized successfully")

			return nil
		},
	}

	return cmd
}
