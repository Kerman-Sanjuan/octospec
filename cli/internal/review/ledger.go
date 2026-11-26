package review

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// LedgerSchema versions the record so a later reader can tell an old file from a broken one.
const LedgerSchema = 1

// Record is one pull request's verdict, kept in the tree so it survives the merge. The comment is
// for the person reading that pull request; the record is for the questions that span many.
type Record struct {
	Schema    int       `json:"schema"`
	PR        int       `json:"pr"`
	Issue     int       `json:"issue,omitempty"`
	Head      string    `json:"head,omitempty"`
	Score     float64   `json:"score"`
	Threshold float64   `json:"threshold"`
	Risk      string    `json:"risk"`
	RiskWhy   []string  `json:"risk_why,omitempty"`
	Verdict   string    `json:"verdict"`
	Findings  int       `json:"findings"`
	At        time.Time `json:"at"`
}

// LedgerPath is the file that holds a pull request's record. One file per pull request, never an
// append-only log: two open branches must not write the same line, and a merge conflict in a ledger
// is silent data loss in the one place history was wanted.
func LedgerPath(dir string, pr int) string {
	return filepath.Join(dir, strconv.Itoa(pr)+".json")
}

// WriteRecord stores rec as ledger/<pr>.json, rewriting whatever was there. A record with no pull
// request number has no key and is refused.
func WriteRecord(dir string, rec Record) (string, error) {
	if rec.PR <= 0 {
		return "", errors.New("a ledger record needs a pull request number")
	}
	if rec.Schema == 0 {
		rec.Schema = LedgerSchema
	}
	if rec.At.IsZero() {
		rec.At = time.Now().UTC()
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	b, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return "", err
	}
	p := LedgerPath(dir, rec.PR)
	tmp, err := os.CreateTemp(dir, ".record-*")
	if err != nil {
		return "", err
	}
	name := tmp.Name()
	if _, err := tmp.Write(append(b, '\n')); err != nil {
		tmp.Close()
		os.Remove(name)
		return "", err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(name)
		return "", err
	}
	if err := os.Chmod(name, 0o644); err != nil {
		os.Remove(name)
		return "", err
	}
	if err := os.Rename(name, p); err != nil {
		os.Remove(name)
		return "", err
	}
	return p, nil
}

// ReadAll loads every record in dir, newest first. A file that does not parse is skipped and
// counted: one corrupt record must not hide the rest of the series.
func ReadAll(dir string) ([]Record, int, error) {
	entries, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	var out []Record
	skipped := 0
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			skipped++
			continue
		}
		var r Record
		if err := json.Unmarshal(b, &r); err != nil || r.PR <= 0 {
			skipped++
			continue
		}
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].At.After(out[j].At) })
	return out, skipped, nil
}

// Summary is the series in five numbers. Anything more belongs in a dashboard, which is out of
// scope on purpose.
type Summary struct {
	Total     int
	ByRisk    map[string]int
	ByVerdict map[string]int
	MeanScore float64
}

// Summarise counts the records by risk and by verdict, and averages the scores.
func Summarise(rs []Record) Summary {
	s := Summary{Total: len(rs), ByRisk: map[string]int{}, ByVerdict: map[string]int{}}
	if len(rs) == 0 {
		return s
	}
	var sum float64
	for _, r := range rs {
		s.ByRisk[r.Risk]++
		s.ByVerdict[r.Verdict]++
		sum += r.Score
	}
	s.MeanScore = sum / float64(len(rs))
	return s
}

// Render formats the series as a table plus the summary, in the order a person reads it: what we
// merged, then how risky it was, then the counts.
func Render(rs []Record, skipped int) string {
	var b strings.Builder
	if len(rs) == 0 {
		b.WriteString("no records\n")
		if skipped > 0 {
			fmt.Fprintf(&b, "%d file(s) skipped as unreadable\n", skipped)
		}
		return b.String()
	}
	fmt.Fprintf(&b, "%-6s %-6s %-7s %-6s %-10s %s\n", "PR", "SCORE", "RISK", "VERDICT", "WHEN", "WHY")
	for _, r := range rs {
		why := strings.Join(r.RiskWhy, "; ")
		if runes := []rune(why); len(runes) > 60 {
			// Truncate by runes: a path with an accent is two bytes in UTF-8, and cutting it in
			// half prints a replacement character in the one column that explains the row.
			why = string(runes[:57]) + "..."
		}
		fmt.Fprintf(&b, "%-6d %-6.1f %-7s %-6s %-10s %s\n",
			r.PR, r.Score, r.Risk, r.Verdict, r.At.Format("2006-01-02"), why)
	}
	s := Summarise(rs)
	fmt.Fprintf(&b, "\n%d record(s), mean score %.1f\n", s.Total, s.MeanScore)
	fmt.Fprintf(&b, "by risk: %s\n", counts(s.ByRisk, []string{RiskLow, RiskMedium, RiskHigh}))
	fmt.Fprintf(&b, "by verdict: %s\n", counts(s.ByVerdict, []string{VerdictPass, VerdictHuman}))
	if skipped > 0 {
		fmt.Fprintf(&b, "%d file(s) skipped as unreadable\n", skipped)
	}
	return b.String()
}

func counts(m map[string]int, order []string) string {
	var parts []string
	seen := map[string]bool{}
	for _, k := range order {
		seen[k] = true
		parts = append(parts, fmt.Sprintf("%s %d", k, m[k]))
	}
	var rest []string
	for k := range m {
		if !seen[k] {
			rest = append(rest, k)
		}
	}
	sort.Strings(rest)
	for _, k := range rest {
		parts = append(parts, fmt.Sprintf("%s %d", k, m[k]))
	}
	return strings.Join(parts, ", ")
}
