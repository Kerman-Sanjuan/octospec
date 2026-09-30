## Context

The release workflow runs `goreleaser` on a `v*` tag with no linker flags, so the tag alone does not set the reported version. The version lives in the `const version` of the CLI. The latest release is v1.0.0, which predates the agentic layer. The `release-process` spec hardcodes `1.0.0` in its install scenario, so it goes stale when a new version ships.

## Goals / Non-Goals

**Goals:**
- Publish v1.1.0 with the agentic layer.
- Make the install check version-agnostic.

**Non-Goals:**
- The onboarding acceptance run (#32).
- Installer version pinning (#37).

## Decisions

- **Bump the const, not the linker flags.** `goreleaser` sets no ldflags, so the tag alone does not change what `octospec version` prints. Bump `const version`. *Alternative rejected:* add ldflags, which is a larger change to the build.
- **Land the version and the changelog through a pull request.** The releasing doc requires it, and `main` requires the `cli` and `gates` checks. *Alternative rejected:* push to `main` directly.
- **Make the spec version-agnostic.** The install scenario names the released version instead of `1.0.0`. *Alternative rejected:* bump the number each release, which goes stale again.

## Risks / Trade-offs

- [The tag and the const disagree] -> The release checklist runs `octospec version` after the install.
- [The release workflow fails after the tag] -> Re-run the workflow, or delete and re-push the tag after a fix.

## Migration Plan

None. Minor release, no breaking change.

## Open Questions

- None.
