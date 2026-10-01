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
	if len(cmds) != 6 {
		t.Fatalf("want 6 commands, got %d", len(cmds))
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
	want := []string{"idea", "spec", "apply", "ship", "archive"}
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
