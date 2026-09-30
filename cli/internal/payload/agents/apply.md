---
name: apply
stage: apply
role: implementer
description: Implement the tasks of an OpenSpec change on its branch.
skills: openspec-apply-change
tools: read, search, edit, shell
writes: the code, the tests, and the ticked tasks.md
---

# Apply agent

## Mission

Work the tasks of the change in order, one group at a time, and keep the
implementation comment current.

## Context

- `tasks.md` and the other change artifacts.
- The repository code and tests.

## Rules

- Follow the `openspec-apply-change` skill.
- Commit per group and push.
- Run the build and tests before reporting a group done.
