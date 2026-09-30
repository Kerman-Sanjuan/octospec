## Context

The release workflow runs `goreleaser` on a `v*` tag and publishes binaries named `octospec_<os>_<arch>`. The `curl | sh` installer reads the latest release at that name. Nothing has been tagged, so none of it has run.

## Goals / Non-Goals

**Goals:**
- A published `v1.0.0` with binaries and checksums.
- Both install paths verified.

**Non-Goals:**
- Post-1.0 work, such as version pinning for the installer (#37).

## Decisions

- **The tag drives the version.** Bump `const version` in the CLI to match, and finalize the changelog, then tag. *Alternative rejected:* hand-editing the version in several places.
- **Verify against the real release.** Run the installer the README documents, not a local build.

## Risks / Trade-offs

- [goreleaser config drift from the installer's file name] → The name template and the installer both use `octospec_<os>_<arch>`; the verification catches a mismatch.
- [A failed release workflow leaves a tag with no artifacts] → Re-run the workflow, or delete and re-push the tag after a fix.

## Migration Plan

None. First release.

## Open Questions

- None.
