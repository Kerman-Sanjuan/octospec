## Context

`install.sh --repo` seeds `.github/` (issue forms, prompts, CI workflow), `openspec/config.yaml` and `scripts/check-gates.sh`. The workflow's labels are assumed to exist but are never created, so a fresh repo fails at the first label operation. This change makes seeding create them.

## Goals / Non-Goals

**Goals:**
- `install.sh --repo` creates the six workflow labels.
- Re-running the seed is safe (idempotent).
- The label set is defined in one place.

**Non-Goals:**
- Changing the lifecycle or the commands' label behaviour.
- Provisioning labels globally (labels are per-repo).

## Decisions

- **Provision labels from the `--repo` seed step.** Seeding already targets one repository with `gh` available. *Alternative rejected:* a global step — there is nothing global to seed, labels are per-repo.
- **Idempotent with `gh label create --force`.** `--force` updates an existing label and never fails, keeping colours/descriptions in sync. *Alternative rejected:* check-then-create (more code, racy).
- **Put the label set in `scripts/seed-labels.sh`, invoked by `install.sh --repo`.** A dedicated script keeps `install.sh` readable, is runnable standalone, and gives one definition of the label set. *Alternative rejected:* inlining the list in `install.sh` (hidden, harder to test).

## Risks / Trade-offs

- [`--force` overwrites a user's custom label colour/description] → The colours/descriptions are part of the workflow contract and are documented; a reseed intentionally re-syncs them.
- [`gh` unavailable or unauthenticated] → The `--repo` step already assumes `gh`; `scripts/seed-labels.sh` skips with a warning if `gh` is missing rather than failing the whole seed.

## Migration Plan

None. Re-running `install.sh --repo` adds the labels to an existing repo.

## Open Questions

- None.
