# Issue #36: `octospec install` and `seed`: add `--dry-run` and `uninstall`

https://github.com/Kerman-Sanjuan/octospec/issues/36

## User story

As a user, I want to preview an install and undo it.

## Context / problem

`install` and `seed` write to disk immediately, and there is no way to remove what octospec manages.

## Requirements

- `octospec install` and `octospec seed` SHALL support `--dry-run`, which writes nothing and prints what it would do.
- `octospec uninstall` SHALL remove the files octospec manages, using the recorded state, and leave hand-edited files alone.

## Success criteria

- `--dry-run` changes nothing on disk; `uninstall` removes the managed files and reports any it preserved.

## Out of scope

- Removing the OpenSpec schema or global tools it did not install.
