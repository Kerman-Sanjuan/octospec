## Context

The root `.github/workflows/openspec.yml` sets `GH_TOKEN: ${{ github.token }}`
and `issues: read`; the seeded copy under
`cli/internal/seed/repo/.github/workflows/openspec.yml` does not. Commit
`6701cff` fixed the root only. On a freshly seeded repo, `gh issue view` fails
for lack of auth, G1 reads an empty body, and reports a false negative:
`FAIL G1 issue #1 missing section(s): ...`. `check-gates-test.sh` injects a
`gh` that always succeeds, and the e2e runs the gates with no change in the
diff, so G1 never exercises this path.

## Goals / Non-Goals

**Goals:**
- The seeded workflow authenticates `gh`, like the root.
- G1 cannot false-fail when `gh` cannot read the issue.
- A test fails if the seeded workflow drifts again.

**Non-Goals:**
- No change to the gate set or its rules beyond G1's failure handling.
- No new CLI features.

## Decisions

### Sync the seeded workflow and enforce it with a test
Copy the token, the `issues: read` permission, and the `setup-node` major into
the seed, and add a seed test that reads both workflows, strips
`actions/*@` version comments, and fails when the token, the permission, or the
major differ. A drift test is the durable fix; a one-off copy is what regressed.
The alternative, a symlink or a shared file, is not possible across the
embedded seed and the repository's own `.github`.

### Distinguish "cannot read" from "incomplete"
In `scripts/check-gates.sh`, capture `gh issue view` separately from the body:
run it, keep its exit code, and treat a non-zero exit as a SKIP, not a FAIL. An
unread issue is not an incomplete one. The alternative, treating an empty body
as incomplete, is exactly the false negative. The seeded copy stays
byte-identical, guarded by the existing `TestSeededGatesMatchRoot`.

### A gate test for the unauthenticated path
Add a case to `scripts/check-gates-test.sh` whose fake `gh` exits non-zero with
no body, asserting G1 reports SKIP and the run does not fail. This closes the gap
the audit found: the existing fake always succeeds.

## Risks / Trade-offs

- [A genuinely missing issue would now SKIP instead of FAIL] -> A missing issue
  already returns an error from `gh`, so it is a "cannot read" case; failing it
  would blame the user for a link that points at nothing, which the G8/G1
  coverage already surfaces elsewhere. The SKIP message names the reason.
- [The seed test duplicates logic] -> It reads two small files and checks a few
  fields, which is cheaper than letting a consumer hit the false negative.

## Migration Plan

No data migration. Consumer repos pick up the fixed seed on their next
`octospec seed`/`update`. Rollback is a revert.

## Open Questions

- None. The audit reproduces the failure and names the files.
