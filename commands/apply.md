---
description: Implement an OpenSpec change on a branch and sync the issue checklist.
argument-hint: "[change]"
---

Load the `openspec-apply-change` skill and follow it.

1. Resolve the change (`$ARGUMENTS`, or the only unarchived change).
   Read `.openspec.yaml` for `github.issue`; it must be set.
2. Create and switch to the branch `feat/<issue>-<slug>` for a feature or
   `fix/<issue>-<slug>` for a bug. The change's own artifacts
   (`openspec/changes/<change>/`) are expected to be uncommitted from
   `/spec`: carry them onto the branch and commit them as the baseline
   (`chore(#<issue>): add spec artifacts`). Refuse only if the working tree
   has changes outside `openspec/changes/<change>/`.
3. Work the tasks in `tasks.md` group by group. Commit each group with
   `<type>(#<issue>): <summary>`.
4. After each group, tick the tasks and post the updated checklist as a new
   comment on the issue, headed with the commit SHA.
5. Run the project's build and tests before reporting a group done. Report a
   failure as a failure.

Stop on a blocker and ask.
