// Command octospec installs the octospec commands into each supported tool's
// native location from one canonical payload.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kerman-sanjuan/octospec/cli/internal/install"
	"github.com/kerman-sanjuan/octospec/cli/internal/seed"
	"github.com/kerman-sanjuan/octospec/cli/internal/targets"
	"github.com/kerman-sanjuan/octospec/cli/internal/update"
)

const version = "0.1.0"

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		usage()
		os.Exit(2)
	}
	switch args[0] {
	case "install":
		opts, err := parseInstall(args[1:])
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(2)
		}
		if len(opts.Tools) == 0 && isTerminal(os.Stdin) {
			opts.Tools, err = wizard()
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(2)
			}
		}
		if err := install.Run(opts); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	case "update":
		repo := ""
		for i := 1; i < len(args); i++ {
			if args[i] == "--repo" && i+1 < len(args) {
				i++
				repo = args[i]
			}
		}
		if err := update.Run(repo); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	case "seed":
		repo, skipLabels := "", false
		for i := 1; i < len(args); i++ {
			switch args[i] {
			case "--repo":
				if i+1 >= len(args) {
					fmt.Fprintln(os.Stderr, "error: --repo needs a value")
					os.Exit(2)
				}
				i++
				repo = args[i]
			case "--no-labels":
				skipLabels = true
			default:
				fmt.Fprintf(os.Stderr, "error: unknown flag %q\n", args[i])
				os.Exit(2)
			}
		}
		if repo == "" {
			repo, _ = os.Getwd()
		}
		r, err := seed.Run(seed.Options{Repo: repo, SkipLabels: skipLabels})
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Printf("schema -> %s (%d files)\nrepo -> %s (%d files)\nlabels -> %d\n",
			seed.SchemaDir(), r.Schema, repo, r.Repo, r.Labels)
	case "version", "--version", "-v":
		fmt.Println("octospec", version)
	case "help", "--help", "-h":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", args[0])
		usage()
		os.Exit(2)
	}
}

func parseInstall(args []string) (install.Options, error) {
	var opts install.Options
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--tool", "-t":
			i++
			if i >= len(args) {
				return opts, fmt.Errorf("--tool needs a value")
			}
			opts.Tools = append(opts.Tools, args[i])
		case "--repo":
			i++
			if i >= len(args) {
				return opts, fmt.Errorf("--repo needs a value")
			}
			opts.Repo = args[i]
		default:
			return opts, fmt.Errorf("unknown flag %q", args[i])
		}
	}
	return opts, nil
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

func usage() {
	fmt.Print(`octospec - install the octospec workflow into your tools

Usage:
  octospec install [--tool pi|opencode|copilot|claude] [--repo <path>]
  octospec seed [--repo <path>] [--no-labels]
  octospec update [--repo <path>]
  octospec version
  octospec help

With no --tool, an interactive prompt asks which tools to install for (and
defaults to all when there is no terminal). Global targets (pi, opencode)
write to your home directory; repo-local targets (copilot, claude) write into
--repo (default: the current directory).
`)
}
