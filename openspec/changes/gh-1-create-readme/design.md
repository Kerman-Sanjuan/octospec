## Context

The README is a 5-line stub that only links to the design doc. octospec is now specified and implemented (schema override, command shims, seed files, install script, CI gates), but none of that is discoverable from the repo's front door. This change replaces the stub with a proper README that is an entry point, not a docs site.

## Goals / Non-Goals

**Goals:**
- Explain what octospec is and the problem it solves.
- Describe the workflow loop and the command surface.
- Document the architecture and installation for pi, opencode, and GitHub Copilot.
- Point to the design doc and implementation plan for depth.

**Non-Goals:**
- Replacing the design doc or implementation plan.
- A multi-page docs site.
- Tutorials or onboarding guides beyond install + workflow loop.

## Decisions

- **README structure follows the workflow, not the repo layout.** Sections ordered for a new reader: what/why, workflow loop, architecture, install, commands, deeper docs. *Alternative rejected:* mirroring the file tree, which optimizes for contributors over new users.
- **README is an entry point that links for depth.** It summarizes and links to `docs/plans/*` rather than duplicating them. *Alternative rejected:* inlining the full design, which would go stale and duplicate the source of truth.
- **Install section covers the three target agents** (pi, opencode, GitHub Copilot), matching `install.sh` and the multi-tool layout. *Alternative rejected:* a single-agent install that ignores the canonical repo's multi-tool goal.

## Risks / Trade-offs

- [README drifts from the design doc as the workflow evolves] → Keep summaries high-level and link to `docs/plans/*` as the authoritative detail; do not restate specifics that change.
- [Command list goes stale when commands are added] → Derive the command table from `commands/` and the design doc's command surface table; keep it a fixed, reviewed list.

## Migration Plan

None. Documentation-only change to `README.md`; no deploy or rollback.

## Open Questions

- None.
