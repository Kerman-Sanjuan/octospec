## ADDED Requirements

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
