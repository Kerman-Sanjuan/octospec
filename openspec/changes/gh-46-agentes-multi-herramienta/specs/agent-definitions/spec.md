## ADDED Requirements

### Requirement: A canonical agent per stage
octospec SHALL define one agent per stage of the loop (idea, spec, apply, ship, archive), each with a scoped mission, a model role, the skills it may load, and the artifact it writes.

#### Scenario: Every stage has an agent
- **WHEN** the canonical payload is listed
- **THEN** there is exactly one agent for each stage: idea, spec, apply, ship, and archive.

#### Scenario: A stage agent is scoped
- **WHEN** the spec agent is rendered
- **THEN** its file references only the spec stage command and the skills the spec stage needs.

### Requirement: One source rendered for every tool
octospec SHALL render the agent definitions for every supported tool from the single canonical source, and the repository SHALL NOT store per-tool agent copies.

#### Scenario: Render for a tool
- **WHEN** the agent definitions are rendered for a tool
- **THEN** they go to that tool's agent directory in that tool's native format.

#### Scenario: No per-tool copy
- **WHEN** a canonical agent changes
- **THEN** every tool is regenerated from it and there is no per-tool copy to edit.

### Requirement: Adding a tool is a mapping
Adding a supported tool SHALL add a target mapping only and SHALL NOT touch the agent definitions.

#### Scenario: Add a tool
- **WHEN** a new tool is added
- **THEN** only a target mapping is added and the agent definitions are unchanged.

### Requirement: The CLI orchestrates the stage
The CLI SHALL stay the orchestrator: it validates the input, selects the agent for the stage, and the agent writes the stage artifact.

#### Scenario: Run a stage
- **WHEN** a stage command runs
- **THEN** the CLI selects that stage's agent and the agent writes that stage's artifact.

### Requirement: The agents do not weaken the gates
Adding the agents SHALL NOT change the gate rules, and the gates SHALL stay green.

#### Scenario: Gates stay green
- **WHEN** the agents are added
- **THEN** `openspec validate --all --strict` passes and `cli` and `gates` stay green.

### Requirement: The developing repo follows the standard
The repository that develops octospec SHALL follow the agent standard in its own instructions, `AGENTS.md` and the workflow commands.

#### Scenario: Contributor reads the rules
- **WHEN** a contributor reads `AGENTS.md`
- **THEN** it points to the canonical agent definitions and the model configuration.
