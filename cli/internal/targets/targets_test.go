package targets

import (
	"strings"
	"testing"

	"github.com/kerman-sanjuan/octospec/cli/internal/payload"
)

func TestRenderDefault(t *testing.T) {
	tp, _ := Get("pi")
	c := payload.Command{Name: "spec", Description: "d", ArgumentHint: "<x>", Body: "BODY"}
	got := tp.Render(c)
	if !strings.HasPrefix(got, "---\ndescription: d\nargument-hint: <x>\n---\n\nBODY") {
		t.Fatalf("unexpected render:\n%s", got)
	}
}

func TestRenderCopilot(t *testing.T) {
	tp, _ := Get("copilot")
	c := payload.Command{Name: "spec", Description: "d", ArgumentHint: "<x>", Body: "BODY"}
	got := tp.Render(c)
	if !strings.HasPrefix(got, "---\nmode: agent\ndescription: d\n---\n\nBODY") {
		t.Fatalf("unexpected render:\n%s", got)
	}
	if strings.Contains(got, "argument-hint") {
		t.Fatal("copilot should not carry argument-hint")
	}
}

func TestFilenames(t *testing.T) {
	want := map[string]string{"pi": "spec.md", "opencode": "spec.md", "copilot": "spec.prompt.md", "claude": "spec.md"}
	for name, w := range want {
		tp, _ := Get(name)
		if got := tp.Filename("spec"); got != w {
			t.Errorf("%s filename = %q, want %q", name, got, w)
		}
	}
}

func TestScopes(t *testing.T) {
	pi, _ := Get("pi")
	if pi.Scope != Global {
		t.Error("pi should be global")
	}
	cp, _ := Get("copilot")
	if cp.Scope != RepoLocal {
		t.Error("copilot should be repo-local")
	}
}

func TestRenderAgentClaude(t *testing.T) {
	tp, _ := Get("claude")
	a := payload.Agent{Name: "spec", Description: "d", Tools: []string{"read", "edit", "shell"}, Body: "BODY"}
	got := tp.RenderAgent(a, "sonnet")
	if !strings.HasPrefix(got, "---\nname: spec\ndescription: d\nmodel: sonnet\ntools: Read, Glob, Grep, Edit, Write, Bash\n---\n\nBODY") {
		t.Fatalf("unexpected render:\n%s", got)
	}
}

func TestRenderAgentCopilot(t *testing.T) {
	tp, _ := Get("copilot")
	a := payload.Agent{Name: "ship", Description: "d", Tools: []string{"read", "shell"}, Body: "BODY"}
	got := tp.RenderAgent(a, "")
	if !strings.HasPrefix(got, "---\nmode: agent\ndescription: d\ntools: read, execute\n---\n\nBODY") {
		t.Fatalf("unexpected render:\n%s", got)
	}
	if strings.Contains(got, "model:") {
		t.Fatal("an empty model should be omitted")
	}
}

func TestRenderAgentOpencode(t *testing.T) {
	tp, _ := Get("opencode")
	a := payload.Agent{Name: "apply", Description: "d", Tools: []string{"read", "edit", "shell"}, Body: "BODY"}
	got := tp.RenderAgent(a, "")
	if !strings.Contains(got, "permissions:\n  - read\n  - edit\n  - bash\n") {
		t.Fatalf("unexpected render:\n%s", got)
	}
	if strings.Contains(got, "tools:") {
		t.Fatal("opencode should not carry a tools field")
	}
}

func TestRenderAgentAdvisory(t *testing.T) {
	tp, _ := Get("pi")
	a := payload.Agent{Name: "idea", Description: "d", Tools: []string{"read", "search", "shell"}, Body: "BODY"}
	got := tp.RenderAgent(a, "")
	if !strings.Contains(got, "Tools: read, search, shell\n\nBODY") {
		t.Fatalf("unexpected render:\n%s", got)
	}
	if strings.Contains(got, "tools:") {
		t.Fatal("pi has no tool field, so the scope must stay in the body")
	}
}

func TestAgentFilenames(t *testing.T) {
	want := map[string]string{"pi": "spec.md", "opencode": "spec.md", "copilot": "spec.agent.md", "claude": "spec.md"}
	for name, w := range want {
		tp, _ := Get(name)
		if got := tp.AgentFilename("spec"); got != w {
			t.Errorf("%s agent filename = %q, want %q", name, got, w)
		}
	}
}
