## Why

`/spec` produces the proposed approach (proposal/specs/design/tasks), but it does not commit: the artifacts float uncommitted in the working tree of the current branch. The human's approval gate is on the approach, so the spec must be reviewable and revisable — which requires a version-controlled surface, not issue comments. This also caused #2 (the dirty-tree deadlock) and forces `/apply` to commit a baseline.

## What Changes

- `/spec` creates the change branch (`feat|fix/<issue>-<slug>`, derived from the issue's `type:*` label) and commits the artifacts on it.
- `/spec` pushes the branch and heads each published issue comment with the artifact commit SHA.
- Re-running `/spec` reuses the existing change and branch and updates the published comments in place.
- `/apply` works on the branch `/spec` created; it no longer creates the branch or commits a baseline.
- "Approach approved" becomes a human action.

## Capabilities

### New Capabilities
- `spec-branching`: where a change's plan artifacts live, how they are versioned, and how the approval gate is expressed.

### Modified Capabilities
<!-- None. -->

## Impact

- `commands/spec.md` and `commands/apply.md` (and the `.github/prompts/` + `repo-template/` mirrors).
- `schema/octospec/schema.yaml` (the `apply` instruction).
- Revisits the #2 fix (the `/apply` dirty-tree / baseline special case).
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/13
