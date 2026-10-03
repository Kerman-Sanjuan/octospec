## Why

`release-acceptance` is red on `main`. On a `workflow_run` event from a tag
push, `github.event.workflow_run.head_branch` is the tag name (`v1.0.0`), which
the workflow passes as `HEAD_REF`. Gate G3 then demands the branch pattern and
fails. The `acceptance` job passed on both OSes; only `gates` failed. The README
also pins a `--version v1.1.0` tag that no longer exists.

## What Changes

- Treat a non-`feat`/`fix` ref as long-lived in `scripts/check-gates.sh`, so G3
  skips a tag or any long-lived ref instead of failing it.
- Keep the seeded copy of the gate script in sync.
- Add a gate test for the tag-ref case.
- Pin the README install example to `v1.0.0`.

## Capabilities

### New Capabilities

### Modified Capabilities
- `gates`: G3 skips a long-lived or tag ref, not only `main`/`master`.

## Impact

- `scripts/check-gates.sh`, `cli/internal/seed/repo/scripts/check-gates.sh`,
  `scripts/check-gates-test.sh`, `README.md`.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/100

## Changelog

### Fixed
- Gate G3 skips a tag ref, so the release acceptance check is green on a tag.
- The README pins the install example to the published `v1.0.0` tag.
