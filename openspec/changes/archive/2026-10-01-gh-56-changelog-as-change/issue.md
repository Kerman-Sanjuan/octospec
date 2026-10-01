# Issue #56: CHANGELOG.md should live inside an OpenSpec change

https://github.com/Kerman-Sanjuan/octospec/issues/56

## User story

As a maintainer, I want the changelog to live inside an OpenSpec change so that it stays in sync with the actual changes and follows the same review process.

## Context / problem

`CHANGELOG.md` is a flat file in the repo root. It can drift from the actual changes because it is not reviewed through the same pull request as the change it describes. Every other artifact (specs, design, tasks) lives under `openspec/changes/`; the changelog should follow the same pattern.

## Requirements

- The changelog content shall live inside an OpenSpec change under `openspec/changes/`.
- The release process shall update the changelog by copying from the latest archived change.
- `CHANGELOG.md` in the root shall be a summary generated from the latest archived change.

## Success criteria

- `openspec/changes/archive/` contains a change whose `proposal.md` holds the changelog for the release.
- `CHANGELOG.md` is generated from the latest archived change.
- The release process uses the changelog from the archived change.

## Out of scope

- Back-filling the changelog into existing archived changes.
- A tool to export the changelog. Manual copy is acceptable for now.
