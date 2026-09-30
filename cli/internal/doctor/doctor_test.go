package doctor

import (
	"testing"

	"github.com/kerman-sanjuan/octospec/cli/internal/install"
)

func TestRunReportsMissingCommands(t *testing.T) {
	checks, err := Run(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, c := range checks {
		if c.Name == "commands" {
			found = true
			if c.OK {
				t.Fatal("an empty repo should report commands missing")
			}
		}
	}
	if !found {
		t.Fatal("expected a commands check")
	}
}

func TestRunReportsAgentsAndModels(t *testing.T) {
	dir := t.TempDir()
	if err := install.Run(install.Options{Tools: []string{"claude"}, Repo: dir}); err != nil {
		t.Fatal(err)
	}
	checks, err := Run(dir)
	if err != nil {
		t.Fatal(err)
	}
	var agents, models bool
	for _, c := range checks {
		if c.Name == "agents (claude)" {
			agents = c.OK
		}
		if c.Name == "models" {
			models = c.OK
		}
	}
	if !agents || !models {
		t.Fatalf("agents = %v, models = %v", agents, models)
	}
}
