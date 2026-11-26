package review

import (
	"strings"
	"testing"
)

func f(id string, axis Axis, typ Type, file string, line int) Finding {
	return Finding{ID: id, Axis: axis, Type: typ, File: file, Line: line}
}

func TestParseAcceptsBothShapes(t *testing.T) {
	doc := `{"findings":[{"id":"a","axis":"tests","type":"test","file":"x.go","line":3}]}`
	list := `[{"axis":"ui","type":"style","file":"y.tsx"}]`
	for _, in := range []string{doc, list} {
		fs, err := Parse(strings.NewReader(in))
		if err != nil {
			t.Fatalf("%s: %v", in, err)
		}
		if len(fs) != 1 || fs[0].ID == "" {
			t.Fatalf("%s -> %+v", in, fs)
		}
	}
}

func TestParseRejectsUnknownAxisOrType(t *testing.T) {
	for _, in := range []string{
		`[{"axis":"vibes","type":"test","file":"x.go"}]`,
		`[{"axis":"tests","type":"smell","file":"x.go"}]`,
		`[{"axis":"tests","type":"test"}]`,
	} {
		if _, err := Parse(strings.NewReader(in)); err == nil {
			t.Fatalf("%s was accepted", in)
		}
	}
}

func TestCleanChangeScoresTenAndPasses(t *testing.T) {
	r := Score(nil, DefaultPolicy(), RiskInput{Familiar: true})
	if r.Score != 10 || r.Verdict != VerdictPass || r.Risk != RiskLow {
		t.Fatalf("result = %+v", r)
	}
}

func TestScoreIsDeterministic(t *testing.T) {
	fs := []Finding{
		f("a", AxisTests, TypeTest, "x_test.go", 10),
		f("b", AxisArchitecture, TypeDuplication, "y.go", 2),
		f("c", AxisCorrectness, TypeBlocking, "z.go", 7),
	}
	in := RiskInput{ChangedFiles: []string{"z.go"}, DiffLines: 40, Familiar: true}
	first := Score(fs, DefaultPolicy(), in)
	for i := 0; i < 5; i++ {
		if got := Score(fs, DefaultPolicy(), in); got.Score != first.Score || got.Verdict != first.Verdict {
			t.Fatalf("run %d differs: %+v vs %+v", i, got, first)
		}
	}
	// 10 - 2 - 1.5 - 4 = 2.5
	if first.Score != 2.5 {
		t.Fatalf("score = %v, want 2.5", first.Score)
	}
	if first.ByAxis[AxisTests] != 2 || first.ByAxis[AxisCorrectness] != 4 {
		t.Fatalf("by axis = %+v", first.ByAxis)
	}
}

func TestScoreNeverGoesBelowZero(t *testing.T) {
	fs := []Finding{
		f("a", AxisCorrectness, TypeBlocking, "a.go", 1),
		f("b", AxisCorrectness, TypeBlocking, "b.go", 1),
		f("c", AxisCorrectness, TypeBlocking, "c.go", 1),
	}
	if got := Score(fs, DefaultPolicy(), RiskInput{Familiar: true}).Score; got != 0 {
		t.Fatalf("score = %v, want 0", got)
	}
}

func TestThresholdIsInclusive(t *testing.T) {
	p := DefaultPolicy()
	p.Threshold = 9 // one docs finding costs 0.5 -> 9.5 passes; two -> 9.0 blocks
	one := []Finding{f("a", AxisUI, TypeDocs, "README.md", 0)}
	if got := Score(one, p, RiskInput{Familiar: true}); got.Verdict != VerdictPass {
		t.Fatalf("9.5 against 9.0 should pass: %+v", got)
	}
	two := append(one, f("b", AxisUI, TypeDocs, "README.md", 0))
	if got := Score(two, p, RiskInput{Familiar: true}); got.Verdict != VerdictHuman {
		t.Fatalf("9.0 against 9.0 should ask for a human: %+v", got)
	}
}

func TestPolicyOverridesWeights(t *testing.T) {
	p := DefaultPolicy()
	p.Weights[TypeTest] = 0.2
	got := Score([]Finding{f("a", AxisTests, TypeTest, "x_test.go", 1)}, p, RiskInput{Familiar: true})
	if got.Score != 9.8 {
		t.Fatalf("score = %v, want 9.8", got.Score)
	}
}

func TestRiskTable(t *testing.T) {
	p := DefaultPolicy()
	p.SensitivePaths = []string{"internal/payments", "migrations"}
	p.LargeDiffLines = 400
	cases := []struct {
		name string
		in   RiskInput
		want string
	}{
		{"sensitive always asks", RiskInput{ChangedFiles: []string{"internal/payments/card.go"}, DiffLines: 5, Familiar: true}, RiskHigh},
		{"large and unfamiliar", RiskInput{ChangedFiles: []string{"a.go"}, DiffLines: 900, Familiar: false}, RiskHigh},
		{"large but familiar", RiskInput{ChangedFiles: []string{"a.go"}, DiffLines: 900, Familiar: true}, RiskMedium},
		{"small but unfamiliar", RiskInput{ChangedFiles: []string{"a.go"}, DiffLines: 20, Familiar: false}, RiskMedium},
		{"small and familiar", RiskInput{ChangedFiles: []string{"a.go"}, DiffLines: 20, Familiar: true}, RiskLow},
	}
	for _, c := range cases {
		if got, _ := Risk(c.in, p); got != c.want {
			t.Errorf("%s: risk %q, want %q", c.name, got, c.want)
		}
	}
}

func TestHighRiskBlocksWhateverTheScore(t *testing.T) {
	p := DefaultPolicy()
	p.SensitivePaths = []string{"internal/payments"}
	got := Score(nil, p, RiskInput{ChangedFiles: []string{"internal/payments/x.go"}, Familiar: true})
	if got.Score != 10 || got.Verdict != VerdictHuman || got.Risk != RiskHigh {
		t.Fatalf("result = %+v", got)
	}
}

func TestLoadPolicyFallsBackToDefaults(t *testing.T) {
	p, err := LoadPolicy(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if p.Source != "defaults" || p.Threshold != 7 || p.LargeDiffLines != 400 {
		t.Fatalf("policy = %+v", p)
	}
}
