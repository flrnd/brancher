package cli

import (
	"os"
	"strings"
	"testing"

	"github.com/flrnd/brancher/internal/cli/input"
	"github.com/flrnd/brancher/internal/git"
	"github.com/spf13/cobra"
)

func TestInitCommand(t *testing.T) {
	// isolate filesystem
	tmpDir := t.TempDir()

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get wd: %v", err)
	}

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	defer func() {
		if err := os.Chdir(oldWd); err != nil {
			t.Fatalf("failed to restore wd: %v", err)
		}
	}()

	// --- mock input (use askFn, not Ask) ---
	inputs := []string{
		"", // provider → default github
		"", // owner → detected
		"", // repo → detected
	}

	i := 0
	originalAsk := input.AskFn
	input.AskFn = func(cmd *cobra.Command, q string) string {
		val := inputs[i]
		i++
		return val
	}
	defer func() { input.AskFn = originalAsk }()

	cmd := NewInitCommand(func() (git.RemoteReader, error) {
		return fakeInitDriver{}, nil
	})

	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(".brancher/config.yml"); err != nil {
		t.Fatalf("config file not created")
	}

	data, err := os.ReadFile(".brancher/config.yml")
	if err != nil {
		t.Fatalf("failed reading config: %v", err)
	}

	content := string(data)

	if !strings.Contains(content, "myorg") {
		t.Fatalf("expected owner in config, got: %s", content)
	}

	if !strings.Contains(content, "myrepo") {
		t.Fatalf("expected repo in config, got: %s", content)
	}
}

type fakeInitDriver struct{}

func (f fakeInitDriver) OriginURL() (string, error) {
	return "git@github.com:myorg/myrepo.git", nil
}
