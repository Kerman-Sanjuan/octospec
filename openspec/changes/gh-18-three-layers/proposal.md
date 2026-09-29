## Why

A change floods its issue with comments: `/spec` posts one per artifact (proposal, design, one per capability, tasks) and `/apply` posts one per task group, and re-runs add more. The timeline stops being readable. An issue should carry exactly three structured layers: the issue (requirement), one plan comment, one implementation comment.

## What Changes

- `/spec` posts a single **plan comment** consolidating proposal, capabilities/specs, design, and tasks under fixed headings, refreshed in place on re-run.
- `/apply` posts a single **implementation comment** with the task checklist and a short insight per task, refreshed in place after each group.
- Each layer gets a predefined, recognizable structure.
- If the issue is ambiguous, `/spec` asks and routes the answer to the issue body; it does not invent requirements.

## Capabilities

### New Capabilities
- `issue-timeline`: the fixed, three-layer comment structure an issue carries and how it is kept in sync.

### Modified Capabilities
<!-- None. -->

## Impact

- `commands/spec.md`, `commands/apply.md` (and the `.github/prompts/` + `repo-template/` mirrors).
- `schema/octospec/schema.yaml` (the publish/apply instructions).
- `README.md`.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/18
