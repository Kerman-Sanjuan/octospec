## Why

Today every stage of the octospec loop (`idea`, `spec`, `apply`, `ship`, `archive`) needs a human to invoke it and approve the handoff to the next stage, so the loop cannot run unattended, for example overnight (#109).

## What Changes

- Add an autonomous mode that drives an issue through all five stages (idea → spec → apply → ship → archive) without asking for human confirmation between stages.
- The mode selects issues to work on from the GitHub backlog's active statuses (`status:in-progress` and similar) and also picks up bugs discovered while working, always against the existing product definition in `openspec/specs/`.
- On a stage failure or a failed gate, the mode retries that stage until it passes before moving on.
- The mode keeps using GitHub as the single record of work: comments, commits, pull requests, and label updates, exactly as the manual stages do today.

## Capabilities

### New Capabilities
- `autonomous-mode`: an unattended orchestration loop that runs the five stage commands in sequence for one issue at a time, retries a failing stage, and leaves the full trail on GitHub.

### Modified Capabilities

(none - the existing stage commands and their artifacts are unchanged; the new capability calls them as-is)

## Impact

- Affected code: a new orchestration entry point in `cli/` (and/or a new payload command) that sequences the existing `idea`, `spec`, `apply`, `ship`, `archive` stages for a single issue.
- Affected systems: GitHub Issues (status labels, comments), GitHub Actions gates (G2a, G2b, G3, G4, G6) remain the pass/fail contract for retries.
- No change to the gates themselves, and no change to the requirement that work lands through a pull request.
- Issue: https://github.com/Kerman-Sanjuan/octospec/issues/109

## Changelog

### Added
- An autonomous mode that runs the idea → spec → apply → ship → archive loop for an issue without human confirmation between stages, retrying a stage until its gate passes, with the full trail recorded on GitHub.
