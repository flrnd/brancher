// Package provider
package provider

import (
	"context"
	"os"

	"github.com/flrnd/brancher/internal/config"
	"github.com/flrnd/brancher/internal/provider"
	"github.com/flrnd/brancher/internal/task"
)

type GitHubProvider struct {
	token string
	owner string
	repo  string
}

func init() {
	task.RegisterProvider(task.Definition{
		Name: provider.GitHub,
		Required: []string{
			EnvToken,
		},
		New: New,
	})
}

func New(cfg *config.Config) (task.Provider, error) {
	token := os.Getenv(EnvToken)

	return &GitHubProvider{
		token: token,
		owner: cfg.Project.Owner,
		repo:  cfg.Project.Repo,
	}, nil
}

func (p *GitHubProvider) Name() provider.Name {
	return provider.GitHub
}

func (p *GitHubProvider) RequiredEnv() []string {
	return []string{EnvToken}
}

func (p *GitHubProvider) ListTasks(ctx context.Context) ([]task.Task, error) {
	return nil, nil
}

func (p *GitHubProvider) GetTask(ctx context.Context, id string) (task.Task, error) {
	return task.Task{}, nil
}
