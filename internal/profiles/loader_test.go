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
	// Use a name that cannot exist in any installed package.
	p, err := Resolve("zzz-nonexistent-profile-grimoire-test", dir)
	if err != nil {
		t.Fatalf("missing file should not error: %v", err)
	}
	if p.Source != "" {
		t.Errorf("expected empty Source for missing profile, got %q", p.Source)
	}
	if p.Name != "zzz-nonexistent-profile-grimoire-test" {
		t.Errorf("Name = %q, want %q", p.Name, "zzz-nonexistent-profile-grimoire-test")
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

func TestResolve_NoFallbackToDefault(t *testing.T) {
	dir := t.TempDir()
	writeProfileFile(t, dir, ".grimoire/profiles/default.toml", `
name = "default"
[[skills]]
name = "default-skill"
`)

	// Profile resolution no longer falls back to default.toml.
	// An unmatched name returns an empty Profile.
	p, err := Resolve("zzz-nonexistent-for-fallback-test", dir)
	if err != nil {
		t.Fatal(err)
	}
	if p.Source != "" {
		t.Errorf("expected empty Source (no fallback), got %q", p.Source)
	}
}

func TestResolveAll_MultipleProfiles(t *testing.T) {
	dir := t.TempDir()
	// Use unlikely names to avoid matching real package profiles.
	writeProfileFile(t, dir, ".grimoire/profiles/zzz-grimoire-test-solid.toml", `
[[skills]]
name = "apply-solid"
`)
	writeProfileFile(t, dir, ".grimoire/profiles/zzz-grimoire-test-tdd.toml", `
[[skills]]
name = "apply-tdd"
`)

	got, err := ResolveAll([]string{"zzz-grimoire-test-solid", "zzz-grimoire-test-tdd"}, dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(got))
	}
	if got[0].Skills[0].Name != "apply-solid" {
		t.Errorf("got[0] skill = %q", got[0].Skills[0].Name)
	}
	if got[1].Skills[0].Name != "apply-tdd" {
		t.Errorf("got[1] skill = %q", got[1].Skills[0].Name)
	}
}

func TestProfilesMatchingGlob_DepthAnywhere(t *testing.T) {
	root := t.TempDir()
	write := func(rel string) {
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(`name = "test"`), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// Nested layout — no top-level profiles/ directory
	write("engineering/tdd.toml")
	write("engineering/bdd.toml")
	write("security/owasp.toml")

	t.Run("name matches at any depth", func(t *testing.T) {
		got, err := ProfilesMatchingGlob(root, "tdd")
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 {
			t.Errorf("expected 1 profile, got %d", len(got))
		}
	})

	t.Run("dir name returns all under dir", func(t *testing.T) {
		got, err := ProfilesMatchingGlob(root, "engineering")
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 2 {
			t.Errorf("expected 2 profiles under engineering/, got %d", len(got))
		}
	})

	t.Run("empty glob returns all", func(t *testing.T) {
		got, err := ProfilesMatchingGlob(root, "")
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 3 {
			t.Errorf("expected 3 profiles, got %d", len(got))
		}
	})
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

	sources := []skills.SkillsPackage{{Name: "test", Root: root}}
	refs := ResolveByTags("oop", sources)
	if len(refs) != 2 {
		t.Fatalf("expected 2 oop-tagged skills, got %d: %v", len(refs), refs)
	}
}

func TestResolveWithOptions_TagFallback(t *testing.T) {
	root := t.TempDir()
	writeTaggedSkill(t, root, "eng", "apply-solid", []string{"oop"})

	sources := []skills.SkillsPackage{{Name: "test", Root: root}}
	p, err := ResolveWithOptions("oop", t.TempDir(), ResolveOptions{Packages: sources})
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
	sources := []skills.SkillsPackage{{Name: "test", Root: t.TempDir()}}
	p, err := ResolveWithOptions("oop", projectDir, ResolveOptions{Packages: sources})
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

// --- composition tests ---

func TestResolveSkills_Extends_InheritsParentSkills(t *testing.T) {
	dir := t.TempDir()
	writeProfileFile(t, dir, ".grimoire/profiles/oop.toml", `
[[skills]]
name = "apply-solid"
[[skills]]
name = "apply-lod"
`)
	writeProfileFile(t, dir, ".grimoire/profiles/my-team.toml", `
extends = ["oop"]
[[skills]]
name = "apply-internal"
`)

	p, err := Resolve("my-team", dir)
	if err != nil {
		t.Fatal(err)
	}
	resolved := ResolveSkills(&p, dir, nil, nil)
	names := skillNames(resolved)
	if !contains(names, "apply-solid") || !contains(names, "apply-lod") {
		t.Errorf("missing inherited skills: %v", names)
	}
	if !contains(names, "apply-internal") {
		t.Errorf("missing explicit skill: %v", names)
	}
}

func TestResolveSkills_Tags_BulkActivates(t *testing.T) {
	dir := t.TempDir()
	root := t.TempDir()
	writeTaggedSkill(t, root, "eng", "apply-solid", []string{"oop"})
	writeTaggedSkill(t, root, "eng", "apply-lod", []string{"oop"})
	writeTaggedSkill(t, root, "eng", "apply-tdd", []string{"tdd"})

	writeProfileFile(t, dir, ".grimoire/profiles/oop-team.toml", `
tags = ["oop"]
`)
	p, err := Resolve("oop-team", dir)
	if err != nil {
		t.Fatal(err)
	}
	sources := []skills.SkillsPackage{{Name: "test", Root: root}}
	resolved := ResolveSkills(&p, dir, sources, nil)
	names := skillNames(resolved)
	if len(names) != 2 {
		t.Errorf("expected 2 oop-tagged skills, got %v", names)
	}
	if contains(names, "apply-tdd") {
		t.Error("apply-tdd should not be included (tagged tdd not oop)")
	}
}

func TestResolveSkills_Exclude_RemovesSkills(t *testing.T) {
	dir := t.TempDir()
	writeProfileFile(t, dir, ".grimoire/profiles/oop.toml", `
[[skills]]
name = "apply-solid"
[[skills]]
name = "apply-lod"
`)
	writeProfileFile(t, dir, ".grimoire/profiles/my-team.toml", `
extends = ["oop"]
exclude = ["apply-lod"]
`)
	p, err := Resolve("my-team", dir)
	if err != nil {
		t.Fatal(err)
	}
	resolved := ResolveSkills(&p, dir, nil, nil)
	names := skillNames(resolved)
	if contains(names, "apply-lod") {
		t.Error("apply-lod should be excluded")
	}
	if !contains(names, "apply-solid") {
		t.Error("apply-solid should still be present")
	}
}

func TestResolveSkills_Priority_LowerWins(t *testing.T) {
	dir := t.TempDir()
	writeProfileFile(t, dir, ".grimoire/profiles/base.toml", `
[[skills]]
name = "apply-solid"
priority = 30
[[skills]]
name = "apply-lod"
priority = 10
`)
	p, err := Resolve("base", dir)
	if err != nil {
		t.Fatal(err)
	}
	resolved := ResolveSkills(&p, dir, nil, nil)
	if len(resolved) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(resolved))
	}
	// priority 10 (apply-lod) should come first
	if resolved[0].Name != "apply-lod" {
		t.Errorf("expected apply-lod first (priority 10), got %s", resolved[0].Name)
	}
	if resolved[1].Name != "apply-solid" {
		t.Errorf("expected apply-solid second (priority 30), got %s", resolved[1].Name)
	}
}

