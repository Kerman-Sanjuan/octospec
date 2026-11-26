# Proposal — Grill the intake

## Why

`/idea` interviews for five form fields, one question at a time, and stops when the fields are filled —
not when the outcome is testable. Everything that was not decided in that conversation is guessed later
by `/spec` and `/apply`, where a gap costs ten times more. G1 then accepts the result, because it only
checks that four headings exist.

## What Changes

- The intake works a **decision frontier**: rounds of numbered questions, each with a recommended answer,
  asking only the decisions whose prerequisites are settled.
- The agent **looks up facts** it can read and asks only for **decisions**.
- The grill ends when the frontier is empty; answers land in the issue under `## Grill`.
- An unattended run stops with `needs_human` instead of inventing answers.
- G1 grows a second half: `## Grill` present, every requirement carries SHALL/MUST, success criteria
  non-empty, `## Open questions` empty.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `issue-intake`: the grill, its frontier, its rules and where the answers live.
- `gates`: G1 checks that the intake produced a decidable issue, not only that the headings exist.

## Impact

- `cli/internal/payload/commands/idea.md`, `bug.md`, `spec.md`
- `cli/internal/payload/agents/idea.md`
- `scripts/check-gates.sh` and `cli/internal/seed/repo/scripts/check-gates.sh` (they must stay identical)
- `cli/internal/seed/repo/.github/ISSUE_TEMPLATE/feature.yml`, `bug.yml`
- `README.md`, `docs/commands.md`, `docs/getting-started.md`

## Changelog

### Changed

- `/idea` and `/bug` grill the decision frontier in rounds before filing, with a recommended answer per
  question, and record the closed decisions in the issue under `## Grill`.
- G1 now fails an issue whose requirements state no SHALL/MUST, whose success criteria are empty, or
  that still carries open questions.
