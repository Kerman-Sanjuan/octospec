## Context

`octospec install` decides which tools to target in `cli/cmd/octospec/main.go`:
when no `--tool` flag is passed and stdin is a terminal, it calls `wizard()`,
which prints a list and reads a comma-or-space separated line with
`bufio.Reader`. `cli/internal/install/install.go` then treats an empty tool
slice as "all tools". The models command already runs a huh form
(`cli/internal/models/models.go`), with one `huh.NewInput` per role. The huh
dependency (`github.com/charmbracelet/huh v1.0.0`) is present, and the targets
list in `cli/internal/targets/targets.go` is the single source of tool names.

## Goals / Non-Goals

**Goals:**
- Replace the blank-line tool prompt with a huh multi-select, seeded from the
  previous install.
- Replace the model text inputs with a select from a static known-model catalog
  per tool, plus a free-form fallback.
- Keep flags and non-terminal runs byte-for-byte the same.

**Non-Goals:**
- No new tools, no change to `targets.Targets`.
- No change to the `.octospec/octospec.json` schema.
- No runtime detection of models from the tools.

## Decisions

### Use huh for the tool selector
huh is already a direct dependency and already drives the models form, so the
tool selector uses `huh.NewMultiSelect`. The alternative, extending the
`bufio` prompt, keeps the code smaller but does not fix the typo problem or add
a browse experience, which is the point of the issue.

### Move the install tool selection into the install package
`wizard()` lives in `main.go` because it is I/O. The multi-select is a huh form
like `askScope()`, which already lives in `cli/internal/install/install.go`. Put
the selector beside `askScope()` and call it from `install.Run` on the same
interactive path, so `main.go` only decides interactive versus not. This keeps
all install-time prompts in one package and lets a test call the selection
helper with an injected set of names. The alternative, leaving it in `main.go`,
would keep TUI code in the command layer and make it untestable.

### Seed the selector from `cfg.Tools`
`install.Run` already loads the config before it plans files, so the recorded
`Tools` are in hand. Pre-check those names when they are in the current targets
list. This answers the issue's open question: a re-run starts from the previous
choice, and a first install starts empty.

### A static model catalog in the models package
Add the known models per tool as a literal in `cli/internal/models/` (a new
`catalog.go`). The form offers them through `huh.NewSelect` with an extra
free-form option that reveals a `huh.NewInput`. A static list is what the issue
chose; runtime detection is out of scope. Keeping it in the models package puts
the data next to its only consumer and out of `targets`, which describes
locations, not models.

### Empty selection means all
`huh.NewMultiSelect` can return an empty slice. `install.Run` already treats an
empty tool slice as all tools, so the selector passes the empty slice through
and the existing behaviour holds without a special case.

## Risks / Trade-offs

- [A model catalog goes stale as providers ship models] -> The list is a
  convenience, not a gate: the free-form fallback always accepts any string, so
  a stale list never blocks a valid model.
- [huh forms need a TTY and are hard to unit test] -> Keep the TUI code thin
  (build the options, read the result) and test the pure parts (target names,
  catalog contents, tool-slice handling) without a terminal.
- [Pre-checking could surprise an operator who wants a fresh choice] -> The
  checked boxes are visible and toggleable; the previous install is a sensible
  default, and unchecking all still installs all, which the prompt explains.

## Migration Plan

No data migration. Existing `.octospec/octospec.json` files already carry
`tools`; the selector reads that field. Rollback is a revert: the flag and
non-terminal paths never changed.

## Open Questions

- None. The issue's two open questions are resolved: the catalog lives in the
  models package, and the selector pre-checks `cfg.Tools`.
