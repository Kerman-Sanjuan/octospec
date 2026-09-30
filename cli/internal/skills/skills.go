// Package skills makes the OpenSpec skills available to each tool by asking
// the OpenSpec CLI, which owns them.
package skills

import (
	"fmt"
	"os/exec"
	"strings"
)

// toolMap maps octospec tool names to OpenSpec's --tools names.
var toolMap = map[string]string{
	"pi":       "pi",
	"opencode": "opencode",
	"copilot":  "github-copilot",
	"claude":   "claude",
}

// Names maps octospec tool names to OpenSpec's, dropping unknown ones.
func Names(tools []string) []string {
	var out []string
	for _, t := range tools {
		if n, ok := toolMap[t]; ok && !contains(out, n) {
			out = append(out, n)
		}
	}
	return out
}

// Install runs `openspec init --tools <list>` in repo so the OpenSpec skills
// exist for the given tools. It returns an error when openspec is absent.
func Install(repo string, tools []string) error {
	names := Names(tools)
	if len(names) == 0 {
		return nil
	}
	bin, err := exec.LookPath("openspec")
	if err != nil {
		return fmt.Errorf("openspec not on PATH; skills skipped")
	}
	cmd := exec.Command(bin, "init", "--tools", strings.Join(names, ","), "--force")
	cmd.Dir = repo
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("openspec init failed: %v: %s", err, out)
	}
	return nil
}

func contains(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}
