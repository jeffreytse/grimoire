package detect

import (
	"os"
	"path/filepath"
	"testing"
)

// touch creates an empty file at the given path.
func touch(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("touch mkdir: %v", err)
	}
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatalf("touch write: %v", err)
	}
}

// ── Engineering project markers ─────────────────────────────────────────────

var engineeringMarkers = []string{
	"package.json",
	"pyproject.toml",
	"Cargo.toml",
	"go.mod",
	"pom.xml",
	"build.gradle",
	"Gemfile",
	"requirements.txt",
}

func TestProfile_EngineeringMarker(t *testing.T) {
	for _, marker := range engineeringMarkers {
		t.Run(marker, func(t *testing.T) {
			dir := t.TempDir()
			touch(t, filepath.Join(dir, marker))
			got := Profile(dir)
			if got != "engineering" {
				t.Errorf("Profile(%q) = %q; want engineering", marker, got)
			}
		})
	}
}

// ── Source file globs ────────────────────────────────────────────────────────

func TestProfile_GoSourceFile_ReturnsEngineering(t *testing.T) {
	dir := t.TempDir()
	touch(t, filepath.Join(dir, "main.go"))
	if got := Profile(dir); got != "engineering" {
		t.Errorf("got %q; want engineering", got)
	}
}

func TestProfile_PythonSourceFile_ReturnsEngineering(t *testing.T) {
	dir := t.TempDir()
	touch(t, filepath.Join(dir, "app.py"))
	if got := Profile(dir); got != "engineering" {
		t.Errorf("got %q; want engineering", got)
	}
}

func TestProfile_TypeScriptSourceFile_ReturnsEngineering(t *testing.T) {
	dir := t.TempDir()
	touch(t, filepath.Join(dir, "index.ts"))
	if got := Profile(dir); got != "engineering" {
		t.Errorf("got %q; want engineering", got)
	}
}

func TestProfile_JavaScriptSourceFile_ReturnsEngineering(t *testing.T) {
	dir := t.TempDir()
	touch(t, filepath.Join(dir, "index.js"))
	if got := Profile(dir); got != "engineering" {
		t.Errorf("got %q; want engineering", got)
	}
}

func TestProfile_RustSourceFile_ReturnsEngineering(t *testing.T) {
	dir := t.TempDir()
	touch(t, filepath.Join(dir, "lib.rs"))
	if got := Profile(dir); got != "engineering" {
		t.Errorf("got %q; want engineering", got)
	}
}

func TestProfile_JavaSourceFile_ReturnsEngineering(t *testing.T) {
	dir := t.TempDir()
	touch(t, filepath.Join(dir, "Main.java"))
	if got := Profile(dir); got != "engineering" {
		t.Errorf("got %q; want engineering", got)
	}
}

func TestProfile_RubySourceFile_ReturnsEngineering(t *testing.T) {
	dir := t.TempDir()
	touch(t, filepath.Join(dir, "app.rb"))
	if got := Profile(dir); got != "engineering" {
		t.Errorf("got %q; want engineering", got)
	}
}

// ── Writing profile ──────────────────────────────────────────────────────────

func TestProfile_MoreThan2Markdown_ReturnsWriting(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"a.md", "b.md", "c.md"} {
		touch(t, filepath.Join(dir, name))
	}
	if got := Profile(dir); got != "writing" {
		t.Errorf("got %q; want writing", got)
	}
}

func TestProfile_ExactlyTwoMarkdown_DoesNotReturnWriting(t *testing.T) {
	dir := t.TempDir()
	touch(t, filepath.Join(dir, "a.md"))
	touch(t, filepath.Join(dir, "b.md"))
	got := Profile(dir)
	if got == "writing" {
		t.Error("2 markdown files should NOT trigger writing profile (requires >2)")
	}
}

// ── Empty / unknown ───────────────────────────────────────────────────────────

func TestProfile_EmptyDir_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	if got := Profile(dir); got != "" {
		t.Errorf("empty dir: got %q; want empty string", got)
	}
}

func TestProfile_EmptyString_UsesCurrentDir(t *testing.T) {
	// Just verify it doesn't panic with empty dir string.
	got := Profile("")
	_ = got // whatever the CWD contains, no crash
}

// ── Priority: manifest beats glob ────────────────────────────────────────────

func TestProfile_ManifestTakesPriorityOverSourceGlob(t *testing.T) {
	dir := t.TempDir()
	// Both a manifest and source files present — manifest wins first
	touch(t, filepath.Join(dir, "go.mod"))
	touch(t, filepath.Join(dir, "main.go"))
	if got := Profile(dir); got != "engineering" {
		t.Errorf("got %q; want engineering", got)
	}
}

// ── Writing profile does not override source files ────────────────────────────

func TestProfile_MarkdownPlusGoFile_ReturnsEngineering(t *testing.T) {
	// Source file markers must be checked before the markdown glob.
	dir := t.TempDir()
	for _, name := range []string{"a.md", "b.md", "c.md"} {
		touch(t, filepath.Join(dir, name))
	}
	touch(t, filepath.Join(dir, "main.go"))
	if got := Profile(dir); got != "engineering" {
		t.Errorf("got %q; want engineering (source file should dominate markdown glob)", got)
	}
}

// ── fileExists helper ─────────────────────────────────────────────────────────

func TestFileExists_ExistingFile_ReturnsTrue(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "exists.txt")
	touch(t, p)
	if !fileExists(p) {
		t.Error("fileExists returned false for an existing file")
	}
}

func TestFileExists_MissingFile_ReturnsFalse(t *testing.T) {
	p := filepath.Join(t.TempDir(), "nope.txt")
	if fileExists(p) {
		t.Error("fileExists returned true for a non-existent file")
	}
}
