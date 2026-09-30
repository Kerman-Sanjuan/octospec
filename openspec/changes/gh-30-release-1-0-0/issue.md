# Issue #30: Cut and verify v1.0.0

https://github.com/Kerman-Sanjuan/octospec/issues/30

## User story

As a user, I want a released version I can install with `curl | sh` or `go install`.

## Context / problem

The installer points at `releases/latest/download/octospec_<os>_<arch>`, but no release exists, so the documented install path has never run. `go install` is also unproven.

## Requirements

- A `v1.0.0` tag SHALL publish release binaries through the release workflow.
- `curl | sh` and `go install .../octospec/cli/cmd/octospec@latest` SHALL install a working binary.
- The release SHALL publish checksums.

## Success criteria

- On a clean machine, both install paths produce a runnable `octospec`, and `octospec version` reports `1.0.0`.

## Out of scope

- Post-1.0 features.
