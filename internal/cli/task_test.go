package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/flrnd/brancher/internal/config"
	"github.com/flrnd/brancher/internal/provider"
	"github.com/flrnd/brancher/internal/task"
)

type fakeTasksProvider struct{}

func (f fakeTasksProvider) Name() provider.Name {
	return provider.Name("fake")
}

func (f fakeTasksProvider) RequiredEnv() []string { return nil }

func (f fakeTasksProvider) ListTasks(ctx context.Context) ([]task.Task, error) {
	return []task.Task{
		{ID: "1", Title: "First task"},
		{ID: "2", Title: "Second task"},
	}, nil
}

func (f fakeTasksProvider) GetTask(ctx context.Context, id string) (task.Task, error) {
	return task.Task{}, nil
}

func TestTasksCommand(t *testing.T) {
	origLoad := loadConfig
	origProvider := newProvider

	defer func() {
		loadConfig = origLoad
		newProvider = origProvider
	}()

	loadConfig = func() (*config.Config, error) {
		return &config.Config{
			Provider: "github",
		}, nil
	}

	newProvider = func(cfg *config.Config) (task.Provider, error) {
		return fakeTasksProvider{}, nil
	}

	cmd := NewTasksCommand()

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()

	if output == "" {
		t.Fatalf("expected output, got empty")
	}

	if !bytes.Contains([]byte(output), []byte("First task")) {
		t.Fatalf("expected task output, got: %s", output)
	}
}
