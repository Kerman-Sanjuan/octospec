// Package install writes the canonical commands into each selected target.
package install

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kerman-sanjuan/octospec/cli/internal/payload"
	"github.com/kerman-sanjuan/octospec/cli/internal/targets"
)

// Options controls an install run.
type Options struct {
	Tools []string // empty means all targets
	Repo  string   // repo root for repo-local targets; defaults to the cwd
}

// Run installs the canonical commands into every selected target.
func Run(opts Options) error {
	cmds, err := payload.Commands()
	if err != nil {
		return err
	}
	repo := opts.Repo
	if repo == "" {
		repo, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	selected := opts.Tools
	if len(selected) == 0 {
		for _, t := range targets.Targets {
			selected = append(selected, t.Name)
		}
	}
	for _, name := range selected {
		t, ok := targets.Get(name)
		if !ok {
			return fmt.Errorf("unknown tool %q (want one of pi, opencode, copilot, claude)", name)
		}
		dir := t.Expand(repo)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		for _, c := range cmds {
			path := filepath.Join(dir, t.Filename(c.Name))
			if err := os.WriteFile(path, []byte(t.Render(c)), 0o644); err != nil {
				return err
			}
		}
		fmt.Printf("%-9s -> %s (%d commands)\n", name, dir, len(cmds))
	}
	return nil
}
