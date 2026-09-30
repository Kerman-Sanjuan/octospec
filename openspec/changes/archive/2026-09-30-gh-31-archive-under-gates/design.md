## Context

`/archive` moves the change into `openspec/changes/archive/`, syncs `openspec/specs/`, fills the new `Purpose`s, closes the issue, and deletes the merged branch. It commits and pushes to `main`. Branch protection now requires pull requests and the `cli` and `gates` checks, and enforces admins, so the push is rejected.

## Goals / Non-Goals

**Goals:**
- The archive lands on `main` with no manual step.
- The archive commits and pushes in every path.

**Non-Goals:**
- Changing the gate definitions for pull requests.
- Reviewing the archive, which is deterministic maintenance.

## Decisions

- **The owner pushes the archive directly; the protection does not block maintenance.** The archive is a deterministic step (move, sync, close, delete), not a reviewed change, so it does not need a pull request or the checks. Implemented by setting `enforce_admins` to false in the protection. *Alternative rejected:* an archive pull request, which adds a second PR per change and was explicitly called inefficient.
- **A committed `scripts/protect-main.sh`** makes the protection reproducible, so the settings are not ad hoc API calls. It sets the required checks (`cli`, `gates`), strict mode, `enforce_admins: false`, and blocks force pushes and deletions.

## Risks / Trade-offs

- [With `enforce_admins: false`, the owner can bypass the checks on any branch] → Acceptable for a single-maintainer repo. The checks still gate every pull request for contributors, and the owner only bypasses for maintenance.
- [The protection drifts from the script] → The script is the source of truth and can be re-run.

## Migration Plan

Run `scripts/protect-main.sh` once to apply `enforce_admins: false`. Existing settings (required checks, strict) stay.

## Open Questions

- None.
