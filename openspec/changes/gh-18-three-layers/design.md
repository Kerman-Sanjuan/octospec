## Context

`/spec` and `/apply` currently mirror every artifact and every task group into its own issue comment, so a change's timeline grows unbounded and loses any fixed shape. #13 added refresh-in-place for `/spec` re-runs; this change makes each phase publish exactly one structured comment.

## Goals / Non-Goals

**Goals:**
- Exactly one plan comment (`/spec`) and one implementation comment (`/apply`) per issue.
- A fixed, recognizable structure per layer.
- Refresh in place on every re-run.

**Non-Goals:**
- `/idea`'s interview clarity bar.
- Reducing the number of PRs or reviews.

## Decisions

- **One comment per phase, not one per artifact/group.** *Alternative rejected:* keep one-per-artifact and only refresh on re-run (#13) — still leaves N comments on the first run.
- **Fixed headings define the layers**: `## Plan` (with `### Proposal`, `### Capabilities`, `### Design`, `### Tasks`) and `## Implementation` (with the checklist and per-task insights). *Alternative rejected:* free-form comments — not machine- or eye-scannable.
- **Comment ids live in `.openspec.yaml`** (`github.comments.plan`, `github.comments.implementation`), reused across runs. This generalises the mechanism #13 introduced for the per-artifact comments. *Alternative rejected:* `--edit-last` (ambiguous once both phases post).
- **Per-task insight is required, not optional**, so the implementation comment is worth reading. *Alternative rejected:* a bare checklist that only repeats `tasks.md`.
- **Ambiguity is asked, not guessed.** `/spec` may decide design; it may not invent requirements. Answers route to the issue body, the single source of truth.

## Risks / Trade-offs

- [One big plan comment can get long] → Keep artifact summaries tight; the full files remain in the repo and on the branch.
- [Comment id lost or a comment deleted] → If the recorded id is missing, create the comment and re-record it.
- [A run that fails midway leaves a partial comment] → Refresh on the next run overwrites it; the header carries the commit SHA so staleness is visible.

## Migration Plan

None. The next `/spec`/`/apply` run on a fresh change uses the new structure; existing archived issues keep their history.

## Open Questions

- Should the plan comment link the branch and the artifact file paths, or is the commit SHA enough?
