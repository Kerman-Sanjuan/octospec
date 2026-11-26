// Command octospec installs the octospec commands into each supported tool's
// native location from one canonical payload.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/kerman-sanjuan/octospec/cli/internal/doctor"
	"github.com/kerman-sanjuan/octospec/cli/internal/install"
	"github.com/kerman-sanjuan/octospec/cli/internal/lock"
	"github.com/kerman-sanjuan/octospec/cli/internal/models"
	"github.com/kerman-sanjuan/octospec/cli/internal/result"
	"github.com/kerman-sanjuan/octospec/cli/internal/review"
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
	root.AddCommand(newInstallCmd(), newSeedCmd(), newUpdateCmd(), newUninstallCmd(), newDoctorCmd(), newModelsCmd(), newLockCmd(), newReviewCmd(), newResultCmd(), newVersionCmd())
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
	var wait, acquireTTL time.Duration
	acquire := &cobra.Command{
		Use:   "acquire",
		Short: "Take the archive lock for a change",
		RunE: func(c *cobra.Command, args []string) error {
			root, err := repoRoot(repo)
			if err != nil {
				return err
			}
			rec, err := lock.Acquire(root, acquireChange, wait, acquireTTL)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(lockExitCode(err))
			}
			fmt.Printf("lock held by %s on %s, lease expires in %s\n",
				rec.Change, rec.Host, time.Until(rec.Expiry()).Round(time.Minute))
			return nil
		},
	}
	acquire.Flags().StringVar(&acquireChange, "change", "", "change that will hold the lock")
	acquire.Flags().DurationVar(&wait, "wait", 0, "how long to wait for a free lock, e.g. 2m")
	acquire.Flags().DurationVar(&acquireTTL, "ttl", 0, "how long the lease lasts without a heartbeat (default 30m)")

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
				fmt.Printf("stale: %s; held by %s on %s; clear it with `octospec lock release --force --change %s`\n",
					st.Why, st.Record.Change, st.Record.Host, st.Record.Change)
			default:
				fmt.Printf("held by %s since %s on %s (pid %d), lease expires in %s\n",
					st.Record.Change, st.Record.AcquiredAt.Format(time.RFC3339), st.Record.Host, st.Record.PID,
					time.Until(st.Record.Expiry()).Round(time.Minute))
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

	var heartbeatChange string
	var heartbeatTTL time.Duration
	heartbeat := &cobra.Command{
		Use:   "heartbeat",
		Short: "Extend the lease you already hold",
		RunE: func(c *cobra.Command, args []string) error {
			root, err := repoRoot(repo)
			if err != nil {
				return err
			}
			rec, err := lock.Heartbeat(root, heartbeatChange, heartbeatTTL)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(lockExitCode(err))
			}
			fmt.Printf("lease extended for %s, expires in %s\n",
				rec.Change, time.Until(rec.Expiry()).Round(time.Minute))
			return nil
		},
	}
	heartbeat.Flags().StringVar(&heartbeatChange, "change", "", "change that holds the lock")
	heartbeat.Flags().DurationVar(&heartbeatTTL, "ttl", 0, "how long the lease lasts from now (default 30m)")

	cmd.AddCommand(acquire, status, release, heartbeat)
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

// reviewExitCode is what a CI job blocks on. 5 means "the review asks for a human", which is not
// the same as the command failing (1) and is deliberately distinct from the lock's 3 and 4.
const reviewExitCode = 5

