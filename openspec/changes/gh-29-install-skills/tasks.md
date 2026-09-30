## 1. Install the skills

- [ ] 1.1 Add a helper that maps the tool names and runs `openspec init --tools <list>` (`cli/internal/skills/skills.go`)
- [ ] 1.2 Call it from `install` and report the result, skipping with a warning when `openspec` is absent (`cli/internal/install/install.go`)

## 2. Tests

- [ ] 2.1 Test the tool-name mapping (`cli/internal/skills/skills_test.go`)
- [ ] 2.2 Test that `install` still succeeds when `openspec` is absent (`cli/internal/install/install_test.go`)

## 3. Docs

- [ ] 3.1 Note that `install` also brings the OpenSpec skills (`README.md`)

## 4. Verify

- [ ] 4.1 Install into a temp repo and confirm the four skill sets exist for each tool (`cli/`)
