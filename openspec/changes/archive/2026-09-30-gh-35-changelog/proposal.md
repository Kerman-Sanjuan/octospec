## Why

A release needs a human-facing record and a repeatable procedure. Today the binaries are configured to publish on a tag, but nothing says how to cut one or what changed.

## What Changes

- Add `CHANGELOG.md`, grouped by impact.
- Add `docs/releasing.md` with the release steps and the tag convention.

## Capabilities

### New Capabilities
- `release-process`: how a version is cut and recorded.

### Modified Capabilities
<!-- None. -->

## Impact

- `CHANGELOG.md`, `docs/releasing.md` (new), `README.md`.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/35