// newResultCmd is the contract between a stage and whatever runs it: one object per execution, with
// a closed vocabulary, so a scheduler decides retry or park without reading an agent's prose.
func newResultCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "result",
		Short: "Write or check the record a stage leaves behind",
	}
	var dir string
	var r result.Result
	var sha string
	var asJSON bool
	var artifacts []string
	write := &cobra.Command{
		Use:   "write",
		Short: "Write the result of one stage execution",
		RunE: func(c *cobra.Command, args []string) error {
			r.Artifacts = artifacts
			r.Normalize()
			if err := r.Validate(); err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(1)
			}
			p, err := result.Write(dir, r, sha)
			if err != nil {
				return err
			}
			if asJSON {
				b, err := json.MarshalIndent(r, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(b))
				return nil
			}
			fmt.Println(p)
			return nil
		},
	}
	write.Flags().StringVar(&dir, "dir", ".octospec/runs", "directory for the records")
	write.Flags().StringVar(&r.Stage, "stage", "", "stage that ran: "+strings.Join(result.Stages, ", "))
	write.Flags().IntVar(&r.Issue, "issue", 0, "issue the stage worked on")
	write.Flags().StringVar(&r.Change, "change", "", "change the stage worked on")
	write.Flags().StringVar(&r.Outcome, "outcome", "", "ok or stalled")
	write.Flags().StringVar(&r.Reason, "reason", "", "why it stalled: "+strings.Join(result.Reasons, ", "))
	write.Flags().StringArrayVar(&artifacts, "artifact", nil, "URL or path a person can open; repeatable")
	write.Flags().Int64Var(&r.DurationMS, "duration-ms", 0, "how long the stage took")
	write.Flags().StringVar(&sha, "sha", "", "commit the stage ran on (default: HEAD)")
	write.Flags().BoolVar(&asJSON, "json", false, "print the object instead of the path")

	show := &cobra.Command{
		Use:   "show <file>",
		Short: "Print one record, and fail when it does not validate",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			rec, err := result.Read(args[0])
			if err != nil {
				return err
			}
			b, err := json.MarshalIndent(rec, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(b))
			return nil
		},
	}

	cmd.AddCommand(write, show)
	return cmd
}

// newReviewCmd scores a pull request from a list of findings. The command reports by default and
// blocks only with --enforce, so a threshold can be set from a real distribution of scores.
func newReviewCmd() *cobra.Command {
	var repo, findingsPath string
	var asJSON, enforce, familiar bool
	var changedFiles []string
	var diffLines int

	cmd := &cobra.Command{
		Use:   "review",
		Short: "Score a pull request from a list of findings and decide whether it needs a human",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			root, err := repoRoot(repo)
			if err != nil {
				return err
			}
			var in io.Reader = os.Stdin
			if findingsPath != "" && findingsPath != "-" {
				f, err := os.Open(findingsPath)
				if err != nil {
					return err
				}
				defer f.Close()
				in = f
			}
			findings, err := review.Parse(in)
			if err != nil {
				return err
			}
			policy, err := review.LoadPolicy(root)
			if err != nil {
				return err
			}
			res := review.Score(findings, policy, review.RiskInput{
				ChangedFiles: changedFiles,
				DiffLines:    diffLines,
				Familiar:     familiar,
			})

			if asJSON {
				b, err := json.MarshalIndent(res, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(b))
			} else {
				printReview(res)
			}

			if enforce && res.Verdict != review.VerdictPass {
				os.Exit(reviewExitCode)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&findingsPath, "findings", "", "file with the findings; `-` or empty reads stdin")
	cmd.Flags().BoolVar(&asJSON, "json", false, "print the result as JSON")
	cmd.Flags().BoolVar(&enforce, "enforce", false, "exit 5 when the verdict asks for a human")
	cmd.Flags().StringSliceVar(&changedFiles, "changed-files", nil, "paths the pull request touches, for the risk table")
	cmd.Flags().IntVar(&diffLines, "diff-lines", 0, "added plus removed lines, for the risk table")
	cmd.Flags().BoolVar(&familiar, "familiar", false, "the author has touched these paths before")
	cmd.Flags().StringVar(&repo, "repo", "", "repository root (default: cwd)")
	return cmd
}

// printReview renders the result for a person: the number, where it was lost, and why the risk
// is what it is.
func printReview(r review.Result) {
	fmt.Printf("score %.1f/10 (threshold %.1f) · risk %s · %s\n", r.Score, r.Threshold, r.Risk, r.Verdict)
	fmt.Printf("policy %s\n", r.PolicySource)
	if len(r.Deductions) > 0 {
		fmt.Println("\ndeductions")
		for _, d := range r.Deductions {
			where := d.File
			if d.Line > 0 {
				where = fmt.Sprintf("%s:%d", d.File, d.Line)
			}
			note := d.Note
			if note == "" {
				note = string(d.Type)
			}
			fmt.Printf("  -%.1f  %-13s %s  %s\n", d.Weight, d.Axis, where, note)
		}
	}
	if len(r.ByAxis) > 0 {
		fmt.Println("\nby axis")
		for _, axis := range review.Axes {
			if v, ok := r.ByAxis[axis]; ok {
				fmt.Printf("  %-13s -%.1f\n", axis, v)
			}
		}
	}
	if len(r.RiskWhy) > 0 {
		fmt.Println("\nrisk")
		for _, why := range r.RiskWhy {
			fmt.Printf("  %s\n", why)
		}
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
