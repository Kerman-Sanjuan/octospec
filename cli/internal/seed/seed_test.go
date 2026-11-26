package seed

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnsureGitignoreAppendsOnce(t *testing.T) {
	repo := t.TempDir()
	if err := ensureGitignore(repo); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(repo, ".gitignore")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), ".octospec/locks/") {
		t.Fatalf("the lock directory is not ignored: %q", b)
	}
	if err := ensureGitignore(repo); err != nil {
		t.Fatal(err)
	}
	b, err = os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(string(b), ".octospec/locks/"); n != 1 {
		t.Fatalf("the block was appended %d times: %q", n, b)
	}
}

func TestEnsureGitignoreKeepsExistingContent(t *testing.T) {
	repo := t.TempDir()
	path := filepath.Join(repo, ".gitignore")
	if err := os.WriteFile(path, []byte("node_modules"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ensureGitignore(repo); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(b), "node_modules") {
		t.Fatalf("the existing content was rewritten: %q", b)
	}
	if !strings.Contains(string(b), ".octospec/locks/") {
		t.Fatalf("the block was not appended: %q", b)
	}
}
