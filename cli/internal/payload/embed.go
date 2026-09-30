// Package payload holds the canonical octospec command and agent files,
// embedded in the binary so there is exactly one source and no per-tool copies
// in the repo.
package payload

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

//go:embed commands/*.md
var commands embed.FS

//go:embed agents/*.md
var agents embed.FS

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

// Agent is one canonical stage agent: its identity, its model role, the skills
// it may load, and its mission body.
type Agent struct {
	Name        string
	Stage       string
	Role        string
	Description string
	Skills      []string
	Writes      string
	Body        string
}

// stageOrder fixes the order of the loop.
var stageOrder = map[string]int{"idea": 0, "spec": 1, "apply": 2, "ship": 3, "archive": 4}

// agentRoles is the set of roles the model configuration understands. It
// mirrors config.Roles, which cannot be imported here without a cycle.
var agentRoles = map[string]bool{"thinking": true, "implementer": true, "reviewer": true}

// Agents returns every canonical agent, ordered by stage.
func Agents() ([]Agent, error) {
	entries, err := fs.ReadDir(agents, "agents")
	if err != nil {
		return nil, err
	}
	var out []Agent
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") || e.Name() == "README.md" {
			continue
		}
		b, err := agents.ReadFile("agents/" + e.Name())
		if err != nil {
			return nil, err
		}
		a := ParseAgent(strings.TrimSuffix(e.Name(), ".md"), string(b))
		if err := a.Validate(); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool {
		if stageOrder[out[i].Stage] != stageOrder[out[j].Stage] {
			return stageOrder[out[i].Stage] < stageOrder[out[j].Stage]
		}
		return out[i].Name < out[j].Name
	})
	return out, nil
}

// ParseAgent splits a canonical agent into front-matter fields and body.
func ParseAgent(name, content string) Agent {
	a := Agent{Name: name}
	lines := strings.Split(content, "\n")
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		i := 1
		for ; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) == "---" {
				break
			}
			k, v, ok := strings.Cut(lines[i], ":")
			if !ok {
				continue
			}
			v = strings.TrimSpace(v)
			switch strings.TrimSpace(k) {
			case "name":
				a.Name = v
			case "stage":
				a.Stage = v
			case "role":
				a.Role = v
			case "description":
				a.Description = v
			case "skills":
				a.Skills = splitList(v)
			case "writes":
				a.Writes = v
			}
		}
		if i < len(lines) {
			body := lines[i+1:]
			if len(body) > 0 && strings.TrimSpace(body[0]) == "" {
				body = body[1:]
			}
			a.Body = strings.Join(body, "\n")
			return a
		}
	}
	a.Body = content
	return a
}

// Validate reports whether the agent has every required field.
func (a Agent) Validate() error {
	if a.Name == "" || a.Stage == "" || a.Role == "" || a.Description == "" || a.Writes == "" {
		return fmt.Errorf("agent %q is missing a required field", a.Name)
	}
	if _, ok := stageOrder[a.Stage]; !ok {
		return fmt.Errorf("agent %q has unknown stage %q", a.Name, a.Stage)
	}
	if !agentRoles[a.Role] {
		return fmt.Errorf("agent %q has unknown role %q", a.Name, a.Role)
	}
	if strings.TrimSpace(a.Body) == "" {
		return fmt.Errorf("agent %q has an empty body", a.Name)
	}
	return nil
}

func splitList(v string) []string {
	var out []string
	for _, p := range strings.Split(v, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
