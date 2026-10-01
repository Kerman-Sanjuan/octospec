## Context

`CHANGELOG.md` is a flat file in the repo root. It is edited by hand at release time, so it can drift from the changes that shipped. Every other artifact lives under `openspec/changes/`; the changelog should follow the same review path.

## Goals / Non-Goals

**Goals:**
- The changelog entry is written next to the change, so the same pull request reviews both.
- `/archive` folds the entry into `CHANGELOG.md`, so the file cannot drift.
- The release only names the version and the date.

**Non-Goals:**
- A tool that renders the changelog from the archive. The archive step is enough.
- Back-filling the entry into the changes that already shipped.

## Decisions

- **The entry lives in the proposal.** Each `proposal.md` carries a `## Changelog` section. *Alternative rejected:* a separate `changelog.md` artifact per change, because it duplicates the proposal review.
- **`/archive` copies the entry.** The step runs before the archive commit, so one commit holds both. *Alternative rejected:* a release-time scrape of the archive, which needs a parser and can miss entries.
- **`None` means no entry.** A change with no user-visible effect writes `None`, and the step skips it. *Alternative rejected:* an empty section, which is ambiguous.

## Risks / Trade-offs

- [A change forgets its `## Changelog` section] -> The section is required by the schema and the archive step stops when it is missing.
- [Two changes ship in one release] -> Both entries land under the same unreleased heading; the release moves them together.

## Migration Plan

Add the `## Changelog` section to the proposal template and to the archive command. Update `docs/releasing.md`.

## Open Questions

- None.
