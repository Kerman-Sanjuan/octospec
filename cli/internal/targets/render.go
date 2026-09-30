// Rendering turns a canonical command into the file a target tool expects.
package targets

import (
	"fmt"
	"strings"

	"github.com/kerman-sanjuan/octospec/cli/internal/payload"
)

const (
	formatDefault = "default"
	formatCopilot = "copilot"
)

// RenderAgent formats a canonical agent for this target, stamping the
// configured model, or omitting it when empty, and mapping the tool surface to
// the target's syntax.
func (t Target) RenderAgent(a payload.Agent, model string) string {
	return renderAgent(t.AgentFormat, t.ToolFormat, a, model)
}

// render formats a canonical command for a target format.
func render(format string, c payload.Command) string {
	if format == formatCopilot {
		return fmt.Sprintf("---\nmode: agent\ndescription: %s\n---\n\n", c.Description) + c.Body
	}
	var fm strings.Builder
	fmt.Fprintf(&fm, "---\ndescription: %s\n", c.Description)
	if c.ArgumentHint != "" {
		fmt.Fprintf(&fm, "argument-hint: %s\n", c.ArgumentHint)
	}
	fm.WriteString("---\n\n")
	return fm.String() + c.Body
}

// renderAgent formats a canonical agent for a target format.
func renderAgent(format, toolFormat string, a payload.Agent, model string) string {
	body := a.Body
	if toolFormat == toolAdvisory {
		body = "Tools: " + strings.Join(mapTools(toolFormat, a.Tools), ", ") + "\n\n" + body
	}
	var fm strings.Builder
	if format == formatCopilot {
		fm.WriteString("---\nmode: agent\n")
		fmt.Fprintf(&fm, "description: %s\n", a.Description)
		if model != "" {
			fmt.Fprintf(&fm, "model: %s\n", model)
		}
		writeTools(&fm, toolFormat, a.Tools)
		fm.WriteString("---\n\n")
		return fm.String() + body
	}
	fmt.Fprintf(&fm, "---\nname: %s\ndescription: %s\n", a.Name, a.Description)
	if model != "" {
		fmt.Fprintf(&fm, "model: %s\n", model)
	}
	writeTools(&fm, toolFormat, a.Tools)
	fm.WriteString("---\n\n")
	return fm.String() + body
}

// writeTools writes the tool surface in the harness's front-matter syntax. The
// advisory format has no field, so the body carries the scope instead.
func writeTools(fm *strings.Builder, format string, caps []string) {
	names := mapTools(format, caps)
	switch format {
	case toolClaude, toolCopilot:
		fmt.Fprintf(fm, "tools: %s\n", strings.Join(names, ", "))
	case toolOpencode:
		// Provisional: opencode documents a `permissions` field for tool
		// access control, but not its exact shape. Confirm against the
		// opencode docs.
		fm.WriteString("permissions:\n")
		for _, n := range names {
			fmt.Fprintf(fm, "  - %s\n", n)
		}
	}
}

// toolTable maps a canonical capability to the harness's tool names.
var toolTable = map[string]map[string][]string{
	toolClaude: {
		"read":   {"Read", "Glob", "Grep"},
		"edit":   {"Edit", "Write"},
		"search": {"Grep", "Glob"},
		"shell":  {"Bash"},
		"web":    {"WebFetch", "WebSearch"},
		"agent":  {"Task"},
	},
	toolCopilot: {
		"read":   {"read"},
		"edit":   {"edit"},
		"search": {"search"},
		"shell":  {"execute"},
		"web":    {"web"},
		"agent":  {"agent"},
	},
	toolOpencode: {
		"read":   {"read"},
		"edit":   {"edit"},
		"search": {"grep"},
		"shell":  {"bash"},
		"web":    {"webfetch"},
		"agent":  {"task"},
	},
	toolAdvisory: {
		"read":   {"read"},
		"edit":   {"edit"},
		"search": {"search"},
		"shell":  {"shell"},
		"web":    {"web"},
		"agent":  {"agent"},
	},
}

// mapTools expands canonical capabilities into harness tool names, in order
// and without duplicates.
func mapTools(format string, caps []string) []string {
	table := toolTable[format]
	seen := map[string]bool{}
	var out []string
	for _, c := range caps {
		for _, n := range table[c] {
			if !seen[n] {
				seen[n] = true
				out = append(out, n)
			}
		}
	}
	return out
}
