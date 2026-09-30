## 1. Multi-tool standard and canonical definitions

- [x] 1.1 Write the canonical agent format and the per-tool render rules (`cli/internal/payload/agents/README.md`)
- [x] 1.2 Add the five stage agent definitions, one per stage (`cli/internal/payload/agents/*.md`)
- [x] 1.3 Embed and parse the agent definitions beside the commands (`cli/internal/payload/embed.go`)
- [x] 1.4 Extend each target with an agent directory and an agent format (`cli/internal/targets/targets.go`)
- [x] 1.5 Render the agents into planned files (`cli/internal/plan/plan.go`)
- [x] 1.6 Add the `tools` field to the canonical format and map it per harness (`cli/internal/payload/agents/*.md`, `cli/internal/targets/render.go`)
- [x] 1.7 Declare the per-harness support levels in the standard doc (`cli/internal/payload/agents/README.md`)

## 2. Model configuration and TUI

- [x] 2.1 Add the `models` role map and its defaults to the state (`cli/internal/config/config.go`)
- [x] 2.2 Stamp the role model into a rendered agent, or omit it when empty (`cli/internal/targets/render.go`)
- [x] 2.3 Add the interactive role-to-model TUI (`cli/internal/models/models.go`)
- [x] 2.4 Add the `octospec models` command with a non-interactive flag (`cli/cmd/octospec/main.go`)

## 3. Install, update, and doctor

- [x] 3.1 Write the agent definitions on install (`cli/internal/install/install.go`)
- [x] 3.2 Re-render the agents on update and preserve local edits (`cli/internal/update/update.go`)
- [x] 3.3 Report the agents and the model configuration in doctor (`cli/internal/doctor/doctor.go`)

## 4. Self-hosting and docs

- [x] 4.1 Migrate `AGENTS.md` to the agent standard (`AGENTS.md`)
- [x] 4.2 Point the workflow commands at the stage agents (`cli/internal/payload/commands/*.md`)
- [x] 4.3 Document the agents per stage and the model choice (`README.md`)
- [x] 4.4 Document the model TUI in the command reference and getting started (`docs/commands.md`, `docs/getting-started.md`)

## 5. Tests and gates

- [ ] 5.1 Add tests for the agent renderers (`cli/internal/targets/targets_test.go`, `cli/internal/payload/embed_test.go`)
- [ ] 5.2 Add tests for the model config, install, and update (`cli/internal/install/install_test.go`, `cli/internal/update/update_test.go`)
- [ ] 5.3 Add tests for doctor (`cli/internal/doctor/doctor_test.go`)
- [ ] 5.4 Run gofmt, vet, tests, `openspec validate --all --strict`, and the gates (`scripts/check-gates.sh`)
- [ ] 5.5 Add tests for the tool surface render and the declared support levels (`cli/internal/targets/targets_test.go`)
