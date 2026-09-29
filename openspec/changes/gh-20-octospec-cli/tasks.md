## 1. Go CLI skeleton

- [x] 1.1 Create the Go module and entrypoint (`cli/go.mod`, `cli/cmd/octospec/main.go`)
- [x] 1.2 Define the target mapping table for pi, opencode, Copilot, Claude (`cli/internal/targets/targets.go`)
- [x] 1.3 Embed the canonical command payload (`cli/internal/payload/embed.go`, `cli/internal/payload/commands/*`)
- [x] 1.4 Implement `install` (wizard plus `--tool`) writing rendered commands into each target (`cli/internal/install/install.go`)

## 2. Targets

- [x] 2.1 Default front matter (`description`, `argument-hint`) for pi, opencode, and Claude (`cli/internal/targets/render.go`)
- [x] 2.2 Copilot front matter (`mode: agent`) (`cli/internal/targets/render.go`)
- [x] 2.3 Global target dirs: pi -> `~/.pi/agent/prompts/`, opencode -> `~/.config/opencode/command/` (`cli/internal/targets/targets.go`)
- [x] 2.4 Repo-local target dirs: Copilot -> `.github/prompts/`, Claude -> `.claude/commands/` (`cli/internal/targets/targets.go`)

## 3. Update

- [x] 3.1 Persist saved wizard answers (`cli/internal/config/config.go`)
- [x] 3.2 Implement `update` with a managed-file hash manifest that preserves local edits (`cli/internal/update/update.go`)

## 4. Replace install.sh

- [x] 4.1 Port schema install into the CLI (`cli/internal/seed/schema.go`)
- [x] 4.2 Port label provisioning into the CLI (`cli/internal/seed/labels.go`)
- [x] 4.3 Port the repo seed: issue forms, CI, `openspec/config.yaml` (`cli/internal/seed/repo.go`)
- [x] 4.4 Delete `install.sh` and the command copies (`install.sh`, `.github/prompts/`); `schema/` and `repo-template/` moved under `cli/internal/seed/`

## 5. Tests

- [x] 5.1 Table tests for every renderer (`cli/internal/targets/*_test.go`)
- [x] 5.2 Install and update tests into temp dirs, including edit preservation (`cli/internal/install/*_test.go`, `cli/internal/update/*_test.go`)

## 6. CI

- [x] 6.1 Add a CLI workflow: build, `gofmt`, `go vet`, `go test ./...` (`.github/workflows/cli.yml`)
- [x] 6.2 Run the CLI workflow and the octospec gates on every pull request (`.github/workflows/openspec.yml`)

## 7. CD

- [x] 7.1 Add a release workflow that produces tagged binaries (`.github/workflows/release.yml`, `.goreleaser.yaml`)
- [x] 7.2 Add the `curl | sh` installer that fetches a release (`install.sh`)

## 8. Ship feedback loop

- [x] 8.1 Make `/ship` run/observe CI and loop failures back to `/apply` (`cli/internal/payload/commands/ship.md`)
- [x] 8.2 No mirror: the payload is the single source, rendered per tool by the CLI

## 9. Docs and verify

- [ ] 9.1 Document the CLI in the README (`README.md`)
- [ ] 9.2 Dry-run: install per target, confirm `update` preserves edits, and a red CI loops back (`cli/`, `commands/ship.md`)
