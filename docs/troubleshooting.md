# Troubleshooting

## A command fails with "skill not found"

The command loads an OpenSpec skill (`openspec-propose`, and the others) that is not present for your tool.

Run `octospec install` again. It calls `openspec init --tools`, which writes the skills per tool. Then confirm with `octospec doctor`.

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
