## 1. Fix the intake commands

- [x] 1.1 Make `/idea` create the issue from the body and apply the labels (`cli/internal/payload/commands/idea.md`)
- [x] 1.2 Make `/bug` create the issue from the body and apply the labels (`cli/internal/payload/commands/bug.md`)

## 2. Guard and document

- [x] 2.1 Add a test that no command uses `--template` with `--body` (`cli/internal/payload/embed_test.go`)
- [x] 2.2 Keep the command reference accurate (`docs/commands.md`)

## 3. Checks

- [x] 3.1 Run gofmt, vet, tests, `openspec validate --all --strict`, and the gates (`scripts/check-gates.sh`)
