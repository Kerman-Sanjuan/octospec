## 1. Model catalog

- [x] 1.1 Add `cli/internal/models/catalog.go` with a known-model list per tool and a lookup that returns the models for the installed tools.
- [x] 1.2 Add a `catalog_test.go` case that every tool in `targets.Targets` has a catalog entry and every entry is non-empty.

## 2. Model selector

- [ ] 2.1 Rewrite `form` in `cli/internal/models/models.go` to use a per-role `huh.NewSelect` over the catalog plus a free-form fallback that reveals a `huh.NewInput`.
- [ ] 2.2 Keep the existing empty-model-clears-the-role behaviour and the flag path unchanged in `cli/internal/models/models.go`.

## 3. Tool selector

- [ ] 3.1 Add a huh multi-select helper in `cli/internal/install/install.go` that lists the targets with a short description and preselects the names in `cfg.Tools`.
- [ ] 3.2 Call the selector from the interactive path in `install.Run`, passing the selection as the tool set and letting an empty selection mean all tools.
- [ ] 3.3 Add a one-line description per tool, from the target name, for the selector labels.
- [ ] 3.4 Remove `wizard()` from `cli/cmd/octospec/main.go` and let `install.Run` own the interactive choice.

## 4. Tests and docs

- [ ] 4.1 Add tests in `cli/internal/install/install_test.go` for the preselection and the empty-means-all handling, without a terminal.
- [ ] 4.2 Add tests in `cli/internal/models/models_test.go` for the catalog lookup and the free-form fallback.
- [ ] 4.3 Update the README `Agents and models` section and the install example to describe the selector.
- [ ] 4.4 Run `gofmt -l`, `go vet ./...`, `go test ./...`, `openspec validate --all --strict`, and `sh scripts/check-gates.sh`.
