## Why

The distribution path has no safety net. The installer trusts whatever `releases/latest` serves, it accepts no version, and nothing checks that a published binary runs. Delivery depends on a manual step that no longer exists.

## What Changes

- The `curl | sh` installer verifies the release `checksums.txt` before it installs the binary.
- The installer accepts a version, so a user can pin.
- A release workflow runs an acceptance test against the published binary on linux and macOS.
- `goreleaser` generates the release notes from the commits.
- Add a `CONTRIBUTING.md` that teaches how to contribute.
- The docs list the supported platforms and the install channels.

## Capabilities

### New Capabilities
- `contributing-guide`: a contributing guide that teaches the loop, the branch naming, the tests, and the release.

### Modified Capabilities
- `cli-distribution`: the installer verifies the checksum and accepts a version.
- `release-pipeline`: the release runs an acceptance test against the published binary and publishes generated release notes.
- `readme`: the install section lists the supported platforms and the channels.

## Impact

- `install.sh`, `.github/workflows/release.yml`, `.goreleaser.yaml`, and a new acceptance workflow.
- `CONTRIBUTING.md`, `README.md`, and `docs/releasing.md`.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/50
