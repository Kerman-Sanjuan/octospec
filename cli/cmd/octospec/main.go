// Command octospec installs the octospec commands into each supported tool's
// native location from one canonical payload.
package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kerman-sanjuan/octospec/cli/internal/doctor"
	"github.com/kerman-sanjuan/octospec/cli/internal/install"
	"github.com/kerman-sanjuan/octospec/cli/internal/lock"
	"github.com/kerman-sanjuan/octospec/cli/internal/models"
	"github.com/kerman-sanjuan/octospec/cli/internal/seed"
	"github.com/kerman-sanjuan/octospec/cli/internal/uninstall"
	"github.com/kerman-sanjuan/octospec/cli/internal/update"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// version is the release the binary was built from. The release build injects
// the tag with -X main.version=<tag>; a plain build keeps "dev".
var version = "dev"

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "octospec",
		Short:         "Install the octospec spec-driven workflow into your tools",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(newInstallCmd(), newSeedCmd(), newUpdateCmd(), newUninstallCmd(), newDoctorCmd(), newModelsCmd(), newLockCmd(), newVersionCmd())
	return root
}

// newLockCmd exposes the lock that guards openspec/specs during /archive, so that a scheduler can
// take it and a human can see who holds it or clear a dead one.
func newLockCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock",
		Short: "Inspect or manage the lock that guards openspec/specs during /archive",
	}
	var repo string
	cmd.PersistentFlags().StringVar(&repo, "repo", "", "repository root (default: cwd)")

	var acquireChange string
	var wait time.Duration
	acquire := &cobra.Command{
		Use:   "acquire",
		Short: "Take the archive lock for a change",
		RunE: func(c *cobra.Command, args []string) error {
			root, err := repoRoot(repo)
			if err != nil {
				return err
			}
			rec, err := lock.Acquire(root, acquireChange, wait)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(lockExitCode(err))
			}
			fmt.Printf("lock held by %s (pid %d on %s)\n", rec.Change, rec.PID, rec.Host)
			return nil
		},
	}
	acquire.Flags().StringVar(&acquireChange, "change", "", "change that will hold the lock")
	acquire.Flags().DurationVar(&wait, "wait", 0, "how long to wait for a free lock, e.g. 2m")

	status := &cobra.Command{
		Use:   "status",
		Short: "Report the lock state without modifying it",
		RunE: func(c *cobra.Command, args []string) error {
			root, err := repoRoot(repo)
			if err != nil {
				return err
			}
			st, err := lock.Status(root)
			if err != nil {
				return err
			}
			switch {
			case !st.Held:
				fmt.Println("free")
			case st.Stale:
				fmt.Printf("stale: held by %s on %s (%s); clear it with `octospec lock release --force --change %s`\n",
					st.Record.Change, st.Record.Host, st.Why, st.Record.Change)
			default:
				fmt.Printf("held by %s since %s on %s (pid %d)\n",
					st.Record.Change, st.Record.AcquiredAt.Format(time.RFC3339), st.Record.Host, st.Record.PID)
			}
			return nil
		},
	}

	var releaseChange string
	var force bool
	release := &cobra.Command{
		Use:   "release",
		Short: "Release the archive lock, optionally clearing a dead holder",
		RunE: func(c *cobra.Command, args []string) error {
			root, err := repoRoot(repo)
			if err != nil {
				return err
			}
			if err := lock.Release(root, releaseChange, force); err != nil {
				return err
			}
			fmt.Println("lock free")
			return nil
		},
	}
	release.Flags().StringVar(&releaseChange, "change", "", "change that holds the lock")
	release.Flags().BoolVar(&force, "force", false, "clear the lock even if another change holds it")

	cmd.AddCommand(acquire, status, release)
	return cmd
}

// lockExitCode gives a scheduler a stable code to act on instead of parsing English: 3 means
// somebody else holds the lock and it may retry later, 4 means the holder is gone and a human
// has to decide.
func lockExitCode(err error) int {
	switch {
	case errors.Is(err, lock.ErrStale):
		return 4
	case errors.Is(err, lock.ErrHeld):
		return 3
	default:
		return 1
	}
}

func repoRoot(repo string) (string, error) {
	if repo != "" {
		return repo, nil
	}
	return os.Getwd()
}

