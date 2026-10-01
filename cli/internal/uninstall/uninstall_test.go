package uninstall

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kerman-sanjuan/octospec/cli/internal/install"
)

func TestUninstallRemovesManagedAndKeepsEdits(t *testing.T) {
	dir := t.TempDir()
	if err := install.Run(install.Options{Tools: []string{"claude"}, Repo: dir}); err != nil {
		t.Fatal(err)
	}
	edited := filepath.Join(dir, ".claude/commands/spec.md")
	if err := os.WriteFile(edited, []byte("hand edited\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Run(dir); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(edited)
	if err != nil {
		t.Fatalf("an edited file should be kept: %v", err)
	}
	if string(b) != "hand edited\n" {
		t.Fatalf("edited content = %q", b)
	}
	if _, err := os.Stat(filepath.Join(dir, ".claude/commands/apply.md")); err == nil {
		t.Fatal("a managed file should be removed")
	}
	if _, err := os.Stat(filepath.Join(dir, ".octospec/octospec.json")); err == nil {
		t.Fatal("the state file should be removed")
	}
}

func TestUninstallWithoutState(t *testing.T) {
	if err := Run(t.TempDir()); err != nil {
		t.Fatal(err)
	}
}
