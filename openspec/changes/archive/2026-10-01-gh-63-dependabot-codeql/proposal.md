## Why

The repository is public and used by others, but it does not run dependency or code scanning on itself. Both are free on a public repository and protect users from a known-vulnerable dependency or a code flaw.

## What Changes

- Add a Dependabot config for Go modules (`cli`) and GitHub Actions.
- Add a CodeQL workflow for Go, on pull requests and on a weekly schedule.

## Capabilities

### New Capabilities
- `supply-chain`: dependency updates and code scanning for the repository.

### Modified Capabilities
<!-- None. -->

## Impact

- `.github/dependabot.yml` and `.github/workflows/codeql.yml`.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/63
