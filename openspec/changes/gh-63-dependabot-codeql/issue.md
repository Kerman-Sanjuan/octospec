# Issue #63: Enable Dependabot and CodeQL on octospec

https://github.com/Kerman-Sanjuan/octospec/issues/63

## User story

As a user, I want dependency and code scanning to be on, so a known-vulnerable dependency or a code flaw is caught before it reaches me.

## Context / problem

octospec seeds scanning for other repositories (#39), but does not enable it on its own repository. Both are free on a public repository.

## Requirements

- The repository SHALL enable Dependabot for Go modules and GitHub Actions.
- The repository SHALL run CodeQL for Go on pull requests and on a schedule.

## Success criteria

- Dependabot opens update pull requests.
- CodeQL runs and reports on pull requests.

## Out of scope

- The seed default for other repositories (#39).
