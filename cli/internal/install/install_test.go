package install

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunInstallsFiles(t *testing.T) {
	dir := t.TempDir()
	if err := Run(Options{Tools: []string{"copilot", "claude"}, Repo: dir}); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{
		".github/prompts/spec.prompt.md",
		".claude/commands/spec.md",
		".octospec/octospec.json",
	} {
		if _, err := os.Stat(filepath.Join(dir, p)); err != nil {
			t.Errorf("missing %s: %v", p, err)
		}
	}
}

func TestRunRejectsUnknownTool(t *testing.T) {
	if err := Run(Options{Tools: []string{"nope"}, Repo: t.TempDir()}); err == nil {
		t.Fatal("expected an error for an unknown tool")
	}
}
