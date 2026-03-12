// Package providers
package providers

import (
	"context"

	"github.com/flrnd/brancher/internal/task"
)

type GitHubProvider struct{}

func (p GitHubProvider) Name() string {
	return "github"
}

func (p GitHubProvider) RequiredEnv() []string {
	return []string{"GITHUB_TOKEN"}
}

func (p GitHubProvider) ListTasks(ctx context.Context) ([]task.Task, error) {
	return nil, nil
}

func (p GitHubProvider) GetTask(ctx context.Context, id string) (task.Task, error) {
	return task.Task{}, nil
}
