## 1. Linter

- [ ] 1.1 Add the `golangci-lint` config (`cli/.golangci.yml`)
- [ ] 1.2 Run `golangci-lint` in the `cli` job (`.github/workflows/cli.yml`)

## 2. Race and coverage

- [ ] 2.1 Run `go test -race` with a coverage profile and report the total (`.github/workflows/cli.yml`)

## 3. Checks

- [ ] 3.1 Run gofmt, vet, tests, `openspec validate --all --strict`, and the gates (`scripts/check-gates.sh`)
