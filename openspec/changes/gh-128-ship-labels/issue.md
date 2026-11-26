# Issue #128 — The /ship stage opens the pull request without labels

https://github.com/Kerman-Sanjuan/octospec/issues/128

Read-only mirror of the issue narrative. The issue body stays the human source of truth.

## User story

As the person triaging work, I want pull requests to carry the same kind and mode labels as the
issue they close, so that I can filter the ones a machine made without reading every title.

## Context

#124, #125 and #127 all carry zero labels while their issues carry `type:*`, `status:*` and
sometimes `mode:autonomous`. `agents/ship.md` says "Add `status:in-review` only after the gates
pass" and never says on which object, nor that the kind and the mode belong on the pull request.

Leaving `mode:autonomous` on the issue only means the pull request — the thing that gets merged —
does not say how it was made or what kind of change it is.

## Requirements

- `/ship` SHALL copy the issue's `type:*` labels onto the pull request.
- `/ship` SHALL copy `mode:autonomous` when the issue has it.
- `/ship` SHALL stop and name the reason when the issue has no `type:*` label.
- Copying SHALL be idempotent.
- `status:in-review` SHALL stay on the issue and SHALL NOT be copied.

## Success criteria

- After `/ship` on `type:bug` + `mode:autonomous`, the pull request lists both.
- Re-running `/ship` leaves the label set unchanged.
- An issue without `type:*` stops the stage with a named reason.

## Out of scope

- Deciding which labels a repository wants beyond the ones the issue already carries.
