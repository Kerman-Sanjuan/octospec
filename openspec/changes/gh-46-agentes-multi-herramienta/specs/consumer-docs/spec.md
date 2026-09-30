## ADDED Requirements

### Requirement: Docs cover the agents and the model TUI
The command reference SHALL list the agent definitions and the command that opens the model TUI, and the getting started guide SHALL show how to set a model.

#### Scenario: Command reference lists the model TUI
- **WHEN** a reader opens the command reference
- **THEN** the agents and the model TUI command are listed with their purpose.

#### Scenario: Getting started sets a model
- **WHEN** a new user follows the getting started guide
- **THEN** they set the model for at least one role.
