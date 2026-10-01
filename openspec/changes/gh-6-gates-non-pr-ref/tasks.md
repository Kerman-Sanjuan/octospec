## 1. Fix the gate script

- [ ] 1.1 Skip the PR-only branch gate on a long-lived branch (`scripts/check-gates.sh`)

## 2. Guard and document

- [ ] 2.1 Add a gate test for the skip and the PR check (`scripts/check-gates-test.sh`)
- [ ] 2.2 Run the gate test in the gates workflow (`.github/workflows/openspec.yml`)
- [ ] 2.3 Note the PR-only scope in the README gates table (`README.md`)

## 3. Checks

- [ ] 3.1 Run gofmt, vet, tests, `openspec validate --all --strict`, and the gates (`scripts/check-gates.sh`)
