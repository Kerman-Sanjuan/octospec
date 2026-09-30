## MODIFIED Requirements

### Requirement: Install from one canonical payload
The CLI SHALL install the commands and the agent definitions for a target tool from one canonical payload embedded in the binary, and the repo SHALL NOT store per-tool derived copies.

#### Scenario: One source, many destinations
- **WHEN** the canonical payload changes
- **THEN** every target is regenerated from it with no per-tool copy to edit.

### Requirement: Target mapping
Each supported tool SHALL be defined by a target mapping (destination directory plus file format for commands and for agents); adding a tool SHALL NOT add a source copy.

#### Scenario: Adding a tool
- **WHEN** a new tool is added
- **THEN** only a new mapping is added and the payload is unchanged.

## ADDED Requirements

### Requirement: Install and update render the agents
The CLI SHALL write the agent definitions for the selected tools on install and update, and SHALL preserve local edits to them like any managed file.

#### Scenario: Install writes the agents
- **WHEN** `octospec install --tool <pi|opencode|copilot|claude>` runs
- **THEN** the agent definitions are written to that tool's agent location.

#### Scenario: Hand edit to an agent survives update
- **WHEN** a managed agent file is edited locally and `octospec update` runs
- **THEN** the edit is preserved and reported.
