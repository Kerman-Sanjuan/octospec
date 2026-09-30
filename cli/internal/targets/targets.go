// Package targets maps each supported tool to its destination directory and
// file format. Adding a tool is a new entry here, never a new source copy.
package targets

import (
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

// Target is one tool. It maps where the commands land and where the agents
// land, plus the file format for each.
type Target struct {
	Name        string
	Scope       Scope
	Dir         string // Global: may start with "~/"; RepoLocal: relative to the repo
	Ext         string // file extension for the target's command files
	Format      string // command renderer: formatDefault or formatCopilot
	AgentDir    string // directory for the target's agent files
	AgentExt    string // file extension for the target's agent files
	AgentFormat string // agent renderer: formatDefault or formatCopilot
}

// Targets is the single source of tool mappings.
var Targets = []Target{
	{Name: "pi", Scope: Global, Dir: "~/.pi/agent/prompts", Ext: ".md", Format: formatDefault,
		AgentDir: "~/.pi/agent/agents", AgentExt: ".md", AgentFormat: formatDefault},
	{Name: "opencode", Scope: Global, Dir: "~/.config/opencode/command", Ext: ".md", Format: formatDefault,
		AgentDir: "~/.config/opencode/agent", AgentExt: ".md", AgentFormat: formatDefault},
	{Name: "copilot", Scope: RepoLocal, Dir: ".github/prompts", Ext: ".prompt.md", Format: formatCopilot,
		AgentDir: ".github/agents", AgentExt: ".md", AgentFormat: formatCopilot},
	{Name: "claude", Scope: RepoLocal, Dir: ".claude/commands", Ext: ".md", Format: formatDefault,
		AgentDir: ".claude/agents", AgentExt: ".md", AgentFormat: formatDefault},
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

// AgentExpand resolves the target's agent directory for a run rooted at repo.
func (t Target) AgentExpand(repo string) string {
	if t.Scope == Global {
		return expandHome(t.AgentDir)
	}
	return filepath.Join(repo, t.AgentDir)
}

// Filename returns the file name a command gets on this target.
func (t Target) Filename(name string) string {
	return name + t.Ext
}

// AgentFilename returns the file name an agent gets on this target.
func (t Target) AgentFilename(name string) string {
	return name + t.AgentExt
}

// Render returns the file content for a command on this target.
func (t Target) Render(c payload.Command) string {
	return render(t.Format, c)
}

func expandHome(p string) string {
	if strings.HasPrefix(p, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, p[2:])
		}
	}
	return p
}
