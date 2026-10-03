# gates Specification

## Purpose
The gate set that CI enforces: what each gate checks and how strong it is.
## Requirements
### Requirement: G1 issue sections
A gate SHALL fail when the OpenSpec change in the pull request links an issue that is missing a required section (user story, context, requirements, success criteria). The gate SHALL skip when the pull request touches no OpenSpec change, or when `gh` is unavailable.

#### Scenario: Issue is complete
- **WHEN** a pull request touches `openspec/changes/<name>/` and the linked issue has every required section
- **THEN** G1 passes.

#### Scenario: Issue is incomplete
- **WHEN** the linked issue is missing a required section
- **THEN** G1 fails and names the missing sections.

#### Scenario: No change in the pull request
- **WHEN** the pull request touches no OpenSpec change
- **THEN** G1 skips.

### Requirement: G2 validate and completeness
A gate SHALL fail when `openspec validate --all --strict` fails, or when an unarchived change is incomplete. When `openspec` is unavailable, both the validate check (G2a) and the completeness check (G2b) SHALL report a SKIP or a FAIL, consistently, and SHALL NOT pass silently.

#### Scenario: A change is incomplete
- **WHEN** an unarchived change lacks a required artifact
- **THEN** G2 fails and names the change.

#### Scenario: OpenSpec is unavailable
- **WHEN** `openspec` is not on PATH
- **THEN** G2a and G2b report the same way, so neither passes silently.

### Requirement: G3 branch name
A gate SHALL check the branch name only for `feat/`, `fix/`, and `dependabot/` refs. For `feat/` and `fix/` it SHALL fail unless the ref matches `feat|fix/<issue>-<slug>`. For `dependabot/` it SHALL skip when the author is `dependabot[bot]` and fail otherwise. For any other ref (main, a tag such as `v1.0.0`, or another long-lived ref) it SHALL skip.

#### Scenario: A conforming branch
- **WHEN** the branch is `feat/12-add-thing`
- **THEN** G3 passes.

#### Scenario: A malformed branch
- **WHEN** the branch is `feat/12` and does not match `<issue>-<slug>`
- **THEN** G3 fails.

#### Scenario: A long-lived branch or tag ref
- **WHEN** `HEAD_REF` is `main` or a tag such as `v1.0.0`
- **THEN** G3 skips instead of failing the run.

#### Scenario: A dependency branch
- **WHEN** the branch is `dependabot/go_modules/cli/x` and the author is `dependabot[bot]`
- **THEN** G3 skips.

#### Scenario: A dependency branch from a human
- **WHEN** the branch is `dependabot/go_modules/cli/x` and the author is not `dependabot[bot]`
- **THEN** G3 fails.

### Requirement: G4 pull request links the issue
A gate SHALL fail when the pull request body does not reference the issue with `Closes`, `Fixes`, or `Resolves` and a number. The gate SHALL skip when the body is not provided. The gate SHALL skip when the pull-request author is `dependabot[bot]`.

#### Scenario: A linked pull request
- **WHEN** the body contains `Closes #12`
- **THEN** G4 passes.

#### Scenario: An unlinked pull request
- **WHEN** the body has no `Closes`, `Fixes`, or `Resolves` reference and the author is not `dependabot[bot]`
- **THEN** G4 fails.

#### Scenario: A dependency pull request
- **WHEN** the body has no issue reference and the author is `dependabot[bot]`
- **THEN** G4 skips.

### Requirement: G5 the target stack passes
A gate SHALL fail when the target repository's tests, build, or lint fail. The target repository provides this gate; the shared gate script does not.

#### Scenario: A red build
- **WHEN** the target stack's check fails
- **THEN** G5 fails.

### Requirement: G6 spec delta
A gate SHALL fail when a pull request touches `openspec/changes/<name>/` without a `specs/*.md` delta. The gate SHALL skip when no change is touched. The failure message SHALL name the real control, `skip_specs: true` in the change's `.openspec.yaml`.

#### Scenario: A change without a delta
- **WHEN** the pull request edits a change but adds no `specs/*.md`
- **THEN** G6 fails and the message names `skip_specs: true`.

### Requirement: G7 tasks checked
A gate SHALL warn, without failing, when an unarchived change has unchecked tasks. When `openspec` is unavailable, G7 SHALL report a SKIP or a FAIL, consistently with G2, and SHALL NOT pass silently.

#### Scenario: An unarchived change with open tasks
- **WHEN** an unarchived change has an unchecked task
- **THEN** G7 warns and the pull request can still merge.

#### Scenario: OpenSpec is unavailable
- **WHEN** `openspec` is not on PATH
- **THEN** G7 reports the same way as G2, so it does not pass silently.

### Requirement: G8 one change per issue
A gate SHALL warn, without failing, when more than one unarchived change links the same issue.

#### Scenario: Two changes for one issue
- **WHEN** two unarchived changes carry the same `github.issue`
- **THEN** G8 warns.

### Requirement: Gates are runnable locally and in CI
The shared gate script SHALL run the gates that do not depend on the target stack (G1 to G4, and G6 to G8), and the repository SHALL run it in CI. The repository's gate test SHALL cover the `dependabot[bot]` exemption for G3 and G4.

#### Scenario: Run the gates locally
- **WHEN** a maintainer runs `sh scripts/check-gates.sh`
- **THEN** the stack-independent gates report PASS, FAIL, or SKIP.

#### Scenario: The bot exemption is tested
- **WHEN** a maintainer runs `sh scripts/check-gates-test.sh`
- **THEN** it passes and includes a case that runs the gate script as `dependabot[bot]` and asserts G3 and G4 skip.

### Requirement: Archived changes carry no unchecked tasks
Every archived change SHALL have all of its tasks checked, so an archived change never contradicts gate G7.

#### Scenario: Archived changes have no open tasks
- **WHEN** an archived change is inspected
- **THEN** every checkbox in its `tasks.md` is checked.

