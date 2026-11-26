// Package seed installs the parts of octospec that live outside the commands:
// the OpenSpec schema, the repo seed files, and the workflow labels.
package seed

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed all:schema
var schemaFS embed.FS

//go:embed all:repo
var repoFS embed.FS

// Options controls a seed run.
type Options struct {
	Repo       string // target repository root
	SkipLabels bool
	DryRun     bool // print the plan and write nothing
}

// Result counts what a seed run wrote.
type Result struct {
	Schema int
	Repo   int
	Labels int
}

// Run installs the schema, the repo seed files, and (unless skipped) labels.
func Run(opts Options) (Result, error) {
	var r Result
	var err error
	if opts.DryRun {
		r.Schema = countTree(schemaFS, "schema")
		r.Repo = countTree(repoFS, "repo")
		fmt.Printf("dry run: would install the schema -> %s (%d files)\n", SchemaDir(), r.Schema)
		fmt.Printf("dry run: would write the repo seed -> %s (%d files)\n", opts.Repo, r.Repo)
		if !opts.SkipLabels {
			fmt.Println("dry run: would provision the workflow labels")
		}
		return r, nil
	}
	if r.Schema, err = InstallSchema(); err != nil {
		return r, err
	}
	if r.Repo, err = InstallRepo(opts.Repo); err != nil {
		return r, err
	}
	if err := ensureGitignore(opts.Repo); err != nil {
		return r, err
	}
	if opts.SkipLabels {
		return r, nil
	}
	r.Labels, err = InstallLabels(opts.Repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
	}
	return r, nil
}

// countTree returns the number of files under src in an embedded tree.
func countTree(fsys fs.FS, src string) int {
	n := 0
	_ = fs.WalkDir(fsys, src, func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			n++
		}
		return nil
	})
	return n
}

// copyTree writes every file under src (an embedded dir) into dst, preserving
// relative paths.
func copyTree(fsys fs.FS, src, dst string, count *int) error {
	return fs.WalkDir(fsys, src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, p)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		b, err := fs.ReadFile(fsys, p)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(target, b, 0o644); err != nil {
			return err
		}
		*count++
		return nil
	})
}

// gitignoreBlock keeps the runtime state out of the repository. The archive lock is written by
// /archive inside the tree, so an un-ignored copy would be committed by the next `git add -A`.
const gitignoreBlock = "# octospec runtime state\n.octospec/locks/\n"

// ensureGitignore appends the managed block to .gitignore when it is missing. It never rewrites
// or reorders what the repository already has.
func ensureGitignore(repo string) error {
	path := filepath.Join(repo, ".gitignore")
	b, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if strings.Contains(string(b), ".octospec/locks/") {
		return nil
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	prefix := ""
	if len(b) > 0 && !strings.HasSuffix(string(b), "\n") {
		prefix = "\n"
	}
	if len(b) > 0 {
		prefix += "\n"
	}
	_, err = f.WriteString(prefix + gitignoreBlock)
	return err
}
