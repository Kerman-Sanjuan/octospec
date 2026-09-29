---
mode: agent
description: Archive the change, update the specs, and close the issue.
---


Load the `openspec-archive-change` skill and follow it.

1. Resolve the change (`$ARGUMENTS`, or the only unarchived change) and read
   `github.issue`. Confirm the PR is merged; if it is not, stop and say so.
2. `openspec archive "<change>" --yes`
3. If `Closes #<issue>` did not already close it, close the issue:
   `gh issue close <issue>`
4. Remove `status:*` labels from the issue.
5. Report the archive path and the spec files updated.
