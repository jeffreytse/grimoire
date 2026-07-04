package manifest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jeffreytse/grimoire/internal/config"
)

func writeTOML(t *testing.T, dir, name, content string) string {
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

func TestParseFile_Absent(t *testing.T) {
	mf, err := ParseFile("/nonexistent/grimoire.toml")
	if err != nil {
		t.Fatalf("absent file should not error: %v", err)
	}
	if mf.Package.Name != "" {
		t.Errorf("expected empty package name, got %q", mf.Package.Name)
	}
}

func TestParseFile_Package(t *testing.T) {
	dir := t.TempDir()
	path := writeTOML(t, dir, "grimoire.toml", `
[package]
name = "my-project"
version = "0.1.0"
description = "Test project"
authors = ["Alice <alice@example.com>"]
license = "MIT"
`)
	mf, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if mf.Package.Name != "my-project" {
		t.Errorf("name = %q", mf.Package.Name)
	}
	if mf.Package.Version != "0.1.0" {
		t.Errorf("version = %q", mf.Package.Version)
	}
	if mf.Package.License != "MIT" {
		t.Errorf("license = %q", mf.Package.License)
	}
	if len(mf.Package.Authors) != 1 || mf.Package.Authors[0] != "Alice <alice@example.com>" {
		t.Errorf("authors = %v", mf.Package.Authors)
	}
}

func TestParseFile_DepsStringForm(t *testing.T) {
	dir := t.TempDir()
	path := writeTOML(t, dir, "grimoire.toml", `
[dependencies]
apply-solid = "^1.0.0"
apply-dry = "*"
`)
	mf, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(mf.Deps) != 2 {
		t.Fatalf("want 2 deps, got %d", len(mf.Deps))
	}
	if mf.Deps["apply-solid"].Version != "^1.0.0" {
		t.Errorf("apply-solid version = %q", mf.Deps["apply-solid"].Version)
	}
	if mf.Deps["apply-dry"].Version != "*" {
		t.Errorf("apply-dry version = %q", mf.Deps["apply-dry"].Version)
	}
}

func TestParseFile_DepsPackageInKey(t *testing.T) {
	dir := t.TempDir()
	path := writeTOML(t, dir, "grimoire.toml", `
[dependencies]
apply-solid = "^1.0.0"
"acmecorp/practices:apply-tdd" = "~2.0.0"
"github.com/myteam/skills:domain/skill" = "*"
`)
	mf, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(mf.Deps) != 3 {
		t.Fatalf("want 3 deps, got %d", len(mf.Deps))
	}
	if mf.Deps["apply-solid"].Version != "^1.0.0" {
		t.Errorf("apply-solid version = %q", mf.Deps["apply-solid"].Version)
	}
	if mf.Deps["acmecorp/practices:apply-tdd"].Version != "~2.0.0" {
		t.Errorf("acmecorp dep version = %q", mf.Deps["acmecorp/practices:apply-tdd"].Version)
	}
	if mf.Deps["github.com/myteam/skills:domain/skill"].Version != "*" {
		t.Errorf("github dep version = %q", mf.Deps["github.com/myteam/skills:domain/skill"].Version)
	}
}

func TestParseFile_Standards(t *testing.T) {
	dir := t.TempDir()
	path := writeTOML(t, dir, "grimoire.toml", `
[standards]
profiles = ["engineering"]
report-path = ".grimoire/report.json"
staleness-days = 7

[standards.engineering]
practices = ["apply-solid-principles"]
fallback = "ask"
compliance-threshold = 0.8

[standards.engineering.testing]
practices = ["apply-tdd"]
`)
	mf, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(mf.Standards.Profiles) != 1 || mf.Standards.Profiles[0] != "engineering" {
		t.Errorf("profiles = %v", mf.Standards.Profiles)
	}
	if mf.Standards.ReportPath != ".grimoire/report.json" {
		t.Errorf("report-path = %q", mf.Standards.ReportPath)
	}
	if mf.Standards.StalenessDays != 7 {
		t.Errorf("staleness-days = %d", mf.Standards.StalenessDays)
	}

	eng, ok := mf.DomainSections["engineering"]
	if !ok {
		t.Fatal("expected engineering domain section")
	}
	if len(eng.Practices) != 1 || eng.Practices[0] != "apply-solid-principles" {
		t.Errorf("engineering practices = %v", eng.Practices)
	}
	if eng.Fallback != "ask" {
		t.Errorf("fallback = %q", eng.Fallback)
	}
	if eng.ComplianceThreshold != 0.8 {
		t.Errorf("compliance-threshold = %v", eng.ComplianceThreshold)
	}

	testingSec, ok := mf.DomainSections["engineering.testing"]
	if !ok {
		t.Fatal("expected engineering.testing domain section")
	}
	if len(testingSec.Practices) != 1 || testingSec.Practices[0] != "apply-tdd" {
		t.Errorf("engineering.testing practices = %v", testingSec.Practices)
	}
}

func TestParseFile_SearchPackages(t *testing.T) {
	dir := t.TempDir()
	path := writeTOML(t, dir, "grimoire.toml", `
[core]
search-packages = ["github.com/acmecorp/practices", "github.com/myteam/skills"]
`)
	mf, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(mf.Core.SearchPackages) != 2 {
		t.Fatalf("want 2 search-packages, got %d", len(mf.Core.SearchPackages))
	}
	if mf.Core.SearchPackages[0] != "github.com/acmecorp/practices" {
		t.Errorf("search-packages[0] = %q", mf.Core.SearchPackages[0])
	}
}

func TestMerge_ProjectWins(t *testing.T) {
	projDir := t.TempDir()
	globalDir := t.TempDir()

	writeTOML(t, projDir, "grimoire.toml", `
[package]
name = "my-project"

[dependencies]
apply-solid = "^1.0.0"

[standards]
staleness-days = 3
`)

	writeTOML(t, globalDir, "grimoire.toml", `
[dependencies]
apply-solid = "^2.0.0"
apply-dry = "*"

[standards]
staleness-days = 7
`)

	projMF, _ := ParseFile(filepath.Join(projDir, "grimoire.toml"))
	globalMF, _ := ParseFile(filepath.Join(globalDir, "grimoire.toml"))

	r := merge([]ManifestFile{projMF, globalMF}, []string{"project", "global"})

	// project version of apply-solid wins
	if r.Deps["apply-solid"].Version != "^1.0.0" {
		t.Errorf("apply-solid = %q, want ^1.0.0", r.Deps["apply-solid"].Version)
	}
	// global-only dep still present
	if r.Deps["apply-dry"].Version != "*" {
		t.Errorf("apply-dry = %q, want *", r.Deps["apply-dry"].Version)
	}
	// project staleness-days wins
	if r.Standards.StalenessDays != 3 {
		t.Errorf("staleness-days = %d, want 3", r.Standards.StalenessDays)
	}
}

func TestWriteFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "grimoire.toml")

	mf := ManifestFile{
		Package: PackageMeta{Name: "test", Version: "1.0.0", License: "MIT"},
		Deps: map[string]DepSpec{
			"apply-solid":                  {Version: "^1.0.0"},
			"acmecorp/practices:apply-dry": {Version: "*"},
		},
		DevDeps:        make(map[string]DepSpec),
		DomainSections: make(map[string]config.DomainSection),
		InlineProfiles: make(map[string]config.InlineProfileDef),
		Standards:      StandardsSection{Profiles: []string{"engineering"}, StalenessDays: 7},
	}

	if err := WriteFile(path, &mf); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	got, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile after write: %v", err)
	}
	if got.Package.Name != "test" {
		t.Errorf("name = %q", got.Package.Name)
	}
	if got.Deps["apply-solid"].Version != "^1.0.0" {
		t.Errorf("apply-solid = %q", got.Deps["apply-solid"].Version)
	}
	if got.Deps["acmecorp/practices:apply-dry"].Version != "*" {
		t.Errorf("acmecorp dep version = %q", got.Deps["acmecorp/practices:apply-dry"].Version)
	}
	if got.Standards.StalenessDays != 7 {
		t.Errorf("staleness-days = %d", got.Standards.StalenessDays)
	}
}

