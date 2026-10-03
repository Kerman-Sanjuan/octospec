package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kerman-sanjuan/octospec/cli/internal/targets"
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

func TestValidPreselectKeepsKnownTools(t *testing.T) {
	got := validPreselect([]string{"opencode", "nope", "claude"})
	want := []string{"opencode", "claude"}
	if len(got) != len(want) {
		t.Fatalf("preselect = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("preselect = %v, want %v", got, want)
		}
	}
}

func TestWithAllTools(t *testing.T) {
	if got := withAllTools(nil); len(got) != len(targets.Targets) {
		t.Fatalf("empty should expand to every target, got %v", got)
	}
	if got := withAllTools([]string{"pi"}); len(got) != 1 || got[0] != "pi" {
		t.Fatalf("a non-empty list should pass through, got %v", got)
	}
}

func TestToolOptionsCoverEveryTarget(t *testing.T) {
	if got := len(toolOptions()); got != len(targets.Targets) {
		t.Fatalf("options = %d, want %d", got, len(targets.Targets))
	}
}

func TestRunRejectsUnknownTool(t *testing.T) {
	if err := Run(Options{Tools: []string{"nope"}, Repo: t.TempDir()}); err == nil {
		t.Fatal("expected an error for an unknown tool")
	}
}

func TestRunInstallsAgents(t *testing.T) {
	dir := t.TempDir()
	if err := Run(Options{Tools: []string{"claude"}, Repo: dir}); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dir, ".claude/agents/spec.md"))
	if err != nil {
		t.Fatalf("missing agent: %v", err)
	}
	if !strings.Contains(string(b), "tools:") {
		t.Fatalf("agent missing the tool surface:\n%s", b)
	}
}

func TestRunLocalWritesNothingUnderHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	repo := t.TempDir()
	if err := Run(Options{Tools: []string{"pi"}, Repo: repo}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(repo, ".pi/prompts/spec.md")); err != nil {
		t.Fatalf("expected a repo-local command: %v", err)
	}
	if _, err := os.Stat(filepath.Join(home, ".pi/agent/prompts/spec.md")); err == nil {
		t.Fatal("a repo-local install must not write under home")
	}
	if _, err := os.Stat(filepath.Join(repo, ".pi/agents/spec.md")); err == nil {
		t.Fatal("pi should have no agent files")
	}
}

func TestRunGlobalWritesHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	repo := t.TempDir()
	if err := Run(Options{Tools: []string{"claude"}, Repo: repo, Global: true}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(home, ".claude/commands/spec.md")); err != nil {
		t.Fatalf("expected a global command: %v", err)
	}
	if _, err := os.Stat(filepath.Join(repo, ".claude/commands/spec.md")); err == nil {
		t.Fatal("a global install should not write the command into the repo")
	}
}

func TestRunDryRunWritesNothing(t *testing.T) {
	repo := t.TempDir()
	if err := Run(Options{Tools: []string{"claude"}, Repo: repo, DryRun: true}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(repo, ".claude")); err == nil {
		t.Fatal("a dry run should write nothing")
	}
	if _, err := os.Stat(filepath.Join(repo, ".octospec")); err == nil {
		t.Fatal("a dry run should not write the state")
	}
}
