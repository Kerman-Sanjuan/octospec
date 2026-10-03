## Why

The `octospec install` wizard asks you to type tool names on a blank line, so a
typo fails the run and there is nothing to browse. The models TUI asks you to
type model IDs with no list to pick from. Both should be browsers.

## What Changes

- Replace the blank-line tool prompt with a huh multi-select of the tools in
  `targets.Targets`, each with a short description. Space toggles, enter
  confirms.
- Pre-check the tools recorded in `.octospec/octospec.json` (`cfg.Tools`) so a
  re-run starts from the previous selection.
- An empty selection installs for all tools, matching the blank input today.
- Replace the free-form model inputs with a select per role backed by a static
  known-model catalog, plus a free-form fallback for models not in the list.
- Keep flags and non-terminal runs exactly as they are.

## Capabilities

### New Capabilities
- `interactive-selection`: a multi-select tool browser for `octospec install`,
  seeded from the previous install, that drives the install like `--tool`.

### Modified Capabilities
- `agent-model-config`: the model TUI now picks from a known-model catalog per
  tool, with a free-form fallback, instead of a bare text input.

## Impact

- `cli/internal/models/models.go` and a new model catalog beside it.
- `cli/cmd/octospec/main.go`: the install wizard becomes a huh form.
- `cli/internal/install/install.go`: accepts a preselected tool set.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/84

## Changelog

### Added
- An interactive multi-select for the tools to install, seeded from the previous
  install.

### Changed
- `octospec models` offers a known-model list per tool with a free-form
  fallback instead of a bare text input.
