---
description: Run CI, open a pull request linked to the issue, and iterate failures back to apply.
argument-hint: "[change]"
---

Runs the `ship` stage in the current session. The `ship` agent defines its scope, tools, and model.

Open the PR for the current branch once the gates are green.

1. Read `.openspec.yaml` for `github.issue`.
2. Session check. Run `octospec session list`. If it shows a session for this
   issue, make sure you are in its worktree (the tab opened by
   `octospec session start <issue>`); work there. If `git branch
   --show-current` is a change branch for a different issue
   (`feat|fix/<other>-<slug>`), STOP. This checkout belongs to another task; do
   not switch the branch or push. Tell the user to run
   `octospec session start <issue>` or move to the right checkout. No session is
   fine: continue in the current checkout as before.
3. Push the branch: `git push -u origin HEAD`.
4. Create or update the pull request. The body MUST contain `Closes #<issue>`,
   then:
   - a short summary,
   - the change name and artifact paths,
   - the `openspec validate` result.
   `gh pr create --title "<type>(#<issue>): <summary>" --body-file <file>`
5. Run or observe CI (`gh pr checks --watch`). If any gate fails, hand the
   specific failing check back to `/apply` to fix and iterate - do not report a
   bare failure, and do not continue until it is green.
6. Once every gate is green, add label `status:in-review`; leave
   `status:in-progress` (the `/archive` step clears it).
7. Print the PR URL.
