## 1. Dry-run

- [ ] 1.1 Add `--dry-run` to install (`cli/internal/install/install.go`, `cli/cmd/octospec/main.go`)
- [ ] 1.2 Add `--dry-run` to seed (`cli/internal/seed/seed.go`, `cli/cmd/octospec/main.go`)

## 2. Uninstall

- [ ] 2.1 Add `octospec uninstall`, which removes the unmodified managed files and preserves edits (`cli/internal/uninstall/uninstall.go`, `cli/cmd/octospec/main.go`)

## 3. Docs

- [ ] 3.1 Document the flags and the command (`docs/commands.md`, `README.md`)

## 4. Tests and checks

- [ ] 4.1 Add tests for the dry run and uninstall (`cli/internal/uninstall/uninstall_test.go`)
- [ ] 4.2 Run gofmt, vet, tests, `openspec validate --all --strict`, and the gates (`scripts/check-gates.sh`)
