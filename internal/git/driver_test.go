package git

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	gitpkg "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

func TestCreateAndCheckoutBranch(t *testing.T) {
	driver := newTestDriver(t)

	if err := driver.CreateAndCheckoutBranch("feature-x"); err != nil {
		t.Fatalf("CreateAndCheckoutBranch returned error: %v", err)
	}

	current, err := driver.CurrentBranch()
	if err != nil {
		t.Fatalf("CurrentBranch returned error: %v", err)
	}

	wantRef := plumbing.NewBranchReferenceName("feature-x")
	if current.RefName != wantRef {
		t.Fatalf("expected current branch %s, got %s", wantRef, current.RefName)
	}

	head, err := driver.repo.Head()
	if err != nil {
		t.Fatalf("Head returned error: %v", err)
	}

	if head.Name() != wantRef {
		t.Fatalf("expected HEAD to point to %s, got %s", wantRef, head.Name())
	}

	if _, err := driver.repo.Reference(wantRef, true); err != nil {
		t.Fatalf("expected branch reference to exist: %v", err)
	}
}

func TestCreateAndCheckoutBranchFailsWhenBranchExists(t *testing.T) {
	driver := newTestDriver(t)

	if err := driver.CreateAndCheckoutBranch("feature-x"); err != nil {
		t.Fatalf("CreateAndCheckoutBranch returned error: %v", err)
	}

	if err := driver.CreateAndCheckoutBranch("feature-x"); err == nil {
		t.Fatal("expected error when branch already exists")
	}
}

func TestCreateAndCheckoutBranchKeepsUnstagedChanges(t *testing.T) {
	driver := newTestDriver(t)

	worktree, err := driver.repo.Worktree()
	if err != nil {
		t.Fatalf("Worktree returned error: %v", err)
	}

	file, err := worktree.Filesystem.OpenFile("README.md", os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		t.Fatalf("OpenFile returned error: %v", err)
	}

	if _, err := file.Write([]byte("seed\nlocal change\n")); err != nil {
		_ = file.Close()
		t.Fatalf("Write returned error: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}

	if err := driver.CreateAndCheckoutBranch("feature-x"); err != nil {
		t.Fatalf("CreateAndCheckoutBranch returned error: %v", err)
	}

	current, err := driver.CurrentBranch()
	if err != nil {
		t.Fatalf("CurrentBranch returned error: %v", err)
	}

	wantRef := plumbing.NewBranchReferenceName("feature-x")
	if current.RefName != wantRef {
		t.Fatalf("expected current branch %s, got %s", wantRef, current.RefName)
	}

	reader, err := worktree.Filesystem.Open("README.md")
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll returned error: %v", err)
	}

	if string(data) != "seed\nlocal change\n" {
		t.Fatalf("expected unstaged changes to be preserved, got %q", string(data))
	}
}

func newTestDriver(t *testing.T) *GoGitDriver {
	t.Helper()

	repoPath := t.TempDir()

	repo, err := gitpkg.PlainInit(repoPath, false)
	if err != nil {
		t.Fatalf("PlainInit returned error: %v", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		t.Fatalf("Worktree returned error: %v", err)
	}

	readmePath := filepath.Join(repoPath, "README.md")
	if err := os.WriteFile(readmePath, []byte("seed\n"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	if _, err := worktree.Add("README.md"); err != nil {
		t.Fatalf("Add returned error: %v", err)
	}

	if _, err := worktree.Commit("initial commit", &gitpkg.CommitOptions{
		Author: &object.Signature{
			Name:  "Brancher",
			Email: "brancher@example.com",
		},
	}); err != nil {
		t.Fatalf("Commit returned error: %v", err)
	}

	return &GoGitDriver{repo: repo}
}
