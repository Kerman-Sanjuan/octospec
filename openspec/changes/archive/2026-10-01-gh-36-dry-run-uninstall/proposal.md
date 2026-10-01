## Why

`install` and `seed` write to disk with no preview and no undo. A user who wants to try octospec, or to back out, has to remove files by hand.

## What Changes

- `octospec install` and `octospec seed` gain `--dry-run`, which writes nothing and prints what they would do.
- A new `octospec uninstall` removes the files octospec manages, using the recorded state, and leaves hand-edited files alone.

## Capabilities

### New Capabilities
<!-- None. -->

### Modified Capabilities
- `cli-distribution`: install and seed support a dry run, and an uninstall command removes the managed files.

## Impact

- `cli/internal/install/`, `cli/internal/seed/`, a new `cli/internal/uninstall/`, and `cli/cmd/octospec/main.go`.
- `README.md` and `docs/commands.md`.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/36
