package seed

import (
	"fmt"
	"os/exec"
)

type label struct {
	name  string
	color string
	desc  string
}

var workflowLabels = []label{
	{"type:feature", "a2eeef", "New capability"},
	{"type:bug", "d73a4a", "Something behaves differently than expected"},
	{"status:backlog", "d4c5f9", "Queued for planning"},
	{"status:spec-ready", "0e8a16", "Spec artifacts created and validated"},
	{"status:in-progress", "fbca04", "Work in progress"},
	{"status:in-review", "1d76db", "Under review"},
}

// InstallLabels provisions the workflow labels in repo via gh, idempotently.
func InstallLabels(repo string) (int, error) {
	if _, err := exec.LookPath("gh"); err != nil {
		return 0, fmt.Errorf("gh not found; skipping label provisioning")
	}
	n := 0
	for _, l := range workflowLabels {
		cmd := exec.Command("gh", "label", "create", l.name,
			"--color", l.color, "--description", l.desc, "--force")
		cmd.Dir = repo
		if out, err := cmd.CombinedOutput(); err != nil {
			return n, fmt.Errorf("gh label create %s: %v: %s", l.name, err, out)
		}
		n++
	}
	return n, nil
}
