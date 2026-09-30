package doctor

import "testing"

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
