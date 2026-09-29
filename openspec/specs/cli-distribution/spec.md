# cli-distribution Specification

## Purpose
The Go CLI: one canonical command payload rendered into each tool's directory (pi, opencode, GitHub Copilot, Claude Code), how it is distributed (`curl | sh` and `go install`), and how `update` preserves local edits.
## Requirements
### Requirement: Single Go binary
`octospec/cli` SHALL be a single Go binary with no runtime dependency.

#### Scenario: No runtime dependency
- **WHEN** the binary is installed
- **THEN** it runs without Node or any other runtime.

### Requirement: Install from one canonical payload
The CLI SHALL install the commands for a target tool from one canonical payload embedded in the binary, and the repo SHALL NOT store per-tool derived copies.

#### Scenario: One source, many destinations
- **WHEN** the canonical payload changes
- **THEN** every target is regenerated from it with no per-tool copy to edit.

### Requirement: Target mapping
Each supported tool SHALL be defined by a target mapping (destination directory plus file format); adding a tool SHALL NOT add a source copy.

#### Scenario: Adding a tool
- **WHEN** a new tool is added
- **THEN** only a new mapping is added and the payload is unchanged.

### Requirement: Supported targets
The CLI SHALL support pi, opencode, GitHub Copilot, and Claude Code, with global targets for pi and opencode and repo-local targets for Copilot and Claude.

#### Scenario: Install per tool
- **WHEN** `octospec install --tool <pi|opencode|copilot|claude>` runs
- **THEN** the commands are written to that tool's directory.

### Requirement: Update preserves local edits
The CLI SHALL provide `update`, which re-applies from saved answers and preserves local edits tracked by managed-file hashes.

#### Scenario: Hand edit survives update
- **WHEN** a managed file is edited locally and `octospec update` runs
- **THEN** the edit is preserved and reported.

### Requirement: Replaces install.sh
The CLI SHALL cover the old `install.sh` and its `--repo` seed: schema install, label provisioning, and the `repo-template` seed (issue forms, CI, `openspec/config.yaml`).

#### Scenario: Seed a repo
- **WHEN** the CLI seeds a repository
- **THEN** the schema is installed, the labels provisioned, and the seed files written, as `install.sh --repo` did.

### Requirement: Distribution
The CLI SHALL be installed via a `curl | sh` installer from a release (and `go install` for contributors) and SHALL support non-interactive use.

#### Scenario: One-line install
- **WHEN** a user runs the `curl | sh` installer
- **THEN** the released binary is installed and runnable.

