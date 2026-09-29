# Issue #13: `/spec` owns the change branch and commits the spec artifacts

https://github.com/Kerman-Sanjuan/octospec/issues/13

## User story

As a maintainer, I want `/spec` to commit the change artifacts on the change branch, so that I can review and revise the proposed approach (the *how*) with a diff and history before any implementation.

## Context / problem

The issue is a **requirement** (what/why). `/spec` produces a **proposed approach** — proposal, specs, design, tasks (the *how*). The human's approval gate is on the approach, so "not happy with the approach" is expected, and the spec must survive revision. Today `/spec` does not commit: the artifacts float uncommitted in the working tree of the current branch. That (a) caused #2 (the dirty-tree deadlock), (b) forces `/apply` to commit a baseline, and (c) leaves the published issue comments headed with only a base SHA, not the revision they represent. Issue comments have no diff and no history, so they cannot be the review/revise surface.

## Requirements

- `/spec` SHALL create the change branch (`feat/<issue>-<slug>` for a feature, `fix/<issue>-<slug>` for a bug) and commit the change artifacts on it.
- `/spec` SHALL derive `feat` / `fix` from the issue's `type:*` label (`type:feature` → `feat`, `type:bug` → `fix`).
- `/spec` SHALL push the branch.
- The published issue comments SHALL be headed with the commit SHA of the artifacts they were generated from.
- Re-running `/spec` SHALL update the existing published comments in place and SHALL NOT append new ones.
- `/spec` SHALL be idempotent: re-running reuses the existing change and branch.
- `/apply` SHALL use the branch created by `/spec`; it SHALL NOT create the branch or commit a baseline.
- The "approach approved" signal SHALL be a human action, not something `/spec` sets automatically.

## Success criteria

- After `/spec <issue>`, the change artifacts are committed on `feat|fix/<issue>-<slug>` and pushed.
- The published comments are headed with the artifact commit SHA.
- Editing the artifacts and re-running `/spec` updates the same comments to the new SHA.
- `/apply` starts on the existing branch with a clean working tree.
- `main` is untouched until the PR merges.

## Out of scope

- Redesigning the published-copy mechanism beyond refresh-in-place.
- Changing the label lifecycle (see #3).
