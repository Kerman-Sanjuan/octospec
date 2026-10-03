package models

import "github.com/kerman-sanjuan/octospec/cli/internal/targets"

// catalog holds the model identifiers the TUI offers for each tool. It is a
// convenience, not a gate: the role form always allows a model entered by hand,
// so a stale list never blocks a valid model. Update it here as providers ship
// models.
var catalog = map[string][]string{
	"claude": {
		"sonnet",
		"opus",
		"haiku",
	},
	"opencode": {
		"anthropic/claude-sonnet-4",
		"anthropic/claude-opus-4",
		"openai/gpt-4o",
		"google/gemini-2.5-pro",
	},
	"copilot": {
		"gpt-4o",
		"gpt-4.1",
		"claude-sonnet-4",
	},
	"pi": {
		"anthropic/claude-sonnet-4",
		"anthropic/claude-opus-4",
		"openai/gpt-4o",
		"google/gemini-2.5-pro",
	},
}

// Catalog returns the known models for the given tools, in tool order and
// de-duplicated. An unknown tool is skipped. An empty tool list returns the
// models for every tool, in target order.
func Catalog(tools []string) []string {
	if len(tools) == 0 {
		for _, t := range targets.Targets {
			tools = append(tools, t.Name)
		}
	}
	var out []string
	seen := map[string]bool{}
	for _, tool := range tools {
		for _, m := range catalog[tool] {
			if !seen[m] {
				seen[m] = true
				out = append(out, m)
			}
		}
	}
	return out
}
