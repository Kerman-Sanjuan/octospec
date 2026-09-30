# Commands

octospec has two surfaces: the CLI, and the workflow commands you run inside your agent.

## CLI

| Command | Purpose |
|---|---|
| `octospec install [--tool <name>] [--repo <path>]` | Install the commands, and the OpenSpec skills, for each tool. |
| `octospec seed [--repo <path>] [--no-labels]` | Install the OpenSpec schema, the repo seed files, and the workflow labels. |
| `octospec update [--repo <path>]` | Re-apply an install from the saved state, keeping your local edits. |
| `octospec doctor [--repo <path>]` | Report what is missing, without changing anything. |
| `octospec version` | Print the version. |

## Tools

| Tool | Commands land in | Skills land in |
|---|---|---|
| pi | `~/.pi/agent/prompts/` | `.pi/skills/` |
| opencode | `~/.config/opencode/command/` | `.opencode/skills/` |
| GitHub Copilot | `.github/prompts/` | `.github/skills/` |
| Claude Code | `.claude/commands/` | `.claude/skills/` |

## Workflow commands

| Command | Phase | Behaviour |
|---|---|---|
| `/explore <topic>` | optional | Think through an idea. No files, no issue, no change. |
| `/idea [idea]` | intake | Interview you, then file the issue from the feature form. |
| `/bug [summary]` | intake | File a bug issue from the bug form. |
| `/spec <#\|url>` | plan | Turn the issue into a change, commit it on the branch, publish one `## Plan` comment, and validate. |
| `/apply [change]` | execute | Work the tasks on the branch, keeping one `## Implementation` comment. |
| `/ship [change]` | ship | Run CI, loop failures back to `/apply`, and open the PR once it is green. |
| `/archive [change]` | archive | Archive the change, sync the specs, commit and push, close the issue, and delete the branch. |

## Labels

The workflow uses `type:feature`, `type:bug`, `status:backlog`, `status:spec-ready`, `status:in-progress`, and `status:in-review`. `octospec seed` provisions them.
