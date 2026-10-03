package install

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/charmbracelet/huh"
	"github.com/kerman-sanjuan/octospec/cli/internal/config"
	"github.com/kerman-sanjuan/octospec/cli/internal/targets"
)

// lineReader yields one scripted answer per Read, in order, then EOF. huh's
// accessible mode builds a fresh scanner per prompt, so a plain strings.Reader
// would let the first prompt swallow every line.
type lineReader struct {
	lines []string
	i     int
}

func (r *lineReader) Read(p []byte) (int, error) {
	if r.i >= len(r.lines) {
		return 0, io.EOF
	}
	n := copy(p, r.lines[r.i]+"\n")
	r.i++
	return n, nil
}

// driveAccessible makes the next form run read the given scripted lines in
// huh's accessible mode and discards the prompts.
func driveAccessible(t *testing.T, answers ...string) {
	t.Helper()
	old := runForm
	runForm = func(f *huh.Form) error {
		var out bytes.Buffer
		return f.WithAccessible(true).WithOutput(&out).WithInput(&lineReader{lines: answers}).Run()
	}
	t.Cleanup(func() { runForm = old })
}

// optionIndex returns the 1-based position of a tool in the selector, given the
// order targets declares them.
func optionIndex(name string) string {
	for i, tp := range targets.Targets {
		if tp.Name == name {
			return strconv.Itoa(i + 1)
		}
	}
	return ""
}

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

func TestRunInteractiveSelectsSubset(t *testing.T) {
	repo := t.TempDir()
	// Answer the multi-select with opencode then claude, then confirm (0); the
	// scope confirm defaults to repo-local with a blank line.
	answers := []string{optionIndex("opencode"), optionIndex("claude"), "0", ""}
	driveAccessible(t, answers...)
	if err := Run(Options{Repo: repo, Interactive: true}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(repo, ".opencode/command/spec.md")); err != nil {
		t.Fatalf("opencode was not installed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(repo, ".claude/commands/spec.md")); err != nil {
		t.Fatalf("claude was not installed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(repo, ".pi/prompts/spec.md")); err == nil {
		t.Fatal("pi was not selected but was installed")
	}
	cfg, err := config.Load(repo)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(cfg.Tools, ",") != "opencode,claude" {
		t.Fatalf("recorded tools = %v, want [opencode claude]", cfg.Tools)
	}
}

func TestRunInteractiveNoneInstallsAll(t *testing.T) {
	repo := t.TempDir()
	// Confirm immediately (0) with nothing toggled; the scope confirm is blank.
	driveAccessible(t, "0", "")
	if err := Run(Options{Repo: repo, Interactive: true}); err != nil {
		t.Fatal(err)
	}
	for _, tp := range targets.Targets {
		if _, err := os.Stat(filepath.Join(repo, tp.Dir, tp.Filename("spec"))); err != nil {
			t.Errorf("tool %s was not installed: %v", tp.Name, err)
		}
	}
}

func TestRunInteractivePreselectsPreviousTools(t *testing.T) {
	repo := t.TempDir()
	// A prior install records claude only.
	if err := config.Save(repo, config.Config{
		Tools:  []string{"claude"},
		Scopes: map[string]string{},
		Models: map[string]string{},
		Files:  map[string]string{},
	}); err != nil {
		t.Fatal(err)
	}
	// Confirm with nothing toggled. The preselected claude must survive.
	driveAccessible(t, "0", "")
	if err := Run(Options{Repo: repo, Interactive: true}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(repo, ".claude/commands/spec.md")); err != nil {
		t.Fatalf("the preselected claude was not installed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(repo, ".pi/prompts/spec.md")); err == nil {
		t.Fatal("pi was not preselected but was installed")
	}
}

func TestRunInteractiveGlobalFromConfirm(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	repo := t.TempDir()
	// Select claude by number, confirm, then answer the scope prompt "y".
	answers := []string{optionIndex("claude"), "0", "y"}
	driveAccessible(t, answers...)
	if err := Run(Options{Repo: repo, Interactive: true}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(home, ".claude/commands/spec.md")); err != nil {
		t.Fatalf("expected a global install under home: %v", err)
	}
}

func TestRunNonInteractiveNoForm(t *testing.T) {
	// A non-interactive run must not call the form runner at all.
	called := false
	old := runForm
	runForm = func(f *huh.Form) error { called = true; return nil }
	t.Cleanup(func() { runForm = old })
	if err := Run(Options{Tools: []string{"claude"}, Repo: t.TempDir()}); err != nil {
		t.Fatal(err)
	}
	if called {
		t.Fatal("a non-interactive install opened a form")
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
