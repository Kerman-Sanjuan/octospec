package review

import (
	"fmt"
	"strings"
)

// Risk levels. Three inputs, four outcomes: a risk *score* invites tuning the weights until
// everything is medium, so the table stays small enough to argue about.
const (
	RiskLow    = "low"
	RiskMedium = "medium"
	RiskHigh   = "high"
)

// RiskInput is what the script can read about a change without believing the model.
type RiskInput struct {
	// ChangedFiles are the paths the pull request touches, relative to the repository root.
	ChangedFiles []string
	// DiffLines is added plus removed lines.
	DiffLines int
	// Familiar is true when the person responsible for the merge has touched these paths before.
	Familiar bool
}

// Risk applies the table and returns the level plus the facts that produced it.
func Risk(in RiskInput, p Policy) (string, []string) {
	var why []string
	sensitive := ""
	for _, f := range in.ChangedFiles {
		if matchesAny(f, p.SensitivePaths) {
			sensitive = f
			break
		}
	}
	large := p.LargeDiffLines > 0 && in.DiffLines > p.LargeDiffLines
	if large {
		why = append(why, fmt.Sprintf("diff of %d lines", in.DiffLines))
	}
	if !in.Familiar {
		why = append(why, "the author has not touched these paths before")
	}
	switch {
	case sensitive != "":
		return RiskHigh, append([]string{"touches sensitive path " + sensitive}, why...)
	case large && !in.Familiar:
		return RiskHigh, why
	case large || !in.Familiar:
		return RiskMedium, why
	default:
		return RiskLow, nil
	}
}

func matchesAny(path string, prefixes []string) bool {
	p := strings.TrimPrefix(filepathSlash(path), "./")
	for _, prefix := range prefixes {
		prefix = strings.TrimPrefix(filepathSlash(prefix), "./")
		if prefix == "" {
			continue
		}
		if p == prefix || strings.HasPrefix(p, prefix+"/") || strings.HasPrefix(p, prefix) {
			return true
		}
	}
	return false
}

// filepathSlash keeps the comparison stable when a path arrives from Windows.
func filepathSlash(p string) string { return strings.ReplaceAll(p, "\\", "/") }
