# Proposal — label the pull request

## Why

The pull request is what gets merged, and it is the only object a reviewer ever filters. Ours come
out blank: #124, #125 and #127 all have zero labels while their issues declare a `type:*` and,
since the autonomous work started, `mode:autonomous`. The flag that says "a machine may do this
unattended" currently stops at the issue, which is not the artifact that changes `main`.

`agents/ship.md` mentions a label once — "Add `status:in-review` only after the gates pass" — without
saying on which object, and says nothing about copying the kind or the mode.

## What Changes

- `/ship` reads the issue's labels and copies `type:*` plus `mode:autonomous` onto the pull request.
- An issue with no `type:*` stops the stage: a pull request with no kind cannot be filtered, and
  guessing the kind from the title is how labels stop meaning anything.
- `status:in-review` stays on the issue, and the payload says so explicitly.
- Copying is idempotent because `gh pr edit --add-label` is: re-running `/ship` adds nothing twice.

## Capabilities

- `ship-feedback-loop`: added — the pull request carries the issue's kind and mode.

## Impact

- `cli/internal/payload/commands/ship.md` — a step between opening the pull request and running CI.
- `cli/internal/payload/agents/ship.md` — the rule, with the reason.
- `docs/commands.md` — the `/ship` row says what it labels.

## Changelog

### Fixed

- `/ship` now labels the pull request with the issue's `type:*` and `mode:autonomous`, instead of
  leaving it blank.
