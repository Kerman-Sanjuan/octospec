## Why

The install path is documented but has never run, because there is no release. A 1.0.0 needs a tag, published binaries, and a verified install.

## What Changes

- Bump the CLI version to `1.0.0` and finalize `CHANGELOG.md` for the release.
- Tag `v1.0.0`, which triggers the release workflow, and verify `curl | sh` and `go install`.

## Capabilities

### New Capabilities
<!-- None. -->

### Modified Capabilities
- `release-process`: the first release is cut and the install path is verified end to end.

## Impact

- `cli/cmd/octospec/main.go` (the version), `CHANGELOG.md`, and a `v1.0.0` tag.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/30
