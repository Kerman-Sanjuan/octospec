package seed

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestProposalTemplateCarriesChangelog(t *testing.T) {
	b, err := fs.ReadFile(schemaFS, "schema/templates/proposal.md")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "## Changelog") {
		t.Fatal("proposal template is missing the ## Changelog section")
	}
}

func TestSchemaDocumentsChangelog(t *testing.T) {
	b, err := fs.ReadFile(schemaFS, "schema/schema.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "**Changelog**") {
		t.Fatal("schema is missing the Changelog instruction")
	}
}

// The gate script exists once for this repository and once in the repo seed
// that consumer repositories receive. The two must stay byte-identical.
func TestSeededGatesMatchRoot(t *testing.T) {
	seeded, err := fs.ReadFile(repoFS, "repo/scripts/check-gates.sh")
	if err != nil {
		t.Fatal(err)
	}
	root, err := os.ReadFile(filepath.Join("..", "..", "..", "scripts", "check-gates.sh"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(seeded, root) {
		t.Fatal("scripts/check-gates.sh and the seeded copy have drifted")
	}
}

// The seeded gates workflow must authenticate gh like the repository's own
// copy. A seeded repo without a token makes `gh issue view` fail, which used to
// false-fail G1. This test keeps the token, the issues permission, and the
// setup-node major in sync.
func TestSeededWorkflowAuthenticatesGh(t *testing.T) {
	seeded, err := fs.ReadFile(repoFS, "repo/.github/workflows/openspec.yml")
	if err != nil {
		t.Fatal(err)
	}
	root, err := os.ReadFile(filepath.Join("..", "..", "..", ".github", "workflows", "openspec.yml"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(seeded)

	for _, want := range []string{
		"GH_TOKEN: ${{ github.token }}",
		"issues: read",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("seeded workflow is missing %q", want)
		}
	}

	if got, want := setupNodeMajor(s), setupNodeMajor(string(root)); got != want {
		t.Errorf("seeded workflow uses setup-node@%s, root uses @%s", got, want)
	}
}

// setupNodeMajor returns the major version in the first `setup-node@vN` token.
func setupNodeMajor(yaml string) string {
	re := regexp.MustCompile(`actions/setup-node@v(\d+)`)
	m := re.FindStringSubmatch(yaml)
	if m == nil {
		return ""
	}
	return m[1]
}
