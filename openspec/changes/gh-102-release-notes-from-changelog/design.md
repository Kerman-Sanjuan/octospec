## Context

`.goreleaser.yaml` sets `changelog: {use: github, sort: asc}`, so GoReleaser
generates release notes from the git log. On the v1.0.0 tag that produced 261
lines of raw commits, including `Merge branch ...` entries, with no `### Added
/ Changed / Fixed` structure. `CHANGELOG.md` already holds the curated, grouped
entry per version, folded in by `/archive`. The `release-debut` spec says the
release carries notes from the changelog, but the config says otherwise.

## Goals / Non-Goals

**Goals:**
- The release notes are the `CHANGELOG.md` section for the tag, grouped by
  impact.
- A missing section never fails a release.

**Non-Goals:**
- No retag of v1.0.0 and no rewrite of its existing notes.
- No change to the binaries or the installer.

## Decisions

### Extract the section in a small POSIX script and pass it with --release-notes
`scripts/release-notes.sh <version>` prints the `## <version>` section from
`CHANGELOG.md`, stopping at the next `## ` heading, stripping a leading `v`, and
printing a one-line fallback when the section is missing. The release workflow
runs it with `GITHUB_REF_NAME` and passes the file to
`goreleaser release --release-notes`. This keeps the change small and testable
with a shell test, matching the repo's other scripts. The alternative, a
GoReleaser `changelog` template that rewrites the log, cannot reproduce the
hand-curated grouping and would duplicate the changelog.

### Disable the native changelog
Set `changelog: {disable: true}`. Without it, GoReleaser still computes a commit
changelog; `--release-notes` overrides the published body, but disabling removes
the noise and makes the intent explicit.

### Run the extraction test in CI
Add `scripts/release-notes-test.sh` to the `cli` job. It checks the real v1.0.0
section, the `v` prefix handling, the missing-version fallback, and a synthetic
changelog, so the extraction cannot silently break.

## Risks / Trade-offs

- [The CHANGELOG section and the tag must agree] -> `/archive` folds the entry
  under the version heading before the tag is cut, so they agree by construction;
  the fallback covers a mistag without failing the release.
- [A shell awk parser can be brittle] -> The heading match is anchored
  (`^## <version>( |$)`), and the test covers the boundary and the prefix.

## Migration Plan

No data migration. The next tag picks up the new notes; v1.0.0 keeps its current
notes. Rollback is a revert.

## Open Questions

- None. The issue names the symptom and the fix.
