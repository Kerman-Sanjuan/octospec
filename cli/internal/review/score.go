package review

import (
	"math"
	"sort"
)

// Deduction is one finding as the score saw it: what cost how much, and where.
type Deduction struct {
	ID     string  `json:"id"`
	Axis   Axis    `json:"axis"`
	Type   Type    `json:"type"`
	Weight float64 `json:"weight"`
	File   string  `json:"file"`
	Line   int     `json:"line,omitempty"`
	Note   string  `json:"note,omitempty"`
}

// Result is the whole verdict, ready to be printed as JSON or rendered into a comment.
type Result struct {
	Score        float64          `json:"score"`
	Threshold    float64          `json:"threshold"`
	Verdict      string           `json:"verdict"` // pass | human
	Risk         string           `json:"risk"`
	RiskWhy      []string         `json:"risk_why,omitempty"`
	ByAxis       map[Axis]float64 `json:"by_axis"`
	Deductions   []Deduction      `json:"deductions"`
	PolicySource string           `json:"policy_source"`
}

// Verdicts.
const (
	VerdictPass  = "pass"
	VerdictHuman = "human"
)

// Score starts at ten and subtracts what each finding costs. A high risk asks for a human
// whatever the score: the point of the risk table is that no number of style points buys a
// sensitive change its merge.
func Score(findings []Finding, p Policy, in RiskInput) Result {
	total := 10.0
	byAxis := map[Axis]float64{}
	deductions := make([]Deduction, 0, len(findings))
	for _, f := range findings {
		w, ok := p.Weights[f.Type]
		if !ok {
			w = DefaultWeights[f.Type]
		}
		total -= w
		byAxis[f.Axis] = round1(byAxis[f.Axis] + w)
		deductions = append(deductions, Deduction{
			ID: f.ID, Axis: f.Axis, Type: f.Type, Weight: w,
			File: f.File, Line: f.Line, Note: f.Note,
		})
	}
	if total < 0 {
		total = 0
	}
	total = round1(total)
	sort.SliceStable(deductions, func(i, j int) bool {
		if deductions[i].Weight != deductions[j].Weight {
			return deductions[i].Weight > deductions[j].Weight
		}
		return deductions[i].ID < deductions[j].ID
	})
	risk, why := Risk(in, p)
	verdict := VerdictPass
	if risk == RiskHigh || total <= p.Threshold {
		verdict = VerdictHuman
	}
	return Result{
		Score: total, Threshold: p.Threshold, Verdict: verdict,
		Risk: risk, RiskWhy: why, ByAxis: byAxis, Deductions: deductions,
		PolicySource: p.Source,
	}
}

func round1(v float64) float64 { return math.Round(v*10) / 10 }
