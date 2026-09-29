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
