## Why

`/idea` and `/bug` instruct `gh issue create --template feature.yml --body "<fields>"`, but `gh` rejects `--template` together with `--body`. The documented primary path always errors, so every run depends on the fallback.

## What Changes

- Make "pass the sections as the body" the primary path in `/idea` and `/bug`.
- Apply the labels explicitly on the command path.
- Keep `feature.yml` and `bug.yml` as the human web-form contract.
- Add a regression test that fails if a command uses `--template` with `--body`.

## Capabilities

### New Capabilities
- `issue-intake`: how `/idea` and `/bug` create the issue from the body and apply the labels.

### Modified Capabilities
<!-- None. -->

## Impact

- `cli/internal/payload/commands/idea.md`, `cli/internal/payload/commands/bug.md`.
- `cli/internal/payload/embed_test.go` and `docs/commands.md`.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/9
