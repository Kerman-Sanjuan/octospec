## ADDED Requirements

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
