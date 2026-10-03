# interactive-selection Specification

## Purpose
The multi-select tool browser for `octospec install`: pick tools from the known set, seeded from the previous install, and drive the install like `--tool`.
## Requirements
### Requirement: An interactive tool selector
`octospec install` SHALL open a multi-select list of the supported tools when no `--tool` flag is passed and stdin is a terminal. The list SHALL name every tool in `targets.Targets` with a short description, space SHALL toggle a tool, and enter SHALL confirm.

#### Scenario: Present the tools
- **WHEN** `octospec install` runs with no `--tool` flag on a terminal
- **THEN** a multi-select lists pi, opencode, copilot, and claude, and enter confirms the choice.

#### Scenario: A subset drives the install
- **WHEN** the operator selects opencode and claude
- **THEN** only opencode and claude are installed.

#### Scenario: An empty choice installs all
- **WHEN** the operator confirms with nothing selected
- **THEN** every tool is installed.

### Requirement: Preselect the previous install
The tool selector SHALL pre-check the tools recorded in `.octospec/octospec.json` (`cfg.Tools`). A first install has no recorded tools, so it SHALL start empty.

#### Scenario: Re-run preselects
- **WHEN** `.octospec/octospec.json` records opencode
- **THEN** the selector opens with opencode checked.

#### Scenario: First install
- **WHEN** there is no `.octospec/octospec.json`
- **THEN** the selector opens with nothing checked.

### Requirement: Non-interactive install is unchanged
A run with no terminal, or with a `--tool` flag, SHALL NOT open the selector and SHALL keep the current behaviour.

#### Scenario: No terminal
- **WHEN** `octospec install` runs without a terminal
- **THEN** it installs for all tools and never opens the selector.

#### Scenario: Tool flag skips the selector
- **WHEN** `octospec install --tool pi` runs
- **THEN** it installs pi without opening the selector.

