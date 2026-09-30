## ADDED Requirements

### Requirement: A single model configuration
octospec SHALL store the model for each role in one configuration point, and SHALL NOT require the model to be set per agent or per tool.

#### Scenario: Change in one place
- **WHEN** the model for a role changes in the configuration
- **THEN** every agent for that role uses the new model on the next render.

#### Scenario: No per-tool duplication
- **WHEN** a model is set for a role
- **THEN** no per-tool file needs an edit.

### Requirement: Rendered agents carry the configured model
Each rendered agent SHALL carry the model configured for its role.

#### Scenario: Render carries the model
- **WHEN** an agent is rendered
- **THEN** its file names the model configured for its role, or omits the model when the role has none.

### Requirement: The model roles
octospec SHALL define the roles `thinking`, `implementer`, and `reviewer`, and SHALL map each stage to one role.

#### Scenario: Default role mapping
- **WHEN** the default configuration is used
- **THEN** idea and spec map to thinking, apply maps to implementer, and ship and archive map to reviewer.

### Requirement: An interactive TUI for the models
octospec SHALL provide an interactive TUI that sets the model of each role.

#### Scenario: Set a model in the TUI
- **WHEN** the user opens the model TUI
- **THEN** they pick a model for each role and the choice is saved to the configuration.

### Requirement: Non-interactive use
The model configuration SHALL be settable without a terminal.

#### Scenario: Set a model without a TUI
- **WHEN** the command runs with a flag and no terminal
- **THEN** the configuration is written without opening the TUI.
