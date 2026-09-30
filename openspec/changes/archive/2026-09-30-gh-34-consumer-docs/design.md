## Context

The README serves contributors. Adopters need a shorter path: install, seed, run, and a place to look when something breaks.

## Goals / Non-Goals

**Goals:**
- A getting started, a command reference, and troubleshooting, linked from the README.

**Non-Goals:**
- A hosted docs site, or duplicating the design docs.

## Decisions

- **Three focused documents, in `docs/`.** One per need, so each stays short and linkable. *Alternative rejected:* one long page, which buries the answers.
- **The command reference covers both surfaces.** The CLI verbs (`install`, `seed`, `update`, `doctor`) and the workflow commands (`/idea` to `/archive`), since adopters touch both.

## Risks / Trade-offs

- [Docs drift from the code] → The AGENTS.md definition of done already requires updating the README when commands change; the same applies here.

## Migration Plan

None.

## Open Questions

- None.
