package task

import "context"

type Provider interface {
	Name() string
	RequiredEnv() []string

	ListTasks(ctx context.Context) ([]Task, error)
	GetTask(ctx context.Context, id string) (Task, error)
}
