// Package models sets the model each agent role uses, interactively through a
// TUI or by flag.
package models

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/kerman-sanjuan/octospec/cli/internal/config"
)

// Options controls a model run.
type Options struct {
	Repo        string
	Pairs       map[string]string // role -> model
	Interactive bool              // open the TUI for the roles
}

// Set writes the role-to-model pairs, or opens the TUI on the interactive
// path. An empty model clears a role back to the tool default.
func Set(opts Options) error {
	repo := opts.Repo
	if repo == "" {
		var err error
		if repo, err = os.Getwd(); err != nil {
			return err
		}
	}
	cfg, err := config.Load(repo)
	if err != nil {
		return err
	}
	if cfg.Models == nil {
		cfg.Models = map[string]string{}
	}
	for role, model := range opts.Pairs {
		if !config.ValidRole(role) {
			return fmt.Errorf("unknown role %q (want one of %s)", role, strings.Join(config.Roles, ", "))
		}
		cfg.Models[role] = strings.TrimSpace(model)
	}
	if opts.Interactive {
		if err := form(&cfg); err != nil {
			return err
		}
	}
	return config.Save(repo, cfg)
}

// customChoice is the sentinel select value that reveals the free-form input.
const customChoice = "\x00custom"

// runForm runs a form. It is a variable so tests can drive the real form in
// huh's accessible mode without a terminal.
var runForm = func(f *huh.Form) error { return f.Run() }

// form opens the role-to-model TUI. Each role gets a select over the known
// models, plus a free-form input shown only when the operator picks
// "Custom...". An empty choice clears the role back to the tool default.
func form(cfg *config.Config) error {
	available := Catalog(cfg.Tools)
	roles := config.Roles
	choices := make([]string, len(roles))
	customs := make([]string, len(roles))
	groups := make([]*huh.Group, 0, len(roles)*2)
	for i, role := range roles {
		choices[i], customs[i] = seedChoice(cfg.Models[role], available)

		options := make([]huh.Option[string], 0, len(available)+2)
		options = append(options, huh.NewOption("(tool default)", ""))
		for _, m := range available {
			options = append(options, huh.NewOption(m, m))
		}
		options = append(options, huh.NewOption("Custom...", customChoice))

		selectIdx, customIdx := i, i
		groups = append(groups,
			huh.NewGroup(huh.NewSelect[string]().
				Title(role).
				Description(roleHint(role)).
				Options(options...).
				Value(&choices[selectIdx])),
			huh.NewGroup(huh.NewInput().
				Title(role+" (custom)").
				Description("Type any model identifier.").
				Value(&customs[customIdx])).
				WithHideFunc(func() bool { return choices[selectIdx] != customChoice }),
		)
	}
	if err := runForm(huh.NewForm(groups...)); err != nil {
		return err
	}
	for i, role := range roles {
		if choices[i] == customChoice {
			cfg.Models[role] = strings.TrimSpace(customs[i])
			continue
		}
		cfg.Models[role] = strings.TrimSpace(choices[i])
	}
	return nil
}

// seedChoice maps a stored model to a select value and a custom value. A model
// in the catalog selects itself; anything else selects "Custom..." and seeds
// the input; an empty model selects the tool default.
func seedChoice(model string, available []string) (choice, custom string) {
	switch {
	case model == "":
		return "", ""
	case contains(available, model):
		return model, ""
	default:
		return customChoice, model
	}
}

func contains(list []string, want string) bool {
	for _, s := range list {
		if s == want {
			return true
		}
	}
	return false
}

func roleHint(role string) string {
	switch role {
	case config.RoleThinking:
		return "Used by the idea and spec stages."
	case config.RoleImplementer:
		return "Used by the apply stage."
	case config.RoleReviewer:
		return "Used by the ship and archive stages."
	}
	return ""
}
