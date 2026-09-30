// Package doctor reports the state of an octospec install. It is read-only.
package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kerman-sanjuan/octospec/cli/internal/config"
	"github.com/kerman-sanjuan/octospec/cli/internal/plan"
	"github.com/kerman-sanjuan/octospec/cli/internal/seed"
)

// Check is one line of the report.
type Check struct {
	Name string
	OK   bool
	Note string
}

var workflowLabels = []string{
	"type:feature", "type:bug",
	"status:backlog", "status:spec-ready", "status:in-progress", "status:in-review",
}

// Run reports the state of an install rooted at repo.
func Run(repo string) ([]Check, error) {
	if repo == "" {
		var err error
		if repo, err = os.Getwd(); err != nil {
			return nil, err
		}
	}
	var checks []Check

	if _, err := exec.LookPath("openspec"); err == nil {
		checks = append(checks, Check{"openspec CLI", true, ""})
	} else {
		checks = append(checks, Check{"openspec CLI", false, "not on PATH"})
	}

	sd := seed.SchemaDir()
	if _, err := os.Stat(filepath.Join(sd, "schema.yaml")); err == nil {
		checks = append(checks, Check{"schema", true, sd})
	} else {
		checks = append(checks, Check{"schema", false, "run `octospec seed`"})
	}

	cfg, err := config.Load(repo)
	if err != nil {
		return checks, err
	}
	if len(cfg.Tools) == 0 {
		checks = append(checks, Check{"commands", false, "run `octospec install`"})
	} else {
		files, err := plan.Files(cfg.Tools, repo, nil)
		if err != nil {
			return checks, err
		}
		missing := 0
		for _, f := range files {
			if _, err := os.Stat(f.Path); err != nil {
				missing++
			}
		}
		checks = append(checks, Check{
			Name: "commands (" + strings.Join(cfg.Tools, ",") + ")",
			OK:   missing == 0,
			Note: fmt.Sprintf("%d missing", missing),
		})
	}

	checks = append(checks, checkLabels(repo))

	if _, err := os.Stat(filepath.Join(repo, "scripts", "check-gates.sh")); err == nil {
		checks = append(checks, Check{"check-gates.sh", true, ""})
	} else {
		checks = append(checks, Check{"check-gates.sh", false, "run `octospec seed`"})
	}

	return checks, nil
}

func checkLabels(repo string) Check {
	gh, err := exec.LookPath("gh")
	if err != nil {
		return Check{"labels", false, "gh not on PATH"}
	}
	cmd := exec.Command(gh, "label", "list", "--limit", "100", "--json", "name", "--jq", ".[].name")
	cmd.Dir = repo
	out, err := cmd.Output()
	if err != nil {
		return Check{"labels", false, "gh label list failed (is this a GitHub repo?)"}
	}
	have := map[string]bool{}
	for _, line := range strings.Split(string(out), "\n") {
		have[strings.TrimSpace(line)] = true
	}
	var missing []string
	for _, l := range workflowLabels {
		if !have[l] {
			missing = append(missing, l)
		}
	}
	if len(missing) == 0 {
		return Check{"labels", true, ""}
	}
	return Check{"labels", false, "missing: " + strings.Join(missing, ", ")}
}
