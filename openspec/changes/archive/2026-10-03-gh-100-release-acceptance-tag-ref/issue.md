# Issue #100: release-acceptance fails G3 on a tag ref

https://github.com/Kerman-Sanjuan/octospec/issues/100

## User story

As a maintainer, I want the release acceptance check to pass on a tag, so that
the workflow that validates a release is green on `main` and the README does not
point at a deleted tag.

## Context / problem

`release-acceptance` is red on `main` after the v1.0.0 debut. When the workflow
runs from a `workflow_run` event triggered by a tag push, GitHub sets
`workflow_run.head_branch` to the tag name (`v1.0.0`), not a branch. The workflow
passes that value as `HEAD_REF`, and gate G3 then requires the branch pattern
`feat|fix/<issue>-<slug>` and fails:

```
FAIL G3 branch name: 'v1.0.0' does not match feat|fix/<issue>-<slug>
```

The `acceptance` job (which installs and runs the published binary on ubuntu and
macOS) passed on both OSes. Only the `gates` job failed, and it is the release
validation workflow, so a red `main` is misleading.

Separately, `README.md` line 180 still shows `--version v1.1.0`, a tag that was
deleted in the 1.0.0 reset, so a reader who copies the command gets a 404.

## Requirements

- The release acceptance workflow SHALL NOT fail G3 on a tag ref.
- `scripts/check-gates.sh` SHALL treat a non-`feat`/`fix` ref as a long-lived ref
  and SKIP G3, rather than fail it.
- The README SHALL show an install example pinned to a tag that exists (v1.0.0).
- A gate test SHALL cover the tag-ref case so G3 does not regress.

## Success criteria

- `HEAD_REF=v1.0.0 sh scripts/check-gates.sh` reports `SKIP G3` and does not
  fail the run.
- `release-acceptance` is green on `main` (via `workflow_dispatch` with
  `v1.0.0`, no retag needed).
- The README pins `--version v1.0.0`.
- `scripts/check-gates-test.sh` passes and includes the tag-ref case.

## Out of scope

- No change to the release binaries or a new tag.
- No rewrite of the generated release notes (tracked separately).

