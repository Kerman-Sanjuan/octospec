# octospec

**GitHub-native OpenSpec.** Issues are the backlog and the human layer; the repo
holds the machine artifacts; CI enforces the gates.

octospec is a small, inspectable workflow built on top of
[OpenSpec](https://github.com/Fission-AI/OpenSpec) that makes spec-driven
development GitHub-native. It adds only the GitHub behaviour OpenSpec does not
have: issues as the human-readable backlog, pull requests as the execution
surface, and hard validation in CI.

## Why

Most spec frameworks keep everything in the repository. That is great for the
machine, but it hides the backlog from the people who review it. octospec keeps
the narrative where people already are — GitHub Issues — and the machine
artifacts where OpenSpec can validate and archive them — files in the repo.

There is no dual write. Each artifact has exactly one authoritative home:

| Artifact | Home | Role |
|---|---|---|
| User story, requirements, success criteria, out of scope | GitHub **Issue** body | backlog + human source of truth |
| Proposal, specs, design, tasks | repo `openspec/changes/<name>/` | machine layer, committed with the branch |
| Published copies of the above | Issue **comments** | read-only mirrors, pinned to the commit SHA and refreshed in place on re-run |
| Change ↔ issue link | `.openspec.yaml` → `github.issue` | traceability |
| Living specs | `openspec/specs/` after archive | the contract of record |

## The loop

```
/idea    ──▶ issue in the backlog                       ← human refines on GitHub
              │
              ▼
/spec    ──▶ branch + committed change + comments    ← human approves
              │
              ▼
/apply   ──▶ tasks, commits, checklist synced
              │
              ▼
/ship    ──▶ PR (Closes #n) + CI gates                  ← human reviews and merges
              │
              ▼
/archive ──▶ specs updated, issue closed
```

1. **`/idea`** interviews you for a feature and files the issue. The issue body
   is the human source of truth.
2. **`/spec`** creates the change branch (`feat|fix/<issue>-<slug>`), commits
   the change artifacts on it, publishes them as a single `## Plan` comment
   headed with the commit SHA, and validates. `/spec` does not approve its own
   work.
3. **`/apply`** works the tasks on the branch `/spec` created, group by group,
   keeping a single `## Implementation` comment (the checklist plus a short
   insight per task) up to date.
4. **`/ship`** opens a pull request that closes the issue (the branch was
   pushed by `/spec`).
5. **`/archive`** archives the change, updates `openspec/specs/` (filling each
   capability's Purpose), closes the issue, and deletes the merged branch.

### The issue timeline

An issue carries exactly three structured layers, and nothing else:

| Layer | Source | Header |
|---|---|---|
| Requirement | issue body | — |
| Plan | `/spec` | `## Plan` (`### Proposal`, `### Capabilities`, `### Design`, `### Tasks`) |
| Implementation | `/apply` | `## Implementation` (checklist + a short insight per task) |

The plan and implementation comments are **updated in place** on every re-run
(their ids live in `.openspec.yaml`), so iterating never adds a new comment.

## Architecture

octospec extends OpenSpec rather than wrapping it. OpenSpec's own `/opsx:*`
commands and skills keep working; octospec adds a project-local schema and thin
command shims.

```
octospec/
├── cli/                    # the Go CLI: one canonical payload, rendered per tool
│   └── internal/
│       ├── payload/commands/  # the canonical commands (embedded)
│       ├── targets/           # tool -> directory + front matter
│       └── seed/              # OpenSpec schema + repo seed (embedded)
├── scripts/check-gates.sh  # gate checks, runnable locally and in CI
├── .github/workflows/      # cli (Go), openspec (gates), release
└── openspec/               # this repo's own changes and specs
```

- **`cli/`** is the Go CLI. It embeds one canonical command payload and renders
  it into each tool's directory - there is no per-tool copy in the repo.
- **`cli/internal/payload/commands/`** is the single source of the commands.
- **`cli/internal/targets/`** maps each tool to its directory and front matter.
- **`cli/internal/seed/`** holds the OpenSpec schema and the repo seed (issue
  forms, CI, config), embedded in the binary.
- **`scripts/check-gates.sh`** runs the gates locally and in CI.

### Gates

Hard gates block CI or merge; advisory gates warn. `scripts/check-gates.sh` runs
the gates that do not depend on the target stack (G2–G4, G6); a repository adds
its own build/test gate for G5.

| Gate | Checks | Strength |
|---|---|---|
| G1 | issue has the required sections before `/spec` produces artifacts | hard |
| G2 | `openspec validate --all --strict` and every unarchived change is complete | hard |
| G3 | branch matches `feat\|fix/<issue>-<slug>` | hard |
| G4 | PR body references the issue (`Closes #n`) | hard |
| G5 | tests / build / lint pass | hard |
| G6 | a PR touching `openspec/changes/<name>/` carries a spec delta | hard |
| G7 | all tasks checked before merge | advisory |
| G8 | one OpenSpec change per issue | advisory |

## Install

Requirements: [OpenSpec CLI](https://github.com/Fission-AI/OpenSpec) ≥ 1.3.1 and
[`gh`](https://cli.github.com) ≥ 2.x, authenticated (`gh auth status`).

Install the CLI (fetches the latest release, then runs `octospec install`):

```sh
curl -fsSL https://raw.githubusercontent.com/Kerman-Sanjuan/octospec/main/install.sh | sh
```

Contributors can instead `go install github.com/kerman-sanjuan/octospec/cli/cmd/octospec@latest`.

Install the commands for a tool - repeat `--tool`, or omit it for all:

```sh
octospec install --tool pi --tool opencode --tool copilot --tool claude --repo .
```

| Tool | Destination | Scope |
|---|---|---|
| pi | `~/.pi/agent/prompts/` | global |
| opencode | `~/.config/opencode/command/` | global |
| GitHub Copilot | `.github/prompts/` | repo-local |
| Claude Code | `.claude/commands/` | repo-local |

Seed a repository - OpenSpec schema, issue forms, CI, config, and the workflow
labels:

```sh
octospec seed --repo /path/to/repo
```

Update after a new release - re-applies and keeps your local edits:

```sh
octospec update
```

Adding a tool is a new entry in `cli/internal/targets` - never a new copy.
- OpenSpec skills (`openspec-propose`, `openspec-apply-change`,
  `openspec-archive-change`, `openspec-explore`) ship per tool.

## Commands

| Command | Phase | Behaviour |
|---|---|---|
| `/explore <topic>` | optional | Think through an idea. No files, no issue, no change. |
| `/idea [idea]` | intake | Interview for a feature and file the issue from the feature form. |
| `/bug [summary]` | intake | File a bug issue from the bug form. |
| `/spec <#\|url>` | plan | Turn an issue into an OpenSpec change, publish the artifacts to the issue, validate, and label `status:spec-ready`. |
| `/apply [change]` | execute | Branch, commit the change baseline, work the tasks group by group, and sync the checklist. |
| `/ship [change]` | ship | Push the branch and open a PR whose body closes the issue. |
| `/archive [change]` | archive | Archive the change, update `openspec/specs/`, and close the issue. |

## Deeper documentation

- [Design](docs/plans/2026-09-29-github-openspec-workflow-design.md) — decisions,
  artifact map, gates, and trade-offs.
- [Implementation plan](docs/plans/2026-09-29-github-openspec-workflow-plan.md) —
  the task-by-task build.

## License

MIT
