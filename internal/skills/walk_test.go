package skills

import (
	"os"
	"path/filepath"
	"testing"
)

// ── Test helpers ────────────────────────────────────────────────────────────

// buildSkillsRoot creates a minimal skills root with the given layout:
//
//	root/
//	  <domain>/            (nested domain)
//	    <subdomain>/
//	      skills/
//	        <skill>/
//	          SKILL.md
//	  <flatdomain>/        (flat domain — has skills/ directly)
//	    skills/
//	      <skill>/
//	        SKILL.md
func buildNestedDomain(t *testing.T, root, domain, subdomain, skill string) {
	t.Helper()
	dir := filepath.Join(root, domain, subdomain, "skills", skill)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("buildNestedDomain: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# "+skill), 0o644); err != nil {
		t.Fatalf("buildNestedDomain SKILL.md: %v", err)
	}
}

func buildFlatDomain(t *testing.T, root, domain, skill string) {
	t.Helper()
	dir := filepath.Join(root, domain, "skills", skill)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("buildFlatDomain: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# "+skill), 0o644); err != nil {
		t.Fatalf("buildFlatDomain SKILL.md: %v", err)
	}
}

// ── IsNested ────────────────────────────────────────────────────────────────

func TestIsNested_TrueWhenNoSkillsDir(t *testing.T) {
	dir := t.TempDir()
	if !IsNested(dir) {
		t.Error("expected IsNested=true when skills/ does not exist")
	}
}

func TestIsNested_FalseWhenSkillsDirHasEntries(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "skills", "my-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if IsNested(dir) {
		t.Error("expected IsNested=false when skills/ has entries")
	}
}

func TestIsNested_TrueWhenSkillsDirIsEmpty(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "skills"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if !IsNested(dir) {
		t.Error("expected IsNested=true for empty skills/")
	}
}

