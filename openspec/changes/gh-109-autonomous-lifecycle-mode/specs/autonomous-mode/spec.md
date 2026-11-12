## ADDED Requirements

### Requirement: Autonomous mode runs the full loop unattended
The autonomous mode SHALL drive an issue through all five stages (idea, spec, apply, ship, archive) in sequence without asking for human confirmation between stages.

#### Scenario: Full loop without confirmation
- **WHEN** the autonomous mode starts on an issue
- **THEN** it runs idea, spec, apply, ship, and archive in order for that issue, moving from one stage to the next on its own once each stage completes, until a pull request exists or the issue is archived.

### Requirement: Issue selection follows backlog status and respects the product definition
The autonomous mode SHALL select issues to work from the backlog's active statuses (such as `status:in-progress`), and SHALL identify bugs found while working, always checked against the existing specs in `openspec/specs/` rather than invented requirements.

#### Scenario: Picks an active issue
- **WHEN** the autonomous mode looks for work
- **THEN** it selects an issue already in an active status in the GitHub backlog, instead of an issue still in `status:backlog` awaiting a human decision.

#### Scenario: Files a bug found during work
- **WHEN** the autonomous mode finds behaviour that contradicts an existing spec in `openspec/specs/` while working an issue
- **THEN** it records the discrepancy as a new issue through the normal intake, without inventing requirements beyond what the specs already state.

### Requirement: Stage failures retry until the gate passes
When a stage fails or its gate does not pass, the autonomous mode SHALL retry that stage until it passes, instead of skipping ahead or stopping silently.

#### Scenario: Retry a failing gate
- **WHEN** a stage's gate (G2a, G2b, G3, G4, or G6) fails
- **THEN** the autonomous mode re-runs the stage to address the failure and re-checks the gate, repeating until it passes or a bounded retry limit is reached.

#### Scenario: Give up visibly after repeated failure
- **WHEN** a stage keeps failing past the retry limit
- **THEN** the autonomous mode posts a comment on the issue describing the failure and stops working that issue, rather than leaving it in a silent or ambiguous state.

### Requirement: Every action is recorded on GitHub
The autonomous mode SHALL record every action it takes (comments, commits, pull requests, label changes) on GitHub, in the same way the manual stage commands do.

#### Scenario: Overnight run leaves a visible trail
- **WHEN** a developer leaves the autonomous mode running overnight and checks the next day
- **THEN** every issue it touched shows its status labels, plan/implementation comments, and any pull request, exactly as if a human had run `/idea`, `/spec`, `/apply`, `/ship`, and `/archive` manually.
