---
description: Archive the change, sync the specs, close the issue, and push.
argument-hint: "[change]"
---

Runs the `archive` stage in the current session. The `archive` agent defines its scope, tools, and model.

Load the `openspec-archive-change` skill and follow it.

This is the only stage that writes shared state: it folds your change's deltas into
`openspec/specs/`, which every other change reads. Take the lock before you read or write
anything, and release it even when you fail.

1. Resolve the change (`$ARGUMENTS`, or the only unarchived change) and read
   `github.issue`. Confirm the PR is merged; if it is not, stop and say so.
2. Take the archive lock for this change:
   `octospec lock acquire --change <change> --wait 2m`.
   If it fails, write nothing: report `stalled / needs_human`, name the change that holds the
   lock, and give the command that clears a stale one. Never archive around the lock.
3. Fold the changelog entry. Read the `## Changelog` section from
   `openspec/changes/<change>/proposal.md`. If the section is missing, stop
   and ask the author to add it. If it reads `None`, skip. Otherwise insert
   its bullets under the unreleased heading in `CHANGELOG.md`, creating the
   heading when it is absent.
4. `openspec archive "<change>" --yes`
5. For every `openspec/specs/<capability>/spec.md` whose `## Purpose` still
   reads `TBD`, replace it with a one-line purpose for that capability.
6. Check the scope of the archive before committing: `git status --porcelain
   openspec/specs` must list only the capabilities your change declares in
   `proposal.md`. If it touches anything else, stop, release the lock, and
   name the capability you did not declare. An archive that edits somebody
   else's spec is a bug, not a review comment.
7. Commit the archive and push it:
   `git add -A openspec CHANGELOG.md`
   `git commit -m "chore(#<issue>): archive <change>"`
   `git push origin main`
8. Release the lock: `octospec lock release --change <change>`. If you fail at
   any point after step 2, release it before returning.
9. If `Closes #<issue>` did not already close it, close the issue:
   `gh issue close <issue>`
10. Remove `status:*` labels from the issue.
11. Delete the merged branch: `git branch -d feat|fix/<issue>-<slug>` (the
    prefix comes from the issue's `type:*` label). Optionally delete it
    remotely: `git push origin --delete feat|fix/<issue>-<slug>`.
12. Report the archive path, the spec files updated, and that the lock is free.

The push in step 7 lands because the branch protection lets the owner push
maintenance commits such as an archive (see `scripts/protect-main.sh`).