func TestIsNested_TrueWhenSkillsDirHasOnlyHiddenEntries(t *testing.T) {
	dir := t.TempDir()
	skillsDir := filepath.Join(dir, "skills")
	if err := os.MkdirAll(skillsDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(skillsDir, ".gitkeep"), nil, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if !IsNested(dir) {
		t.Error("expected IsNested=true when skills/ contains only hidden files")
	}
}

// ── ListDomains ─────────────────────────────────────────────────────────────

func TestListDomains_EmptyRoot(t *testing.T) {
	root := t.TempDir()
	domains, err := ListDomains(root)
	if err != nil {
		t.Fatalf("ListDomains: %v", err)
	}
	if len(domains) != 0 {
		t.Errorf("expected 0 domains, got %d", len(domains))
	}
}

func TestListDomains_SkipsFiles(t *testing.T) {
	root := t.TempDir()
	// Regular file at root level — must be ignored
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	buildFlatDomain(t, root, "engineering", "my-skill")

	domains, err := ListDomains(root)
	if err != nil {
		t.Fatalf("ListDomains: %v", err)
	}
	if len(domains) != 1 {
		t.Errorf("expected 1 domain, got %d", len(domains))
	}
	if domains[0].Name != "engineering" {
		t.Errorf("domain name = %q; want %q", domains[0].Name, "engineering")
	}
}

func TestListDomains_SkipsHiddenDirs(t *testing.T) {
	root := t.TempDir()
	buildFlatDomain(t, root, "engineering", "skill-a")
	if err := os.MkdirAll(filepath.Join(root, ".hidden"), 0o755); err != nil {
		t.Fatalf("mkdir hidden: %v", err)
	}

	domains, err := ListDomains(root)
	if err != nil {
		t.Fatalf("ListDomains: %v", err)
	}
	for _, d := range domains {
		if d.Name == ".hidden" {
			t.Error("hidden dir should not appear in domains")
		}
	}
}

func TestListDomains_SkipsClaudePluginDir(t *testing.T) {
	root := t.TempDir()
	buildFlatDomain(t, root, "health", "skill-b")
	if err := os.MkdirAll(filepath.Join(root, ".claude-plugin"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	domains, err := ListDomains(root)
	if err != nil {
		t.Fatalf("ListDomains: %v", err)
	}
	for _, d := range domains {
		if d.Name == ".claude-plugin" {
			t.Error(".claude-plugin should be excluded from domains")
		}
	}
}

func TestListDomains_NestedFlag(t *testing.T) {
	root := t.TempDir()
	// engineering has subdomain layout (nested)
	buildNestedDomain(t, root, "engineering", "development", "write-unit-test")
	// health has flat layout
	buildFlatDomain(t, root, "health", "diagnose-sleep")

	domains, err := ListDomains(root)
	if err != nil {
		t.Fatalf("ListDomains: %v", err)
	}

	byName := map[string]Domain{}
	for _, d := range domains {
		byName[d.Name] = d
	}

	if !byName["engineering"].Nested {
		t.Error("engineering domain should be Nested=true")
	}
	if byName["health"].Nested {
		t.Error("health domain should be Nested=false")
	}
}

// ── ListSubdomains ──────────────────────────────────────────────────────────

func TestListSubdomains_ReturnsSubdomainsWithSkills(t *testing.T) {
	root := t.TempDir()
	buildNestedDomain(t, root, "engineering", "development", "skill-a")
	buildNestedDomain(t, root, "engineering", "testing", "skill-b")

	subs, err := ListSubdomains(filepath.Join(root, "engineering"))
	if err != nil {
		t.Fatalf("ListSubdomains: %v", err)
	}
	if len(subs) != 2 {
		t.Errorf("expected 2 subdomains, got %d", len(subs))
	}
}

func TestListSubdomains_SkipsDirWithoutSkillsFolder(t *testing.T) {
	root := t.TempDir()
	// subdomain with skills/
	buildNestedDomain(t, root, "engineering", "development", "skill-a")
	// subdomain without skills/
	if err := os.MkdirAll(filepath.Join(root, "engineering", "empty-sub"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	subs, err := ListSubdomains(filepath.Join(root, "engineering"))
	if err != nil {
		t.Fatalf("ListSubdomains: %v", err)
	}
	if len(subs) != 1 {
		t.Errorf("expected 1 subdomain, got %d", len(subs))
	}
	if subs[0].Name != "development" {
		t.Errorf("subdomain name = %q; want development", subs[0].Name)
	}
}

// ── ListSkillsInDir ─────────────────────────────────────────────────────────

func TestListSkillsInDir_ReturnsSkillsWithSKILLmd(t *testing.T) {
	root := t.TempDir()
	buildFlatDomain(t, root, "meta", "apply-tdd")
	buildFlatDomain(t, root, "meta", "apply-dry")

	skillList, err := ListSkillsInDir(filepath.Join(root, "meta"), "meta", "")
	if err != nil {
		t.Fatalf("ListSkillsInDir: %v", err)
	}
	if len(skillList) != 2 {
		t.Errorf("expected 2 skills, got %d", len(skillList))
	}
}

func TestListSkillsInDir_SkipsDirWithoutSKILLmd(t *testing.T) {
	root := t.TempDir()
	// Valid skill
	buildFlatDomain(t, root, "meta", "valid-skill")
	// Dir without SKILL.md
	badDir := filepath.Join(root, "meta", "skills", "no-skill-md")
	if err := os.MkdirAll(badDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	skillList, err := ListSkillsInDir(filepath.Join(root, "meta"), "meta", "")
	if err != nil {
		t.Fatalf("ListSkillsInDir: %v", err)
	}
	if len(skillList) != 1 {
		t.Errorf("expected 1 skill (the valid one), got %d", len(skillList))
	}
	if skillList[0].Name != "valid-skill" {
		t.Errorf("skill name = %q; want valid-skill", skillList[0].Name)
	}
}

func TestListSkillsInDir_SkipsHiddenEntries(t *testing.T) {
	root := t.TempDir()
	buildFlatDomain(t, root, "meta", "real-skill")
	skillsDir := filepath.Join(root, "meta", "skills")
	// Hidden dir in skills/
	if err := os.MkdirAll(filepath.Join(skillsDir, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	skillList, err := ListSkillsInDir(filepath.Join(root, "meta"), "meta", "")
	if err != nil {
		t.Fatalf("ListSkillsInDir: %v", err)
	}
	if len(skillList) != 1 {
		t.Errorf("expected 1 skill, got %d", len(skillList))
	}
}

func TestListSkillsInDir_FieldsArePopulated(t *testing.T) {
	root := t.TempDir()
	buildNestedDomain(t, root, "engineering", "testing", "write-unit-test")

	skillList, err := ListSkillsInDir(
		filepath.Join(root, "engineering", "testing"),
		"engineering", "testing",
	)
	if err != nil {
		t.Fatalf("ListSkillsInDir: %v", err)
	}
	if len(skillList) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skillList))
	}
	sk := skillList[0]
	if sk.Domain != "engineering" {
		t.Errorf("Domain = %q; want engineering", sk.Domain)
	}
	if sk.Subdomain != "testing" {
		t.Errorf("Subdomain = %q; want testing", sk.Subdomain)
	}
	if sk.Name != "write-unit-test" {
		t.Errorf("Name = %q; want write-unit-test", sk.Name)
	}
	if sk.Path == "" {
		t.Error("Path must not be empty")
	}
}

// ── ListAllSkills ───────────────────────────────────────────────────────────

func TestListAllSkills_Flat(t *testing.T) {
	root := t.TempDir()
	buildFlatDomain(t, root, "meta", "skill-one")
	buildFlatDomain(t, root, "meta", "skill-two")

	all, err := ListAllSkills(root)
	if err != nil {
		t.Fatalf("ListAllSkills: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 skills, got %d", len(all))
	}
}

func TestListAllSkills_Nested(t *testing.T) {
	root := t.TempDir()
	buildNestedDomain(t, root, "engineering", "development", "skill-a")
	buildNestedDomain(t, root, "engineering", "testing", "skill-b")

	all, err := ListAllSkills(root)
	if err != nil {
		t.Fatalf("ListAllSkills: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 skills across subdomains, got %d", len(all))
	}
}

func TestListAllSkills_MixedDomains(t *testing.T) {
	root := t.TempDir()
	buildNestedDomain(t, root, "engineering", "development", "eng-skill")
	buildFlatDomain(t, root, "health", "health-skill")

	all, err := ListAllSkills(root)
	if err != nil {
		t.Fatalf("ListAllSkills: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 skills total, got %d", len(all))
	}
}

func TestListAllSkills_EmptyRoot(t *testing.T) {
	root := t.TempDir()
	all, err := ListAllSkills(root)
	if err != nil {
		t.Fatalf("ListAllSkills: %v", err)
	}
	if len(all) != 0 {
		t.Errorf("expected 0 skills, got %d", len(all))
	}
}

// ── ResolveSkillPath ────────────────────────────────────────────────────────

func TestResolveSkillPath_ThreePartRef(t *testing.T) {
	root := t.TempDir()
	buildNestedDomain(t, root, "engineering", "development", "write-unit-test")

	got, err := ResolveSkillPath(root, "engineering/development/write-unit-test")
	if err != nil {
		t.Fatalf("ResolveSkillPath: %v", err)
	}
	want := filepath.Join(root, "engineering", "development", "skills", "write-unit-test")
	if got != want {
		t.Errorf("path = %q; want %q", got, want)
	}
}

func TestResolveSkillPath_TwoPartRef(t *testing.T) {
	root := t.TempDir()
	buildFlatDomain(t, root, "meta", "apply-tdd")

	got, err := ResolveSkillPath(root, "meta/apply-tdd")
	if err != nil {
		t.Fatalf("ResolveSkillPath: %v", err)
	}
	want := filepath.Join(root, "meta", "skills", "apply-tdd")
	if got != want {
		t.Errorf("path = %q; want %q", got, want)
	}
}

func TestResolveSkillPath_InvalidRefReturnsError(t *testing.T) {
	root := t.TempDir()
	_, err := ResolveSkillPath(root, "just-one-part")
	if err == nil {
		t.Fatal("expected error for single-part ref, got nil")
	}
	if _, ok := err.(*ErrInvalidSkillRef); !ok {
		t.Errorf("expected ErrInvalidSkillRef, got %T: %v", err, err)
	}
}

func TestResolveSkillPath_NotFoundReturnsError(t *testing.T) {
	root := t.TempDir()
	_, err := ResolveSkillPath(root, "engineering/development/no-such-skill")
	if err == nil {
		t.Fatal("expected error for missing skill, got nil")
	}
	if _, ok := err.(*ErrSkillNotFound); !ok {
		t.Errorf("expected ErrSkillNotFound, got %T: %v", err, err)
	}
}

func TestResolveSkillPath_FourPartRefReturnsError(t *testing.T) {
	root := t.TempDir()
	_, err := ResolveSkillPath(root, "a/b/c/d")
	if err == nil {
		t.Fatal("expected error for 4-part ref")
	}
}

// ── parseSkillMeta ───────────────────────────────────────────────────────────

func writeSkillMD(t *testing.T, dir, content string) string {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestParseSkillMeta_InlineTags(t *testing.T) {
	dir := t.TempDir()
	skillDir := writeSkillMD(t, filepath.Join(dir, "my-skill"), `---
name: my-skill
description: A test skill
tags: [oop, tdd, solid]
---
# My Skill
`)
	meta, _ := parseSkillMeta(skillDir, true)
	if meta.Name != "my-skill" {
		t.Errorf("name = %q, want my-skill", meta.Name)
	}
	if len(meta.Tags) != 3 || meta.Tags[0] != "oop" || meta.Tags[1] != "tdd" || meta.Tags[2] != "solid" {
		t.Errorf("tags = %v, want [oop tdd solid]", meta.Tags)
	}
}

func TestParseSkillMeta_BlockTags(t *testing.T) {
	dir := t.TempDir()
	skillDir := writeSkillMD(t, filepath.Join(dir, "block-skill"), `---
name: block-skill
tags:
  - functional
  - fp
---
`)
	meta, _ := parseSkillMeta(skillDir, true)
	if len(meta.Tags) != 2 || meta.Tags[0] != "functional" || meta.Tags[1] != "fp" {
		t.Errorf("tags = %v, want [functional fp]", meta.Tags)
	}
}

func TestParseSkillMeta_NoFrontmatter(t *testing.T) {
	dir := t.TempDir()
	skillDir := writeSkillMD(t, filepath.Join(dir, "plain"), `# No frontmatter
Just content.
`)
	meta, _ := parseSkillMeta(skillDir, true)
	if meta.Name != "" || len(meta.Tags) != 0 {
		t.Errorf("expected empty for no frontmatter, got name=%q tags=%v", meta.Name, meta.Tags)
	}
}

func TestParseSkillMeta_NoTags(t *testing.T) {
	dir := t.TempDir()
	skillDir := writeSkillMD(t, filepath.Join(dir, "no-tags"), `---
name: no-tags-skill
description: Skill without tags
---
`)
	meta, _ := parseSkillMeta(skillDir, true)
	if meta.Name != "no-tags-skill" {
		t.Errorf("name = %q", meta.Name)
	}
	if len(meta.Tags) != 0 {
		t.Errorf("expected no tags, got %v", meta.Tags)
	}
}

func TestParseSkillMeta_MissingFile(t *testing.T) {
	meta, _ := parseSkillMeta("/nonexistent/skill/dir", true)
	if meta.Name != "" || len(meta.Tags) != 0 {
		t.Errorf("missing SKILL.md should return empty, got name=%q tags=%v", meta.Name, meta.Tags)
	}
}

func TestParseSkillMeta_ExtendedFields(t *testing.T) {
	dir := t.TempDir()
	skillDir := writeSkillMD(t, filepath.Join(dir, "apply-solid"), `---
name: apply-solid-principles
version: 1.2.3
description: Apply SOLID OOP design principles
authors:
  - Jeffrey Tse
license: MIT
compatibility:
  - opencode
  - claude
metadata:
  audience: senior-engineers
dependencies:
  apply-composition-over-inheritance: "*"
  apply-polymorphism: "^1.0.0"
---
`)
	meta, _ := parseSkillMeta(skillDir, true)

	if meta.Version != "1.2.3" {
		t.Errorf("version = %q", meta.Version)
	}
	if meta.Description != "Apply SOLID OOP design principles" {
		t.Errorf("description = %q", meta.Description)
	}
	if len(meta.Authors) != 1 || meta.Authors[0] != "Jeffrey Tse" {
		t.Errorf("authors = %v", meta.Authors)
	}
	if meta.License != "MIT" {
		t.Errorf("license = %q", meta.License)
	}
	if len(meta.Compatibility) != 2 || meta.Compatibility[0] != "opencode" {
		t.Errorf("compatibility = %v", meta.Compatibility)
	}
	if meta.Metadata["audience"] != "senior-engineers" {
		t.Errorf("metadata[audience] = %q", meta.Metadata["audience"])
	}
	if meta.Dependencies["apply-composition-over-inheritance"] != "*" {
		t.Errorf("dep[apply-composition] = %q", meta.Dependencies["apply-composition-over-inheritance"])
	}
	if meta.Dependencies["apply-polymorphism"] != "^1.0.0" {
		t.Errorf("dep[apply-polymorphism] = %q", meta.Dependencies["apply-polymorphism"])
	}
}

func TestListSkillsInDir_PopulatesTags(t *testing.T) {
	root := t.TempDir()
	skillDir := filepath.Join(root, "meta", "skills", "tagged-skill")
	writeSkillMD(t, skillDir, `---
name: tagged-skill
tags: [oop, engineering]
---
`)
	skills, err := ListSkillsInDir(filepath.Join(root, "meta"), "meta", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}
	if len(skills[0].Tags) != 2 || skills[0].Tags[0] != "oop" {
		t.Errorf("Tags = %v, want [oop engineering]", skills[0].Tags)
	}
}
