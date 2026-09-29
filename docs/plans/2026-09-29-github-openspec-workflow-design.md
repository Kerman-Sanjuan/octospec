# GitHub-native OpenSpec workflow: design

- **Date:** 2026-09-29
- **Status:** approved in brainstorming, ready for implementation planning
- **Owner:** kerman
- **Target agents:** pi, opencode, GitHub Copilot
- **Reference:** [GH-600 study guide](https://learn.microsoft.com/en-us/credentials/certifications/resources/study-guides/gh-600) and the *Developing in Agentic AI Systems* learning path (parts [1](https://learn.microsoft.com/en-us/training/paths/gh-developing-agentic-systems-1/) and [2](https://learn.microsoft.com/en-us/training/paths/github-agentic-systems-part-two/))

## Why

The previous workflow (Plain Concepts `agent-harness`, removed) carried a lot of framework for a single developer and kept its GitHub logic in a large set of `pc-*` skills. We want the opposite: a small, inspectable workflow that is **based on OpenSpec** (its CLI, schema and skills) and adds only the GitHub behaviour we need - issues as the human-readable backlog, pull requests as the execution surface, and hard validation in CI.

The goal is to follow the OpenSpec workflow using GitHub's capabilities, and to do it in a way that satisfies GH-600's principles: plan before act, durable inspectable artifacts, human-in-the-loop at risky transitions, least privilege, and guardrails as enforced controls rather than prose.

## Goals

1. GitHub Issues are the backlog and the human-readable home of the narrative artifacts: user story, requirements, contract, success criteria.
2. The repository holds the machine artifacts (`openspec/`), committed with the branch.
3. A single OpenSpec workflow drives everything, extended, not wrapped, so `/opsx:*` and our commands behave the same.
4. The same workflow works in pi, opencode and GitHub Copilot.
5. CI enforces hard gates; nothing is "fixed" by convention alone.

## Non-goals

- Replacing the OpenSpec CLI with `gh`. The CLI owns artifact scaffolding, status, validation and archiving.
- Bidirectional real-time sync. GitHub and the repo are each canonical for a defined set of artifacts; there is no merge.
- GitHub Projects boards. Issue labels and PR links are enough.
- Storing specs in GitHub Issues. Specs are contracts that must be files in the repo for validation and archiving.
- Multi-agent orchestration. This is a single-developer workflow; GH-600's multi-agent material is reference, not scope.

## Principle: extend the baseline, do not wrap it

OpenSpec resolves schemas by first match:

1. `<repo>/openspec/schemas/<name>/schema.yaml` (project-local)
2. `~/.local/share/openspec/schemas/<name>/schema.yaml` (user-global override)
3. the package's built-in `spec-driven`

The generated `/opsx:*` commands and the OpenSpec skills are **generic drivers**: they call `openspec new`, `openspec status`, `openspec instructions <artifact>`, and then follow whatever `instruction` and `template` the schema returns. Editing the schema therefore changes the behaviour of OpenSpec's own drivers too.

Decision: add a project-local schema `octospec` and point `openspec/config.yaml` at it. Do **not** shadow `spec-driven`; it stays as a fallback and a reference. Do **not** edit generated command or skill files - `openspec update` overwrites them.

## Decision: source of truth is split by artifact

- **GitHub owns the narrative.** The issue body is the only home of the user story, requirements, contract and success criteria. It is written and edited on GitHub.
- **The repo owns the machine artifacts.** `openspec/changes/<name>/` holds proposal, specs, design and tasks. They are committed with the branch and are what the CLI validates and archives.
- **Published copies.** The machine artifacts are published to the issue as comments so the human never has to leave GitHub to review. Published copies are read-only mirrors, not a second source of truth.
- **No dual write.** Each artifact has exactly one authoritative home. This removes the staleness risk the previous design accepted.

## Decision: artifact map

| Artifact | Home | Role |
|---|---|---|
| User story, requirements, success criteria, out-of-scope | GitHub **Issue** body (issue form) | backlog + human source of truth |
| Issue snapshot (`issue.md`) | repo `openspec/changes/<name>/` | durable record of the intake, read by later artifacts |
| Proposal / specs / design / tasks | repo `openspec/changes/<name>/` | machine layer, committed with the branch |
| Published copies of the above | Issue **comments** | read-only, for review and discussion |
| Change ↔ issue link | `.openspec.yaml` → `github.issue`, `github.issue_updated_at` | traceability and drift detection |
| Branch | `feat/<issue>-<slug>` or `fix/<issue>-<slug>` | scope |
| PR | body with `Closes #<issue>` + summary + validation | review and merge gate |
| Living specs | `openspec/specs/` after archive | the contract of record |

Change folder names must start with a letter (OpenSpec rule), so the issue number is prefixed: `gh-<issue>-<slug>` (e.g. `gh-42-add-dark-mode`). Branch names have no such rule and stay `feat\|fix/<issue>-<slug>`.

## Decision: command surface

| Command | Phase | Behaviour |
|---|---|---|
| `/idea` | intake | Interview for a feature. Create the issue from the feature form, label it, return the URL. |
| `/bug` | intake | Create a bug issue from the bug form. |
| `/explore <topic>` | optional | Think before an issue exists. No files, no issue. |
| `/spec <#\|url>` | plan | Load `openspec-propose`. Ensure the issue exists (create from the form if the change started without one), snapshot it to `issue.md`, build the artifacts, publish them to the issue, run `openspec validate`, set the status label. |
| `/apply [change]` | execute | Load `openspec-apply-change`. Create the branch, work tasks in order, commit per task group, tick the boxes and sync the checklist to the issue. |
| `/ship` | ship | Push the branch and open the PR with `Closes #<issue>`, a summary and the validation results. |
| `/archive [change]` | archive | Load `openspec-archive-change`. Archive the change, update `openspec/specs/`, close the issue. |

`/opsx:*` commands are **not** generated (`delivery: skills`). The OpenSpec skills remain, and the three alias commands load them. `/idea` and `/ship` are the only logic that is not OpenSpec's, because both sit outside OpenSpec's artifact graph.

### Where logic lives

- **Schema = artifact semantics.** What each artifact must contain, its dependency order, and the intake artifact. Single-sourced here.
- **Commands = phase orchestration.** Creating the issue, publishing copies, branching, opening the PR, closing the issue. Short and phase-scoped.
- **No artifact content is duplicated in commands.** A command never re-describes what a proposal or a task list should contain.

## The loop

```
/idea  ──▶ issue in the backlog          ← human refines on GitHub
              │
              ▼
/spec   ──▶ change files + issue comments + validate    ← human approves (label)
              │
              ▼
/apply  ──▶ branch, tasks, commits, checklist synced
              │
              ▼
/ship   ──▶ PR (Closes #n) + CI gates    ← human reviews and merges
              │
              ▼
/archive ─▶ specs updated, issue closed
```

## Schema: `octospec`

Artifact graph:

```
issue ──▶ proposal ──▶ specs ──▶ design ──▶ tasks ──▶ (apply)
                └──────────────┘
```

- **`issue`**: generates `issue.md`. Instruction: ensure a GitHub issue exists for this change (adopt an existing one, or create it from the issue form), record `github.issue` and `github.issue_updated_at` in `.openspec.yaml`, and snapshot the body to `issue.md`. Requires nothing. This is what lets `/spec` run from either a fresh idea or an existing issue URL.
- **`proposal`**: generates `proposal.md`. Requires `issue`. Instruction: derive the proposal from the issue snapshot; the issue keeps the narrative, the proposal keeps the change summary the CLI archives.
- **`specs`**: generates `specs/**/*.md`. Requires `proposal`. Instruction unchanged from stock `spec-driven` (ADDED/MODIFIED/REMOVED/RENAMED, four-hashtag scenarios).
- **`design`**: generates `design.md`. Requires `proposal`.
- **`tasks`**: generates `tasks.md`. Requires `specs`, `design`. Checkbox format preserved.
- **`apply`**: `requires: [tasks]`, `tracks: tasks.md`, instruction adds branch naming, per-group commits, and checklist sync.

Templates live under `openspec/schemas/octospec/templates/`.

## Issue forms and labels

Issue forms (`.github/ISSUE_TEMPLATE/`) are the contract. Both the human and the agent create issues through them.

- **Feature form**: user story (`As a … I want … so that …`), context/problem, requirements (SHALL/MUST), success criteria, out of scope, open questions. Required fields are those G1 checks.
- **Bug form**: summary, steps to reproduce, expected, actual, impact.

Labels: `type:feature`, `type:bug`, `status:backlog`, `status:spec-ready`, `status:in-progress`, `status:in-review`.

## Gates

Hard gates block CI or merge; advisory gates warn.

| Gate | Description | Strength |
|---|---|---|
| G1 | Issue has the required sections before `/spec` produces artifacts | hard |
| G2 | `openspec validate --all --strict` passes **and** every unarchived change is complete (`openspec status` reports `isComplete: true`) | hard |
| G3 | Branch name matches `feat\|fix/<issue>-<slug>` | hard |
| G4 | PR body references the issue (`Closes #n`) | hard |
| G5 | Tests / build / lint pass | hard |
| G6 | A PR that touches an `openspec/changes/<name>/` folder carries a spec delta (`specs/**/*.md`) | hard |
| G7 | All tasks checked before merge | advisory |
| G8 | One OpenSpec change per issue (no duplicates) | advisory |

G1-G2 are checked locally in `/spec`; G2-G6 are re-checked in a GitHub Actions workflow (`.github/workflows/openspec.yml`) so they cannot be skipped.

## Multi-tool layout

The schema and `openspec/config.yaml` are tool-agnostic and committed once. Commands are per-tool shims:

- **pi**: prompt templates
- **opencode**: `.opencode/command/<name>.md`
- **GitHub Copilot**: `.github/prompts/<name>.prompt.md`

Skills are generated by OpenSpec per tool (`delivery: skills`).

## Risks and trade-offs

- **Schema drift from upstream.** A custom schema means OpenSpec's built-in guidance changes do not reach the workflow automatically. Mitigation: keep `spec-driven` as a reference and re-read it on OpenSpec upgrades.
- **The CLI cannot call `gh`.** GitHub steps are instruction text the agent follows. Mitigation: CI re-checks the critical ones (G4, G6) so a missed step fails the PR rather than shipping silently.
- **Published copies can go stale.** The issue comment mirrors the file and is not authoritative. Mitigation: every published copy carries the commit SHA it was generated from.
- **Tool differences.** pi, opencode and Copilot load skills and commands differently. Mitigation: keep commands to thin shims and put all shared logic in the schema.

## Open questions

1. Should `/idea` accept a rough paragraph and let the agent ask the interview questions, or should it require the fields up front?
2. Do we want the `issue` artifact to snapshot the whole issue body into `issue.md`, or only the stable narrative sections?
3. Should `/ship` be a separate command, or the tail of `/apply` behind a confirmation? (Current decision: separate, to keep the human in the loop.)
4. Do we hide the OpenSpec skills from direct invocation, or leave them visible as the low-level escape hatch?
