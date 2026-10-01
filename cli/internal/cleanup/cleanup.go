// Package cleanup reports and removes the files that an earlier global install
// left behind when a tool moves to a repo-local install.
package cleanup

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/kerman-sanjuan/octospec/cli/internal/config"
)

// Global reports the recorded files under the home directory that the new plan
// no longer writes. On a terminal it offers to remove them; otherwise it lists
// them and does not prompt.
func Global(cfg *config.Config, keep map[string]bool, interactive bool) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return
	}
	var leftovers []string
	for path := range cfg.Files {
		if keep[path] {
			continue
		}
		if strings.HasPrefix(path, home+string(filepath.Separator)) {
			leftovers = append(leftovers, path)
		}
	}
	if len(leftovers) == 0 {
		return
	}
	sort.Strings(leftovers)

	if interactive {
		fmt.Printf("octospec: %d file(s) from a previous global install are no longer managed.\n", len(leftovers))
		remove := false
		confirm := huh.NewConfirm().
			Title("Remove the leftover global files?").
			Description("They sit under your home directory and were installed by an earlier version.").
			Affirmative("Remove").
			Negative("Keep").
			Value(&remove)
		if err := huh.NewForm(huh.NewGroup(confirm)).Run(); err == nil && remove {
			removed := 0
			for _, path := range leftovers {
				if os.Remove(path) == nil {
					removed++
				}
				delete(cfg.Files, path)
			}
			fmt.Printf("removed %d leftover file(s)\n", removed)
			return
		}
	}

	fmt.Printf("note: %d file(s) from a previous global install remain under your home directory. Remove them if you moved to a repo-local install:\n", len(leftovers))
	for _, path := range leftovers {
		fmt.Printf("  %s\n", path)
	}
}
