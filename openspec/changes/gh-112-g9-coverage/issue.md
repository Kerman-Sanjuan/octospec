# Issue #112 — Link every requirement to a test: gate G9

https://github.com/Kerman-Sanjuan/octospec/issues/112

Read-only mirror of the issue narrative. The issue body stays the human source of truth.

## User story

As a person reading a change six months later, I want every requirement to name the check that
proves it, so that I can tell a requirement nobody tests from one that is quietly broken.

## Context

The harness grades a repository 0/4 on acceptance criteria: specs are written, and nothing says which
test proves which sentence. A change can add five requirements, one test, and pass every gate we have.

## Requirements

- A change SHALL declare its mapping in `openspec/changes/<change>/coverage.md`, one line per
  requirement: `- <requirement title> -> <symbol>`.
- The symbol SHALL be a string that exists in the tree outside `openspec/`: a test name, a script, a
  path.
- G9 SHALL warn when a requirement has no line, and when a symbol is nowhere.
- A change with no spec delta SHALL skip G9.

## Success criteria

- A change with a requirement and no coverage line prints a warning naming the requirement.
- A coverage line pointing at a symbol that does not exist prints a warning naming the symbol.
- A fully covered change passes.
