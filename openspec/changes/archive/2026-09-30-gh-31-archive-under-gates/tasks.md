## 1. Archive commits and pushes

- [x] 1.1 Make `/archive` commit the archive (moved change, synced specs, filled Purposes) and push it (`cli/internal/payload/commands/archive.md`)
- [x] 1.2 Confirm `/archive` closes the issue, clears `status:*`, and deletes the merged branch (`cli/internal/payload/commands/archive.md`)

## 2. Protection allows maintenance

- [x] 2.1 Add `scripts/protect-main.sh` with the required checks, strict mode, `enforce_admins: false`, and force-push and deletion blocked (`scripts/protect-main.sh`)
- [x] 2.2 Apply the settings and confirm an owner push to `main` succeeds (`scripts/protect-main.sh`)

## 3. Docs

- [x] 3.1 Note the archive push and the owner-maintenance policy in the README (`README.md`)

## 4. Verify

- [x] 4.1 Dry-run `/archive` on a completed change and confirm the push lands with no manual step (`cli/internal/payload/commands/archive.md`)
