package seed

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
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
