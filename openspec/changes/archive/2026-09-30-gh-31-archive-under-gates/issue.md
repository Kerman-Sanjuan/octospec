# Issue #31: Make `/archive` land under the enforced gates

https://github.com/Kerman-Sanjuan/octospec/issues/31

## User story

As a maintainer, I want `/archive` to finish without the protected branch rejecting the push.

## Context / problem

`/archive` commits the archive and specs and pushes to `main`. With `cli` and `gates` required and admins enforced, the push is rejected, so the last step of the loop is manual. It also does not commit or push at all in some paths.

## Requirements

- `/archive` SHALL land the archive (schema sync, spec sync, the archive move) without a rejected push.
- It SHALL commit and push the result.

## Success criteria

- `/archive` completes on the protected `main` with no manual step.

## Out of scope

- Changing the gate definitions.
