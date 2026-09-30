package skills

import "testing"

func TestNames(t *testing.T) {
	got := Names([]string{"pi", "opencode", "copilot", "claude"})
	want := []string{"pi", "opencode", "github-copilot", "claude"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestNamesDropsUnknown(t *testing.T) {
	got := Names([]string{"nope"})
	if len(got) != 0 {
		t.Fatalf("want empty, got %v", got)
	}
}
