## Context

octospec's loop runs as five discrete commands (`/idea`, `/spec`, `/apply`, `/ship`, `/archive`), each invoked by a human, who also approves the handoff between stages (for example, the human runs `/apply` to approve the approach drafted by `/spec`). This keeps a human in the loop at every transition, which is deliberate for day-to-day work but blocks leaving the system unattended, e.g. overnight, per issue #109.

The stage commands already read and write their state through GitHub (issue labels, plan/implementation comments, pull requests) and through the `openspec/changes/<change>/` artifacts on disk. An autonomous mode can reuse that state machine instead of inventing a new one.

## Goals / Non-Goals

**Goals:**
- Run the full loop (idea → spec → apply → ship → archive) for one issue at a time with no human confirmation between stages.
- Select which issue to work on from the backlog's active statuses, and pick up bugs found along the way, without inventing requirements beyond `openspec/specs/`.
- Retry a stage that fails or whose gate (G2a, G2b, G3, G4, G6) does not pass, instead of stopping or skipping it.
- Leave the same GitHub trail a human run would leave: labels, comments, commits, pull requests.

**Non-Goals:**
- No scheduler, cron, or always-on daemon. The mode is invoked manually and runs until it stops or is told to stop; "leave it overnight" means a long foreground or background run, not a new trigger mechanism.
- No change to the gates (G2a, G2b, G3, G4, G6) or to what counts as passing.
- No bypass of pull requests. Autonomous mode still proposes a PR for a human to merge; it does not merge to `main` on its own.

## Decisions

- **One new orchestration layer, not five rewritten commands.** The autonomous mode is a thin driver that calls the existing `/idea`, `/spec`, `/apply`, `/ship`, `/archive` stage contracts in sequence for a single issue, reusing their artifacts and labels as the state to resume from. Alternative considered: fold autonomy into each stage command individually (e.g. "auto" flag per command); rejected because it would duplicate the retry/selection logic five times and blur the one-canonical-payload invariant.
- **Issue selection is label-driven, not a separate queue.** The mode asks GitHub for issues already in an active status (`status:in-progress` or equivalent) rather than maintaining its own backlog state. Alternative considered: a dedicated "auto-queue" label; rejected because it would create a second source of truth alongside the existing `status:*` labels.
- **Retry means re-running the same stage, not skipping to the next.** When a gate fails, the mode stays on the failing stage and re-runs it (e.g. `/apply` fixes the failing check, `/ship` re-observes CI) up to a bounded number of attempts, matching the existing `ship-feedback-loop` pattern of handing a failure back to `/apply`.
- **No silent abandonment.** If a stage keeps failing past the retry bound, the mode posts a comment on the issue describing the failure and stops working that issue, rather than leaving it in an ambiguous state or looping forever.

## Risks / Trade-offs

- [Risk] Unattended retries could loop indefinitely on a stage that can never pass (e.g. a gate requiring a human decision) → Mitigation: bound retries per stage and post a clear failure comment, then move to the next issue instead of looping forever.
- [Risk] Running unattended across many issues overnight could open or update many pull requests before a human reviews any of them → Mitigation: the mode works one issue at a time to completion (through pull request) before starting the next, so review backlog grows predictably.
- [Risk] Autonomous bug discovery could create issues that duplicate existing ones or invent requirements not in `openspec/specs/` → Mitigation: the mode only files a bug issue through the existing `/idea`-equivalent intake, checked against current specs, never invents acceptance criteria.

## Migration Plan

No migration: this adds a new, opt-in entry point alongside the existing manual commands. Running the manual commands continues to work exactly as before.

## Open Questions

- Which harness-level mechanism invokes the autonomous driver (a new slash command, a CLI subcommand, or both)? Left to `/apply` to decide against `cli/internal/payload/commands/` conventions.
- How many retries per stage is "bounded"? Left to `/apply` to pick a concrete default and make it configurable if the existing config surface (`.octospec/octospec.json`) supports it.
