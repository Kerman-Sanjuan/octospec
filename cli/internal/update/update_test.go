package update

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kerman-sanjuan/octospec/cli/internal/install"
)

func TestUpdatePreservesLocalEdits(t *testing.T) {
	dir := t.TempDir()
	if err := install.Run(install.Options{Tools: []string{"copilot"}, Repo: dir}); err != nil {
		t.Fatal(err)
	}
	edited := filepath.Join(dir, ".github/prompts/spec.prompt.md")
	if err := os.WriteFile(edited, []byte("hand edited\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Run(dir); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(edited)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "hand edited\n" {
		t.Fatalf("edit not preserved: %q", b)
	}
}

func TestUpdateWithoutInstall(t *testing.T) {
	if err := Run(t.TempDir()); err == nil {
		t.Fatal("expected an error when nothing is installed")
	}
}

func TestUpdatePreservesAgentEdits(t *testing.T) {
	dir := t.TempDir()
	if err := install.Run(install.Options{Tools: []string{"claude"}, Repo: dir}); err != nil {
		t.Fatal(err)
	}
	edited := filepath.Join(dir, ".claude/agents/spec.md")
	if err := os.WriteFile(edited, []byte("hand edited agent\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Run(dir); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(edited)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "hand edited agent\n" {
		t.Fatalf("agent edit not preserved: %q", b)
	}
}
