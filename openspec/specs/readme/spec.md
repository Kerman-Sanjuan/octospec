# readme Specification

## Purpose
What the octospec README must convey: what the project is, the workflow loop, the architecture, installation for each agent, and the full command surface.
## Requirements
### Requirement: README explains what octospec is
The README SHALL explain what octospec is and the problem it solves.

#### Scenario: New reader understands the project
- **WHEN** a first-time reader opens the README
- **THEN** they can state what octospec is and the problem it solves without opening the design docs.

### Requirement: README describes the workflow loop
The README SHALL describe the workflow loop (`/idea → /spec → /apply → /ship → /archive`).

#### Scenario: Reader follows the loop
- **WHEN** a reader reads the workflow section
- **THEN** they can identify the `/idea → /spec → /apply → /ship → /archive` loop and each phase's purpose.

### Requirement: README documents the architecture
The README SHALL document the architecture: schema override, command shims, seed files, install script, and CI gates.

#### Scenario: Architecture is explained
- **WHEN** a reader reads the architecture section
- **THEN** they can identify the schema override, command shims, seed files, install script, and CI gates and their roles.

### Requirement: README provides installation instructions
The README SHALL provide installation instructions for pi, opencode, and GitHub Copilot.

#### Scenario: Fresh clone can install
- **WHEN** a developer follows the install instructions from a fresh clone
- **THEN** they can install octospec for pi, opencode, and GitHub Copilot.

### Requirement: README lists every command
The README SHALL list every command with its phase and behavior.

#### Scenario: Command surface is documented
- **WHEN** a reader consults the commands section
- **THEN** every command (`/idea`, `/bug`, `/explore`, `/spec`, `/apply`, `/ship`, `/archive`) is listed with its phase and behavior.

### Requirement: README links to deeper docs
The README SHALL link to the design doc and implementation plan as deeper references, and its prose SHALL be complete and readable, with no truncated sentences or stray backticks.

#### Scenario: Links resolve
- **WHEN** a reader follows the links in the README
- **THEN** each link resolves to an existing file.

#### Scenario: The skills paragraph reads
- **WHEN** a reader reads the OpenSpec skills paragraph
- **THEN** the sentences are complete, with closed backticks and no line that begins with a stray character.

### Requirement: README documents the agents and the model choice
The README SHALL describe the agent that covers each stage of the loop and how to choose the model each role uses.

#### Scenario: Reader chooses a model
- **WHEN** a reader looks for how to set the model an agent uses
- **THEN** the README points to the model configuration and the interactive TUI.

### Requirement: The install section lists platforms and channels
The README SHALL list the supported platforms and every install channel.

#### Scenario: A reader picks a platform
- **WHEN** a reader opens the install section
- **THEN** they can see the supported platforms and the install channels.

### Requirement: The README states the supported platforms
The README SHALL state the supported operating systems and that native Windows is not supported (WSL works).

#### Scenario: A Windows user reads the README
- **WHEN** a Windows user looks for install instructions
- **THEN** the README says native Windows is not supported, that WSL works, and why.

### Requirement: A quickstart at the top
The README SHALL show a copy-paste quickstart before the deep explanation.

#### Scenario: A newcomer tries octospec
- **WHEN** a newcomer opens the README
- **THEN** a copy-paste quickstart shows how to install the CLI, install the commands, and run the loop.

### Requirement: The prerequisites are complete
The README SHALL list every prerequisite: the OpenSpec CLI, `gh`, and Node, and SHALL state the OpenSpec version the gates are tested against without claiming support for untested newer versions.

#### Scenario: The prerequisites are listed
- **WHEN** a reader looks for the prerequisites
- **THEN** the OpenSpec CLI, `gh`, and Node are named.

#### Scenario: The OpenSpec version is accurate
- **WHEN** a reader reads the OpenSpec prerequisite
- **THEN** it names the tested version (1.3.1) and does not say "or newer", because the strict gates fail on OpenSpec 1.14.0.

