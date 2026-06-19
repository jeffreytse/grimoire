package profiles

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jeffreytse/grimoire/internal/skills"
)

func writeProfileFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestResolve_MissingFile(t *testing.T) {
	dir := t.TempDir()
	p, err := Resolve("engineering", dir)
	if err != nil {
		t.Fatalf("missing file should not error: %v", err)
	}
	if p.Source != "" {
		t.Errorf("expected empty Source for missing profile, got %q", p.Source)
	}
	if p.Name != "engineering" {
		t.Errorf("Name = %q, want %q", p.Name, "engineering")
	}
	if len(p.Skills) != 0 {
		t.Errorf("Skills should be empty for missing profile, got %v", p.Skills)
	}
}

func TestResolve_ParsesSkills(t *testing.T) {
	dir := t.TempDir()
	writeProfileFile(t, dir, ".grimoire/profiles/engineering.toml", `
name = "engineering"
description = "Core engineering practices"

[[skills]]
name = "apply-solid-principles"

[[skills]]
name = "apply-kiss-principle"
`)

	p, err := Resolve("engineering", dir)
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "engineering" {
		t.Errorf("Name = %q", p.Name)
	}
	if p.Description != "Core engineering practices" {
		t.Errorf("Description = %q", p.Description)
	}
	if len(p.Skills) != 2 {
		t.Fatalf("Skills len = %d, want 2", len(p.Skills))
	}
	if p.Skills[0].Name != "apply-solid-principles" {
		t.Errorf("Skills[0] = %q", p.Skills[0].Name)
	}
	if p.Skills[1].Name != "apply-kiss-principle" {
		t.Errorf("Skills[1] = %q", p.Skills[1].Name)
	}
	if p.Source == "" {
		t.Error("Source should be set when file is found")
	}
}

func TestResolve_ProjectLevelBeforeUserLevel(t *testing.T) {
	projectDir := t.TempDir()

	// project-level file
	writeProfileFile(t, projectDir, ".grimoire/profiles/oop.toml", `
name = "oop"
[[skills]]
name = "project-level-skill"
`)

	p, err := Resolve("oop", projectDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Skills) != 1 || p.Skills[0].Name != "project-level-skill" {
		t.Errorf("expected project-level skill, got %v", p.Skills)
	}
	if p.Source == "" {
		t.Error("Source should be set")
	}
}

func TestResolve_FallbackToDefault(t *testing.T) {
	dir := t.TempDir()
	writeProfileFile(t, dir, ".grimoire/profiles/default.toml", `
name = "default"
[[skills]]
name = "default-skill"
`)

	// "unknown" profile has no file, should fall back to default.toml
	p, err := Resolve("unknown", dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Skills) != 1 || p.Skills[0].Name != "default-skill" {
		t.Errorf("expected default-skill, got %v", p.Skills)
	}
}

func TestResolveAll_MultipleProfiles(t *testing.T) {
	dir := t.TempDir()
	writeProfileFile(t, dir, ".grimoire/profiles/engineering.toml", `
[[skills]]
name = "apply-solid"
`)
	writeProfileFile(t, dir, ".grimoire/profiles/tdd.toml", `
[[skills]]
name = "apply-tdd"
`)

	profiles, err := ResolveAll([]string{"engineering", "tdd"}, dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(profiles) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(profiles))
	}
	if profiles[0].Skills[0].Name != "apply-solid" {
		t.Errorf("profiles[0] skill = %q", profiles[0].Skills[0].Name)
	}
	if profiles[1].Skills[0].Name != "apply-tdd" {
		t.Errorf("profiles[1] skill = %q", profiles[1].Skills[0].Name)
	}
}

func TestSearchPaths_Order(t *testing.T) {
	paths := SearchPaths("engineering", "/project")
	if len(paths) != 4 {
		t.Fatalf("expected 4 search paths, got %d", len(paths))
	}
	// first path must be project-level
	if filepath.Base(filepath.Dir(paths[0])) != "profiles" {
		t.Errorf("unexpected first path: %s", paths[0])
	}
	// third path must be project default
	if filepath.Base(paths[2]) != "default.toml" {
		t.Errorf("expected default.toml at index 2, got %s", paths[2])
	}
}

func writeTaggedSkill(t *testing.T, root, domain, name string, tags []string) {
	t.Helper()
	skillDir := filepath.Join(root, domain, "skills", name)
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	fm := "---\nname: " + name + "\ntags: ["
	for i, tag := range tags {
		if i > 0 {
			fm += ", "
		}
		fm += tag
	}
	fm += "]\n---\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(fm), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestResolveByTags_MatchesTaggedSkills(t *testing.T) {
	root := t.TempDir()
	writeTaggedSkill(t, root, "engineering", "apply-solid", []string{"oop", "engineering"})
	writeTaggedSkill(t, root, "engineering", "apply-tdd", []string{"tdd"})
	writeTaggedSkill(t, root, "engineering", "apply-lod", []string{"oop"})

	sources := []skills.SkillsSource{{Name: "test", Root: root}}
	refs := ResolveByTags("oop", sources)
	if len(refs) != 2 {
		t.Fatalf("expected 2 oop-tagged skills, got %d: %v", len(refs), refs)
	}
}

func TestResolveWithOptions_TagFallback(t *testing.T) {
	root := t.TempDir()
	writeTaggedSkill(t, root, "eng", "apply-solid", []string{"oop"})

	sources := []skills.SkillsSource{{Name: "test", Root: root}}
	p, err := ResolveWithOptions("oop", t.TempDir(), ResolveOptions{Sources: sources})
	if err != nil {
		t.Fatal(err)
	}
	if p.Source != "(tag query)" {
		t.Errorf("Source = %q, want (tag query)", p.Source)
	}
	if len(p.Skills) != 1 || p.Skills[0].Name != "apply-solid" {
		t.Errorf("Skills = %v", p.Skills)
	}
}

func TestResolveWithOptions_FileBeforeTagQuery(t *testing.T) {
	projectDir := t.TempDir()
	writeProfileFile(t, projectDir, ".grimoire/profiles/oop.toml", `
name = "oop"
[[skills]]
name = "file-based-skill"
`)

	// even with sources provided, file takes priority
	sources := []skills.SkillsSource{{Name: "test", Root: t.TempDir()}}
	p, err := ResolveWithOptions("oop", projectDir, ResolveOptions{Sources: sources})
	if err != nil {
		t.Fatal(err)
	}
	if p.Source == "(tag query)" {
		t.Error("file should take priority over tag query")
	}
	if len(p.Skills) != 1 || p.Skills[0].Name != "file-based-skill" {
		t.Errorf("Skills = %v", p.Skills)
	}
}
