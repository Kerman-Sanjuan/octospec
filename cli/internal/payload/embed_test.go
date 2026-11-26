package payload

import (
	"strings"
	"testing"
)

func TestParseFrontMatter(t *testing.T) {
	c := Parse("x", "---\ndescription: d\nargument-hint: <a>\n---\n\nhello\n")
	if c.Name != "x" || c.Description != "d" || c.ArgumentHint != "<a>" {
		t.Fatalf("front matter: %+v", c)
	}
	if c.Body != "hello\n" {
		t.Fatalf("body = %q", c.Body)
	}
}

func TestCommandsEmbedded(t *testing.T) {
	cmds, err := Commands()
	if err != nil {
		t.Fatal(err)
	}
	if len(cmds) != 7 {
		t.Fatalf("want 7 commands, got %d", len(cmds))
	}
	for _, c := range cmds {
		if c.Description == "" || c.Body == "" {
			t.Errorf("command %q not parsed", c.Name)
		}
	}
}

func TestParseAgent(t *testing.T) {
	a := ParseAgent("spec", "---\nname: spec\nstage: spec\nrole: thinking\ndescription: d\nskills: a, b\ntools: read, shell\nwrites: w\n---\n\nbody\n")
	if a.Name != "spec" || a.Stage != "spec" || a.Role != "thinking" || a.Description != "d" || a.Writes != "w" {
		t.Fatalf("agent: %+v", a)
	}
	if len(a.Skills) != 2 || a.Skills[0] != "a" || a.Skills[1] != "b" {
		t.Fatalf("skills = %v", a.Skills)
	}
	if len(a.Tools) != 2 || a.Tools[0] != "read" || a.Tools[1] != "shell" {
		t.Fatalf("tools = %v", a.Tools)
	}
	if a.Body != "body\n" {
		t.Fatalf("body = %q", a.Body)
	}
}

func TestAgentsEmbedded(t *testing.T) {
	agents, err := Agents()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"idea", "spec", "apply", "ship", "review", "archive"}
	if len(agents) != len(want) {
		t.Fatalf("want %d agents, got %d", len(want), len(agents))
	}
	for i, a := range agents {
		if a.Stage != want[i] {
			t.Errorf("agent %d stage = %q, want %q", i, a.Stage, want[i])
		}
		if a.Name == "" || a.Role == "" || a.Description == "" || a.Writes == "" || len(a.Tools) == 0 || a.Body == "" {
			t.Errorf("agent %q not fully parsed", a.Name)
		}
	}
}

func TestValidateRejectsAgentWithoutTools(t *testing.T) {
	a := Agent{Name: "x", Stage: "idea", Role: "thinking", Description: "d", Writes: "w", Body: "b"}
	if err := a.Validate(); err == nil {
		t.Fatal("expected an error when the tool surface is empty")
	}
}

func TestCommandsDoNotUseTemplateWithBody(t *testing.T) {
	cmds, err := Commands()
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range cmds {
		if strings.Contains(c.Body, "--template") {
			t.Errorf("command %q uses --template; gh rejects it together with --body", c.Name)
		}
	}
}

func commandBody(t *testing.T, name string) string {
	t.Helper()
	cmds, err := Commands()
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range cmds {
		if c.Name == name {
			return c.Body
		}
	}
	t.Fatalf("command %q not found", name)
	return ""
}

func agentBody(t *testing.T, stage string) string {
	t.Helper()
	agents, err := Agents()
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range agents {
		if a.Stage == stage {
			return a.Body
		}
	}
	t.Fatalf("agent for stage %q not found", stage)
	return ""
}

// The intake grills the decision frontier before filing. Without the frontier
// the command is an interview for five form fields, which is what it was.
func TestIdeaCommandGrillsTheFrontier(t *testing.T) {
	body := commandBody(t, "idea")
	for _, want := range []string{"frontier", "next round", "## Grill", "File nothing while a decision"} {
		if !strings.Contains(body, want) {
			t.Errorf("idea command does not carry %q", want)
		}
	}
	if strings.Contains(body, "one question at a time") {
		t.Error("idea command still asks one question at a time")
	}
}

// Every question carries a recommendation, so "yes" is a valid answer.
func TestIdeaCommandRecommendsAnAnswer(t *testing.T) {
	body := commandBody(t, "idea")
	if !strings.Contains(body, "recommended answer") {
		t.Error("idea command asks without a recommendation")
	}
	if !strings.Contains(body, `"yes" accepts it`) {
		t.Error("idea command does not say that yes accepts the recommendation")
	}
}

// Facts are looked up; only decisions are asked.
func TestIdeaAgentLooksFactsUpAndAsksDecisions(t *testing.T) {
	body := agentBody(t, "idea")
	if !strings.Contains(body, "decision tree") || !strings.Contains(body, "frontier") {
		t.Error("idea agent does not describe the decision tree and its frontier")
	}
	if !strings.Contains(body, "Look facts up, ask for decisions") {
		t.Error("idea agent does not separate facts from decisions")
	}
}

// The frontier covers who, the outcome, done, the failures, the boundary, the
// constraints and the evidence.
func TestIdeaAgentCoversTheSevenGroups(t *testing.T) {
	body := agentBody(t, "idea")
	for _, want := range []string{"who", "observable outcome", "done", "failures", "boundary", "constraints", "evidence"} {
		if !strings.Contains(body, want) {
			t.Errorf("idea agent does not cover %q", want)
		}
	}
}

// With nobody to answer, the grill stops instead of inventing an answer.
func TestIdeaAgentStopsUnattended(t *testing.T) {
	body := agentBody(t, "idea")
	if !strings.Contains(body, "needs_human") {
		t.Error("idea agent does not stop with needs_human when nobody answers")
	}
}

// /spec reads the intake's record instead of re-asking it.
func TestSpecCommandReadsTheGrill(t *testing.T) {
	body := commandBody(t, "spec")
	if !strings.Contains(body, "## Grill") {
		t.Error("spec command does not read the Grill section")
	}
	if !strings.Contains(body, "not asked again") {
		t.Error("spec command does not say it must not re-ask a recorded decision")
	}
}

