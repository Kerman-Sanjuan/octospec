# Issue #12: `install.sh` does not provision the `type:*` and `status:*` labels

https://github.com/Kerman-Sanjuan/octospec/issues/12

## User story

As a maintainer seeding octospec into a fresh repository, I want `install.sh --repo` to provision the workflow labels, so that `/idea`, `/spec`, and `/apply` work without creating labels by hand.

## Context / problem

The workflow uses six labels: `type:feature`, `type:bug`, `status:backlog`, `status:spec-ready`, `status:in-progress`, `status:in-review`. Neither `install.sh` nor `repo-template/` creates them, so they had to be created manually in this repository. On a fresh repo the commands fail until the labels exist, so the label lifecycle cannot run.

> **Impact:** broken out-of-the-box experience.

### Steps to reproduce

1. Seed a fresh repo: `./install.sh --repo <path>`.
2. Run `/idea`, `/spec`, or `/apply`.
3. `gh issue edit --add-label status:...` fails with `'status:...' not found`.

## Requirements

- `install.sh --repo` SHALL provision the six workflow labels.
- Provisioning SHALL be idempotent: re-running does not fail when a label already exists.

## Success criteria

- After `./install.sh --repo <path>`, all six labels exist in the target repository.
- `/idea`, `/spec`, and `/apply` run on a fresh repo without a "label not found" error.

## Out of scope

- Changing the label names, colors, or the status lifecycle.
