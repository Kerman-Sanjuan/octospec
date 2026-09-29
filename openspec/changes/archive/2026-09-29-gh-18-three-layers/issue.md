# Issue #18: Collapse each issue to three structured layers

https://github.com/Kerman-Sanjuan/octospec/issues/18

## User story

As a maintainer reviewing an issue, I want exactly three clearly-structured layers — the issue (requirement), one plan comment (`/spec`), and one implementation comment (`/apply`) — so that I can read the plan and progress at a glance without wading through a stream of comments.

## Context / problem

A change today adds many comments: `/spec` posts one per artifact (proposal, design, one per capability, tasks) and `/apply` posts one per task group; re-runs append more. The timeline becomes noise. #13 made `/spec` refresh in place on re-run and #15 proposes one `/apply` checklist comment, but artifacts are still one-comment-each and no predefined structure distinguishes the layers.

## Requirements

- Each issue SHALL have at most one plan comment (`/spec`) and one implementation comment (`/apply`), in addition to the issue body.
- The plan comment SHALL consolidate proposal, capabilities/specs, design, and tasks under fixed headings.
- The implementation comment SHALL contain the task checklist with a short insight per task.
- Re-running `/spec` or `/apply` SHALL update the existing comment in place and SHALL NOT create new comments.
- Each layer SHALL have a predefined, recognizable structure (a stable header per layer).
- If the issue is ambiguous, `/spec` SHALL ask and route the answer to the issue body; it SHALL NOT invent requirements.

## Success criteria

- A change's issue has exactly one plan comment and one implementation comment.
- All artifacts/specs are readable inside the single plan comment under their headings.
- Progress (checklist + insights) is readable in the single implementation comment.
- Iterating the plan or implementation updates the same comments to the new commit SHA.

## Out of scope

- `/idea`'s interview clarity bar (tracked separately).
- The refresh-in-place mechanism itself (from #13).
