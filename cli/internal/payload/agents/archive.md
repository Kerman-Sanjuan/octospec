---
name: archive
stage: archive
role: reviewer
description: Archive the change, sync the specs, and close the issue.
skills: openspec-archive-change
tools: read, search, edit, shell
writes: the archived change, the synced specs, and the changelog entry
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
- Fold the change's `## Changelog` entry into `CHANGELOG.md` before the archive commit.
- Fill a real `## Purpose` for every new spec.
- Close the issue and clear its `status:*` labels.
