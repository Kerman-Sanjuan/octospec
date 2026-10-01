## Why

`CHANGELOG.md` is a flat file that can drift from the actual changes because it is not reviewed through the same pull request. Every other artifact (specs, design, tasks) lives under `openspec/changes/`; the changelog should follow the same pattern.

## What Changes

- The full changelog entry lives in `openspec/changes/archive/<change>/proposal.md`.
- `CHANGELOG.md` is a generated summary derived from the latest archived change.
- The release process copies the changelog entry from the archived change.

## Capabilities

### New Capabilities

- `changelog-as-change`: the changelog content lives inside an OpenSpec change.

## Impact

- `docs/releasing.md` is updated to describe the new process.
- `CHANGELOG.md` becomes a generated file; do not edit it by hand.

## Changelog

### Added
- The change proposal carries a `## Changelog` section that `/archive` folds into `CHANGELOG.md`.
