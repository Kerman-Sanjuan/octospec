## Why

`scripts/check-gates.sh` assumes it runs on a PR branch. On `main`, the branch gate (G3) fails and the script exits 1, so it cannot be run locally on `main` and the output is misleading.

## What Changes

- Skip the PR-only branch gate when the ref is a long-lived branch, and report the skip.
- Add a test for the gate script and run it in the gates workflow.
- Note the PR-only scope in the README gates table.

## Capabilities

### New Capabilities
<!-- None. -->

### Modified Capabilities
- `release-pipeline`: the gate script runs on a non-PR ref by skipping the PR-only branch gate.

## Impact

- `scripts/check-gates.sh`, a new `scripts/check-gates-test.sh`, `.github/workflows/openspec.yml`, and `README.md`.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/6
