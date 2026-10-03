---
description: Archive the change, sync the specs, close the issue, and push.
argument-hint: "[change]"
---

Runs the `archive` stage in the current session. The `archive` agent defines its scope, tools, and model.

Load the `openspec-archive-change` skill and follow it.

1. Resolve the change (`$ARGUMENTS`, or the only unarchived change) and read
   `github.issue`. Confirm the PR is merged; if it is not, stop and say so.
2. Session check. Run `octospec session list`. If it shows a session for this
   issue, make sure you are in its worktree (the tab opened by
   `octospec session start <issue>`); work there. If `git branch
   --show-current` is a change branch for a different issue
   (`feat|fix/<other>-<slug>`), STOP. This checkout belongs to another task; do
   not switch the branch or push. Tell the user to run
   `octospec session start <issue>` or move to the right checkout.
3. Fold the changelog entry. Read the `## Changelog` section from
   `openspec/changes/<change>/proposal.md`. If the section is missing, stop
   and ask the author to add it. If it reads `None`, skip. Otherwise insert
   its bullets under the unreleased heading in `CHANGELOG.md`, creating the
   heading when it is absent.
4. `openspec archive "<change>" --yes`
5. For every `openspec/specs/<capability>/spec.md` whose `## Purpose` still
   reads `TBD`, replace it with a one-line purpose for that capability.
6. Commit the archive:
   `git add -A openspec CHANGELOG.md`
   `git commit -m "chore(#<issue>): archive <change>"`
7. Push the archive to `main` with a bounded rebase-retry, so a second
   archive landing first does not lose this one:
   ```
   for i in 1 2 3; do
     git fetch origin main
     git rebase origin/main || { echo "rebase conflict; resolve by hand"; exit 1; }
     git push origin HEAD:main && break
   done
   ```
   If the rebase conflicts, STOP and hand the conflict to the human. Never
   force-push.
8. If `Closes #<issue>` did not already close it, close the issue:
   `gh issue close <issue>`
9. Remove `status:*` labels from the issue.
10. Delete the merged branch: `git branch -d feat|fix/<issue>-<slug>` (the
    prefix comes from the issue's `type:*` label). Optionally delete it
    remotely: `git push origin --delete feat|fix/<issue>-<slug>`. If a session
    exists for the issue, end it: `octospec session end <issue>`.
11. Report the archive path and the spec files updated.

The push in step 7 lands because the branch protection lets the owner push
maintenance commits such as an archive (see `scripts/protect-main.sh`).
