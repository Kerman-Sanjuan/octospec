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
| `octospec lock status [--repo <path>]` | Report whether the lock that guards `openspec/specs/` is free, held (by which change, on which host, since when, and how long the lease has left) or stale. Read-only. |
| `octospec lock acquire --change <name> [--wait 2m] [--ttl 30m] [--repo <path>]` | Take the archive lock as a lease. Exit 3 when another change holds it and the wait expired, 4 when the lease has expired, 1 on any other error. |
| `octospec lock heartbeat --change <name> [--ttl 30m] [--repo <path>]` | Extend the lease you already hold. Fails on an expired lease instead of reviving it, because by then somebody else may hold the lock. |
| `octospec lock release --change <name> [--force] [--repo <path>]` | Release the lock. `--force` breaks a lease that is still valid, deliberately; nothing expires it early and nothing clears it automatically. |
| `octospec review --findings <file> [--json] [--enforce] [--changed-files a,b] [--diff-lines n] [--familiar]` | Score a pull request from a list of findings: 0-10, a risk level and a verdict. Reports by default; `--enforce` exits 5 when the verdict asks for a human. |
| `octospec result write --stage <s> --issue <n> --outcome ok\|stalled [--reason <r>] [--artifact <url>] [--change <c>] [--duration-ms <n>] [--sha <sha>] [--json]` | Write what one stage execution did, to `.octospec/runs/<stage>-<sha>.json`. The vocabulary is closed: `reason` is one of `timeout`, `crashed`, `empty_diff`, `no_pr_requested`, `needs_human`, `budget_exhausted`, and an `ok` from `ship` or `archive` must carry an artefact. |
| `octospec result show <file>` | Print one record, and exit non-zero when it does not validate. |
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
| `/review <#\|url>` | review | Read the pull request against its own change, write the findings as JSON, and publish one `## Review` comment with the score, the costed findings and the risk. |

Each command runs its stage agent. The agent per stage and its tool surface live
in `cli/internal/payload/agents/`. Set the model each role uses with
`octospec models`.

## The review policy

`octospec review` is arithmetic over findings, so the numbers live in the repository and tuning
them never needs a release. `octospec seed` writes `.octospec/review.json` with the defaults and
never overwrites an existing one.

| Key | Meaning | Default |
|---|---|---|
| `threshold` | Score at or below which the verdict asks for a human | `7` |
| `weights` | Points each finding type costs (`blocking`, `correctness`, `test`, `duplication`, `naming`, `style`, `docs`) | 4 / 3 / 2 / 1.5 / 0.5 / 0.5 / 0.5 |
| `sensitive_paths` | Path prefixes where any change asks for a human, whatever the score | empty |
| `large_diff_lines` | Added plus removed lines that count as a large diff | `400` |

Risk is a table, not a score: a sensitive path is always `high`; large *and* unfamiliar is `high`;
either one alone is `medium`; neither is `low`. Familiarity is a flag, because only the person
merging knows it.

Start in report-only mode — no `--enforce` — and read the first twenty scores before setting a
threshold. A gate configured from a guess is a gate somebody turns off.

## Labels

The workflow uses `type:feature`, `type:bug`, `status:backlog`, `status:spec-ready`, `status:in-progress`, and `status:in-review`. `octospec seed` provisions them.
