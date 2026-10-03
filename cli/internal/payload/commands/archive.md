---
description: Archive the change, sync the specs, close the issue, and push.
argument-hint: "[change]"
---

Runs the `archive` stage in the current session. The `archive` agent defines its scope, tools, and model.

Load the `openspec-archive-change` skill and follow it.

1. Resolve the change (`$ARGUMENTS`, or the only unarchived change) and read
   `github.issue`. Confirm the PR is merged; if it is not, stop and say so.
2. Fold the changelog entry. Read the `## Changelog` section from
   `openspec/changes/<change>/proposal.md`. If the section is missing, stop
   and ask the author to add it. If it reads `None`, skip. Otherwise insert
   its bullets under the unreleased heading in `CHANGELOG.md`, creating the
   heading when it is absent.
3. `openspec archive "<change>" --yes`
4. For every `openspec/specs/<capability>/spec.md` whose `## Purpose` still
   reads `TBD`, replace it with a one-line purpose for that capability.
5. Commit the archive and push it:
   `git add -A openspec CHANGELOG.md`
   `git commit -m "chore(#<issue>): archive <change>"`
   `git push origin main`
6. If `Closes #<issue>` did not already close it, close the issue:
   `gh issue close <issue>`
7. Remove `status:*` labels from the issue.
8. Delete the merged branch: `git branch -d feat|fix/<issue>-<slug>` (the
   prefix comes from the issue's `type:*` label). Optionally delete it
   remotely: `git push origin --delete feat|fix/<issue>-<slug>`.
9. Report the archive path and the spec files updated.

The push in step 4 lands because the branch protection lets the owner push
maintenance commits such as an archive (see `scripts/protect-main.sh`).
