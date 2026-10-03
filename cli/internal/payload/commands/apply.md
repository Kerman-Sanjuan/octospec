---
description: Implement an OpenSpec change on its branch and sync the implementation comment.
argument-hint: "[change]"
---

Runs the `apply` stage in the current session. The `apply` agent defines its scope, tools, and model.

Load the `openspec-apply-change` skill and follow it.

1. Resolve the change (`$ARGUMENTS`, or the only unarchived change).
   Read `.openspec.yaml` for `github.issue`; it must be set.
2. Session check. Run `octospec session list`. If it shows a session for this
   issue, make sure you are in its worktree (the tab opened by
   `octospec session start <issue>`); work there. If `git branch
   --show-current` is a change branch for a different issue
   (`feat|fix/<other>-<slug>`), STOP. This checkout belongs to another task; do
   not switch the branch or commit. Tell the user to run
   `octospec session start <issue>` or move to the right checkout. No session is
   fine: continue in the current checkout as before.
3. Ensure you are on the change branch `/spec` created:
   `feat|fix/<issue>-<slug>` (prefix from the issue's `type:*` label). The
   branch already holds the committed, pushed artifacts and a clean working
   tree. Do not create a branch and do not commit a baseline.
4. Set the issue label: add `status:in-progress`, remove `status:spec-ready`.
5. Work the tasks in `tasks.md` group by group. Commit each group with
   `<type>(#<issue>): <summary>` and push after each commit.
6. Keep a single **implementation comment** on the issue, headed
   `## Implementation` and the group commit SHA. It lists every task, with a
   short insight on each done task - what changed and how it was verified - so
   the comment is worth reading, not just a copy of `tasks.md`. Record its id
   in `.openspec.yaml` as `github.comments.implementation` and update it in
   place after each group instead of posting a new comment:
   `gh api --method PATCH /repos/{owner}/{repo}/issues/comments/<id> -F body=@<file>`
7. Run the project's build and tests before reporting a group done. Report a
   failure as a failure.

Stop on a blocker and ask.
