## 1. `/spec` owns the branch

- [x] 1.1 Create the change branch (`feat`/`fix` from the `type:*` label), commit the artifacts, and push (`commands/spec.md`)
- [x] 1.2 Head published comments with the artifact commit SHA and refresh them in place, recording ids in `.openspec.yaml` (`commands/spec.md`)
- [x] 1.3 Make re-run idempotent: reuse the existing change and branch (`commands/spec.md`)

## 2. `/apply` executes on it

- [x] 2.1 Remove branch creation, the dirty-tree special case, and the baseline commit; assume the `/spec` branch (`commands/apply.md`)
- [x] 2.2 Update the schema's `apply` instruction to match (`schema/octospec/schema.yaml`)

## 3. Approval

- [x] 3.1 Make "approach approved" a human action; `/spec` stops at `status:spec-ready` (`commands/spec.md`)

## 4. Mirrors and docs

- [x] 4.1 Mirror the `spec` change into `.github/prompts/spec.prompt.md` and `repo-template/.github/prompts/spec.prompt.md` (`.github/prompts/spec.prompt.md`, `repo-template/.github/prompts/spec.prompt.md`)
- [x] 4.2 Mirror the `apply` change into `.github/prompts/apply.prompt.md` and `repo-template/.github/prompts/apply.prompt.md` (`.github/prompts/apply.prompt.md`, `repo-template/.github/prompts/apply.prompt.md`)
- [x] 4.3 Update the README workflow loop to show `/spec` creating the branch (`README.md`)

## 5. Verify

- [ ] 5.1 Dry-run `/spec` on a sample issue: confirm branch + commit + push and in-place comment refresh (`commands/spec.md`)
- [ ] 5.2 Confirm `/apply` starts clean on the branch with no baseline commit (`commands/apply.md`)
