## MODIFIED Requirements

### Requirement: G1 issue sections
A gate SHALL fail when the OpenSpec change in the pull request links an issue that is missing a required section (user story, context, requirements, success criteria). The gate SHALL skip when the pull request touches no OpenSpec change, when `gh` is unavailable, or when `gh` cannot read the issue body (an authentication, permission, or network failure). The gate SHALL NOT fail on a read error.

#### Scenario: Issue is complete
- **WHEN** a pull request touches `openspec/changes/<name>/` and the linked issue has every required section
- **THEN** G1 passes.

#### Scenario: Issue is incomplete
- **WHEN** the linked issue is read and is missing a required section
- **THEN** G1 fails and names the missing sections.

#### Scenario: gh cannot read the issue
- **WHEN** `gh issue view` fails (unauthenticated, missing permission, or network error)
- **THEN** G1 skips with a clear message and does not fail the run, because an unread issue is not an incomplete one.

#### Scenario: No change in the pull request
- **WHEN** the pull request touches no OpenSpec change
- **THEN** G1 skips.
