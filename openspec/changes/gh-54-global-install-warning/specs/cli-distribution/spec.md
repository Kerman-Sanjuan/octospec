## MODIFIED Requirements

### Requirement: Target mapping
Each supported tool SHALL be defined by a target mapping with a repo-local location and a global location (directory plus file format for commands and for agents); adding a tool SHALL NOT add a source copy.

#### Scenario: Adding a tool
- **WHEN** a new tool is added
- **THEN** only a new mapping is added and the payload is unchanged.

### Requirement: Supported targets
The CLI SHALL support pi, opencode, GitHub Copilot, and Claude Code, and SHALL offer a repo-local location and a global location for each of them.

#### Scenario: Install per tool
- **WHEN** `octospec install --tool <pi|opencode|copilot|claude>` runs
- **THEN** the commands are written to that tool's chosen directory.

## ADDED Requirements

### Requirement: A repo-local install writes nothing under home
A repo-local install SHALL NOT write under the user's home directory.

#### Scenario: Nothing under home
- **WHEN** every tool is installed repo-local
- **THEN** no file is written under the home directory.

### Requirement: pi has no agent files
octospec SHALL NOT install agent files for pi, because pi has no agent mechanism.

#### Scenario: pi install
- **WHEN** pi is installed
- **THEN** its commands and skills are written and no agent file is written.
