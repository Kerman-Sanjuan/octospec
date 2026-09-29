// Package targets maps each supported tool to its destination directory and
// file format. Adding a tool is a new entry here, never a new source copy.
package targets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kerman-sanjuan/octospec/cli/internal/payload"
)

// Scope says whether a target writes to a user directory or into a repo.
type Scope int

const (
	// Global targets write to a user directory (no repo needed).
	Global Scope = iota
	// RepoLocal targets write into the repository.
	RepoLocal
)

// Target is one tool.
type Target struct {
	Name  string
	Scope Scope
	Dir   string // Global: may start with "~/"; RepoLocal: relative to the repo
}

// Targets is the single source of tool mappings.
var Targets = []Target{
	{Name: "pi", Scope: Global, Dir: "~/.pi/agent/prompts"},
	{Name: "opencode", Scope: Global, Dir: "~/.config/opencode/command"},
	{Name: "copilot", Scope: RepoLocal, Dir: ".github/prompts"},
	{Name: "claude", Scope: RepoLocal, Dir: ".claude/commands"},
}

// Get returns the target with the given name.
func Get(name string) (Target, bool) {
	for _, t := range Targets {
		if t.Name == name {
			return t, true
		}
	}
	return Target{}, false
}

// Expand resolves the target's directory for a run rooted at repo.
func (t Target) Expand(repo string) string {
	if t.Scope == Global {
		return expandHome(t.Dir)
	}
	return filepath.Join(repo, t.Dir)
}

// Filename returns the file name a command gets on this target.
func (t Target) Filename(name string) string {
	if t.Name == "copilot" {
		return name + ".prompt.md"
	}
	return name + ".md"
}

// Render returns the file content for a command on this target.
func (t Target) Render(c payload.Command) string {
	var fm strings.Builder
	switch t.Name {
	case "copilot":
		fmt.Fprintf(&fm, "---\nmode: agent\ndescription: %s\n---\n\n", c.Description)
	default:
		fmt.Fprintf(&fm, "---\ndescription: %s\n", c.Description)
		if c.ArgumentHint != "" {
			fmt.Fprintf(&fm, "argument-hint: %s\n", c.ArgumentHint)
		}
		fm.WriteString("---\n\n")
	}
	return fm.String() + c.Body
}

func expandHome(p string) string {
	if strings.HasPrefix(p, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, p[2:])
		}
	}
	return p
}
