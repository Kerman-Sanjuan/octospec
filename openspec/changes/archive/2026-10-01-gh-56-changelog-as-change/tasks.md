## 1. Carry the entry in the change

- [x] 1.1 Add a `## Changelog` section to the proposal template (`cli/internal/seed/schema/templates/proposal.md`)
- [x] 1.2 Document the section in the proposal artifact instruction (`cli/internal/seed/schema/schema.yaml`)

## 2. Fold the entry at archive

- [x] 2.1 Make `/archive` copy the change's `## Changelog` entry into `CHANGELOG.md` (`cli/internal/payload/commands/archive.md`)
- [x] 2.2 Stop the archive when the section is missing, and skip it when it reads `None`

## 3. Document and check

- [x] 3.1 Update `docs/releasing.md` to describe the generated changelog
- [x] 3.2 Add a test that the proposal template carries the section (`cli/internal/seed/`)
- [x] 3.3 Run gofmt, vet, tests, `openspec validate --all --strict`, and the gates (`scripts/check-gates.sh`)
