// Command octospec installs the octospec commands into each supported tool's
// native location from one canonical payload.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kerman-sanjuan/octospec/cli/internal/doctor"
	"github.com/kerman-sanjuan/octospec/cli/internal/install"
	"github.com/kerman-sanjuan/octospec/cli/internal/seed"
	"github.com/kerman-sanjuan/octospec/cli/internal/targets"
	"github.com/kerman-sanjuan/octospec/cli/internal/update"
	"github.com/spf13/cobra"
)

const version = "1.0.0"

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
	root.AddCommand(newInstallCmd(), newSeedCmd(), newUpdateCmd(), newDoctorCmd(), newVersionCmd())
	return root
}

func newInstallCmd() *cobra.Command {
	var tools []string
	var repo string
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install the commands into each tool's native location",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(tools) == 0 && isTerminal(os.Stdin) {
				t, err := wizard()
				if err != nil {
					return err
				}
				tools = t
			}
			return install.Run(install.Options{Tools: tools, Repo: repo})
		},
	}
	cmd.Flags().StringArrayVarP(&tools, "tool", "t", nil, "tool to install for (pi, opencode, copilot, claude); repeatable")
	cmd.Flags().StringVar(&repo, "repo", "", "repository root for repo-local tools (default: cwd)")
	return cmd
}

func newSeedCmd() *cobra.Command {
	var repo string
	var noLabels bool
	cmd := &cobra.Command{
		Use:   "seed",
		Short: "Install the OpenSpec schema, the repo seed files, and the labels",
		RunE: func(cmd *cobra.Command, args []string) error {
			if repo == "" {
				repo, _ = os.Getwd()
			}
			r, err := seed.Run(seed.Options{Repo: repo, SkipLabels: noLabels})
			if err != nil {
				return err
			}
			fmt.Printf("schema -> %s (%d files)\nrepo -> %s (%d files)\nlabels -> %d\n",
				seed.SchemaDir(), r.Schema, repo, r.Repo, r.Labels)
			return nil
		},
	}
	cmd.Flags().StringVar(&repo, "repo", "", "repository root (default: cwd)")
	cmd.Flags().BoolVar(&noLabels, "no-labels", false, "skip label provisioning")
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
	fi, err := f.Stat()
	return err == nil && (fi.Mode()&os.ModeCharDevice) != 0
}
