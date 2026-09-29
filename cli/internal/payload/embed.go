// Package payload holds the canonical octospec command files, embedded in the
// binary so there is exactly one source and no per-tool copies in the repo.
package payload

import (
	"embed"
	"io/fs"
	"sort"
	"strings"
)

//go:embed commands/*.md
var commands embed.FS

// Command is one canonical command: its name (without extension), the parsed
// front matter, and the body.
type Command struct {
	Name         string
	Description  string
	ArgumentHint string
	Body         string
}

// Commands returns every canonical command, sorted by name.
func Commands() ([]Command, error) {
	entries, err := fs.ReadDir(commands, "commands")
	if err != nil {
		return nil, err
	}
	var out []Command
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		b, err := commands.ReadFile("commands/" + e.Name())
		if err != nil {
			return nil, err
		}
		out = append(out, Parse(strings.TrimSuffix(e.Name(), ".md"), string(b)))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// Parse splits a canonical command into front-matter fields and body.
func Parse(name, content string) Command {
	c := Command{Name: name}
	lines := strings.Split(content, "\n")
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		i := 1
		for ; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) == "---" {
				break
			}
			if k, v, ok := strings.Cut(lines[i], ":"); ok {
				switch strings.TrimSpace(k) {
				case "description":
					c.Description = strings.TrimSpace(v)
				case "argument-hint":
					c.ArgumentHint = strings.TrimSpace(v)
				}
			}
		}
		if i < len(lines) {
			body := lines[i+1:]
			if len(body) > 0 && strings.TrimSpace(body[0]) == "" {
				body = body[1:]
			}
			c.Body = strings.Join(body, "\n")
			return c
		}
	}
	c.Body = content
	return c
}
