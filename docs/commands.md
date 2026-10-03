# Commands

octospec has two surfaces: the CLI, and the workflow commands you run inside your agent.

## CLI

| Command | Purpose |
|---|---|
| `octospec install [--tool <name>] [--repo <path>] [--global] [--dry-run]` | Install the commands, the agents, and the OpenSpec skills for each tool. With no `--tool` on a terminal it opens a multi-select, pre-checked with the last install. Repo-local by default; `--global` warns and installs into the home directory. `--dry-run` prints the plan and changes nothing. |
| `octospec seed [--repo <path>] [--no-labels] [--dry-run]` | Install the OpenSpec schema, the repo seed files, and the workflow labels. `--dry-run` prints the plan and changes nothing. |
| `octospec update [--repo <path>]` | Re-apply an install from the saved state, keeping your local edits. |
| `octospec uninstall [--repo <path>]` | Remove the files octospec manages, keeping any you edited by hand. |
| `octospec doctor [--repo <path>]` | Report what is missing, without changing anything. |
| `octospec models [--set <role>=<model>] [--repo <path>]` | Set the model each agent role uses. Opens a TUI when no `--set` is given, listing the known models for the installed tools with a free-form fallback. |
| `octospec version` | Print the version. |

## Tools

| Tool | Repo-local (default) | Global | Skills |
|---|---|---|---|
| pi | `.pi/prompts/` | `~/.pi/agent/prompts/` | `.pi/skills/` |
| opencode | `.opencode/command/`, `.opencode/agents/` | `~/.config/opencode/command/`, `~/.config/opencode/agents/` | `.opencode/skills/` |
| GitHub Copilot | `.github/prompts/`, `.github/agents/` | `~/.copilot/prompts/`, `~/.copilot/agents/` | `.github/skills/` |
| Claude Code | `.claude/commands/`, `.claude/agents/` | `~/.claude/commands/`, `~/.claude/agents/` | `.claude/skills/` |

Repo-local is the default. A global install warns that the commands, agents, and
skills appear in every project. pi has no agent mechanism, so it installs no
agents. The scope per tool is recorded in `.octospec/octospec.json`.

## Workflow commands

| Command | Phase | Behaviour |
|---|---|---|
| `/idea [idea]` | intake | Interview you, then file the issue with those sections as the body. |
| `/bug [summary]` | intake | Interview you, then file the bug issue with those sections as the body. |
| `/spec <#\|url>` | plan | Turn the issue into a change, commit it on the branch, publish one `## Plan` comment, and validate. |
| `/apply [change]` | execute | Work the tasks on the branch, keeping one `## Implementation` comment. |
| `/ship [change]` | ship | Run CI, loop failures back to `/apply`, and open the PR once it is green. |
| `/archive [change]` | archive | Archive the change, sync the specs, commit and push, close the issue, and delete the branch. |

Each command runs its stage agent. The agent per stage and its tool surface live
in `cli/internal/payload/agents/`. Set the model each role uses with
`octospec models`.

## Labels

The workflow uses `type:feature`, `type:bug`, `status:backlog`, `status:spec-ready`, `status:in-progress`, and `status:in-review`. `octospec seed` provisions them.
