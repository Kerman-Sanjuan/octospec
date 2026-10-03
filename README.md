<p align="center">
  <img src="docs/logo.svg" alt="octospec logo" width="140" height="140">
</p>

# octospec

[![cli](https://github.com/Kerman-Sanjuan/octospec/actions/workflows/cli.yml/badge.svg)](https://github.com/Kerman-Sanjuan/octospec/actions/workflows/cli.yml)
[![openspec](https://github.com/Kerman-Sanjuan/octospec/actions/workflows/openspec.yml/badge.svg)](https://github.com/Kerman-Sanjuan/octospec/actions/workflows/openspec.yml)
[![release](https://img.shields.io/github/v/release/Kerman-Sanjuan/octospec)](https://github.com/Kerman-Sanjuan/octospec/releases)
[![license](https://img.shields.io/github/license/Kerman-Sanjuan/octospec)](LICENSE)
[![go](https://img.shields.io/github/go-mod/go-version/Kerman-Sanjuan/octospec/main?filename=cli/go.mod)](cli/go.mod)

**GitHub-native OpenSpec.** Issues are the backlog and the human layer; the repo
holds the machine artifacts; CI enforces the gates.

octospec is a small, inspectable workflow built on top of
[OpenSpec](https://github.com/Fission-AI/OpenSpec) that makes spec-driven
development GitHub-native. It adds only the GitHub behaviour OpenSpec does not
have: issues as the human-readable backlog, pull requests as the execution
surface, and hard validation in CI.

## Quickstart

Prerequisites: the [OpenSpec CLI](https://github.com/Fission-AI/OpenSpec) 1.3.1
(installed with npm, so Node.js is needed), and
the [`gh` CLI](https://cli.github.com) 2.x, authenticated. octospec is tested
against 1.3.1; newer versions reject the `github:` key that octospec writes to
`.openspec.yaml`.

```sh
# 1. install the CLI
curl -fsSL https://raw.githubusercontent.com/Kerman-Sanjuan/octospec/main/install.sh | sh

# 2. in your repository: install the commands and seed the workflow
octospec install --tool claude    # or pi, opencode, copilot
octospec seed

# 3. in your agent, run the loop: /idea -> /spec -> /apply -> /ship -> /archive
```

See [Install](#install) for every tool and every channel.

## Why

Most spec frameworks keep everything in the repository. That is great for the
machine, but it hides the backlog from the people who review it. octospec keeps
the narrative where people already are (GitHub Issues) and the machine
artifacts where OpenSpec can validate and archive them (files in the repo).

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
| Requirement | issue body | (none) |
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
│       ├── payload/agents/    # the canonical stage agents (embedded)
│       ├── targets/           # tool -> directory, front matter, tools
│       ├── models/            # the role-to-model TUI
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
the gates that do not depend on the target stack (G1 to G4, and G6 to G8); a
repository adds its own build/test gate for G5.

| Gate | Checks | Strength |
|---|---|---|
| G1 | issue has the required sections before `/spec` produces artifacts | hard |
| G2 | `openspec validate --all --strict` and every unarchived change is complete | hard |
| G3 | branch matches `feat\|fix/<issue>-<slug>` (PR branches only) | hard |
| G4 | PR body references the issue (`Closes #n`) | hard |
| G5 | tests / build / lint pass | hard |
| G6 | a PR touching `openspec/changes/<name>/` carries a spec delta | hard |
| G7 | all tasks checked before merge | advisory |
| G8 | one OpenSpec change per issue | advisory |

## Install

Requirements: [OpenSpec CLI](https://github.com/Fission-AI/OpenSpec) 1.3.1
(installed with npm, so Node.js is needed) and
[`gh`](https://cli.github.com) ≥ 2.x, authenticated (`gh auth status`).

Install the CLI (fetches the latest release, then runs `octospec install`):

```sh
curl -fsSL https://raw.githubusercontent.com/Kerman-Sanjuan/octospec/main/install.sh | sh
```

Contributors can instead `go install github.com/kerman-sanjuan/octospec/cli/cmd/octospec@latest`.

### Platforms and channels

| Platform | `curl \| sh` | `go install` | Release binary |
|---|---|---|---|
| linux amd64 | yes | yes | `octospec_linux_amd64` |
| linux arm64 | yes | yes | `octospec_linux_arm64` |
| darwin amd64 | yes | yes | `octospec_darwin_amd64` |
| darwin arm64 | yes | yes | `octospec_darwin_arm64` |

These are the channels: the `curl | sh` installer, `go install`, and the release binaries. There is no package manager channel yet.

Windows is not supported natively. WSL works, because it is Linux, so the linux binary and the `sh` installer run there. There is no Windows release binary or PowerShell installer.

The installer verifies the release checksum and accepts a version to pin:

```sh
curl -fsSL https://raw.githubusercontent.com/Kerman-Sanjuan/octospec/main/install.sh | sh -s -- --version v1.1.0
```

Install the commands for a tool - repeat `--tool`, or omit it for all. They
install repo-local by default:

```sh
octospec install --tool pi --tool opencode --tool copilot --tool claude --repo .
```

Run `octospec install` with no `--tool` on a terminal and it opens a
multi-select of the tools, pre-checked with the last install. Space toggles a
tool, enter confirms, and choosing none installs for all of them.

| Tool | Repo-local (default) | Global |
|---|---|---|
| pi | `.pi/prompts/` | `~/.pi/agent/prompts/` |
| opencode | `.opencode/command/`, `.opencode/agents/` | `~/.config/opencode/command/`, `~/.config/opencode/agents/` |
| GitHub Copilot | `.github/prompts/`, `.github/agents/` | `~/.copilot/prompts/`, `~/.copilot/agents/` |
| Claude Code | `.claude/commands/`, `.claude/agents/` | `~/.claude/commands/`, `~/.claude/agents/` |

A repo-local install writes nothing under your home directory. A global install
warns that the commands, agents, and skills then appear in every project. Choose
it with `--global`, or answer the prompt on a terminal:

```sh
octospec install --tool opencode --global
```

pi has no agent mechanism, so it installs commands and skills only. The scope per
tool is recorded in `.octospec/octospec.json`, so `octospec update` re-applies it,
and it offers to remove the files that earlier versions installed globally.

Seed a repository - OpenSpec schema, issue forms, CI, config, and the workflow
labels:

```sh
octospec seed --repo /path/to/repo
```

Update after a new release - re-applies and keeps your local edits:

```sh
octospec update
```

Check an install, and see what is missing:

```sh
octospec doctor
```

Adding a tool is a new entry in `cli/internal/targets` - never a new copy.

`install` also brings the OpenSpec skills (`openspec-propose`, and the others)
into each tool's skill directory by calling `openspec init --tools`.

OpenSpec ships its own skills (`openspec-propose`, `openspec-apply-change`,
`openspec-archive-change`, `openspec-explore`). Run `openspec update` to
(re)generate them per tool. They are not committed.

## Agents and models

Each stage runs on a scoped agent. The agents live once in
`cli/internal/payload/agents/` and render into each tool's agent directory.

| Stage | Agent | Role | Tool surface |
|---|---|---|---|
| intake | `idea` | thinking | read, search, shell |
| plan | `spec` | thinking | read, search, edit, shell |
| execute | `apply` | implementer | read, search, edit, shell |
| ship | `ship` | reviewer | read, search, shell |
| archive | `archive` | reviewer | read, search, edit, shell |

The tool surface maps to each tool's native permission field. Claude, opencode,
and Copilot enforce it. pi carries it as advice because it has no permission
field, and Codex is not a target because it has no agent file.

Pick the model each role uses with the TUI, or set one directly:

```sh
octospec models
octospec models --set thinking=sonnet --set implementer=sonnet
```

The TUI lists the known models for your installed tools, and `Custom...` lets
you type any other identifier. `(tool default)` leaves the role without a
model.

The model per role lives in `.octospec/octospec.json`. An empty model means the
tool default, and `octospec update` re-renders the agents with the new value.

## Commands

| Command | Phase | Behaviour |
|---|---|---|
| `/idea [idea]` | intake | Interview for a feature and file the issue with those sections as the body. |
| `/bug [summary]` | intake | Interview for a bug and file the issue with those sections as the body. |
| `/spec <#\|url>` | plan | Turn an issue into a change, commit it on `feat\|fix/<issue>-<slug>`, publish one `## Plan` comment, validate, and label `status:spec-ready`. |
| `/apply [change]` | execute | Work the tasks on the branch `/spec` created, group by group, keeping one `## Implementation` comment. |
| `/ship [change]` | ship | Run CI; on failure, loop back to `/apply`; open the PR once the gates are green. |
| `/archive [change]` | archive | Archive the change, sync `openspec/specs/` (filling each Purpose), commit and push, close the issue, and delete the merged branch. |

## Deeper documentation

- [Getting started](docs/getting-started.md) - install octospec and run one change.
- [Contributing](CONTRIBUTING.md) - how to send a change to octospec.
- [Commands](docs/commands.md) - the CLI and the workflow commands.
- [Troubleshooting](docs/troubleshooting.md) - the common failures.
- [Releasing](docs/releasing.md) - how a version is cut.
- [Changelog](CHANGELOG.md) - what changed, grouped by impact.
- [About](https://github.com/Kerman-Sanjuan/octospec) - the project at a glance:
  description, topics, and releases.
- [Design](docs/plans/2026-09-29-github-openspec-workflow-design.md): decisions,
  artifact map, gates, and trade-offs.
- [Implementation plan](docs/plans/2026-09-29-github-openspec-workflow-plan.md):
  the task-by-task build.

## License

MIT
