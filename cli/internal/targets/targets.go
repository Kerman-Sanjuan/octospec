// Package targets maps each supported tool to its repo-local and global
// directories and its file format. Adding a tool is a new entry here, never a
// new source copy.
package targets

import (
	"os"
	"path/filepath"
	"strings"
)

// Scopes for an install. Repo-local is the default.
const (
	ScopeLocal  = "local"
	ScopeGlobal = "global"
)

// Target is one tool. It maps a repo-local and a global directory for the
// commands and for the agents, plus the file format for each. An empty agent
// directory means the tool has no agent mechanism.
type Target struct {
	Name           string
	Dir            string // repo-local commands
	Ext            string // file extension for the command files
	Format         string // command renderer: formatDefault or formatCopilot
	AgentDir       string // repo-local agents; empty means no agent files
	AgentExt       string // file extension for the agent files
	AgentFormat    string // agent renderer: formatDefault or formatCopilot
	GlobalDir      string // global commands
	GlobalAgentDir string // global agents; empty means no agent files
	ToolFormat     string // how the agent's tool surface is rendered
}

// Tool surface renderers. A harness with no permission field is toolAdvisory:
// the scope goes in the body.
const (
	toolClaude   = "claude"
	toolCopilot  = "copilot"
	toolOpencode = "opencode"
	toolAdvisory = "advisory"
)

// Targets is the single source of tool mappings.
var Targets = []Target{
	{Name: "pi", Dir: ".pi/prompts", Ext: ".md", Format: formatDefault,
		AgentDir: "", AgentExt: "", AgentFormat: formatDefault,
		GlobalDir: "~/.pi/agent/prompts", GlobalAgentDir: "",
		ToolFormat: toolAdvisory},
	{Name: "opencode", Dir: ".opencode/command", Ext: ".md", Format: formatDefault,
		AgentDir: ".opencode/agents", AgentExt: ".md", AgentFormat: formatDefault,
		GlobalDir: "~/.config/opencode/command", GlobalAgentDir: "~/.config/opencode/agents",
		ToolFormat: toolOpencode},
	{Name: "copilot", Dir: ".github/prompts", Ext: ".prompt.md", Format: formatCopilot,
		AgentDir: ".github/agents", AgentExt: ".agent.md", AgentFormat: formatCopilot,
		GlobalDir: "~/.copilot/prompts", GlobalAgentDir: "~/.copilot/agents",
		ToolFormat: toolCopilot},
	{Name: "claude", Dir: ".claude/commands", Ext: ".md", Format: formatDefault,
		AgentDir: ".claude/agents", AgentExt: ".md", AgentFormat: formatDefault,
		GlobalDir: "~/.claude/commands", GlobalAgentDir: "~/.claude/agents",
		ToolFormat: toolClaude},
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

// HasAgents reports whether the target installs agent files.
func (t Target) HasAgents() bool { return t.AgentDir != "" }

// Expand resolves the command directory for a scope.
func (t Target) Expand(repo, scope string) string {
	if scope == ScopeGlobal {
		return expandHome(t.GlobalDir)
	}
	return filepath.Join(repo, t.Dir)
}

// AgentExpand resolves the agent directory for a scope. It is empty when the
// target has no agent mechanism.
func (t Target) AgentExpand(repo, scope string) string {
	if !t.HasAgents() {
		return ""
	}
	if scope == ScopeGlobal {
		return expandHome(t.GlobalAgentDir)
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

func expandHome(p string) string {
	if strings.HasPrefix(p, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, p[2:])
		}
	}
	return p
}
