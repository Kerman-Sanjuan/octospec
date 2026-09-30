## 1. Doctor

- [ ] 1.1 Implement the checks: openspec CLI, schema, installed commands, labels, gate script (`cli/internal/doctor/doctor.go`)
- [ ] 1.2 Add the `doctor` command with `--repo` and a non-zero exit when a check fails (`cli/cmd/octospec/main.go`)

## 2. Tests

- [ ] 2.1 Test a healthy repo and a repo missing the schema (`cli/internal/doctor/doctor_test.go`)

## 3. Docs

- [ ] 3.1 Add `octospec doctor` to the README install section (`README.md`)

## 4. Verify

- [ ] 4.1 Run `doctor` on a seeded repo and on a bare repo, and confirm the exit codes (`cli/`)
