# Issue #102: Release notes should come from the changelog

https://github.com/Kerman-Sanjuan/octospec/issues/102

## User story

As a visitor, I want the v1.0.0 release page to read like the changelog, so that the release notes are a clean, structured summary instead of raw commit noise.

## Context / problem

The v1.0.0 release notes are GoReleaser's native git changelog: 261 lines of
raw commits, including `Merge branch ...` entries, with none of the
`### Added / Changed / Fixed` structure the repository keeps in `CHANGELOG.md`.
The `release-debut` spec says the release carries "release notes from the
changelog", but the configured behaviour (`release-pipeline`: notes from the
commits) generates exactly the noise. The binary is fine; the release page, the
shop window, is not.

## Requirements

- The release SHALL use the matching section of `CHANGELOG.md` as its notes.
- GoReleaser SHALL NOT generate the native git commit changelog.
- A helper SHALL extract a version's section from `CHANGELOG.md`, and it SHALL be
  tested.
- If a version has no section in `CHANGELOG.md`, the release SHALL still work,
  with a sensible fallback, and SHALL not fail.

## Success criteria

- Cutting the next tag produces a release whose notes mirror the `CHANGELOG.md`
  section for that version, grouped by impact.
- `sh scripts/release-notes-test.sh` passes.
- The extraction helper prints the correct section for a given tag on the
  current `CHANGELOG.md`.
- `openspec validate --all --strict`, `scripts/check-gates.sh`, and CI pass.

## Out of scope

- No retag of v1.0.0 and no rewrite of its existing notes.
- No change to the release binaries or the installer.

