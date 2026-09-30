## MODIFIED Requirements

### Requirement: `octospec doctor` reports the install state
`octospec doctor` SHALL report the state of the OpenSpec CLI, the schema, the installed commands, the installed agents, the model configuration, the workflow labels, and the gate script.

#### Scenario: Healthy install
- **WHEN** `octospec doctor` runs on a healthy install
- **THEN** every check reports OK and the command exits zero.

#### Scenario: Missing piece
- **WHEN** the schema or a workflow label is missing
- **THEN** `doctor` reports it and exits non-zero.

#### Scenario: Missing agent
- **WHEN** an installed agent definition is missing
- **THEN** `doctor` reports it and exits non-zero.
