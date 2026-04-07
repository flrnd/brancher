package branch

import (
	"testing"

	"github.com/flrnd/brancher/internal/task"
)

func TestGeneratorGenerate(t *testing.T) {
	strategy, err := ResolveStrategy(StrategyTitle)
	if err != nil {
		t.Fatalf("unexpected error resolving strategy: %v", err)
	}

	g := NewGenerator(strategy)

	task := task.Task{
		ID:    "42",
		Title: "Something does not work!!!",
	}

	got := g.Generate(task)

	want := "42-something-does-not-work"

	if got != want {
		t.Fatalf("Generate() = %q, want %q", got, want)
	}
}

func TestGeneratorGenerateWithPrefixedID(t *testing.T) {
	strategy, err := ResolveStrategy(StrategyTitle)
	if err != nil {
		t.Fatalf("unexpected error resolving strategy: %v", err)
	}

	g := NewGenerator(strategy)

	task := task.Task{
		ID:    "PROJ-123",
		Title: "Implement GitHub task provider",
	}

	got := g.Generate(task)

	want := "proj-123-implement-github-task-provider"

	if got != want {
		t.Fatalf("Generate() = %q, want %q", got, want)
	}
}

func TestGeneratorGenerateFallsBackToTitleWhenIDIsEmpty(t *testing.T) {
	strategy, err := ResolveStrategy(StrategyTitle)
	if err != nil {
		t.Fatalf("unexpected error resolving strategy: %v", err)
	}

	g := NewGenerator(strategy)

	task := task.Task{
		Title: "Implement GitHub task provider",
	}

	got := g.Generate(task)

	want := "implement-github-task-provider"

	if got != want {
		t.Fatalf("Generate() = %q, want %q", got, want)
	}
}

func TestGeneratorGenerateFallsBackToIDWhenTitleIsEmpty(t *testing.T) {
	strategy, err := ResolveStrategy(StrategyTitle)
	if err != nil {
		t.Fatalf("unexpected error resolving strategy: %v", err)
	}

	g := NewGenerator(strategy)

	task := task.Task{
		ID: "PROJ-123",
	}

	got := g.Generate(task)

	want := "proj-123"

	if got != want {
		t.Fatalf("Generate() = %q, want %q", got, want)
	}
}
