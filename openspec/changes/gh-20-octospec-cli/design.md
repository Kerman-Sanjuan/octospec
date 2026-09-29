## Context

octospec currently distributes the same Markdown command file several times and installs with `install.sh`. This change introduces a Go CLI that renders one canonical payload into each tool's directory, removes the copies and the shell script, and gives the new code artifact its own tests, CI/CD, and a `/ship` feedback loop.

## Goals / Non-Goals

**Goals:**
- One canonical command payload; per-tool output is generated, never stored.
- A single Go binary, installable via `curl | sh` (and `go install`).
- Targets: pi, opencode, GitHub Copilot, Claude Code — adding one is a mapping.
- `update` that preserves local edits.
- Go tests, PR CI (build/format/vet/test + octospec gates), tagged release binaries.
- `/ship` feeds CI failures back to `/apply` and proceeds only on green.

**Non-Goals:**
- Changing what the commands do, or the G1-G8 gate definitions.
- Reimplementing the OpenSpec CLI.
- Auto-migrating existing `install.sh` installs.

## Decisions

- **Go, single binary.** No runtime dependency, trivial cross-compilation, one artifact. *Alternative rejected:* Node/npm (the OpenSpec CLI's ecosystem) — drags a runtime in.
- **Cobra for the command tree.** Subcommands, flags, help, and shell completion are standard Go-CLI ergonomics, and the CLI grows with new verbs (`install`, `seed`, `update`). *Alternative rejected:* hand-rolled `os.Args` parsing — no help or completion, more code to maintain.
- **Canonical payload embedded with `go:embed`.** One binary, no data directory to keep in sync. *Alternative rejected:* a payload dir next to the binary — a second thing to ship and drift.
- **A single target table** (`internal/targets`) maps tool -> directory + front-matter. Global for pi/opencode, repo-local for Copilot/Claude. *Alternative rejected:* per-tool folders in the repo (the status quo).
- **`curl | sh` + goreleaser releases.** Standard for Go CLIs; the installer fetches a release binary. `go install` stays for contributors.
- **Managed-file hash manifest** (`octospec-managed.json`) so `update` re-applies and preserves hand edits. *Alternative rejected:* blind overwrite.
- **The CLI absorbs `install.sh`**: schema install, label provisioning, and the repo seed move in. `install.sh` is deleted.
- **`/ship` is the CI gate.** It runs/observes CI; on failure it hands the specific check back to `/apply`; it opens/updates the PR only on green.

## Risks / Trade-offs

- [Contributors now need Go] → release binaries + `curl | sh` for users; `go install` and a documented toolchain for contributors.
- [Embedding means a payload change needs a rebuild] → acceptable; the release is the distribution unit.
- [Target formats differ (pi front matter, Copilot `mode: agent`)] → the renderer owns per-target front matter; tests pin each format.
- [Removing `install.sh` and the copies breaks existing users] → documented manual cleanup; no auto-migration (out of scope).
- [CI cost doubles (Go + octospec gates)] → both are fast; run them as separate jobs in one workflow.

## Migration Plan

1. Land the CLI and its tests/CI alongside the current files.
2. Switch the README and docs to the CLI.
3. Delete `install.sh`, `.github/prompts/`, `repo-template/.github/prompts/`, and the tool mirrors in the same change.
4. Existing users run the CLI once; a note explains removing the old files.

## Open Questions

- Module layout: a `cli/` subdirectory module, or the repo root as the Go module?
- Does the OpenSpec-generated skills (`openspec-propose`, …) install also move into the CLI, or stay with `openspec`?
