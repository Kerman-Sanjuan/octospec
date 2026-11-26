// Package doctor reports the state of an octospec install. It is read-only.
package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/kerman-sanjuan/octospec/cli/internal/budget"
	"github.com/kerman-sanjuan/octospec/cli/internal/config"
	"github.com/kerman-sanjuan/octospec/cli/internal/lock"
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
		checks = append(checks, Check{"agents", false, "run `octospec install`"})
	} else {
		files, err := plan.Files(cfg.Tools, repo, cfg.Models, cfg.Scopes)
		if err != nil {
			return checks, err
		}
		missingCmd, missingAgent := 0, 0
		for _, f := range files {
			if _, err := os.Stat(f.Path); err != nil {
				if f.Kind == plan.KindAgent {
					missingAgent++
				} else {
					missingCmd++
				}
			}
		}
		tools := strings.Join(cfg.Tools, ",")
		checks = append(checks,
			Check{"commands (" + tools + ")", missingCmd == 0, fmt.Sprintf("%d missing", missingCmd)},
			Check{"agents (" + tools + ")", missingAgent == 0, fmt.Sprintf("%d missing", missingAgent)},
		)
	}

	checks = append(checks, checkModels(cfg))

	checks = append(checks, checkLabels(repo))

	if _, err := os.Stat(filepath.Join(repo, "scripts", "check-gates.sh")); err == nil {
		checks = append(checks, Check{"check-gates.sh", true, ""})
	} else {
		checks = append(checks, Check{"check-gates.sh", false, "run `octospec seed`"})
	}

	checks = append(checks, checkLock(repo))

	checks = append(checks, checkBudget(repo))

	return checks, nil
}

// checkBudget prints the ceilings a run will be held to. A repository cannot pretend it has budgets:
// the answer is either its file or the defaults it never chose.
func checkBudget(repo string) Check {
	policy, err := budget.Load(repo)
	if err != nil {
		return Check{"budget", false, err.Error()}
	}
	l := policy.With("apply")
	detail := fmt.Sprintf("%d min, %d tokens, %d files, %d diff lines, %d tasks",
		l.MaxMinutes, l.MaxTokens, l.MaxFiles, l.MaxDiffLines, l.MaxTasks)
	source := "defaults"
	if _, err := os.Stat(filepath.Join(repo, filepath.FromSlash(budget.PolicyPath))); err == nil {
		source = budget.PolicyPath
	}
	return Check{"budget", true, detail + " (" + source + ")"}
}

// checkLock reports who holds the archive lock. A free lock is healthy; an expired lease is not,
// because the next /archive will refuse to run until somebody clears it.
func checkLock(repo string) Check {
	st, err := lock.Status(repo)
	if err != nil {
		return Check{"archive lock", false, err.Error()}
	}
	switch {
	case !st.Held:
		return Check{"archive lock", true, "free"}
	case st.Stale:
		return Check{"archive lock", false, fmt.Sprintf("stale: %s; held by %s; clear with `octospec lock release --force --change %s`",
			st.Why, st.Record.Change, st.Record.Change)}
	default:
		return Check{"archive lock", true, fmt.Sprintf("held by %s since %s on %s, lease expires in %s",
			st.Record.Change, st.Record.AcquiredAt.Format(time.RFC3339), st.Record.Host,
			time.Until(st.Record.Expiry()).Round(time.Minute))}
	}
}

// checkModels reports the configured role-to-model map. An empty map means the
// tool defaults, which is healthy.
func checkModels(cfg config.Config) Check {
	var set []string
	for _, role := range config.Roles {
		if cfg.Models[role] != "" {
			set = append(set, role+"="+cfg.Models[role])
		}
	}
	if len(set) == 0 {
		return Check{"models", true, "tool defaults"}
	}
	return Check{"models", true, strings.Join(set, " ")}
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
