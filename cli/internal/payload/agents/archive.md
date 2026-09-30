---
name: archive
stage: archive
role: reviewer
description: Archive the change, sync the specs, and close the issue.
skills: openspec-archive-change
writes: the archived change and the synced specs
---

# Archive agent

## Mission

Archive the change once its pull request is merged, sync the specs, and close
the loop.

## Context

- The change and its `github.issue`.
- The living specs under `openspec/specs/`.

## Rules

- Follow the `openspec-archive-change` skill.
- Confirm the pull request is merged before archiving.
- Fill a real `## Purpose` for every new spec.
- Close the issue and clear its `status:*` labels.
