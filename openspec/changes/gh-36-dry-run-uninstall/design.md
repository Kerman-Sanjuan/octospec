## Context

`install` and `seed` write immediately. The state in `.octospec/octospec.json` records every managed file and its content hash, which is what `update` uses to tell a managed file from a hand edit.

## Goals / Non-Goals

**Goals:**
- A dry run for install and seed.
- An uninstall that removes only what octospec manages and was not hand-edited.

**Non-Goals:**
- Removing the OpenSpec schema or a tool octospec did not install.

## Decisions

- **`--dry-run` prints the plan and returns before any write.** Both commands already compute what they would write, so the dry run reuses that and writes nothing. *Alternative rejected:* a separate `plan` command, which duplicates the surface.
- **Uninstall uses the recorded hashes.** A file whose content still matches its recorded hash is removed; a file that differs is kept and reported, the same rule as `update`. *Alternative rejected:* removing every recorded path, which would delete a user's edits.
- **Uninstall removes the state file when it is done.** With nothing managed, the state would be misleading. *Alternative rejected:* keep the state, which leaves a stale record.

## Risks / Trade-offs

- [Uninstall could remove a file the user edited to match the recorded content] -> That is the same content octospec wrote, so removal is expected.
- [Uninstall leaves empty directories] -> It prunes the directories it empties.

## Migration Plan

None.

## Open Questions

- Should uninstall also remove the OpenSpec skills it installed? The issue says out of scope, so no.
