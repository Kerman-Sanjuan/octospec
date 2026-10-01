## Context

`/idea` and `/bug` collect the sections, then pass them to `gh issue create`. They use `--template` with `--body`, which `gh` rejects (`--template is not supported when using --body or --body-file`). YAML issue forms exist for the human web UI, and they carry the labels and the required fields, but they cannot be pre-filled from the CLI.

## Goals / Non-Goals

**Goals:**
- The documented CLI path works on the first try.
- The labels are applied on the CLI path.
- The YAML forms stay as the human contract.

**Non-Goals:**
- Change the issue sections.
- Change the web forms.

## Decisions

- **Pass the sections as the body.** `/idea` and `/bug` run `gh issue create --title ... --body-file ...`. *Alternative rejected:* `--template` with `--body`, which always errors.
- **Apply the labels on the command.** Add `--label type:feature --label status:backlog` (or `type:bug`). The form's `labels:` apply only to the web form. *Alternative rejected:* parse the YAML for its labels, which adds a YAML parser to a command.
- **Keep the YAML forms.** They remain the human web-form contract. *Alternative rejected:* delete them, which loses the UI contract.

## Risks / Trade-offs

- [The command body and the YAML form drift] -> Both use the same sections, and a test pins the command path.

## Migration Plan

None.

## Open Questions

- None.