func TestResolveSkills_CycleDetection(t *testing.T) {
	dir := t.TempDir()
	// A extends B, B extends A — should not infinite loop
	writeProfileFile(t, dir, ".grimoire/profiles/a.toml", `
extends = ["b"]
[[skills]]
name = "skill-a"
`)
	writeProfileFile(t, dir, ".grimoire/profiles/b.toml", `
extends = ["a"]
[[skills]]
name = "skill-b"
`)
	p, err := Resolve("a", dir)
	if err != nil {
		t.Fatal(err)
	}
	// Should complete without hanging; exact output not critical — just no panic/loop
	visited := map[string]bool{"a": true}
	resolved := ResolveSkills(&p, dir, nil, visited)
	names := skillNames(resolved)
	if !contains(names, "skill-a") {
		t.Errorf("skill-a missing from %v", names)
	}
}

func TestResolveSkills_ExplicitOverridesPriority(t *testing.T) {
	dir := t.TempDir()
	// Parent has apply-solid at default priority; child explicitly sets priority=1
	writeProfileFile(t, dir, ".grimoire/profiles/parent.toml", `
[[skills]]
name = "apply-solid"
[[skills]]
name = "apply-lod"
`)
	writeProfileFile(t, dir, ".grimoire/profiles/child.toml", `
extends = ["parent"]
[[skills]]
name = "apply-solid"
priority = 1
`)
	p, err := Resolve("child", dir)
	if err != nil {
		t.Fatal(err)
	}
	resolved := ResolveSkills(&p, dir, nil, nil)
	// apply-solid should appear with priority=1 and come first
	if len(resolved) == 0 || resolved[0].Name != "apply-solid" {
		t.Errorf("expected apply-solid first with priority=1, got %v", resolved)
	}
	if resolved[0].Priority != 1 {
		t.Errorf("expected priority=1, got %d", resolved[0].Priority)
	}
}

func TestResolveSkills_BackwardCompat_NoNewFields(t *testing.T) {
	dir := t.TempDir()
	writeProfileFile(t, dir, ".grimoire/profiles/oop.toml", `
name = "oop"
[[skills]]
name = "apply-solid-principles"
[[skills]]
name = "apply-law-of-demeter"
`)
	p, err := Resolve("oop", dir)
	if err != nil {
		t.Fatal(err)
	}
	// No extends/tags/exclude — should behave identically to pre-redesign
	resolved := ResolveSkills(&p, dir, nil, nil)
	if len(resolved) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(resolved))
	}
	if resolved[0].Name != "apply-solid-principles" {
		t.Errorf("order changed: %v", skillNames(resolved))
	}
}

func TestParseProfileRef(t *testing.T) {
	known := []string{"official", "acmecorp/standards", "gitlab.com/acmecorp/standards"}

	cases := []struct {
		ref      string
		wantReg  string
		wantName string
	}{
		{"engineering", "", "engineering"},
		{"official/engineering", "official", "engineering"},
		{"acmecorp/standards/engineering", "acmecorp/standards", "engineering"},
		{"gitlab.com/acmecorp/standards/engineering", "gitlab.com/acmecorp/standards", "engineering"},
		// longest match wins: gitlab.com/acmecorp/standards beats acmecorp/standards
		{"gitlab.com/acmecorp/standards/go-service", "gitlab.com/acmecorp/standards", "go-service"},
		// unknown package → unqualified
		{"unknown/org/profile", "", "unknown/org/profile"},
	}

	for _, c := range cases {
		reg, name := ParseProfileRef(c.ref, known)
		if reg != c.wantReg || name != c.wantName {
			t.Errorf("ParseProfileRef(%q) = (%q, %q), want (%q, %q)",
				c.ref, reg, name, c.wantReg, c.wantName)
		}
	}
}

// helpers

func skillNames(refs []SkillRef) []string {
	names := make([]string, len(refs))
	for i, r := range refs {
		names[i] = r.Name
	}
	return names
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
