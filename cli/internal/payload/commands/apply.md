---
description: Implement an OpenSpec change on its branch and sync the implementation comment.
argument-hint: "[change]"
---

Runs the `apply` stage in the current session. The `apply` agent defines its scope, tools, and model.

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

7. Keep `openspec/changes/<change>/coverage.md` current: one line per requirement of the change,
   pointing at the thing that proves it.

   ```
   - <requirement title, exactly as in the spec> -> <test name or script>
   ```

   A requirement with no line, or a line pointing at a test you renamed an hour ago, is a requirement
   nobody can prove. G9 warns, and the review will ask.
8. Leave the record of what happened:

   ```
   octospec result write --stage apply --issue <issue> --change <change> --outcome ok
   octospec result write --stage apply --issue <issue> --change <change> \
     --outcome stalled --reason <needs_human|empty_diff|budget_exhausted|crashed|timeout>
   ```

   The reason is the part that matters: a stage that ran out of budget and a stage that found
   nothing to do need different retries, and `failed` for both teaches the scheduler to retry a dead
   end.

Stop on a blocker and ask.
