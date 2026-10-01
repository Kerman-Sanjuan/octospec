## Context

`scripts/check-gates.sh` runs G2a, G2b, G3, G4, and G6. G3 checks the branch name against `feat|fix/<issue>-<slug>`. That check only makes sense on a PR branch, but the script runs it on any ref, so on `main` it fails and the script exits 1. G4 and G6 already skip when they do not apply; G3 does not.

## Goals / Non-Goals

**Goals:**
- The script runs on `main` and reports the branch gate as skipped.
- G3 still runs on a PR branch.

**Non-Goals:**
- Change the other gates.
- Change the gate rules.

## Decisions

- **Skip G3 on a long-lived branch.** Treat `main` and `master` as not a PR ref. *Alternative rejected:* a `--pr-only` flag, which every caller must pass and which is easy to forget.
- **Report the skip.** Print `SKIP G3 ...` so the output is honest. *Alternative rejected:* pass silently, which hides that the gate did not run.
- **Test the skip.** A shell test asserts the skip on `main` and the check on a PR branch, and it runs in the gates workflow. *Alternative rejected:* no test, which lets the behavior drift.

## Risks / Trade-offs

- [A new default branch name is not covered] -> The list is `main` and `master`; extend it if the default branch changes.
- [The gate test depends on the repo state] -> It asserts only the G3 line, not the whole script.

## Migration Plan

None.

## Open Questions

- None.
