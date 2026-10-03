## Context

The repository is on `main` after the session work was reverted (#92). There are
two GitHub releases, v1.0.0 and v1.1.0 (both 2026-09-30), cut before the current
polish. The README opens with a bare `# octospec` heading and no visual. The
release workflow (`.github/workflows/release.yml`) runs GoReleaser on a `v*`
tag, and `docs/releasing.md` documents the steps. The README spec already
requires a quickstart, prerequisites, platforms, and links, so this change adds
only the visual layer and the release reset.

## Goals / Non-Goals

**Goals:**
- Give the project a logo and status badges so its landing page looks credible.
- Reset the release history to one polished v1.0.0 with a single changelog
  section.
- Record an end-to-end check of the published binary.

**Non-Goals:**
- No CLI, command, agent, or schema changes.
- No new tools.
- No rewrite of the README's prose beyond the visual layer and the links.

## Decisions

### A hand-authored SVG logo with a committed PNG
The logo is an SVG in the repo so it stays crisp and reviewable, plus a
rendered PNG for contexts that need a raster. The alternative, an external image
host, adds a network dependency and a third-party account for no benefit.

### Badges from shields.io with live endpoints
The CI, release, license, and Go badges use shields.io endpoints that read the
repository and the release feed, so they track reality instead of a static
string. The alternative, checked-in badge images, would go stale.

### Reset the release history rather than stacking a third release
The user wants the first public release to be the polished one. The old v1.0.0
and v1.1.0 releases and tags are removed, the `CHANGELOG.md` 1.1.0 section is
folded into a single 1.0.0 section, and a fresh `v1.0.0` tag is cut. This is a
deliberate rewrite of a public history while the project is still pre-adoption;
the alternative, cutting v1.2.0, leaves the unpolished releases in place, which
the user rejected.

### Verify the published binary, not just the local build
The end-to-end check runs against the installed release (the `curl | sh` path),
because the local `go test` suite already covers the code. This catches
packaging and installer problems that never appear in unit tests.

## Risks / Trade-offs

- [Deleting a tag that consumers may have pinned] -> The project is
  pre-adoption and the only consumers are the maintainer's own machines; the
  re-cut tag is documented in the release notes.
- [A logo authored by an agent may look generic] -> Keep it simple and
  geometric, review it in the PR, and iterate; it is one SVG, cheap to change.
- [A badge endpoint can be flaky or rate-limited] -> Badges are decoration, so a
  transient failure never blocks the release or CI; the README text stands alone.

## Migration Plan

1. Land the logo, badges, and README links through a pull request.
2. Delete the v1.1.0 and old v1.0.0 releases and tags.
3. Fold `CHANGELOG.md` into one 1.0.0 section and merge it.
4. Tag `v1.0.0`, push, and let the release workflow publish.
5. Run the end-to-end install check and fill in the About.
Rollback: re-tag from the prior commit if the release workflow fails; no data
is migrated.

## Open Questions

- None. The issue answers the version (reset to 1.0.0, remove the rest), the
  scope (logo, badges, verify, release cut, cleanup), and the boundary (no
  feature changes).
