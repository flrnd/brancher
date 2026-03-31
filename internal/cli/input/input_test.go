package input

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestAsk(t *testing.T) {
	original := AskFn
	defer func() { AskFn = original }()

	AskFn = func(cmd *cobra.Command, q string) string {
		return "hello"
	}

	result := Ask(nil, "question")

	if result != "hello" {
		t.Fatalf("expected hello, got %s", result)
	}
}
