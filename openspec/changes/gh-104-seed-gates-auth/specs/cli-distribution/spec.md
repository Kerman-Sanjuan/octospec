## MODIFIED Requirements

### Requirement: Replaces install.sh
The CLI SHALL cover the old `install.sh` and its `--repo` seed: schema install, label provisioning, and the `repo-template` seed (issue forms, CI, `openspec/config.yaml`). The seeded CI workflow SHALL carry the same `gh` authentication as the repository's own gates workflow (`GH_TOKEN`, `issues: read`, and the same `actions/setup-node` major), and a test SHALL fail when the two drift.

#### Scenario: Seed a repo
- **WHEN** the CLI seeds a repository
- **THEN** the schema is installed, the labels provisioned, and the seed files written, as `install.sh --repo` did.

#### Scenario: The seeded gates workflow authenticates gh
- **WHEN** a freshly seeded repository's gates workflow runs
- **THEN** it sets `GH_TOKEN` and requests `issues: read`, so the gates can read the linked issue.

#### Scenario: The seeded workflow does not drift
- **WHEN** the seeded CI workflow differs from the root on the token, the permissions, or the `setup-node` major
- **THEN** the seed test fails.
