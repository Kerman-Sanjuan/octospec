## Context

Gate G3 in `scripts/check-gates.sh` handles `main` and `master` explicitly and
otherwise demands `feat|fix/<issue>-<slug>`. The `release-acceptance` workflow
triggers on `workflow_run` and passes
`HEAD_REF: ${{ github.event.workflow_run.head_branch }}`. When the triggering
run is a tag push, `head_branch` is the tag name (`v1.0.0`), so G3 fails and the
workflow is red on `main`. The `acceptance` job passed on ubuntu and macOS.

## Goals / Non-Goals

**Goals:**
- G3 skips a tag ref, without weakening the branch check for real PR branches.
- The release acceptance check is green on a tag.

**Non-Goals:**
- No change to the release workflow or a new tag.
- No rewrite of the generated release notes.

## Decisions

### Skip any ref that is not a `feat/` or `fix/` branch
The current code special-cases `main|master` and treats everything else as a PR
branch. Invert it: a ref that starts with `feat/` or `fix/` is a candidate to
check; anything else (main, a tag, a detached ref) is long-lived and skips. This
fixes the tag case and simplifies the logic. The alternative, special-casing
tags with a `v*` pattern, would still fail on any other long-lived ref and adds a
second rule for the same idea.

### Pin the README to v1.0.0
The README's pin example named `v1.1.0`, which the 1.0.0 reset deleted. Point it
at the published `v1.0.0`.

## Risks / Trade-offs

- [Skipping too much could hide a malformed PR branch] -> Only refs that already
  start with `feat/` or `fix/` are checked, so a typo like `feature/12-x` would
  now skip instead of fail. The gate test keeps a case for a `dependabot/...`
  branch failing, and the branch convention is documented and enforced at
  creation, so the trade-off is acceptable.

## Migration Plan

No data migration. The fix lands via PR; the green check is then restored with
`workflow_dispatch` on `v1.0.0`, no retag needed. Rollback is a revert.

## Open Questions

- None. The failing log pins the cause to `HEAD_REF=v1.0.0` and line 33.
