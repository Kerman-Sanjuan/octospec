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
// configured model, or omitting it when empty.
func (t Target) RenderAgent(a payload.Agent, model string) string {
	return renderAgent(t.AgentFormat, a, model)
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
func renderAgent(format string, a payload.Agent, model string) string {
	var fm strings.Builder
	if format == formatCopilot {
		fm.WriteString("---\nmode: agent\n")
		fmt.Fprintf(&fm, "description: %s\n", a.Description)
		if model != "" {
			fmt.Fprintf(&fm, "model: %s\n", model)
		}
		fm.WriteString("---\n\n")
		return fm.String() + a.Body
	}
	fmt.Fprintf(&fm, "---\nname: %s\ndescription: %s\n", a.Name, a.Description)
	if model != "" {
		fmt.Fprintf(&fm, "model: %s\n", model)
	}
	fm.WriteString("---\n\n")
	return fm.String() + a.Body
}
