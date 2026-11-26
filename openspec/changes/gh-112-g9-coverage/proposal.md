# Proposal — G9 requirement coverage

## Why

The gates check that a change has a spec (G6) and that its tasks are checked (G7). Nothing checks that
a requirement is proved by anything. That is the field where every harness we compare against scores
badly, and where ours scores zero: a requirement can be written, accepted, merged and never executed.

## What Changes

- A change declares `coverage.md`: one line per requirement, `- <title> -> <symbol>`.
- G9 reads the requirement titles from the change's spec deltas, and warns when one has no line, or
  when a line's symbol cannot be found outside `openspec/`.
- Advisory at first, like G7 and G8 were: a warning that names the requirement is useful on day one,
  and a hard failure on day one stops every change in flight for a reason nobody has practised.

## Capabilities

- `gates`: added — G9, its artifact and its conditions.

## Impact

- `scripts/check-gates.sh`, `README.md`, `cli/internal/payload/commands/apply.md`.
- Each change directory grows one file.

## Changelog

### Added

- Gate G9 (advisory): a change must map each of its requirements to the check that proves it, and a
  mapping that points at nothing is a warning.
