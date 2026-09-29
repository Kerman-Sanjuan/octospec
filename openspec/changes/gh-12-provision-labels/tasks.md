## 1. Add the label seeder

- [ ] 1.1 Create `scripts/seed-labels.sh` defining the six labels and creating them idempotently with `gh label create --force` (`scripts/seed-labels.sh`)
- [ ] 1.2 Resolve the target repo (from `$1`, else cwd) and skip with a warning when `gh` is unavailable (`scripts/seed-labels.sh`)

## 2. Wire into the seed

- [ ] 2.1 Invoke `scripts/seed-labels.sh` from the `--repo` step of `install.sh` (`install.sh`)
- [ ] 2.2 Update the README install section to state that `--repo` provisions the labels (`README.md`)

## 3. Verify

- [ ] 3.1 Run `./install.sh --repo <fresh-repo>` twice and confirm the six labels exist and the second run exits 0 (`scripts/seed-labels.sh`, `install.sh`)
- [ ] 3.2 Confirm the labels match the names used by `commands/{idea,bug,spec,apply,ship,archive}.md` (`commands/`)
