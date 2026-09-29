## 1. Plan comment

- [x] 1.1 Make `/spec` post one plan comment with `### Proposal`, `### Capabilities`, `### Design`, `### Tasks`, recording the id in `.openspec.yaml` (`commands/spec.md`)
- [x] 1.2 Refresh the plan comment in place on re-run instead of posting per artifact (`commands/spec.md`)
- [x] 1.3 Add the ask-on-ambiguity rule: questions to the issue body, never invent requirements (`commands/spec.md`)

## 2. Implementation comment

- [x] 2.1 Make `/apply` post one implementation comment with the checklist and a short insight per task, recording the id in `.openspec.yaml` (`commands/apply.md`)
- [x] 2.2 Refresh the implementation comment in place after each group (`commands/apply.md`)

## 3. Schema and mirrors

- [ ] 3.1 Update the schema's `apply` instruction to describe the single implementation comment (`schema/octospec/schema.yaml`)
- [ ] 3.2 Regenerate the `spec` and `apply` mirrors from the canonical commands (`repo-template/.github/prompts/spec.prompt.md`, `.github/prompts/spec.prompt.md`, `repo-template/.github/prompts/apply.prompt.md`, `.github/prompts/apply.prompt.md`)

## 4. Docs

- [ ] 4.1 Document the three-layer structure in the README (`README.md`)

## 5. Verify

- [ ] 5.1 Dry-run `/spec` and `/apply` on a sample issue: confirm one plan comment and one implementation comment, refreshed in place on re-run (`commands/spec.md`, `commands/apply.md`)
