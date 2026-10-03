## MODIFIED Requirements

### Requirement: G3 branch name
A gate SHALL fail when the pull-request branch does not match `feat|fix/<issue>-<slug>`. The gate SHALL skip on a long-lived branch. The gate SHALL skip when the pull-request author is `dependabot[bot]`.

#### Scenario: A conforming branch
- **WHEN** the branch is `feat/12-add-thing`
- **THEN** G3 passes.

#### Scenario: A long-lived branch
- **WHEN** the branch is `main`
- **THEN** G3 skips.

#### Scenario: A non-conforming branch
- **WHEN** the branch is `dependabot/go_modules/cli/x` and the author is not `dependabot[bot]`
- **THEN** G3 fails.

#### Scenario: A dependency branch
- **WHEN** the branch is `dependabot/go_modules/cli/x` and the author is `dependabot[bot]`
- **THEN** G3 skips.

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

### Requirement: Gates are runnable locally and in CI
The shared gate script SHALL run the gates that do not depend on the target stack (G1 to G4, and G6 to G8), and the repository SHALL run it in CI. The repository's gate test SHALL cover the `dependabot[bot]` exemption for G3 and G4.

#### Scenario: Run the gates locally
- **WHEN** a maintainer runs `sh scripts/check-gates.sh`
- **THEN** the stack-independent gates report PASS, FAIL, or SKIP.

#### Scenario: The bot exemption is tested
- **WHEN** a maintainer runs `sh scripts/check-gates-test.sh`
- **THEN** it passes and includes a case that runs the gate script as `dependabot[bot]` and asserts G3 and G4 skip.
