## 1. Model catalog

- [x] 1.1 Add `cli/internal/models/catalog.go` with a known-model list per tool and a lookup that returns the models for the installed tools.
- [x] 1.2 Add a `catalog_test.go` case that every tool in `targets.Targets` has a catalog entry and every entry is non-empty.

## 2. Model selector

- [x] 2.1 Rewrite `form` in `cli/internal/models/models.go` to use a per-role `huh.NewSelect` over the catalog plus a free-form fallback that reveals a `huh.NewInput`.
- [x] 2.2 Keep the existing empty-model-clears-the-role behaviour and the flag path unchanged in `cli/internal/models/models.go`.

## 3. Tool selector

- [x] 3.1 Add a huh multi-select helper in `cli/internal/install/install.go` that lists the targets with a short description and preselects the names in `cfg.Tools`.
- [x] 3.2 Call the selector from the interactive path in `install.Run`, passing the selection as the tool set and letting an empty selection mean all tools.
- [x] 3.3 Add a one-line description per tool, from the target name, for the selector labels.
- [x] 3.4 Remove `wizard()` from `cli/cmd/octospec/main.go` and let `install.Run` own the interactive choice.

## 4. Tests and docs

- [x] 4.1 Add tests in `cli/internal/install/install_test.go` for the preselection and the empty-means-all handling, without a terminal.
- [x] 4.2 Add tests in `cli/internal/models/models_test.go` for the catalog lookup and the free-form fallback.
- [x] 4.3 Update the README `Agents and models` section and the install example to describe the selector.
- [x] 4.4 Run `gofmt -l`, `go vet ./...`, `go test ./...`, `openspec validate --all --strict`, and `sh scripts/check-gates.sh`.

## 5. Test hardening

- [x] 5.1 Add a `runForm` seam in `cli/internal/models/models.go` and `cli/internal/install/install.go` so tests drive the real huh forms in accessible mode.
- [x] 5.2 Integration-test the model form in `cli/internal/models/models_test.go`: catalog choice, custom model, and a seeded custom model.
- [x] 5.3 Integration-test `models.Set` on the interactive path, including a `--set` pair combined with the form.
- [x] 5.4 Integration-test `install.Run` on the interactive path in `cli/internal/install/install_test.go`: subset, none-means-all, preselection, global confirm, and no form without a terminal.
- [x] 5.5 Run the full suite with `-race` and confirm the coverage of the touched packages.
