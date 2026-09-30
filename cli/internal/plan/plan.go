// Package plan turns a set of tools into the concrete files they would write.
package plan

import (
	"fmt"
	"path/filepath"

	"github.com/kerman-sanjuan/octospec/cli/internal/payload"
	"github.com/kerman-sanjuan/octospec/cli/internal/targets"
)

// File is one rendered command or agent for one tool.
type File struct {
	Tool    string
	Kind    string // KindCommand or KindAgent
	Path    string
	Content string
}

// Kinds of planned files.
const (
	KindCommand = "command"
	KindAgent   = "agent"
)

// Files returns every file the given tools would write for a repo. An empty
// tools slice means every target. The models map is keyed by role; an empty or
// missing entry means the tool default.
func Files(tools []string, repo string, models map[string]string) ([]File, error) {
	cmds, err := payload.Commands()
	if err != nil {
		return nil, err
	}
	agents, err := payload.Agents()
	if err != nil {
		return nil, err
	}
	if len(tools) == 0 {
		for _, t := range targets.Targets {
			tools = append(tools, t.Name)
		}
	}
	var out []File
	for _, name := range tools {
		t, ok := targets.Get(name)
		if !ok {
			return nil, fmt.Errorf("unknown tool %q (want one of pi, opencode, copilot, claude)", name)
		}
		dir := t.Expand(repo)
		for _, c := range cmds {
			out = append(out, File{
				Tool:    name,
				Kind:    KindCommand,
				Path:    filepath.Join(dir, t.Filename(c.Name)),
				Content: t.Render(c),
			})
		}
		agentDir := t.AgentExpand(repo)
		for _, a := range agents {
			out = append(out, File{
				Tool:    name,
				Kind:    KindAgent,
				Path:    filepath.Join(agentDir, t.AgentFilename(a.Name)),
				Content: t.RenderAgent(a, models[a.Role]),
			})
		}
	}
	return out, nil
}
