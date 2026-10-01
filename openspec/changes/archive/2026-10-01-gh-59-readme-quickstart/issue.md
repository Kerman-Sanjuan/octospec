# Issue #59: README: quickstart, real prerequisites, and supported platforms

https://github.com/Kerman-Sanjuan/octospec/issues/59

## User story

As a newcomer, I want a copy-paste quickstart and honest prerequisites, so I can try octospec in a minute.

## Context / problem

The README opens with a strong "why", but the install path is far down the page, so a visitor has to hunt for it. The requirements list the OpenSpec CLI and `gh`, but not Node, which the OpenSpec CLI needs. Windows is unsupported because the installer is `sh`, and the README does not say so.

## Requirements

- The README SHALL show a copy-paste quickstart before the deep explanation.
- It SHALL list every prerequisite: the OpenSpec CLI, `gh`, and Node.
- It SHALL state the supported platforms (linux and macOS) and that Windows is not supported yet.

## Success criteria

- A newcomer can install and run one `/idea` by following the top of the README.
- The prerequisites are complete and correct.
- The supported platforms are explicit.

## Out of scope

- A demo recording (tracked separately).
