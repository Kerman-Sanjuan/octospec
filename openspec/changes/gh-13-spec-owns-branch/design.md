## Context

`/spec` builds the approach artifacts but leaves them uncommitted. That produced #2, forces `/apply` to commit a baseline, and gives the published comments only a base SHA. Because the human may reject the approach (the issue is only the requirement), the spec needs a diffable, editable, versioned surface. Spec Kit creates the feature branch at specification time; superpowers creates the workspace at execution time. This change adopts the spec-time model.

## Goals / Non-Goals

**Goals:**
- Commit the artifacts at `/spec` time on the change branch.
- Real SHA pinning for the published comments.
- Idempotent re-run with in-place comment refresh.
- Simpler `/apply` (no branch creation, no baseline commit).

**Non-Goals:**
- Changing the label lifecycle.
- A second "spec PR"; one PR per change, opened at `/ship`.

## Decisions

- **`/spec` creates the branch and commits** (Spec Kit model). *Alternative rejected:* commit at `/apply` (superpowers model) — leaves the spec uncommitted until execution, so it has no reviewable diff or history.
- **Derive `feat` / `fix` from the issue's `type:*` label.** *Alternative rejected:* ask the user — the label already carries the type.
- **Push from `/spec`.** The published SHA must resolve remotely for review. *Alternative rejected:* commit locally and push at `/apply` — the SHA is unresolvable until then.
- **Refresh comments in place, tracking the comment id in `.openspec.yaml`** (`github.comments.<artifact>`). *Alternative rejected:* append a new comment each run (spam) or edit the last comment (ambiguous across artifacts).
- **Approval is a human action.** `/spec` stops at `status:spec-ready` (drafted); the human approves by running `/apply` or setting an approval label. *Alternative rejected:* `/spec` approving its own work — the label would lie.
- **Revisit #2:** `/apply` drops the "allow the change dir / commit baseline" special case, because `/spec` now leaves a clean tree on the branch.

## Risks / Trade-offs

- [A rejected approach leaves a branch behind] → Delete the branch; no PR was opened.
- [Re-running `/spec` must not re-scaffold] → Detect the existing change and branch from `.openspec.yaml` and the branch name, and reuse them.
- [Comment-id tracking adds state] → Store ids in `.openspec.yaml`; if an id is missing, create the comment then record it.
- [Pushing before approval] → No PR is opened at `/spec`; the branch is only a durable draft.

## Migration Plan

None. There are no in-flight changes at the time of writing.

## Open Questions

- None. Resolved: approval is the human running `/apply`; no separate
  `status:approved` label is added, so `/spec` stops at `status:spec-ready`
  (drafted) and never marks its own work approved.
