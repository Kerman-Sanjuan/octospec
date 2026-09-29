# Issue #20: `octospec/cli` — a Go CLI to install octospec from one canonical payload

https://github.com/Kerman-Sanjuan/octospec/issues/20

## User story

As a maintainer distributing octospec, I want a single Go CLI (`octospec/cli`) that installs octospec's commands into each tool's native location (pi, opencode, GitHub Copilot, Claude Code) from one canonical payload, so that the repo stops storing per-tool mirrored folders and adding a tool becomes a path mapping, not a new copy.

## Context / problem

octospec is distributed as near-identical copies per tool plus a shell script: `commands/*.md` (canonical), `.github/prompts/*.prompt.md` (Copilot), `repo-template/.github/prompts/*` (seed), and `install.sh` (copies `commands/*.md` into the global pi and opencode directories and seeds `repo-template/`). Every tool follows the same source pattern (a Markdown command file) and differs only by directory and a small front-matter tweak. So the repo carries N copies that drift — e.g. the `/spec` fix had to be hand-applied to 3 files — the repo root is cluttered with them, and installation is a shell script tied to a checked-out repo rather than an installable artifact.

A new code artifact (a Go CLI) also needs its own quality gates, and today `/ship` only reports a CI failure instead of feeding it back into `/apply` so the change can iterate until green.

## Requirements

**Distribution and install**

- The CLI SHALL be a single Go binary (`octospec/cli`) with no Node/runtime dependency.
- The CLI SHALL be distributed via a `curl | sh` installer (and `go install .../octospec/cli@latest` for contributors).
- The CLI SHALL install octospec's commands into the directory each supported tool expects, from one canonical payload embedded in the binary.
- The repo SHALL NOT store per-tool derived copies (`.github/prompts/`, `repo-template/.github/prompts/`, and the `.opencode/` / `.pi/` mirrors); the CLI SHALL generate them.
- The CLI SHALL support pi, opencode, GitHub Copilot, and Claude Code; adding a tool SHALL be a target mapping (directory + file format), not a new source copy.
- The CLI SHALL support both global targets (pi, opencode — user directories) and repo-local targets (Copilot -> `.github/prompts/`, Claude Code -> `.claude/commands/`).
- The CLI SHALL provide an `update` that re-applies from saved answers and preserves local edits via managed-file hashes.
- The CLI SHALL replace `install.sh` and its `--repo` seed: schema install, label provisioning, and the `repo-template` seed (issue forms, CI, `openspec/config.yaml`).
- The CLI SHALL support non-interactive use (flags/config) in addition to an interactive wizard.

**Tests, CI and CD**

- The CLI SHALL ship with Go tests (unit and integration) runnable as `go test ./...`.
- CI SHALL run on every pull request: build, `gofmt` check, `go vet`, `go test ./...`, and the existing octospec gates (`scripts/check-gates.sh`, `openspec validate`).
- CD SHALL publish versioned release binaries on tag, and the `curl | sh` installer SHALL install from a published release.
- Gates SHALL be enforced controls — a red CI blocks the PR — not prose, per GH-600.

**Ship feedback loop**

- `/ship` SHALL run/observe CI; when it fails it SHALL hand the specific failure back to `/apply` to fix and iterate, instead of only reporting a failure.
- `/ship` SHALL proceed (open/update the PR) only once the gates are green.

## Success criteria

- A user can install the CLI with a single `curl | sh` command — no repo clone, no Node.
- `octospec install --tool <pi|opencode|copilot|claude>` writes the commands into that tool's directory, and the tool picks them up.
- Adding a new tool requires only a new target mapping — no new source copies anywhere.
- The repo root no longer contains `install.sh` or the per-tool copies (`.github/prompts/`, `repo-template/.github/prompts/`, and the `.opencode/` / `.pi/` mirrors).
- `octospec update` re-applies from saved answers, preserves local edits, and reports what it preserved.
- The old `install.sh` + `--repo` behavior (schema install, label provisioning, `repo-template` seed) is fully covered by the CLI.
- `go install .../octospec/cli@latest` gives contributors a working binary.
- `go test ./...` is green locally and in CI, and a red CI blocks the PR.
- A failing CI causes `/ship` to hand back to `/apply` with the specific failure, and the loop repeats until green.
- Tagging a release publishes binaries; `curl | sh` installs the released version.

## Out of scope

- Changing what the commands do — this is about how octospec is distributed and installed, not the `/idea -> ... -> /archive` semantics.
- Changing the gate definitions (G1-G8) themselves.
- Tools beyond pi, opencode, GitHub Copilot, and Claude Code in the first release (the target mapping is extensible).
- Reimplementing the OpenSpec CLI — octospec keeps delegating to it.
- Auto-migrating machines or repos that used `install.sh` — a documented manual cleanup is fine.
- A hosted service, telemetry, or a docs site — a binary and its `--help` are enough.
