## Context

The stage agents exist and are installed with their tool surface and model, but nothing dispatches to them. The commands carry a line, "Runs the `<stage>` agent", that reads as a spawn. Real dispatch differs per harness (Claude Code, opencode, and Copilot have subagents; pi does not), and it cannot be verified from this repository.

## Goals / Non-Goals

**Goals:**
- The commands and the standard describe the current behaviour honestly.
- The `agent-definitions` spec stops implying a dispatch.

**Non-Goals:**
- Wiring real dispatch, which is a follow-up with per-harness verification.
- Subagents in pi.

## Decisions

- **Soften the claim now, wire the dispatch later.** The honest wording is small and correct; real dispatch needs a harness to test against, and pi has none. *Alternative rejected:* wire dispatch blind, which would ship an unverified path and break the claim on pi anyway.
- **The agent is the stage contract, the command is the runner.** The agent defines the scope, tools, and model; the command runs the stage in the current session. *Alternative rejected:* keep the "selects the agent" wording, which describes something that does not happen.

## Risks / Trade-offs

- [The softened wording hides a real gap] -> The follow-up issue tracks the dispatch, so the gap is visible in the backlog, not in a misleading command.

## Migration Plan

None.

## Open Questions

- Which harness to implement dispatch for first (Claude Code and opencode are the likely candidates).
