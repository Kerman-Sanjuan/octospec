// Package review turns a list of findings into a score, a risk level and a verdict.
//
// The model detects; this package does the arithmetic. Splitting it that way is the whole point:
// two runs over the same findings must produce the same number, which is what makes the score
// something a gate can use and a person can argue with.
package review

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Axis is what a finding is about. The five axes are the contract with the comment: the reader
// sees where the score was lost, not only how much.
type Axis string

const (
	AxisCorrectness  Axis = "correctness"
	AxisArchitecture Axis = "architecture"
	AxisDomain       Axis = "domain"
	AxisTests        Axis = "tests"
	AxisUI           Axis = "ui"
)

// Axes is the closed vocabulary, in display order.
var Axes = []Axis{AxisCorrectness, AxisArchitecture, AxisDomain, AxisTests, AxisUI}

func (a Axis) valid() bool {
	for _, x := range Axes {
		if x == a {
			return true
		}
	}
	return false
}

// Type is what kind of problem a finding is, and it is what the policy weighs.
type Type string

const (
	TypeBlocking    Type = "blocking"
	TypeCorrectness Type = "correctness"
	TypeTest        Type = "test"
	TypeDuplication Type = "duplication"
	TypeNaming      Type = "naming"
	TypeStyle       Type = "style"
	TypeDocs        Type = "docs"
)

// DefaultWeights is what each finding type costs when a repository does not say otherwise. They
// are deliberately blunt: the point of the first twenty pull requests is to replace them.
var DefaultWeights = map[Type]float64{
	TypeBlocking:    4,
	TypeCorrectness: 3,
	TypeTest:        2,
	TypeDuplication: 1.5,
	TypeNaming:      0.5,
	TypeStyle:       0.5,
	TypeDocs:        0.5,
}

func (t Type) valid() bool {
	_, ok := DefaultWeights[t]
	return ok
}

// Finding is one thing the model noticed, cited so a human can check it.
type Finding struct {
	ID   string `json:"id,omitempty"`
	Axis Axis   `json:"axis"`
	Type Type   `json:"type"`
	File string `json:"file"`
	Line int    `json:"line,omitempty"`
	Note string `json:"note,omitempty"`
}

// Document is the shape the command prefers. A bare array is accepted too, because an agent
// writing JSON by hand will produce both.
type Document struct {
	Findings []Finding `json:"findings"`
}

// Parse reads findings from a reader. An unknown axis or type is an error, not a silent "other":
// a finding nobody can classify is a finding the score cannot defend.
func Parse(r io.Reader) ([]Finding, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var doc Document
	if err := json.Unmarshal(b, &doc); err == nil && doc.Findings != nil {
		return validate(doc.Findings)
	}
	var list []Finding
	if err := json.Unmarshal(b, &list); err != nil {
		return nil, fmt.Errorf("findings must be a JSON array or {\"findings\": [...]}: %w", err)
	}
	return validate(list)
}

func validate(fs []Finding) ([]Finding, error) {
	for i := range fs {
		f := &fs[i]
		if f.ID == "" {
			f.ID = fmt.Sprintf("f%d", i+1)
		}
		if !f.Axis.valid() {
			return nil, fmt.Errorf("finding %s: unknown axis %q", f.ID, f.Axis)
		}
		if !f.Type.valid() {
			return nil, fmt.Errorf("finding %s: unknown type %q", f.ID, f.Type)
		}
		if strings.TrimSpace(f.File) == "" {
			return nil, fmt.Errorf("finding %s: a file is required", f.ID)
		}
		if f.Line < 0 {
			return nil, fmt.Errorf("finding %s: line must not be negative", f.ID)
		}
	}
	return fs, nil
}
