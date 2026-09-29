---
description: Implement an OpenSpec change on its branch and sync the issue checklist.
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
5. After each group, tick the tasks and post the updated checklist as a new
   comment on the issue, headed with the commit SHA.
6. Run the project's build and tests before reporting a group done. Report a
   failure as a failure.

Stop on a blocker and ask.
