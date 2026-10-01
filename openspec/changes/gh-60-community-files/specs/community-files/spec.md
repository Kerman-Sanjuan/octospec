## ADDED Requirements

### Requirement: A security policy
The repository SHALL have a `SECURITY.md` that says how to report a vulnerability privately.

#### Scenario: A reporter finds the policy
- **WHEN** a reporter looks for how to report a vulnerability
- **THEN** `SECURITY.md` names the private channel.

### Requirement: A code of conduct
The repository SHALL have a `CODE_OF_CONDUCT.md` that adopts the Contributor Covenant.

#### Scenario: A contributor finds the code of conduct
- **WHEN** a contributor opens the repository
- **THEN** `CODE_OF_CONDUCT.md` states the expected behaviour and the enforcement contact.

### Requirement: Issue and pull request templates
The repository SHALL have a pull request template and an `ISSUE_TEMPLATE/config.yml` that links the security policy.

#### Scenario: A new pull request
- **WHEN** a contributor opens a pull request
- **THEN** the body starts from the template.

#### Scenario: The issue chooser links the policy
- **WHEN** a visitor opens the new-issue chooser
- **THEN** it links the security policy.
