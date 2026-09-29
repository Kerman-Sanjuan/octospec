// Package seed installs the parts of octospec that live outside the commands:
// the OpenSpec schema, the repo seed files, and the workflow labels.
package seed

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed all:schema
var schemaFS embed.FS

//go:embed all:repo
var repoFS embed.FS

// Options controls a seed run.
type Options struct {
	Repo       string // target repository root
	SkipLabels bool
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
	if r.Schema, err = InstallSchema(); err != nil {
		return r, err
	}
	if r.Repo, err = InstallRepo(opts.Repo); err != nil {
		return r, err
	}
	if opts.SkipLabels {
		return r, nil
	}
	r.Labels, err = InstallLabels(opts.Repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
		err = nil
	}
	return r, nil
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
