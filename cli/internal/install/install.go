// Package install writes the canonical commands into each selected target and
// records the result so `update` can re-apply it.
package install

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/kerman-sanjuan/octospec/cli/internal/cleanup"
	"github.com/kerman-sanjuan/octospec/cli/internal/config"
	"github.com/kerman-sanjuan/octospec/cli/internal/plan"
	"github.com/kerman-sanjuan/octospec/cli/internal/skills"
	"github.com/kerman-sanjuan/octospec/cli/internal/targets"
)

// Options controls an install run.
type Options struct {
	Tools       []string // empty means all targets
	Repo        string   // repo root for repo-local targets; defaults to the cwd
	Global      bool     // install every selected tool globally
	Interactive bool     // a terminal is available, so ask for the scope
	DryRun      bool     // print the plan and write nothing
}

// Run installs the canonical commands into every selected target.
func Run(opts Options) error {
	repo := opts.Repo
	if repo == "" {
		var err error
		if repo, err = os.Getwd(); err != nil {
			return err
		}
	}
	cfg, err := config.Load(repo)
	if err != nil {
		return err
	}
	if cfg.Scopes == nil {
		cfg.Scopes = map[string]string{}
	}

	tools := opts.Tools
	if len(tools) == 0 && opts.Interactive {
		if tools, err = selectTools(cfg.Tools); err != nil {
			return err
		}
	}
	tools = withAllTools(tools)

	scope := targets.ScopeLocal
	switch {
	case opts.Global:
		scope = targets.ScopeGlobal
	case opts.Interactive:
		if scope, err = askScope(); err != nil {
			return err
		}
	}
	if scope == targets.ScopeGlobal {
		fmt.Println("warning: a global install writes under your home directory, so the octospec commands, agents, and skills appear in every project. Choose repo-local to keep them in this repository.")
	}
	for _, tool := range tools {
		cfg.Scopes[tool] = scope
	}

	files, err := plan.Files(tools, repo, cfg.Models, cfg.Scopes)
	if err != nil {
		return err
	}
	if opts.DryRun {
		fmt.Printf("dry run: would write %d files for %s (%s)\n", len(files), strings.Join(tools, ", "), scope)
		for _, f := range files {
			fmt.Printf("  %s\n", f.Path)
		}
		return nil
	}
	if cfg.Files == nil {
		cfg.Files = map[string]string{}
	}
	cfg.Tools = tools
	keep := make(map[string]bool, len(files))
	for _, f := range files {
		keep[f.Path] = true
	}
	cleanup.Global(&cfg, keep, opts.Interactive)
	for _, f := range files {
		if err := os.MkdirAll(filepath.Dir(f.Path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(f.Path, []byte(f.Content), 0o644); err != nil {
			return err
		}
		cfg.Files[f.Path] = config.Hash([]byte(f.Content))
	}
	if err := config.Save(repo, cfg); err != nil {
		return err
	}
	if err := skills.Install(repo, tools); err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
	}
	fmt.Printf("installed %d files for %s (%s)\n", len(files), strings.Join(tools, ", "), scope)
	return nil
}

// selectTools opens a multi-select of the supported tools, pre-checking the
// names in preselect that are still targets. It returns the chosen names; an
// empty choice is returned as-is, and the caller reads it as "all tools".
func selectTools(preselect []string) ([]string, error) {
	options := toolOptions()
	chosen := validPreselect(preselect)
	if err := huh.NewForm(huh.NewGroup(
		huh.NewMultiSelect[string]().
			Title("Install octospec for which tools?").
			Description("Space toggles a tool, enter confirms. Choose none to install for all tools.").
			Options(options...).
			Value(&chosen),
	)).Run(); err != nil {
		return nil, err
	}
	return chosen, nil
}

// toolOptions builds the multi-select options from the targets, in order.
func toolOptions() []huh.Option[string] {
	options := make([]huh.Option[string], 0, len(targets.Targets))
	for _, t := range targets.Targets {
		options = append(options, huh.NewOption(t.Name+" ("+toolDescription(t.Name)+")", t.Name))
	}
	return options
}

// validPreselect keeps the names that are still targets, dropping the rest.
func validPreselect(preselect []string) []string {
	known := make(map[string]bool, len(targets.Targets))
	for _, t := range targets.Targets {
		known[t.Name] = true
	}
	var chosen []string
	for _, name := range preselect {
		if known[name] {
			chosen = append(chosen, name)
		}
	}
	return chosen
}

// withAllTools turns an empty tool slice into every target, matching the
// selector's "choose none to install all" contract.
func withAllTools(tools []string) []string {
	if len(tools) > 0 {
		return tools
	}
	all := make([]string, 0, len(targets.Targets))
	for _, t := range targets.Targets {
		all = append(all, t.Name)
	}
	return all
}

// toolDescription is the one-line description for a tool in the selector.
func toolDescription(name string) string {
	switch name {
	case "pi":
		return "prompt templates"
	case "opencode":
		return "commands and agents"
	case "copilot":
		return "prompts and agents"
	case "claude":
		return "commands and agents"
	}
	return "commands"
}

// askScope warns about a global install and asks which one to use. It defaults
// to repo-local.
func askScope() (string, error) {
	global := false
	confirm := huh.NewConfirm().
		Title("Install octospec globally?").
		Description("A global install makes the commands, agents, and skills appear in every project. Repo-local keeps them in this repository.").
		Affirmative("Global").
		Negative("Repo-local (recommended)").
		Value(&global)
	if err := huh.NewForm(huh.NewGroup(confirm)).Run(); err != nil {
		return "", err
	}
	if global {
		return targets.ScopeGlobal, nil
	}
	return targets.ScopeLocal, nil
}
