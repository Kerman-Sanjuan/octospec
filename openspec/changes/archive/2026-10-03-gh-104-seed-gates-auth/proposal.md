## Why

The root gates workflow authenticates `gh`, but the copy the CLI seeds does not,
so on a freshly seeded repository `gh` cannot read the linked issue and G1 false
fails on a complete issue. The drift is untested, so it recurred.

## What Changes

- Sync the seeded `.github/workflows/openspec.yml` with the root: add `GH_TOKEN`,
  `issues: read`, and the same `setup-node` major.
- Make `scripts/check-gates.sh` SKIP G1 when `gh` cannot read the issue body
  (auth, network, or a missing issue), instead of failing it.
- Keep the seeded `check-gates.sh` byte-identical to the root.
- Add a gate test for the unauthenticated-`gh` path.
- Add a seed test that fails when the seeded workflow drifts from the root on the
  token, the permissions, and the `setup-node` major.

## Capabilities

### New Capabilities

### Modified Capabilities
- `gates`: G1 skips when `gh` cannot read the issue, so a false negative is
  impossible.
- `cli-distribution`: the seeded CI workflow carries the same `gh`
  authentication as the root, enforced by a drift test.

## Impact

- `scripts/check-gates.sh`, `cli/internal/seed/repo/scripts/check-gates.sh`,
  `cli/internal/seed/repo/.github/workflows/openspec.yml`,
  `scripts/check-gates-test.sh`, `cli/internal/seed/schema_test.go`.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/104

## Changelog

### Fixed
- A seeded repository's gates workflow now authenticates `gh`, so G1 no longer false-fails on a complete issue.
- G1 skips, instead of failing, when `gh` cannot read the issue body.
