## Why

The workflow depends on the human to write each artifact well, so quality is uneven. A single monolithic agent loads the whole manual, and there is no way to pick a model per stage. octospec can split the work into scoped agents, render them for every tool from one source, and let the user choose the model each role uses.

## What Changes

- Define a canonical agent per stage of the loop (`idea`, `spec`, `apply`, `ship`, `archive`), each with a scoped mission, the skills it may load, the context it receives, and the artifact it writes.
- Define the multi-tool agent standard: the canonical fields and how each tool renders them.
- Add a single model configuration that maps each role (thinking, implementer, reviewer) to a model, and stamp that model into every rendered agent.
- Add an interactive TUI to set the model per role without editing files by hand.
- Render the agent definitions from the same canonical payload that holds the commands, and extend `install`, `update`, `plan`, and the target mappings to cover them.
- Keep the CLI as the orchestrator: it validates the input, selects the stage agent, and the agent writes the artifact.
- Migrate the instructions to develop octospec (`AGENTS.md` and the workflow commands of this repo) to the new rules.
- Document the agents per stage and how to pick a model in the README and the docs.

## Capabilities

### New Capabilities
- `agent-definitions`: the canonical per-stage agent definitions, their scoped context and skills, and the render for each supported tool.
- `agent-model-config`: the single model-role configuration and the interactive TUI that edits it.

### Modified Capabilities
- `cli-distribution`: the canonical payload, `plan`, `install`, and `update` also cover the agent definitions, and every target maps an agent location.
- `diagnostics`: `octospec doctor` also reports the agent definitions and the model configuration.
- `readme`: the README documents the agent per stage and how to choose the model each role uses.
- `consumer-docs`: the command reference and the getting started guide document the agents and the model TUI.

## Impact

- `cli/internal/payload/` (new agent definitions), `cli/internal/targets/`, `cli/internal/plan/`, `cli/internal/install/`, `cli/internal/update/`, `cli/internal/config/`.
- New TUI package and a new CLI subcommand in `cli/cmd/octospec/main.go`.
- `AGENTS.md`, `README.md`, and `docs/`.
- Originating issue: https://github.com/Kerman-Sanjuan/octospec/issues/46
