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

// form opens the role-to-model TUI.
func form(cfg *config.Config) error {
	values := make([]string, len(config.Roles))
	fields := make([]huh.Field, 0, len(config.Roles))
	for i, role := range config.Roles {
		values[i] = cfg.Models[role]
		fields = append(fields, huh.NewInput().
			Title(role).
			Description(roleHint(role)).
			Value(&values[i]))
	}
	if err := huh.NewForm(huh.NewGroup(fields...)).Run(); err != nil {
		return err
	}
	for i, role := range config.Roles {
		cfg.Models[role] = strings.TrimSpace(values[i])
	}
	return nil
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
