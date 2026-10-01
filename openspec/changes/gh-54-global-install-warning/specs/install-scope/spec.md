## ADDED Requirements

### Requirement: Warn on a global install
octospec SHALL warn when a tool is installed into a global user directory, and SHALL explain that its commands, agents, and skills then appear in every project.

#### Scenario: A global install warns
- **WHEN** a tool is installed globally
- **THEN** octospec prints the warning that the agents and skills appear in every project.

#### Scenario: A repo-local install is silent
- **WHEN** the install is repo-local
- **THEN** no global warning is printed.

### Requirement: Choose the scope per tool
octospec SHALL let the user choose a repo-local or a global location for every tool that has a global config, and SHALL default to repo-local.

#### Scenario: Default repo-local
- **WHEN** `octospec install` runs with no scope choice
- **THEN** every tool installs repo-local.

#### Scenario: Choose global
- **WHEN** the user chooses global for a tool
- **THEN** that tool installs globally and the warning is printed.

### Requirement: Non-interactive scope
A non-interactive install SHALL NOT hang and SHALL default to repo-local; an explicit global flag SHALL select the global location.

#### Scenario: No terminal
- **WHEN** `octospec install` runs without a terminal
- **THEN** it installs repo-local and exits.

#### Scenario: Global flag
- **WHEN** `octospec install --global --tool opencode` runs
- **THEN** opencode installs globally and prints the warning.

### Requirement: Record and re-apply the scope
octospec SHALL record each tool's scope in `.octospec/octospec.json`, and `octospec update` SHALL re-apply the recorded scope.

#### Scenario: Update keeps the scope
- **WHEN** `octospec update` runs
- **THEN** each tool writes to its recorded location.

### Requirement: Offer to remove globally installed files
octospec SHALL detect the files that earlier versions installed globally and SHALL offer to remove them, as an optional, warned step.

#### Scenario: Offer cleanup
- **WHEN** a repo-local install runs and earlier global files are recorded
- **THEN** octospec offers to remove them and does not hang.
