package skills

import (
	"os"
	"path/filepath"
	"testing"
)

// ── ListAllSkillsFromRegistries ─────────────────────────────────────────────────

func TestListAllSkillsFromRegistries_NoConflict(t *testing.T) {
	a := t.TempDir()
	b := t.TempDir()
	buildNestedDomain(t, a, "engineering", "development", "apply-solid")
	buildNestedDomain(t, b, "design", "patterns", "observer")

	srcs := []SkillsRegistry{{Name: "official", Root: a}, {Name: "myteam", Root: b}}
	skills, conflicts, err := ListAllSkillsFromRegistries(srcs)
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

func TestListAllSkillsFromRegistries_ConflictNestedPath(t *testing.T) {
	a := t.TempDir()
	b := t.TempDir()
	buildNestedDomain(t, a, "engineering", "development", "apply-solid")
	buildNestedDomain(t, b, "engineering", "development", "apply-solid")

	srcs := []SkillsRegistry{{Name: "official", Root: a}, {Name: "myteam", Root: b}}
	skills, conflicts, err := ListAllSkillsFromRegistries(srcs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d: %+v", len(conflicts), conflicts)
	}
	if conflicts[0].WinnerRegistry != "official" {
		t.Errorf("winner = %q, want official", conflicts[0].WinnerRegistry)
	}
	if conflicts[0].LoserRegistry != "myteam" {
		t.Errorf("loser = %q, want myteam", conflicts[0].LoserRegistry)
	}
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}
	if skills[0].Registry != "official" {
		t.Errorf("installed registry = %q, want official", skills[0].Registry)
	}
}

func TestListAllSkillsFromRegistries_SameLeafDifferentDomain(t *testing.T) {
	a := t.TempDir()
	b := t.TempDir()
	buildFlatDomain(t, a, "engineering", "apply-solid")
	buildFlatDomain(t, b, "design", "apply-solid")

	srcs := []SkillsRegistry{{Name: "official", Root: a}, {Name: "myteam", Root: b}}
	skills, conflicts, err := ListAllSkillsFromRegistries(srcs)
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

func TestListAllSkillsFromRegistries_ThreeWayConflict(t *testing.T) {
	a, b, c := t.TempDir(), t.TempDir(), t.TempDir()
	buildNestedDomain(t, a, "engineering", "development", "apply-solid")
	buildNestedDomain(t, b, "engineering", "development", "apply-solid")
	buildNestedDomain(t, c, "engineering", "development", "apply-solid")

	srcs := []SkillsRegistry{
		{Name: "official", Root: a},
		{Name: "teamA", Root: b},
		{Name: "teamB", Root: c},
	}
	skills, conflicts, err := ListAllSkillsFromRegistries(srcs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conflicts) != 2 {
		t.Fatalf("expected 2 conflicts, got %d: %+v", len(conflicts), conflicts)
	}
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill (winner only), got %d", len(skills))
	}
	if skills[0].Registry != "official" {
		t.Errorf("winner registry = %q, want official", skills[0].Registry)
	}
}

func TestListAllSkillsFromRegistries_SingleSourceNoConflicts(t *testing.T) {
	a := t.TempDir()
	buildNestedDomain(t, a, "engineering", "development", "apply-solid")
	buildNestedDomain(t, a, "engineering", "development", "write-tests")

	srcs := []SkillsRegistry{{Name: "official", Root: a}}
	skills, conflicts, err := ListAllSkillsFromRegistries(srcs)
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

func TestListAllSkillsFromRegistries_FrontmatterNameDiffers_SameDirName(t *testing.T) {
	// Two registries have apply-solid/ with different frontmatter name: values.
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

	srcs := []SkillsRegistry{{Name: "official", Root: a}, {Name: "myteam", Root: b}}
	skills, conflicts, err := ListAllSkillsFromRegistries(srcs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conflicts) != 1 {
		t.Fatalf("expected 1 conflict (same dir name, different frontmatter), got %d: %+v", len(conflicts), conflicts)
	}
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill (winner only), got %d", len(skills))
	}
	if skills[0].Registry != "official" {
		t.Errorf("winner = %q, want official", skills[0].Registry)
	}
}
