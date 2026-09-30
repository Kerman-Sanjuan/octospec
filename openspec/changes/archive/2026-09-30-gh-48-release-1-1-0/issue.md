# Issue #48: Cut and verify v1.1.0 (agentic layer)

https://github.com/Kerman-Sanjuan/octospec/issues/48

## User story

As a maintainer, I want to publish the agentic layer as v1.1.0, so the installer serves the current version and the onboarding acceptance run can use the real install path.

## Context / problem

PR #47 shipped the agentic layer to `main`, but the latest release is v1.0.0, which has no agents. `curl | sh` and `go install @latest` still install the old version. The onboarding acceptance test (#32) needs the released installer to carry the agents, the model configuration, and the tool surface.

The `release-process` spec also hardcodes `1.0.0` in its install scenario, so it goes stale the moment a new version ships.

## Requirements

- `const version` SHALL read `1.1.0` in `cli/cmd/octospec/main.go`.
- `CHANGELOG.md` SHALL have a dated `1.1.0` entry grouped by impact.
- The `release-process` spec SHALL describe the install check without hardcoding `1.0.0`.
- Tagging `v1.1.0` SHALL publish binaries and checksums.
- `curl | sh` SHALL install `1.1.0`.

## Success criteria

- `octospec version` prints `octospec 1.1.0` after a `curl | sh` install.
- The release has linux and darwin binaries on amd64 and arm64, plus a checksum file.
- `openspec validate --all --strict` passes and the `cli` and `gates` checks stay green.

## Out of scope

- The onboarding acceptance run itself (#32).
- Installer version pinning (#37).
