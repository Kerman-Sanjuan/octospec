## 1. Targets and state

- [x] 1.1 Give every tool a repo-local and a global location for commands and agents (`cli/internal/targets/targets.go`)
- [x] 1.2 Give pi no agent location (`cli/internal/targets/targets.go`)
- [x] 1.3 Record each tool's scope in the install state (`cli/internal/config/config.go`)
- [x] 1.4 Plan the files for the chosen scope and skip the pi agents (`cli/internal/plan/plan.go`)

## 2. The warning and the choice

- [x] 2.1 Warn on a global install and let the user choose the scope (`cli/internal/install/install.go`)
- [x] 2.2 Add the `--global` flag for non-interactive installs (`cli/cmd/octospec/main.go`)

## 3. Update and cleanup

- [x] 3.1 Re-apply the recorded scope on update (`cli/internal/update/update.go`)
- [x] 3.2 Offer to remove the files earlier versions installed globally (`cli/internal/install/install.go`)

## 4. Docs

- [x] 4.1 Document the scopes and the global warning (`README.md`, `docs/commands.md`)

## 5. Tests and gates

- [x] 5.1 Add tests for the scope planning, the warning, and the pi agent skip (`cli/internal/targets/targets_test.go`, `cli/internal/plan/plan_test.go`)
- [x] 5.2 Run gofmt, vet, tests, `openspec validate --all --strict`, and the gates (`scripts/check-gates.sh`)
