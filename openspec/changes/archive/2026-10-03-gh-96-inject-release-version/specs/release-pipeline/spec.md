## MODIFIED Requirements

### Requirement: Release binaries
CD SHALL publish versioned release binaries on tag, installable by the `curl | sh` installer, and the binary SHALL report the tag it was built from, injected at build time rather than written by hand.

#### Scenario: Tag a release
- **WHEN** a version tag is pushed
- **THEN** release binaries are published and the `curl | sh` installer can install them.

#### Scenario: The binary reports the tag
- **WHEN** the release build injects the tag through `ldflags`
- **THEN** `octospec version` prints that version without the leading `v`.

#### Scenario: A local build
- **WHEN** a developer runs `go build` without the injected value
- **THEN** the binary still runs and `octospec version` prints the `dev` fallback.
