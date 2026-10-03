// Command octospec installs the octospec commands into each supported tool's
// native location from one canonical payload.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kerman-sanjuan/octospec/cli/internal/doctor"
	"github.com/kerman-sanjuan/octospec/cli/internal/install"
	"github.com/kerman-sanjuan/octospec/cli/internal/models"
	"github.com/kerman-sanjuan/octospec/cli/internal/seed"
	"github.com/kerman-sanjuan/octospec/cli/internal/session"
	"github.com/kerman-sanjuan/octospec/cli/internal/targets"
	"github.com/kerman-sanjuan/octospec/cli/internal/uninstall"
	"github.com/kerman-sanjuan/octospec/cli/internal/update"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const version = "1.1.0"

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
	root.AddCommand(newInstallCmd(), newSeedCmd(), newUpdateCmd(), newUninstallCmd(), newDoctorCmd(), newModelsCmd(), newSessionCmd(), newVersionCmd())
	return root
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
			if len(tools) == 0 && isTerminal(os.Stdin) && !dryRun {
				t, err := wizard()
				if err != nil {
					return err
				}
				tools = t
			}
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

func newSessionCmd() *cobra.Command {
	var repo string
	cmd := &cobra.Command{
		Use:   "session",
		Short: "Run independent agents in parallel, one worktree per issue",
	}
	cmd.PersistentFlags().StringVar(&repo, "repo", "", "repository root (default: cwd)")

	var tool string
	start := &cobra.Command{
		Use:   "start <issue>",
		Short: "Create an isolated worktree for an issue and launch its agent",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			issue, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("issue must be a number, got %q", args[0])
			}
			s, err := session.Start(repo, issue, tool)
			if err != nil {
				return err
			}
			launched, err := session.Launch(s, tool)
			if err != nil {
				return err
			}
			fmt.Printf("session for #%d: %s on %s\n", s.Issue, s.Worktree, s.Branch)
			if launched {
				fmt.Printf("launched %s in %s\n", tool, s.Worktree)
			} else {
				fmt.Printf("open your agent in %s, then run /spec %d\n", s.Worktree, s.Issue)
			}
			return nil
		},
	}
	start.Flags().StringVar(&tool, "tool", "", "tool to launch in the worktree (opencode, claude, copilot, pi)")

	list := &cobra.Command{
		Use:   "list",
		Short: "Show the active sessions",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			sessions, err := session.List(repo)
			if err != nil {
				return err
			}
			if len(sessions) == 0 {
				fmt.Println("no active sessions")
				return nil
			}
			for _, s := range sessions {
				fmt.Printf("#%d  %s  %s  pid %d@%s  %s\n",
					s.Issue, s.Branch, s.Worktree, s.PID, s.Host, age(s.StartedAt))
			}
			return nil
		},
	}

	end := &cobra.Command{
		Use:   "end <issue>",
		Short: "Remove an issue's worktree and unregister it",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			issue, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("issue must be a number, got %q", args[0])
			}
			if err := session.End(repo, issue); err != nil {
				return err
			}
			fmt.Printf("session for #%d ended\n", issue)
			return nil
		},
	}

	cmd.AddCommand(start, list, end)
	return cmd
}

// age renders how long ago a session started.
func age(start time.Time) string {
	d := time.Since(start)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	default:
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
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

// wizard asks which tools to install for. An empty answer means all of them.
func wizard() ([]string, error) {
	names := make([]string, 0, len(targets.Targets))
	for _, t := range targets.Targets {
		names = append(names, t.Name)
	}
	fmt.Printf("Install octospec for which tools? [%s] (blank = all)\n", strings.Join(names, " "))
	fmt.Print("> ")
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return nil, nil // no input available: install for all tools
	}
	line = strings.TrimSpace(line)
	if line == "" || line == "all" {
		return nil, nil
	}
	var out []string
	for _, p := range strings.FieldsFunc(line, func(r rune) bool { return r == ',' || r == ' ' }) {
		if _, ok := targets.Get(p); !ok {
			return nil, fmt.Errorf("unknown tool %q", p)
		}
		out = append(out, p)
	}
	return out, nil
}

func isTerminal(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}