func TestParseFile_InlineProfiles(t *testing.T) {
	dir := t.TempDir()
	path := writeTOML(t, dir, "grimoire.toml", `
[profiles.strict-oop]
description = "Strict OOP profile"
tags = ["engineering", "oop"]
extends = ["engineering"]
exclude = ["apply-procedural-style"]
compliance-threshold = 0.9

[[profiles.strict-oop.skills]]
name = "apply-solid-principles"
priority = 80
`)
	mf, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(mf.InlineProfiles) != 1 {
		t.Fatalf("want 1 profile, got %d", len(mf.InlineProfiles))
	}
	p, ok := mf.InlineProfiles["strict-oop"]
	if !ok {
		t.Fatal("expected strict-oop profile")
	}
	if p.Description != "Strict OOP profile" {
		t.Errorf("description = %q", p.Description)
	}
	if len(p.Tags) != 2 || p.Tags[0] != "engineering" {
		t.Errorf("tags = %v", p.Tags)
	}
	if len(p.Extends) != 1 || p.Extends[0] != "engineering" {
		t.Errorf("extends = %v", p.Extends)
	}
	if len(p.Exclude) != 1 || p.Exclude[0] != "apply-procedural-style" {
		t.Errorf("exclude = %v", p.Exclude)
	}
	if p.ComplianceThreshold != 0.9 {
		t.Errorf("compliance-threshold = %v", p.ComplianceThreshold)
	}
	if len(p.Skills) != 1 || p.Skills[0].Name != "apply-solid-principles" || p.Skills[0].Priority != 80 {
		t.Errorf("skills = %+v", p.Skills)
	}
}

