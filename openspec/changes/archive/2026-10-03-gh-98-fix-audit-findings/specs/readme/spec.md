## MODIFIED Requirements

### Requirement: The prerequisites are complete
The README SHALL list every prerequisite: the OpenSpec CLI, `gh`, and Node, and SHALL state the OpenSpec version the gates are tested against without claiming support for untested newer versions.

#### Scenario: The prerequisites are listed
- **WHEN** a reader looks for the prerequisites
- **THEN** the OpenSpec CLI, `gh`, and Node are named.

#### Scenario: The OpenSpec version is accurate
- **WHEN** a reader reads the OpenSpec prerequisite
- **THEN** it names the tested version (1.3.1) and does not say "or newer", because the strict gates fail on OpenSpec 1.14.0.

### Requirement: README links to deeper docs
The README SHALL link to the design doc and implementation plan as deeper references, and its prose SHALL be complete and readable, with no truncated sentences or stray backticks.

#### Scenario: Links resolve
- **WHEN** a reader follows the links in the README
- **THEN** each link resolves to an existing file.

#### Scenario: The skills paragraph reads
- **WHEN** a reader reads the OpenSpec skills paragraph
- **THEN** the sentences are complete, with closed backticks and no line that begins with a stray character.
