# Issue #64: Wire the stage agents into the commands, or soften the claim

https://github.com/Kerman-Sanjuan/octospec/issues/64

## User story

As a user, I want each stage to actually run on its stage agent, so the agentic model is real and not just definitions on disk.

## Context / problem

The stage agents (idea, spec, apply, ship, archive) are rendered for every tool, with their tool surface and model, but nothing dispatches to them. A command such as `/spec` runs in the current session and its "Runs the `spec` agent" line is a description, not a spawn. On harnesses with subagents (Claude Code, opencode, Copilot) the dispatch could work; pi has no subagent mechanism, so it would run inline.

## Requirements

- The commands SHALL dispatch to the stage agent where the harness supports subagents.
- The commands SHALL degrade to inline execution where the harness does not, for example pi.
- The rendered agent files SHALL remain the single source of the stage's scope and tools.

## Success criteria

- On a harness with subagents, `/spec` runs in the spec agent with only its tools.
- On pi, the command still works and runs inline.
- `octospec doctor` and the gates stay green.

## Out of scope

- Subagents in pi, which pi does not support.
- Automating the full loop end to end.
