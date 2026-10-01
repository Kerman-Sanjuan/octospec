# Changelog

Grouped by impact. Newest first. Versions follow semantic versioning.

## Unreleased

### Added
- G1 in `scripts/check-gates.sh`: the linked issue must carry its required sections.
- A living `gates` spec documenting G1 through G8.
- A `## Changelog` section in every change proposal that `/archive` folds into `CHANGELOG.md`.
- `scripts/e2e-test.sh`, a hermetic end-to-end test of seed, install, and the gates.

### Changed
- The release acceptance workflow runs the gates before it verifies the binary.

### Removed
- `/explore`, which duplicated the OpenSpec explore skill.

## 1.1.0 - 2026-09-30

### Added
- Stage agents: one scoped agent per stage (idea, spec, apply, ship, archive), rendered for pi, opencode, GitHub Copilot, and Claude Code from one canonical payload.
- A tool surface per agent, mapped to each tool's permission field. Claude and Copilot use `tools`, opencode uses `permissions`, and pi carries the surface in the body.
- `octospec models`, a role-to-model TUI or `--set` flag, backed by one model configuration in `.octospec/octospec.json`.

### Changed
- `octospec install` and `octospec update` also write the agent definitions.
- `octospec doctor` also reports the agents and the model configuration.
- The build uses Go 1.23.

## 1.0.0 - 2026-09-30

### Added
- A Go CLI, `octospec`, that installs the commands and the OpenSpec skills into pi, opencode, GitHub Copilot, and Claude Code from one canonical payload.
- `octospec seed` for the OpenSpec schema, the repo seed files, and the workflow labels.
- `octospec update`, which re-applies an install and preserves local edits.
- `octospec doctor`, a read-only check of an install.
- A `curl | sh` installer and a release workflow.
- A three-layer issue timeline: one `## Plan` comment and one `## Implementation` comment.
- Consumer docs: getting started, a command reference, and troubleshooting.

### Changed
- `/spec` creates the change branch, commits the artifacts, and publishes a single plan comment.
- `/apply` executes on the branch `/spec` created and keeps a single implementation comment.
- `/ship` runs CI and loops failures back to `/apply`.
- `/archive` commits and pushes the archive.
- `main` requires the `cli` and `gates` checks.

### Removed
- `install.sh` (the old seeding script) and the per-tool command copies.
