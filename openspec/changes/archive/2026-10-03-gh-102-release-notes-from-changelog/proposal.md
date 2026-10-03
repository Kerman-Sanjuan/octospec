## Why

The v1.0.0 release notes are GoReleaser's native git changelog: 261 lines of raw
commits, including `Merge branch ...`, with none of the `### Added / Changed /
Fixed` structure the repository keeps. The `release-debut` spec says the release
carries notes from the changelog, but the configured behaviour generates exactly
the noise.

## What Changes

- Add `scripts/release-notes.sh`, which extracts a version's section from
  `CHANGELOG.md`, with a fallback when the section is missing.
- Make the release workflow pass that section to GoReleaser with
  `--release-notes`.
- Disable GoReleaser's native commit changelog.
- Add `scripts/release-notes-test.sh` and run it in CI.

## Capabilities

### New Capabilities

### Modified Capabilities
- `release-pipeline`: the release notes come from `CHANGELOG.md`, not the commits.

## Impact

- `.goreleaser.yaml`, `.github/workflows/release.yml`,
  `.github/workflows/cli.yml`, `scripts/release-notes.sh`,
  `scripts/release-notes-test.sh`.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/102

## Changelog

### Fixed
- Release notes now come from the changelog section for the version, not the raw git commit log.
