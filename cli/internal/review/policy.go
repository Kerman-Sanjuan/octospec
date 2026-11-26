package review

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// PolicyPath is where a repository keeps the numbers. octospec ships defaults and never owns
// them: tuning the threshold must not need a release.
const PolicyPath = ".octospec/review.json"

// Policy is the repository's say in the score.
type Policy struct {
	Threshold      float64          `json:"threshold,omitempty"`
	Weights        map[Type]float64 `json:"weights,omitempty"`
	SensitivePaths []string         `json:"sensitive_paths,omitempty"`
	LargeDiffLines int              `json:"large_diff_lines,omitempty"`

	// Source says where these numbers came from, so the output can admit it used defaults.
	Source string `json:"-"`
}

// DefaultPolicy is what an adopting repository starts from: a threshold low enough to pass a
// decent change, a diff size that is large for a solo maintainer, and no sensitive paths.
func DefaultPolicy() Policy {
	w := make(map[Type]float64, len(DefaultWeights))
	for k, v := range DefaultWeights {
		w[k] = v
	}
	return Policy{Threshold: 7, Weights: w, LargeDiffLines: 400, Source: "defaults"}
}

// LoadPolicy reads .octospec/review.json, filling anything absent from the defaults.
func LoadPolicy(repo string) (Policy, error) {
	p := DefaultPolicy()
	path := filepath.Join(repo, filepath.FromSlash(PolicyPath))
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return p, nil
	}
	if err != nil {
		return p, err
	}
	var file Policy
	if err := json.Unmarshal(b, &file); err != nil {
		return p, fmt.Errorf("%s: %w", PolicyPath, err)
	}
	if file.Threshold != 0 {
		if file.Threshold < 0 || file.Threshold > 10 {
			return p, fmt.Errorf("%s: threshold must be between 0 and 10, got %v", PolicyPath, file.Threshold)
		}
		p.Threshold = file.Threshold
	}
	for t, w := range file.Weights {
		if !t.valid() {
			return p, fmt.Errorf("%s: unknown finding type %q", PolicyPath, t)
		}
		if w < 0 {
			return p, fmt.Errorf("%s: weight for %q must not be negative", PolicyPath, t)
		}
		p.Weights[t] = w
	}
	if len(file.SensitivePaths) > 0 {
		p.SensitivePaths = file.SensitivePaths
	}
	if file.LargeDiffLines > 0 {
		p.LargeDiffLines = file.LargeDiffLines
	}
	p.Source = PolicyPath
	return p, nil
}
