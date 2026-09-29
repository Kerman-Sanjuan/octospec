package seed

// InstallRepo writes the embedded seed files (issue forms, CI workflow,
// openspec/config.yaml, check-gates.sh) into repo.
func InstallRepo(repo string) (int, error) {
	n := 0
	if err := copyTree(repoFS, "repo", repo, &n); err != nil {
		return n, err
	}
	return n, nil
}
