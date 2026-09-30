# consumer-docs Specification

## Purpose
The adopter-facing docs: getting started, a command reference, troubleshooting, and the README links to them.
## Requirements
### Requirement: Getting started for adopters
The repository SHALL provide a getting started guide that installs the CLI, seeds a repository, and runs the loop.

#### Scenario: A new user follows the guide
- **WHEN** a new user follows the getting started guide
- **THEN** they install octospec, seed a repository, and run one change through the loop.

### Requirement: A command reference
The repository SHALL provide a command reference that lists the CLI commands and the workflow commands with their purpose.

#### Scenario: Commands are documented
- **WHEN** a reader opens the command reference
- **THEN** every CLI command and every workflow command is listed with its purpose.

### Requirement: Troubleshooting
The repository SHALL provide a troubleshooting section covering the common failures: a missing skill, a missing label, a rejected push, and a failing gate.

#### Scenario: A known failure is explained
- **WHEN** a reader hits a common failure
- **THEN** the troubleshooting section explains the cause and the fix.

### Requirement: The README links the docs
The README SHALL link the getting started, command reference, and troubleshooting documents.

#### Scenario: Docs reachable from the README
- **WHEN** a reader opens the README
- **THEN** they can reach all three documents from it.

### Requirement: Docs cover the agents and the model TUI
The command reference SHALL list the agent definitions and the command that opens the model TUI, and the getting started guide SHALL show how to set a model.

#### Scenario: Command reference lists the model TUI
- **WHEN** a reader opens the command reference
- **THEN** the agents and the model TUI command are listed with their purpose.

#### Scenario: Getting started sets a model
- **WHEN** a new user follows the getting started guide
- **THEN** they set the model for at least one role.

