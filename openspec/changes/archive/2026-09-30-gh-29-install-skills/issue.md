# Issue #29: CLI must install the OpenSpec skills

https://github.com/Kerman-Sanjuan/octospec/issues/29

## User story

As a maintainer onboarding a repository, I want the CLI to install the OpenSpec skills for each tool, so that the commands that load them actually work.

## Context / problem

Every command says "load the `openspec-propose` skill", but the seed writes only commands, issue forms, config, and labels. On a fresh repo, `/spec` fails at the first skill load. This is the biggest blocker to a working MVP.

## Requirements

- The CLI SHALL make the OpenSpec skills (`openspec-propose`, `openspec-apply-change`, `openspec-archive-change`, `openspec-explore`) available for each supported tool after `install` or `seed`.
- It MAY do this by rendering the skills per tool, or by running `openspec update`, or by removing the commands' dependency on them.

## Success criteria

- After `octospec install` and `octospec seed` on a fresh repo, `/spec` runs without a missing-skill error.

## Out of scope

- Changing what the skills do.
