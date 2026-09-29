---
mode: agent
description: Archive the change, update the specs, and close the issue.
---


Load the `openspec-archive-change` skill and follow it.

1. Resolve the change (`$ARGUMENTS`, or the only unarchived change) and read
   `github.issue`. Confirm the PR is merged; if it is not, stop and say so.
2. `openspec archive "<change>" --yes`
3. For every `openspec/specs/<capability>/spec.md` whose `## Purpose` still
   reads `TBD`, replace it with a one-line purpose for that capability.
4. If `Closes #<issue>` did not already close it, close the issue:
   `gh issue close <issue>`
5. Remove `status:*` labels from the issue.
6. Delete the merged branch: `git branch -d feat|fix/<issue>-<slug>` (the
   prefix comes from the issue's `type:*` label). Optionally delete it
   remotely: `git push origin --delete feat|fix/<issue>-<slug>`.
7. Report the archive path and the spec files updated.
