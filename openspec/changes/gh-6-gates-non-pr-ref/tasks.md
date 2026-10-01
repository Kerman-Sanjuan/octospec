## 1. Fix the gate script

- [x] 1.1 Skip the PR-only branch gate on a long-lived branch (`scripts/check-gates.sh`)

## 2. Guard and document

- [x] 2.1 Add a gate test for the skip and the PR check (`scripts/check-gates-test.sh`)
- [x] 2.2 Run the gate test in the gates workflow (`.github/workflows/openspec.yml`)
- [x] 2.3 Note the PR-only scope in the README gates table (`README.md`)

## 3. Checks

- [x] 3.1 Run gofmt, vet, tests, `openspec validate --all --strict`, and the gates (`scripts/check-gates.sh`)
