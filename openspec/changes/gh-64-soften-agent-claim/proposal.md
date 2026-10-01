## Why

The commands say "Runs the `<stage>` agent", but nothing dispatches to it; the command runs in the current session. The claim is misleading, and real cross-harness dispatch is not verifiable here because pi has no subagents. This change makes the claim honest and leaves real dispatch to a follow-up.

## What Changes

- Reword the stage line in the workflow commands so it states what happens.
- State the current behaviour in the agent standard.
- Update the `agent-definitions` spec so the stage contract is the agent's scope, tools, and model, not a dispatch.

## Capabilities

### New Capabilities
<!-- None. -->

### Modified Capabilities
- `agent-definitions`: the agent is the stage's contract for scope, tools, and model, and the command runs the stage in the current session.

## Impact

- `cli/internal/payload/commands/*.md` and `cli/internal/payload/agents/README.md`.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/64
