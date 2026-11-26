# Troubleshooting

## A command fails with "skill not found"

The command loads an OpenSpec skill (`openspec-propose`, and the others) that is not present for your tool.

Run `octospec install` again. It calls `openspec init --tools`, which writes the skills per tool. Then confirm with `octospec doctor`.

## `/archive` says the archive lock is held

Another archive is writing `openspec/specs/`. That is the lock doing its job: two changes folding
deltas into the same capability at once lose requirements silently.

The lock is a **lease**: it carries the time it was taken and the time it expires, and that expiry
is the only thing that decides whether the holder is still working. The pid in the file is a hint
for reading, not a liveness test — `octospec lock acquire` is a one-shot call, so its process is
always gone by the time you look.

- `octospec lock status` says who holds it, on which host, since when, and how long the lease has
  left.
- If the lease is valid, wait or let it finish. Do not clear the lock. If your own archive is the
  slow one, extend it: `octospec lock heartbeat --change <name>`.
- If `status` says **stale**, the lease expired: the holder died mid-archive. Clear it deliberately
  with `octospec lock release --force --change <name>`.
- Nothing clears an expired lease automatically, on purpose: a slow archive must never lose to an
  impatient one.

## "label not found"

The workflow labels are missing. Run `octospec seed --repo .`, which provisions them. Confirm with `octospec doctor`.

## A push to `main` is rejected

`main` is protected and requires the `cli` and `gates` checks through a pull request. Open a pull request, as `/ship` does. If you are the owner and are pushing maintenance such as an archive, re-run `scripts/protect-main.sh`, which keeps the checks but lets the owner push maintenance.

## A gate fails on a pull request

`gates` runs `scripts/check-gates.sh`: the OpenSpec validate (G2a), change completeness (G2b), the branch name (G3), the PR body links the issue (G4), and the spec delta (G6). Read the failing line, fix it, push again. `/ship` hands the failure back to `/apply`.

If the branch name is wrong, use `feat|fix/<issue>-<slug>`.

## `octospec doctor` says the schema is missing

Run `octospec seed`. It writes the schema to `${XDG_DATA_HOME:-$HOME/.local/share}/openspec/schemas/octospec/`.

## Nothing is installed, or the state is wrong

Inspect `.octospec/octospec.json`. It lists the tools and the hash of every file octospec manages. `octospec update` re-applies from it and preserves your edits.
