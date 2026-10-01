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
A gate SHALL fail when `openspec validate --all --strict` fails, or when an unarchived change is incomplete.

#### Scenario: A change is incomplete
- **WHEN** an unarchived change lacks a required artifact
- **THEN** G2 fails and names the change.

### Requirement: G3 branch name
A gate SHALL fail when the pull-request branch does not match `feat|fix/<issue>-<slug>`. The gate SHALL skip on a long-lived branch.

#### Scenario: A conforming branch
- **WHEN** the branch is `feat/12-add-thing`
- **THEN** G3 passes.

#### Scenario: A long-lived branch
- **WHEN** the branch is `main`
- **THEN** G3 skips.

### Requirement: G4 pull request links the issue
A gate SHALL fail when the pull request body does not reference the issue with `Closes`, `Fixes`, or `Resolves` and a number. The gate SHALL skip when the body is not provided.

#### Scenario: A linked pull request
- **WHEN** the body contains `Closes #12`
- **THEN** G4 passes.

### Requirement: G5 the target stack passes
A gate SHALL fail when the target repository's tests, build, or lint fail. The target repository provides this gate; the shared gate script does not.

#### Scenario: A red build
- **WHEN** the target stack's check fails
- **THEN** G5 fails.

### Requirement: G6 spec delta
A gate SHALL fail when a pull request touches `openspec/changes/<name>/` without a `specs/*.md` delta. The gate SHALL skip when no change is touched.

#### Scenario: A change without a delta
- **WHEN** the pull request edits a change but adds no `specs/*.md`
- **THEN** G6 fails.

### Requirement: G7 tasks checked
A gate SHALL warn, without failing, when an unarchived change has unchecked tasks.

#### Scenario: An unarchived change with open tasks
- **WHEN** an unarchived change has an unchecked task
- **THEN** G7 warns and the pull request can still merge.

### Requirement: G8 one change per issue
A gate SHALL warn, without failing, when more than one unarchived change links the same issue.

#### Scenario: Two changes for one issue
- **WHEN** two unarchived changes carry the same `github.issue`
- **THEN** G8 warns.

### Requirement: Gates are runnable locally and in CI
The shared gate script SHALL run the gates that do not depend on the target stack (G1 to G4, and G6 to G8), and the repository SHALL run it in CI.

#### Scenario: Run the gates locally
- **WHEN** a maintainer runs `sh scripts/check-gates.sh`
- **THEN** the stack-independent gates report PASS, FAIL, or SKIP.

