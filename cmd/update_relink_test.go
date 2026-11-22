package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// update must relink only skills dirs that already hold grimoire links; an
// agent skills dir created by another tool (or in MCP mode) stays untouched.
func TestHasLinksInto(t *testing.T) {
	pkg := t.TempDir()
	skill := filepath.Join(pkg, "skills", "meta", "skills", "s")
	if err := os.MkdirAll(skill, 0o755); err != nil {
		t.Fatal(err)
	}
	other := t.TempDir()

	empty := t.TempDir()
	foreign := t.TempDir()
	if err := os.Symlink(other, filepath.Join(foreign, "context7-mcp")); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(foreign, "plain-dir"), 0o755); err != nil {
		t.Fatal(err)
	}
	linked := t.TempDir()
	if err := os.Symlink(skill, filepath.Join(linked, "s")); err != nil {
		t.Fatal(err)
	}

	for name, tc := range map[string]struct {
		dir  string
		want bool
	}{
		"missing":       {filepath.Join(empty, "nope"), false},
		"empty":         {empty, false},
		"foreign links": {foreign, false},
		"grimoire link": {linked, true},
	} {
		if got := hasLinksInto(tc.dir, pkg); got != tc.want {
			t.Errorf("%s: hasLinksInto = %v; want %v", name, got, tc.want)
		}
	}
}
