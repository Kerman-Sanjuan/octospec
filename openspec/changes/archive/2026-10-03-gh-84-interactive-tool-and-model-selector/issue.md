# Issue #84: Interactive tool and model selector for octospec install

https://github.com/Kerman-Sanjuan/octospec/issues/84

## User story

As an operator installing octospec, I want a multi-select browser of the
supported tools, so that I choose tools without remembering or typing exact
names.

## Context / problem

The install wizard asks users to type tool names on a blank line:

```
Install octospec for which tools? [pi opencode copilot claude]
```

A mistyped name fails the run, and there is no way to browse what exists. The
models command already runs a huh TUI (`cli/internal/models/models.go`), so the
stack is present. We want the same treatment for tool selection, and a browser
for model roles instead of typing model IDs.

## Requirements

- The `octospec install` command SHALL open a multi-select list of the
  supported tools when no `--tool` flag is passed and stdin is a terminal.
- The selector MUST list every tool in `targets.Targets` by name with a short
  description.
- Space MUST toggle a tool and enter MUST confirm the selection.
- The confirmed selection MUST drive the install exactly as `--tool` values do
  today.
- An empty selection MUST install for all tools, matching the current blank
  input behaviour.
- The `octospec models` command SHALL open a select per role when no `--set`
  flags are passed and stdin is a terminal.
- The models selector MUST offer a static list of known models per tool, plus a
  free-form fallback.
- Non-terminal runs and runs with flags MUST keep the current behaviour.

## Success criteria

- `octospec install` with no flags on a terminal shows a multi-select of pi,
  opencode, copilot, claude.
- Selecting a subset installs only those tools; selecting none installs all.
- `octospec install --tool pi` and the other flag paths are unchanged.
- `octospec models` with no flags lets the operator pick a model per role from a
  known list or type one.
- Existing tests pass, and new tests cover the selection logic where it does not
  need a live terminal.

## Out of scope

- No new tools and no change to the mappings in
  `cli/internal/targets/targets.go`.
- No change to the config format in `cli/internal/config/config.go`.
- No detection of models from the tools at runtime, static list only.
- No change to the install scope (local vs global) logic.
