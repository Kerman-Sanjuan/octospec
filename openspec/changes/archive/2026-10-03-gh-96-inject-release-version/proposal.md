## Why

`octospec version` prints a hardcoded `const version = "1.1.0"`, so a release
built from a `v1.0.0` tag still reports `1.1.0` and the release acceptance check
fails. The tag must drive the version.

## What Changes

- Make the version a variable with a `dev` fallback, set at build time.
- Make GoReleaser inject the tag through `ldflags`, stripping the leading `v`.
- Cover the fallback and the injected value with a test.

## Capabilities

### New Capabilities

### Modified Capabilities
- `release-pipeline`: the release build injects the tag into the binary, so
  `octospec version` reports the release.

## Impact

- `cli/cmd/octospec/main.go`, `.goreleaser.yaml`, and a new test.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/96

## Changelog

### Fixed
- `octospec version` now reports the release tag instead of a hardcoded value.
