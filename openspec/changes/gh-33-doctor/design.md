## Context

An octospec install touches the global schema, the per-tool commands, the repo labels, and the gate script. A missing piece surfaces late, at the point of use. A single check command makes the state visible.

## Goals / Non-Goals

**Goals:**
- One command that reports each piece and exits non-zero on a problem.
- Read-only.

**Non-Goals:**
- Fixing anything, which stays with `install`, `seed`, and `update`.

## Decisions

- **Report, do not repair.** A doctor that mutates is a second installer. It only reads and exits non-zero. *Alternative rejected:* an auto-fix flag, which hides the cause.
- **Check labels through `gh` in the repo.** Labels are per-repo, so the check runs `gh label list` in the target and compares against the six workflow labels.

## Risks / Trade-offs

- [The labels check needs `gh` and a GitHub repo] → When `gh` is absent or the repo has no remote, the check reports it instead of failing hard.

## Migration Plan

None. New command.

## Open Questions

- None.
