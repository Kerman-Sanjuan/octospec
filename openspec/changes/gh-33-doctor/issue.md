# Issue #33: Add `octospec doctor`

https://github.com/Kerman-Sanjuan/octospec/issues/33

## User story

As a user, I want to check my octospec install and see what is missing.

## Context / problem

There is no way to see the state of an install: whether the schema is present, the labels exist, the targets are wired, or the checks are configured.

## Requirements

- `octospec doctor` SHALL report the state of the schema, the workflow labels, the configured targets, and the check configuration.
- It SHALL exit non-zero when something required is missing.

## Success criteria

- Removing the schema or a label makes `doctor` report it clearly.

## Out of scope

- Auto-fixing what it finds.