func TestWriteFile_InlineProfilesRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "grimoire.toml")

	mf := ManifestFile{
		Deps:           make(map[string]DepSpec),
		DevDeps:        make(map[string]DepSpec),
		DomainSections: make(map[string]config.DomainSection),
		InlineProfiles: map[string]config.InlineProfileDef{
			"strict-oop": {
				Description:              "Strict OOP",
				Tags:                     []string{"engineering"},
				Extends:                  []string{"engineering"},
				ComplianceThreshold:      0.9,
				ComplianceThresholdError: -1,
				Skills:                   []config.InlineSkillRef{{Name: "apply-solid-principles", Priority: 80}},
			},
		},
	}

	if err := WriteFile(path, &mf); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	got, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	p, ok := got.InlineProfiles["strict-oop"]
	if !ok {
		t.Fatal("expected strict-oop after round-trip")
	}
	if p.Description != "Strict OOP" {
		t.Errorf("description = %q", p.Description)
	}
	if p.ComplianceThreshold != 0.9 {
		t.Errorf("compliance-threshold = %v", p.ComplianceThreshold)
	}
	if len(p.Skills) != 1 || p.Skills[0].Name != "apply-solid-principles" {
		t.Errorf("skills = %+v", p.Skills)
	}
}

func TestMerge_InlineProfiles(t *testing.T) {
	proj := ManifestFile{
		Deps:           make(map[string]DepSpec),
		DevDeps:        make(map[string]DepSpec),
		DomainSections: make(map[string]config.DomainSection),
		InlineProfiles: map[string]config.InlineProfileDef{
			"strict-oop": {Description: "from-project", ComplianceThresholdError: -1},
		},
	}
	global := ManifestFile{
		Deps:           make(map[string]DepSpec),
		DevDeps:        make(map[string]DepSpec),
		DomainSections: make(map[string]config.DomainSection),
		InlineProfiles: map[string]config.InlineProfileDef{
			"strict-oop": {Description: "from-global", ComplianceThresholdError: -1},
			"relaxed":    {Description: "from-global-only", ComplianceThresholdError: -1},
		},
	}

	r := merge([]ManifestFile{proj, global}, []string{"project", "global"})

	if r.InlineProfiles["strict-oop"].Description != "from-project" {
		t.Errorf("project should win: got %q", r.InlineProfiles["strict-oop"].Description)
	}
	if r.InlineProfiles["relaxed"].Description != "from-global-only" {
		t.Errorf("global-only profile missing: got %q", r.InlineProfiles["relaxed"].Description)
	}
}
