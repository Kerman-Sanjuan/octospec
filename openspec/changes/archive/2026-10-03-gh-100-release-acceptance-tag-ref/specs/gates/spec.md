## MODIFIED Requirements

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
