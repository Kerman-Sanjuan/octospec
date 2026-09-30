## Why

PR #47 shipped the agentic layer to `main`, but the latest release is v1.0.0, which has no agents. The installer still serves the old version, and the onboarding acceptance test (#32) needs the released installer to carry the agents, the model configuration, and the tool surface.

## What Changes

- Bump `const version` in the CLI to `1.1.0`.
- Add the `1.1.0` changelog entry, grouped by impact.
- Make the `release-process` install requirement version-agnostic.
- Tag `v1.1.0` to publish the binaries and checksums.

## Capabilities

### New Capabilities
<!-- None. -->

### Modified Capabilities
- `release-process`: the install check describes the released version instead of hardcoding `1.0.0`.

## Impact

- `cli/cmd/octospec/main.go`, `CHANGELOG.md`, and the `release-process` spec.
- A `v1.1.0` tag.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/48
