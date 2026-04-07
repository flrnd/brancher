package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/flrnd/brancher/internal/config"
	"github.com/flrnd/brancher/internal/git"
	"github.com/flrnd/brancher/internal/provider"
	"github.com/flrnd/brancher/internal/task"
)

type fakeProvider struct{}

func (f fakeProvider) Name() provider.Name {
	return provider.Name("fake")
}

func (f fakeProvider) RequiredEnv() []string { return nil }
func (f fakeProvider) ListTasks(ctx context.Context) ([]task.Task, error) {
	return nil, nil
}

func (f fakeProvider) GetTask(ctx context.Context, id string) (task.Task, error) {
	return task.Task{
		ID:    id,
		Title: "Test task",
	}, nil
}

type fakeDriver struct{}

func (f fakeDriver) CreateBranch(name string) error            { return nil }
func (f fakeDriver) DeleteBranch(name string) error            { return nil }
func (f fakeDriver) ListLocalBranches() ([]git.Branch, error)  { return nil, nil }
func (f fakeDriver) ListRemoteBranches() ([]git.Branch, error) { return nil, nil }
func (f fakeDriver) ListAllBranches() ([]git.Branch, error)    { return nil, nil }
func (f fakeDriver) CurrentBranch() (git.Branch, error)        { return git.Branch{}, nil }

func TestStartCommand(t *testing.T) {
	loadConfig = func() (*config.Config, error) {
		return &config.Config{
			Provider: "github",
			Branch: config.BranchConfig{
				Strategy: "title",
			},
		}, nil
	}

	newProvider = func(cfg *config.Config) (task.Provider, error) {
		return fakeProvider{}, nil
	}

	newGitDriver = func() (git.Driver, error) {
		return fakeDriver{}, nil
	}

	cmd := NewStartCommand()

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	cmd.SetArgs([]string{"42"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()

	if output == "" {
		t.Fatalf("expected output, got empty")
	}

	if !bytes.Contains([]byte(output), []byte("Test task")) {
		t.Fatalf("expected task title in output, got: %s", output)
	}

	if !bytes.Contains([]byte(output), []byte("Created branch")) {
		t.Fatalf("expected branch creation output, got: %s", output)
	}

	if !bytes.Contains([]byte(output), []byte("42-test-task")) {
		t.Fatalf("expected id-prefixed branch name in output, got: %s", output)
	}
}
