// Package github
package github

import (
	"context"
	"fmt"

	"github.com/flrnd/brancher/internal/config"
	"github.com/flrnd/brancher/internal/env"
	"github.com/flrnd/brancher/internal/provider"
	"github.com/flrnd/brancher/internal/task"
)

type GitHubProvider struct {
	token  string
	owner  string
	repo   string
	client *client
}

func init() {
	task.RegisterProvider(task.Definition{
		Name: provider.GitHub,
		Required: []string{
			env.ProviderToken(provider.GitHub),
		},
		New: New,
	})
}

func New(cfg *config.Config) (task.Provider, error) {
	token := env.Get(env.ProviderToken(provider.GitHub))
	c := newClient(token, cfg.Project.Owner, cfg.Project.Repo)

	return &GitHubProvider{
		token:  token,
		owner:  cfg.Project.Owner,
		repo:   cfg.Project.Repo,
		client: c,
	}, nil
}

func (p *GitHubProvider) Name() provider.Name {
	return provider.GitHub
}

func (p *GitHubProvider) RequiredEnv() []string {
	return []string{env.ProviderToken(provider.GitHub)}
}

func (p *GitHubProvider) ListTasks(ctx context.Context) ([]task.Task, error) {
	issues, err := p.client.listIssues(ctx)
	if err != nil {
		return nil, err
	}

	tasks := make([]task.Task, 0, len(issues))
	for _, it := range issues {
		tasks = append(tasks, mapIssueToTask(it))
	}

	return tasks, nil
}

func (p *GitHubProvider) GetTask(ctx context.Context, id string) (task.Task, error) {
	if id == "" {
		return task.Task{}, fmt.Errorf("task id is required")
	}

	it, err := p.client.getIssue(ctx, id)
	if err != nil {
		return task.Task{}, err
	}

	return mapIssueToTask(it), nil
}

func mapIssueToTask(it issue) task.Task {
	labels := make([]string, 0, len(it.Labels))
	for _, l := range it.Labels {
		labels = append(labels, l.Name)
	}

	return task.Task{
		ID:     fmt.Sprintf("%d", it.Number),
		Title:  it.Title,
		State:  mapState(it.State),
		Labels: labels,
		URL:    it.HTMLURL,
	}
}

func mapState(state string) task.State {
	switch state {
	case "closed":
		return task.StateClosed
	default:
		return task.StateOpen
	}
}
