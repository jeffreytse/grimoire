package skills

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// ── ListAllSkillsFromPackages ─────────────────────────────────────────────────

func TestListAllSkillsFromPackages_NoConflict(t *testing.T) {
	a := t.TempDir()
	b := t.TempDir()
	buildNestedDomain(t, a, "engineering", "development", "apply-solid")
	buildNestedDomain(t, b, "design", "patterns", "observer")

	srcs := []SkillsPackage{{Name: "official", Root: a}, {Name: "myteam", Root: b}}
	skills, conflicts, err := ListAllSkillsFromPackages(srcs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conflicts) != 0 {
		t.Fatalf("expected 0 conflicts, got %d: %+v", len(conflicts), conflicts)
	}
	if len(skills) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(skills))
	}
}

func TestListAllSkillsFromPackages_ConflictNestedPath(t *testing.T) {
	a := t.TempDir()
	b := t.TempDir()
	buildNestedDomain(t, a, "engineering", "development", "apply-solid")
	buildNestedDomain(t, b, "engineering", "development", "apply-solid")

	srcs := []SkillsPackage{{Name: "official", Root: a}, {Name: "myteam", Root: b}}
	skills, conflicts, err := ListAllSkillsFromPackages(srcs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d: %+v", len(conflicts), conflicts)
	}
	if conflicts[0].WinnerPackage != "official" {
		t.Errorf("winner = %q, want official", conflicts[0].WinnerPackage)
	}
	if conflicts[0].LoserPackage != "myteam" {
		t.Errorf("loser = %q, want myteam", conflicts[0].LoserPackage)
	}
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}
	if skills[0].Package != "official" {
		t.Errorf("installed package = %q, want official", skills[0].Package)
	}
}

func TestListAllSkillsFromPackages_SameLeafDifferentDomain(t *testing.T) {
	a := t.TempDir()
	b := t.TempDir()
	buildFlatDomain(t, a, "engineering", "apply-solid")
	buildFlatDomain(t, b, "design", "apply-solid")

	srcs := []SkillsPackage{{Name: "official", Root: a}, {Name: "myteam", Root: b}}
	skills, conflicts, err := ListAllSkillsFromPackages(srcs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conflicts) != 0 {
		t.Fatalf("expected 0 conflicts (different domains), got %d: %+v", len(conflicts), conflicts)
	}
	if len(skills) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(skills))
	}
}

func TestListAllSkillsFromPackages_ThreeWayConflict(t *testing.T) {
	a, b, c := t.TempDir(), t.TempDir(), t.TempDir()
	buildNestedDomain(t, a, "engineering", "development", "apply-solid")
	buildNestedDomain(t, b, "engineering", "development", "apply-solid")
	buildNestedDomain(t, c, "engineering", "development", "apply-solid")

	srcs := []SkillsPackage{
		{Name: "official", Root: a},
		{Name: "teamA", Root: b},
		{Name: "teamB", Root: c},
	}
	skills, conflicts, err := ListAllSkillsFromPackages(srcs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conflicts) != 2 {
		t.Fatalf("expected 2 conflicts, got %d: %+v", len(conflicts), conflicts)
	}
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill (winner only), got %d", len(skills))
	}
	if skills[0].Package != "official" {
		t.Errorf("winner package = %q, want official", skills[0].Package)
	}
}

func TestListAllSkillsFromPackages_SingleSourceNoConflicts(t *testing.T) {
	a := t.TempDir()
	buildNestedDomain(t, a, "engineering", "development", "apply-solid")
	buildNestedDomain(t, a, "engineering", "development", "write-tests")

	srcs := []SkillsPackage{{Name: "official", Root: a}}
	skills, conflicts, err := ListAllSkillsFromPackages(srcs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conflicts) != 0 {
		t.Fatalf("expected 0 conflicts for single source, got %d", len(conflicts))
	}
	if len(skills) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(skills))
	}
}

func TestListAllSkillsFromPackages_FrontmatterNameDiffers_SameDirName(t *testing.T) {
	// Two packages have apply-solid/ with different frontmatter name: values.
	// Conflict detection must key on the directory name, not the frontmatter name.
	a := t.TempDir()
	b := t.TempDir()

	dirA := filepath.Join(a, "engineering", "development", "skills", "apply-solid")
	if err := os.MkdirAll(dirA, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dirA, "SKILL.md"),
		[]byte("---\nname: Apply SOLID Principles\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	dirB := filepath.Join(b, "engineering", "development", "skills", "apply-solid")
	if err := os.MkdirAll(dirB, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dirB, "SKILL.md"),
		[]byte("---\nname: SOLID Design Patterns\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	srcs := []SkillsPackage{{Name: "official", Root: a}, {Name: "myteam", Root: b}}
	skills, conflicts, err := ListAllSkillsFromPackages(srcs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conflicts) != 1 {
		t.Fatalf("expected 1 conflict (same dir name, different frontmatter), got %d: %+v", len(conflicts), conflicts)
	}
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill (winner only), got %d", len(skills))
	}
	if skills[0].Package != "official" {
		t.Errorf("winner = %q, want official", skills[0].Package)
	}
}

// Libraries come from the system file (an organization's), the user's, and
// the project's [dependencies]; a project's own [[package]] is ignored.
func TestAllPackages_Layers(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(root, "xdg"))
	t.Setenv("HOME", filepath.Join(root, "home"))
	sys := filepath.Join(root, "etc", "grimoire.toml")
	t.Setenv("GRIMOIRE_SYSTEM_CONFIG", sys)
	proj := filepath.Join(root, "proj")
	t.Setenv("GRIMOIRE_PROJECT_DIR", proj)
	write := func(path, body string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	pkg := func(name string) string {
		return "[[package]]\nname = \"" + name + "\"\nurl = \"https://example.com/" + name + ".git\"\n"
	}
	names := func() []string {
		var out []string
		for _, p := range AllPackages() {
			out = append(out, p.Name)
		}
		return out
	}

	if got := names(); len(got) != 1 || got[0] != OfficialPackageDerivedName() {
		t.Fatalf("fresh: %v", got)
	}

	// A project's [[package]] is ignored; its [dependencies] bring a library.
	write(filepath.Join(proj, "grimoire.toml"), pkg("sneaky-skills")+"\n[dependencies]\n\"acme/team-skills\" = \"*\"\n")
	got := names()
	if slices.Contains(got, "sneaky-skills") {
		t.Errorf("project [[package]] honored: %v", got)
	}
	if !slices.ContainsFunc(got, func(n string) bool { return strings.Contains(n, "acme/team-skills") }) || !slices.Contains(got, OfficialPackageDerivedName()) {
		t.Errorf("project dependency library: %v", got)
	}

	// An organization's library from the system file, a user's from theirs.
	write(sys, pkg("org-skills"))
	write(filepath.Join(root, "xdg", "grimoire", "grimoire.toml"), pkg("my-skills"))
	got = names()
	for _, want := range []string{"org-skills", "my-skills"} {
		if !slices.Contains(got, want) {
			t.Errorf("layers: %v lacks %s", got, want)
		}
	}
	if slices.Contains(got, OfficialPackageDerivedName()) {
		t.Errorf("official included though user/system list packages: %v", got)
	}

	// The user wins per name: they can disable a library the organization lists.
	write(filepath.Join(root, "xdg", "grimoire", "grimoire.toml"), pkg("my-skills")+pkg("org-skills")+"enabled = false\n")
	if got := names(); slices.Contains(got, "org-skills") {
		t.Errorf("user could not disable org-skills: %v", got)
	}
}
