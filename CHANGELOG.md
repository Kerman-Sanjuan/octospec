# Changelog

Grouped by impact. Newest first. Versions follow semantic versioning.

## Unreleased

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

## 1.0.0

The first minimum viable release: install, onboard, and ship through the loop.
