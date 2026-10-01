## 1. Implement G1

- [x] 1.1 Add the G1 check to the gate script (`scripts/check-gates.sh`)
- [x] 1.2 Cover G1 and G3 in the gate tests (`scripts/check-gates-test.sh`)

## 2. Run the gates at release

- [x] 2.1 Add a `gates` job and keep the OS matrix for acceptance (`.github/workflows/release-acceptance.yml`)

## 3. Add the end-to-end test

- [x] 3.1 Write `scripts/e2e-test.sh`: build, seed, install, and run the gates in a temporary repository
- [x] 3.2 Run it in CI (`.github/workflows/cli.yml`)

## 4. Remove explore and document

- [x] 4.1 Remove `/explore` from the payload and the command count test (`cli/internal/payload/commands/explore.md`, `cli/internal/payload/embed_test.go`)
- [x] 4.2 Drop `/explore` from the README and the command reference (`README.md`, `docs/commands.md`)

## 5. Checks

- [x] 5.1 Run gofmt, vet, tests, `openspec validate --all --strict`, and the gates
