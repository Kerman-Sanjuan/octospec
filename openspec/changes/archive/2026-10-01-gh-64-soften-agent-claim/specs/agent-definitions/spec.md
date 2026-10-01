## MODIFIED Requirements

### Requirement: The CLI orchestrates the stage
The CLI SHALL stay the orchestrator: it validates the input and runs the stage in the current session. The stage agent definition SHALL be the stage's contract for its scope, tools, and model; the command SHALL NOT claim to dispatch to a subagent.

#### Scenario: Run a stage
- **WHEN** a stage command runs
- **THEN** it runs the stage in the current session under the stage agent's scope and tools.

#### Scenario: No dispatch claim
- **WHEN** a command and the agent standard are read
- **THEN** they do not claim that the stage is dispatched to a subagent.
