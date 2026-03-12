package task

import (
	"context"

	"github.com/flrnd/brancher/internal/provider"
)

type Provider interface {
	Name() provider.Name
	RequiredEnv() []string
	ListTasks(ctx context.Context) ([]Task, error)
	GetTask(ctx context.Context, id string) (Task, error)
}
