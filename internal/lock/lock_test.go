package lock

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFile_Absent(t *testing.T) {
	lf, err := ParseFile("/nonexistent/grimoire.lock")
	if err != nil {
		t.Fatalf("absent file should not error: %v", err)
	}
	if len(lf.Skills) != 0 {
		t.Errorf("expected empty skills, got %d", len(lf.Skills))
	}
}

func TestParseFile_Basic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "grimoire.lock")
	content := `[[skills]]
name = "apply-solid-principles"
version = "1.2.3"
source = "jeffreytse/grimoire-core"
resolved = "https://github.com/jeffreytse/grimoire-core.git"
commit = "abc123"
checksum = "sha256:deadbeef"

[[skills]]
name = "apply-dry-principle"
version = "1.0.0"
source = "jeffreytse/grimoire-core"
resolved = "https://github.com/jeffreytse/grimoire-core.git"
commit = "abc123"
checksum = "sha256:cafebabe"
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	lf, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(lf.Skills) != 2 {
		t.Fatalf("want 2 skills, got %d", len(lf.Skills))
	}
	if lf.Skills[0].Name != "apply-solid-principles" {
		t.Errorf("name = %q", lf.Skills[0].Name)
	}
	if lf.Skills[0].Version != "1.2.3" {
		t.Errorf("version = %q", lf.Skills[0].Version)
	}
	if lf.Skills[0].Commit != "abc123" {
		t.Errorf("commit = %q", lf.Skills[0].Commit)
	}
}

func TestWriteFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "grimoire.lock")

	lf := LockFile{Skills: []Entry{
		{Name: "apply-solid", Version: "1.0.0", Source: "jeffreytse/grimoire-core", Resolved: "https://github.com/jeffreytse/grimoire-core.git", Commit: "deadbeef", Checksum: "sha256:abc"},
	}}

	if err := WriteFile(path, lf); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	got, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(got.Skills) != 1 || got.Skills[0].Name != "apply-solid" {
		t.Errorf("round-trip failed: %+v", got.Skills)
	}
}

func TestUpsert(t *testing.T) {
	lf := LockFile{}
	lf.Upsert(&Entry{Name: "apply-solid", Version: "1.0.0"})
	lf.Upsert(&Entry{Name: "apply-dry", Version: "2.0.0"})
	lf.Upsert(&Entry{Name: "apply-solid", Version: "1.1.0"}) // update

	if len(lf.Skills) != 2 {
		t.Fatalf("want 2, got %d", len(lf.Skills))
	}
	e := lf.Find("apply-solid")
	if e == nil || e.Version != "1.1.0" {
		t.Errorf("upsert did not update: %+v", e)
	}
}

func TestRemove(t *testing.T) {
	lf := LockFile{Skills: []Entry{
		{Name: "a"}, {Name: "b"}, {Name: "c"},
	}}
	lf.Remove("b")
	if len(lf.Skills) != 2 {
		t.Fatalf("want 2, got %d", len(lf.Skills))
	}
	if lf.Find("b") != nil {
		t.Error("b should be removed")
	}
}

func TestFind_Missing(t *testing.T) {
	lf := LockFile{}
	if lf.Find("nope") != nil {
		t.Error("Find on empty lock should return nil")
	}
}
