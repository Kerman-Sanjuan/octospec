# supply-chain Specification

## Purpose
Dependency updates and code scanning for the repository.
## Requirements
### Requirement: Dependency updates
The repository SHALL have a Dependabot config for Go modules and GitHub Actions.

#### Scenario: Dependabot opens updates
- **WHEN** a dependency has a newer version
- **THEN** Dependabot opens a pull request against the repository.

### Requirement: Code scanning
The repository SHALL run CodeQL for Go on pull requests and on a schedule.

#### Scenario: CodeQL runs on a pull request
- **WHEN** a pull request is opened
- **THEN** CodeQL analyzes the Go code and reports.

