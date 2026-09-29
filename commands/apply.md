---
description: Implement an OpenSpec change on its branch and sync the implementation comment.
argument-hint: "[change]"
---

Load the `openspec-apply-change` skill and follow it.

1. Resolve the change (`$ARGUMENTS`, or the only unarchived change).
   Read `.openspec.yaml` for `github.issue`; it must be set.
2. Ensure you are on the change branch `/spec` created:
   `feat|fix/<issue>-<slug>` (prefix from the issue's `type:*` label). The
   branch already holds the committed, pushed artifacts and a clean working
   tree. Do not create a branch and do not commit a baseline.
3. Set the issue label: add `status:in-progress`, remove `status:spec-ready`.
4. Work the tasks in `tasks.md` group by group. Commit each group with
   `<type>(#<issue>): <summary>` and push after each commit.
5. Keep a single **implementation comment** on the issue, headed
   `## Implementation` and the group commit SHA. It lists every task, with a
   short insight on each done task - what changed and how it was verified - so
   the comment is worth reading, not just a copy of `tasks.md`. Record its id
   in `.openspec.yaml` as `github.comments.implementation` and update it in
   place after each group instead of posting a new comment:
   `gh api --method PATCH /repos/{owner}/{repo}/issues/comments/<id> -F body=@<file>`
6. Run the project's build and tests before reporting a group done. Report a
   failure as a failure.

Stop on a blocker and ask.
