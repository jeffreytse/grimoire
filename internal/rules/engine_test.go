package rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jeffreytse/grimoire/internal/skills"
)

func makeSkillDir(t *testing.T, root, domain, skillName string, withSkillMd bool, frontmatter string) {
	t.Helper()
	dir := filepath.Join(root, domain, "skills", skillName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if withSkillMd {
		content := frontmatter
		if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func testSource(root string) []skills.SkillsSource {
	return []skills.SkillsSource{{Name: "test", Root: root}}
}

func TestCheckSkillHasSkillMd_Missing(t *testing.T) {
	root := t.TempDir()
	makeSkillDir(t, root, "engineering", "apply-solid", false, "")

	diags := checkSkillHasSkillMd(testSource(root))
	if len(diags) != 1 {
		t.Fatalf("want 1 diagnostic, got %d", len(diags))
	}
	if diags[0].Code != "skill-has-skill-md" {
		t.Errorf("code = %q", diags[0].Code)
	}
	if diags[0].Severity != 1 {
		t.Errorf("severity = %d, want 1 (Error)", diags[0].Severity)
	}
	if diags[0].Source != "grimoire-rules" {
		t.Errorf("source = %q", diags[0].Source)
	}
}

func TestCheckSkillHasSkillMd_Present(t *testing.T) {
	root := t.TempDir()
	fm := "---\nname: apply-solid\ntags: [engineering]\n---\n# Apply SOLID\n"
	makeSkillDir(t, root, "engineering", "apply-solid", true, fm)

	diags := checkSkillHasSkillMd(testSource(root))
	if len(diags) != 0 {
		t.Errorf("want 0 diagnostics, got %d: %+v", len(diags), diags)
	}
}

func TestCheckSkillMdFrontmatter_MissingFrontmatter(t *testing.T) {
	root := t.TempDir()
	makeSkillDir(t, root, "engineering", "apply-solid", true, "# No frontmatter here\n")

	diags := checkSkillMdFrontmatter(testSource(root))
	if len(diags) == 0 {
		t.Fatal("want diagnostics for missing frontmatter, got none")
	}
	found := false
	for _, d := range diags {
		if d.Code == "skill-md-has-frontmatter" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected skill-md-has-frontmatter diagnostic, got: %+v", diags)
	}
}

func TestCheckSkillMdFrontmatter_MissingTags(t *testing.T) {
	root := t.TempDir()
	fm := "---\nname: apply-solid\n---\n# No tags\n"
	makeSkillDir(t, root, "engineering", "apply-solid", true, fm)

	diags := checkSkillMdFrontmatter(testSource(root))
	found := false
	for _, d := range diags {
		if d.Code == "skill-md-has-tags" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected skill-md-has-tags diagnostic, got: %+v", diags)
	}
}

func TestCheckSkillMdFrontmatter_Complete(t *testing.T) {
	root := t.TempDir()
	fm := "---\nname: apply-solid\ntags: [engineering]\n---\n# SOLID\n"
	makeSkillDir(t, root, "engineering", "apply-solid", true, fm)

	diags := checkSkillMdFrontmatter(testSource(root))
	if len(diags) != 0 {
		t.Errorf("want 0 diagnostics for complete SKILL.md, got %d: %+v", len(diags), diags)
	}
}

func TestCheckSettingsParseable_Absent(t *testing.T) {
	dir := t.TempDir()
	diags := checkSettingsParseable(dir)
	if len(diags) != 0 {
		t.Errorf("absent settings.toml should produce no diagnostics, got %d", len(diags))
	}
}

func TestCheckSettingsParseable_Valid(t *testing.T) {
	dir := t.TempDir()
	settingsDir := filepath.Join(dir, ".grimoire")
	if err := os.MkdirAll(settingsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(settingsDir, "settings.toml"), []byte(`
[standards.engineering]
practices = ["apply-solid-principles"]
`), 0o644); err != nil {
		t.Fatal(err)
	}

	diags := checkSettingsParseable(dir)
	if len(diags) != 0 {
		t.Errorf("valid settings.toml should produce no diagnostics, got %d", len(diags))
	}
}

func TestCheckSettingsParseable_Invalid(t *testing.T) {
	dir := t.TempDir()
	settingsDir := filepath.Join(dir, ".grimoire")
	if err := os.MkdirAll(settingsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(settingsDir, "settings.toml"), []byte(`
[[not valid toml
`), 0o644); err != nil {
		t.Fatal(err)
	}

	diags := checkSettingsParseable(dir)
	if len(diags) != 1 {
		t.Fatalf("invalid settings.toml should produce 1 diagnostic, got %d", len(diags))
	}
	if diags[0].Code != "settings-toml-parseable" {
		t.Errorf("code = %q", diags[0].Code)
	}
	if diags[0].Severity != 1 {
		t.Errorf("severity = %d, want 1 (Error)", diags[0].Severity)
	}
}

func TestEngine_Run_Empty(t *testing.T) {
	eng := &Engine{
		SkillsSources: []skills.SkillsSource{},
		ProjectDir:    t.TempDir(),
	}
	diags := eng.Run()
	// no sources, no project files → no findings (agent symlink checks may find nothing)
	for _, d := range diags {
		if d.Source != "grimoire-rules" {
			t.Errorf("all diagnostics must have source=grimoire-rules, got %q", d.Source)
		}
	}
}
