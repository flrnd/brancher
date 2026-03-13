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

	want := "something-does-not-work"

	if got != want {
		t.Fatalf("Generate() = %q, want %q", got, want)
	}
}
