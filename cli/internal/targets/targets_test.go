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
