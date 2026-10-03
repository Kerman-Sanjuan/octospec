package models

import (
	"bytes"
	"io"
	"strconv"
	"testing"

	"github.com/charmbracelet/huh"
	"github.com/kerman-sanjuan/octospec/cli/internal/config"
)

// lineReader yields one scripted answer per Read, in order, then EOF. huh's
// accessible mode builds a fresh scanner for each prompt, so a plain
// strings.Reader would let the first prompt swallow every line.
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
// huh's accessible mode and discards the prompts. It restores the real runner
// when the test ends.
func driveAccessible(t *testing.T, answers ...string) {
	t.Helper()
	old := runForm
	runForm = func(f *huh.Form) error {
		var out bytes.Buffer
		return f.WithAccessible(true).WithOutput(&out).WithInput(&lineReader{lines: answers}).Run()
	}
	t.Cleanup(func() { runForm = old })
}

func TestSetWritesModels(t *testing.T) {
	dir := t.TempDir()
	if err := Set(Options{Repo: dir, Pairs: map[string]string{config.RoleReviewer: "opus"}}); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Models[config.RoleReviewer] != "opus" {
		t.Fatalf("models = %v", cfg.Models)
	}
}

func TestSetRejectsUnknownRole(t *testing.T) {
	if err := Set(Options{Repo: t.TempDir(), Pairs: map[string]string{"nope": "x"}}); err == nil {
		t.Fatal("expected an error for an unknown role")
	}
}

func TestSetInteractiveWritesChosenModels(t *testing.T) {
	dir := t.TempDir()
	// A prior install of claude makes the catalog real: sonnet, opus, haiku.
	if err := config.Save(dir, config.Config{
		Tools:  []string{"claude"},
		Scopes: map[string]string{},
		Models: map[string]string{},
		Files:  map[string]string{},
	}); err != nil {
		t.Fatal(err)
	}
	catalog := Catalog([]string{"claude"})
	answers := concat(
		answersForRole("sonnet", catalog), // thinking
		answersForRole("haiku", catalog),  // implementer
		answersForRole("opus", catalog),   // reviewer
	)
	driveAccessible(t, answers...)
	if err := Set(Options{Repo: dir, Interactive: true}); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		config.RoleThinking:    "sonnet",
		config.RoleImplementer: "haiku",
		config.RoleReviewer:    "opus",
	}
	for role, model := range want {
		if cfg.Models[role] != model {
			t.Errorf("role %s = %q, want %q", role, cfg.Models[role], model)
		}
	}
}

func TestSetPairAndInteractiveCombine(t *testing.T) {
	// A --set pair is applied first, then the form runs. The form's seeded
	// choice should reflect the pair, so accepting blanks keeps it.
	dir := t.TempDir()
	if err := config.Save(dir, config.Config{
		Tools:  []string{"claude"},
		Scopes: map[string]string{},
		Models: map[string]string{},
		Files:  map[string]string{},
	}); err != nil {
		t.Fatal(err)
	}
	// thinking seeded to "opus" (known), implementer and reviewer default.
	driveAccessible(t, "", "", "", "", "", "")
	if err := Set(Options{
		Repo:        dir,
		Pairs:       map[string]string{config.RoleThinking: "opus"},
		Interactive: true,
	}); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Models[config.RoleThinking] != "opus" {
		t.Fatalf("thinking = %q, want opus preserved through the form", cfg.Models[config.RoleThinking])
	}
}

func TestSeedChoice(t *testing.T) {
	catalog := []string{"sonnet", "opus"}

	if choice, custom := seedChoice("", catalog); choice != "" || custom != "" {
		t.Fatalf("empty model should select the default, got %q %q", choice, custom)
	}
	if choice, custom := seedChoice("sonnet", catalog); choice != "sonnet" || custom != "" {
		t.Fatalf("a known model should select itself, got %q %q", choice, custom)
	}
	if choice, custom := seedChoice("gpt-4o", catalog); choice != customChoice || custom != "gpt-4o" {
		t.Fatalf("an unknown model should seed the custom input, got %q %q", choice, custom)
	}
}

// optionIndex returns the 1-based position huh's accessible select expects for
// a value, given the option order form builds: (tool default), the catalog,
// then Custom...
func optionIndex(model string, catalog []string) int {
	if model == "" {
		return 1
	}
	for i, m := range catalog {
		if m == model {
			return i + 2
		}
	}
	return len(catalog) + 2
}

// answersForRole renders the two accessible answers a role consumes: the
// select choice and the custom input. The free-form input runs even when the
// select is not "Custom..." because huh does not honor hide functions in
// accessible mode, so a blank line keeps it inert.
func answersForRole(model string, catalog []string) []string {
	return []string{itoa(optionIndex(model, catalog)), ""}
}

func itoa(n int) string {
	return strconv.Itoa(n)
}

func concat(parts ...[]string) []string {
	var out []string
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

func TestFormChoosesFromCatalog(t *testing.T) {
	catalog := Catalog([]string{"claude"})
	answers := concat(
		answersForRole("sonnet", catalog), // thinking
		answersForRole("opus", catalog),   // implementer
		answersForRole("", catalog),       // reviewer: tool default
	)
	driveAccessible(t, answers...)
	cfg := config.Config{Tools: []string{"claude"}, Models: map[string]string{}}
	if err := form(&cfg); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		config.RoleThinking:    "sonnet",
		config.RoleImplementer: "opus",
		config.RoleReviewer:    "",
	}
	for role, model := range want {
		if cfg.Models[role] != model {
			t.Errorf("role %s = %q, want %q", role, cfg.Models[role], model)
		}
	}
}

func TestFormAcceptsCustomModel(t *testing.T) {
	catalog := Catalog([]string{"claude"})
	custom := []string{itoa(len(catalog) + 2), "my-local-model"} // Custom... then the value
	answers := concat(
		custom,                           // thinking
		answersForRole("haiku", catalog), // implementer
		answersForRole("", catalog),      // reviewer
	)
	driveAccessible(t, answers...)
	cfg := config.Config{Tools: []string{"claude"}, Models: map[string]string{}}
	if err := form(&cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Models[config.RoleThinking] != "my-local-model" {
		t.Fatalf("thinking = %q, want the custom model", cfg.Models[config.RoleThinking])
	}
	if cfg.Models[config.RoleImplementer] != "haiku" {
		t.Fatalf("implementer = %q, want haiku", cfg.Models[config.RoleImplementer])
	}
}

func TestFormSeedsExistingCustomModel(t *testing.T) {
	// An existing model outside the catalog opens the form on Custom..., so
	// accepting the seeded value with a blank line keeps it.
	answers := concat(
		[]string{"", ""}, // thinking: accept the seeded Custom... value
		[]string{"", ""}, // implementer: accept (tool default)
		[]string{"", ""}, // reviewer: accept (tool default)
	)
	driveAccessible(t, answers...)
	cfg := config.Config{
		Tools:  []string{"claude"},
		Models: map[string]string{config.RoleThinking: "my-local-model"},
	}
	if err := form(&cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Models[config.RoleThinking] != "my-local-model" {
		t.Fatalf("thinking = %q, want the seeded custom model", cfg.Models[config.RoleThinking])
	}
}
