## Context

`octospec install` writes the commands. The commands load OpenSpec skills, which the install does not write, so a fresh repo cannot run `/spec`. OpenSpec generates those skills through `openspec init --tools`, which is non-interactive and writes one skill set per tool.

## Goals / Non-Goals

**Goals:**
- After `octospec install`, the OpenSpec skills exist for the selected tools.
- No drift: OpenSpec owns its skills.

**Non-Goals:**
- Re-shipping the skills from octospec.
- Changing the skills.

## Decisions

- **Call `openspec init --tools <list>` from `install`.** OpenSpec owns its skills and knows each tool's skill directory, so octospec should ask it rather than embed a copy that drifts with the OpenSpec version. *Alternative rejected:* embedding the skills in the CLI, which duplicates another project's files and rots.
- **Map the tool names.** `pi`, `opencode`, and `claude` map directly; `copilot` maps to OpenSpec's `github-copilot`.
- **Skip with a warning when `openspec` is absent.** The commands still install, and `install` says the skills were skipped.

## Risks / Trade-offs

- [`openspec init` also writes its own workspace and config] → It does not overwrite `openspec/config.yaml` in non-interactive mode; the run only adds the skill files.
- [An OpenSpec CLI change to the tool names] → The mapping lives in one place.

## Migration Plan

None. Re-run `octospec install` to add the skills.

## Open Questions

- None.
