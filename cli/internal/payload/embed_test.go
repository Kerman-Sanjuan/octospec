package payload

import "testing"

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
