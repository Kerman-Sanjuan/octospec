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
The CLI SHALL install the commands and the agent definitions for a target tool from one canonical payload embedded in the binary, and the repo SHALL NOT store per-tool derived copies.

#### Scenario: One source, many destinations
- **WHEN** the canonical payload changes
- **THEN** every target is regenerated from it with no per-tool copy to edit.

### Requirement: Target mapping
Each supported tool SHALL be defined by a target mapping with a repo-local location and a global location (directory plus file format for commands and for agents); adding a tool SHALL NOT add a source copy.

#### Scenario: Adding a tool
- **WHEN** a new tool is added
- **THEN** only a new mapping is added and the payload is unchanged.

### Requirement: Supported targets
The CLI SHALL support pi, opencode, GitHub Copilot, and Claude Code, and SHALL offer a repo-local location and a global location for each of them.

#### Scenario: Install per tool
- **WHEN** `octospec install --tool <pi|opencode|copilot|claude>` runs
- **THEN** the commands are written to that tool's chosen directory.

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

### Requirement: Install the OpenSpec skills
`octospec install` SHALL make the OpenSpec skills available to each selected tool, so that commands that load them run.

#### Scenario: Skills present after install
- **WHEN** `octospec install --tool <pi|opencode|copilot|claude>` runs in a repository
- **THEN** the OpenSpec skills exist in that tool's skill directory.

#### Scenario: Missing OpenSpec is reported
- **WHEN** the `openspec` CLI is not on the PATH
- **THEN** `install` reports that the skills were skipped, and the commands still install.

### Requirement: Install and update render the agents
The CLI SHALL write the agent definitions for the selected tools on install and update, and SHALL preserve local edits to them like any managed file.

#### Scenario: Install writes the agents
- **WHEN** `octospec install --tool <pi|opencode|copilot|claude>` runs
- **THEN** the agent definitions are written to that tool's agent location.

#### Scenario: Hand edit to an agent survives update
- **WHEN** a managed agent file is edited locally and `octospec update` runs
- **THEN** the edit is preserved and reported.

### Requirement: The installer verifies the checksum
The `curl | sh` installer SHALL verify the release `checksums.txt` against the downloaded binary before it installs it, and SHALL fail without installing on a mismatch.

#### Scenario: A good download
- **WHEN** the binary matches its checksum
- **THEN** the installer installs it.

#### Scenario: A corrupted download
- **WHEN** the binary does not match its checksum
- **THEN** the installer stops and installs nothing.

### Requirement: The installer accepts a version
The `curl | sh` installer SHALL accept a version to install, and SHALL default to the latest release.

#### Scenario: Pin a version
- **WHEN** the installer runs with a version such as `--version v1.1.0`
- **THEN** it installs that release instead of the latest.

#### Scenario: Default to latest
- **WHEN** the installer runs without a version
- **THEN** it installs the latest release.

### Requirement: A repo-local install writes nothing under home
A repo-local install SHALL NOT write under the user's home directory.

#### Scenario: Nothing under home
- **WHEN** every tool is installed repo-local
- **THEN** no file is written under the home directory.

### Requirement: pi has no agent files
octospec SHALL NOT install agent files for pi, because pi has no agent mechanism.

#### Scenario: pi install
- **WHEN** pi is installed
- **THEN** its commands and skills are written and no agent file is written.

### Requirement: Dry-run for install and seed
`octospec install` and `octospec seed` SHALL support `--dry-run`, which writes nothing and prints what they would do.

#### Scenario: Install dry-run writes nothing
- **WHEN** `octospec install --dry-run` runs
- **THEN** it prints the files it would write and writes nothing.

#### Scenario: Seed dry-run writes nothing
- **WHEN** `octospec seed --dry-run` runs
- **THEN** it prints what it would install and writes nothing.

### Requirement: Uninstall removes the managed files
`octospec uninstall` SHALL remove the files octospec manages, using the recorded state, and SHALL leave hand-edited files alone.

#### Scenario: Remove managed files
- **WHEN** `octospec uninstall` runs
- **THEN** the unmodified managed files are removed and reported.

#### Scenario: Preserve a hand edit
- **WHEN** a managed file was edited by hand
- **THEN** `uninstall` keeps it and reports it.

