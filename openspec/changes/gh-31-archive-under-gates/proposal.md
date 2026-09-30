## Why

`/archive` commits the archive and the synced specs, then pushes to `main`. With the `cli` and `gates` checks required and admins enforced, that push is rejected, so the last step of the loop is manual. `/archive` also fails to commit or push in the paths that only move files.

## What Changes

- `/archive` commits the archive (the moved change, the synced specs, and the filled Purposes) and pushes it.
- The branch-protection settings allow the owner to push maintenance commits such as an archive, so the push is not rejected.
- The protection still requires `cli` and `gates` on pull requests.

## Capabilities

### New Capabilities
- `archive-flow`: how a completed change is archived, synced, closed, and pushed.

### Modified Capabilities
<!-- None. -->

## Impact

- `cli/internal/payload/commands/archive.md` (the canonical command).
- `scripts/protect-main.sh` (new) and the branch-protection settings.
- `README.md`.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/31