func newUninstallCmd() *cobra.Command {
	var repo string
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Remove the files octospec manages, keeping local edits",
		RunE: func(cmd *cobra.Command, args []string) error {
			return uninstall.Run(repo)
		},
	}
	cmd.Flags().StringVar(&repo, "repo", "", "repository root (default: cwd)")
	return cmd
}

func newModelsCmd() *cobra.Command {
	var repo string
	var sets []string
	cmd := &cobra.Command{
		Use:   "models",
		Short: "Set the model each agent role uses",
		RunE: func(cmd *cobra.Command, args []string) error {
			pairs := map[string]string{}
			for _, s := range sets {
				role, model, ok := strings.Cut(s, "=")
				if !ok {
					return fmt.Errorf("want role=model, got %q", s)
				}
				pairs[strings.TrimSpace(role)] = strings.TrimSpace(model)
			}
			interactive := len(pairs) == 0 && isTerminal(os.Stdin)
			return models.Set(models.Options{Repo: repo, Pairs: pairs, Interactive: interactive})
		},
	}
	cmd.Flags().StringArrayVar(&sets, "set", nil, "set a role's model (role=model); repeatable")
	cmd.Flags().StringVar(&repo, "repo", "", "repository root (default: cwd)")
	return cmd
}

func newInstallCmd() *cobra.Command {
	var tools []string
	var repo string
	var global bool
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install the commands into each tool's native location",
		RunE: func(cmd *cobra.Command, args []string) error {
			interactive := !global && !dryRun && isTerminal(os.Stdin)
			return install.Run(install.Options{Tools: tools, Repo: repo, Global: global, Interactive: interactive, DryRun: dryRun})
		},
	}
	cmd.Flags().StringArrayVarP(&tools, "tool", "t", nil, "tool to install for (pi, opencode, copilot, claude); repeatable")
	cmd.Flags().StringVar(&repo, "repo", "", "repository root for repo-local tools (default: cwd)")
	cmd.Flags().BoolVar(&global, "global", false, "install into each tool's global directory (warned)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print what would be written and change nothing")
	return cmd
}

func newSeedCmd() *cobra.Command {
	var repo string
	var noLabels bool
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "seed",
		Short: "Install the OpenSpec schema, the repo seed files, and the labels",
		RunE: func(cmd *cobra.Command, args []string) error {
			if repo == "" {
				repo, _ = os.Getwd()
			}
			r, err := seed.Run(seed.Options{Repo: repo, SkipLabels: noLabels, DryRun: dryRun})
			if err != nil {
				return err
			}
			if dryRun {
				return nil
			}
			fmt.Printf("schema -> %s (%d files)\nrepo -> %s (%d files)\nlabels -> %d\n",
				seed.SchemaDir(), r.Schema, repo, r.Repo, r.Labels)
			return nil
		},
	}
	cmd.Flags().StringVar(&repo, "repo", "", "repository root (default: cwd)")
	cmd.Flags().BoolVar(&noLabels, "no-labels", false, "skip label provisioning")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print what would be installed and change nothing")
	return cmd
}

func newUpdateCmd() *cobra.Command {
	var repo string
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Re-apply an install, preserving local edits",
		RunE: func(cmd *cobra.Command, args []string) error {
			return update.Run(repo)
		},
	}
	cmd.Flags().StringVar(&repo, "repo", "", "repository root (default: cwd)")
	return cmd
}

func newDoctorCmd() *cobra.Command {
	var repo string
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Report the state of an octospec install",
		RunE: func(cmd *cobra.Command, args []string) error {
			checks, err := doctor.Run(repo)
			if err != nil {
				return err
			}
			bad := 0
			for _, c := range checks {
				mark := "OK  "
				if !c.OK {
					mark = "MISS"
					bad++
				}
				if c.Note != "" {
					fmt.Printf("%s %-24s %s\n", mark, c.Name, c.Note)
				} else {
					fmt.Printf("%s %s\n", mark, c.Name)
				}
			}
			if bad > 0 {
				return fmt.Errorf("%d check(s) missing", bad)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&repo, "repo", "", "repository root (default: cwd)")
	return cmd
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("octospec", version)
		},
	}
}

func isTerminal(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}
