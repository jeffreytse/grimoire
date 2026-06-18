package settings

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeSettingsFile(t *testing.T, dir, name, content string) string {
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

func TestParseFile_MissingFile(t *testing.T) {
	fs, err := ParseFile("/nonexistent/path/settings.toml")
	if err != nil {
		t.Fatalf("missing file should not error: %v", err)
	}
	if fs.Core.Home != "" || fs.Core.Source != "" || len(fs.Core.Profiles) != 0 {
		t.Error("missing file should return zero CoreSection")
	}
	if len(fs.Sections) != 0 {
		t.Error("missing file should return empty Sections")
	}
}

func TestParseFile_CoreSection(t *testing.T) {
	dir := t.TempDir()
	path := writeSettingsFile(t, dir, "settings.toml", `
[core]
home = "/opt/grimoire"
source = "https://example.com/skills.git"
`)
	fs, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if fs.Core.Home != "/opt/grimoire" {
		t.Errorf("Core.Home = %q, want /opt/grimoire", fs.Core.Home)
	}
	if fs.Core.Source != "https://example.com/skills.git" {
		t.Errorf("Core.Source = %q", fs.Core.Source)
	}
	// profiles must NOT be parsed from [core]
	if len(fs.Core.Profiles) != 0 {
		t.Errorf("Core.Profiles should be empty when not set in [standards], got %v", fs.Core.Profiles)
	}
}

func TestParseFile_StandardsProfiles(t *testing.T) {
	dir := t.TempDir()
	path := writeSettingsFile(t, dir, "settings.toml", `
[standards]
profiles = ["oop", "tdd"]
`)
	fs, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(fs.Core.Profiles) != 2 || fs.Core.Profiles[0] != "oop" || fs.Core.Profiles[1] != "tdd" {
		t.Errorf("Core.Profiles = %v, want [oop tdd]", fs.Core.Profiles)
	}
}

func TestParseFile_DomainSection(t *testing.T) {
	dir := t.TempDir()
	path := writeSettingsFile(t, dir, "settings.toml", `
[standards.engineering]
practices = ["Google Engineering Practices"]
fallback = "ask"
compliance-threshold = 75
compliance-threshold-error = 0
`)
	fs, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	ds, ok := fs.Sections["engineering"]
	if !ok {
		t.Fatal("expected engineering section")
	}
	if len(ds.Practices) != 1 || ds.Practices[0] != "Google Engineering Practices" {
		t.Errorf("Practices = %v", ds.Practices)
	}
	if ds.Fallback != "ask" {
		t.Errorf("Fallback = %q", ds.Fallback)
	}
	if ds.ComplianceThreshold != 75 {
		t.Errorf("ComplianceThreshold = %v", ds.ComplianceThreshold)
	}
	// 0 must be preserved (not treated as unset)
	if ds.ComplianceThresholdError != 0 {
		t.Errorf("ComplianceThresholdError = %v, want 0", ds.ComplianceThresholdError)
	}
}

func TestParseFile_SubdomainSection(t *testing.T) {
	dir := t.TempDir()
	path := writeSettingsFile(t, dir, "settings.toml", `
[standards.engineering]
compliance-threshold = 75

[standards.engineering.architecture]
compliance-threshold = 85
`)
	fs, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if fs.Sections["engineering"].ComplianceThreshold != 75 {
		t.Errorf("domain threshold = %v", fs.Sections["engineering"].ComplianceThreshold)
	}
	if fs.Sections["engineering.architecture"].ComplianceThreshold != 85 {
		t.Errorf("subdomain threshold = %v", fs.Sections["engineering.architecture"].ComplianceThreshold)
	}
}

func TestMerge_ProjectPersonalOverridesProjectShared(t *testing.T) {
	personal := FileSettings{
		Sections: map[string]DomainSection{
			"engineering": {Fallback: "ask", ComplianceThresholdError: -1},
		},
	}
	shared := FileSettings{
		Sections: map[string]DomainSection{
			"engineering": {Fallback: "skip", ComplianceThresholdError: -1},
		},
	}

	r := Merge([]FileSettings{personal, shared}, []string{"local.toml", "shared.toml"})
	ds := r.ResolveSection("engineering")
	if ds.Fallback != "ask" {
		t.Errorf("project-personal should win: Fallback = %q", ds.Fallback)
	}
	if src := r.Sources["engineering.fallback"]; src != "local.toml" {
		t.Errorf("source should be local.toml, got %q", src)
	}
}

func TestMerge_ProjectSharedOverridesGlobal(t *testing.T) {
	shared := FileSettings{
		Sections: map[string]DomainSection{
			"engineering": {ComplianceThreshold: 80, ComplianceThresholdError: -1},
		},
	}
	global := FileSettings{
		Sections: map[string]DomainSection{
			"engineering": {ComplianceThreshold: 60, ComplianceThresholdError: -1},
		},
	}

	r := Merge([]FileSettings{shared, global}, []string{"shared.toml", "global.toml"})
	ds := r.ResolveSection("engineering")
	if ds.ComplianceThreshold != 80 {
		t.Errorf("project-shared should win: threshold = %v", ds.ComplianceThreshold)
	}
}

func TestMerge_IndependentKeysMergeAcrossLayers(t *testing.T) {
	shared := FileSettings{
		Sections: map[string]DomainSection{
			"engineering": {ComplianceThreshold: 80, ComplianceThresholdError: -1},
		},
	}
	global := FileSettings{
		Sections: map[string]DomainSection{
			"engineering": {
				Practices:                []string{"Google"},
				ComplianceThresholdError: -1,
			},
		},
	}

	r := Merge([]FileSettings{shared, global}, []string{"shared.toml", "global.toml"})
	ds := r.ResolveSection("engineering")
	// threshold from shared, practices from global
	if ds.ComplianceThreshold != 80 {
		t.Errorf("threshold = %v", ds.ComplianceThreshold)
	}
	if len(ds.Practices) == 0 || ds.Practices[0] != "Google" {
		t.Errorf("practices = %v", ds.Practices)
	}
}

func TestResolveSection_SubdomainOverridesDomain(t *testing.T) {
	r := Resolved{
		sections: map[string]DomainSection{
			"engineering": {
				ComplianceThreshold:      75,
				Fallback:                 "skip",
				ComplianceThresholdError: -1,
			},
			"engineering.architecture": {
				ComplianceThreshold:      85,
				ComplianceThresholdError: -1,
			},
		},
		Sources: map[string]string{},
	}

	ds := r.ResolveSection("engineering.architecture")
	if ds.ComplianceThreshold != 85 {
		t.Errorf("subdomain should override: threshold = %v", ds.ComplianceThreshold)
	}
	// fallback not set in subdomain → inherit from domain
	if ds.Fallback != "skip" {
		t.Errorf("should inherit domain fallback, got %q", ds.Fallback)
	}
}

func TestResolveSection_SubdomainInheritsDomainKeys(t *testing.T) {
	r := Resolved{
		sections: map[string]DomainSection{
			"engineering": {
				Practices:                []string{"Google"},
				ComplianceThreshold:      75,
				ComplianceThresholdError: 0,
			},
			"engineering.testing": {
				Practices:                []string{"apply-tdd"},
				ComplianceThresholdError: -1, // unset at subdomain level
			},
		},
		Sources: map[string]string{},
	}

	ds := r.ResolveSection("engineering.testing")
	// practices overridden by subdomain
	if len(ds.Practices) != 1 || ds.Practices[0] != "apply-tdd" {
		t.Errorf("practices = %v", ds.Practices)
	}
	// threshold not set in subdomain → inherit from domain
	if ds.ComplianceThreshold != 75 {
		t.Errorf("should inherit domain threshold: %v", ds.ComplianceThreshold)
	}
	// threshold-error not set in subdomain → inherit domain's 0
	if ds.ComplianceThresholdError != 0 {
		t.Errorf("should inherit domain threshold-error 0, got %v", ds.ComplianceThresholdError)
	}
}

func TestComplianceThresholdError_ZeroPreserved(t *testing.T) {
	dir := t.TempDir()
	path := writeSettingsFile(t, dir, "settings.toml", `
[standards.engineering]
compliance-threshold-error = 0
`)
	fs, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	ds := fs.Sections["engineering"]
	if ds.ComplianceThresholdError != 0 {
		t.Errorf("0 must be preserved as 'allow no errors', got %v", ds.ComplianceThresholdError)
	}
}

func TestWriteFile_RoundTrip(t *testing.T) {
	original := FileSettings{
		Core: CoreSection{
			Home:     "/opt/grimoire",
			Profiles: []string{"oop"},
		},
		Sections: map[string]DomainSection{
			"engineering": {
				Practices:                []string{"Google"},
				ComplianceThreshold:      80,
				ComplianceThresholdError: -1,
			},
		},
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "settings.toml")
	if err := WriteFile(path, original); err != nil {
		t.Fatal(err)
	}

	// verify emitted TOML uses [standards.*] namespace and profiles not in [core]
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	rawStr := string(raw)
	if !strings.Contains(rawStr, "[standards.") {
		t.Errorf("emitted TOML must use [standards.*] sections, got:\n%s", rawStr)
	}
	if strings.Contains(rawStr, "profiles") && strings.Contains(rawStr, "[core]") {
		// profiles should be under [standards], not [core]
		coreIdx := strings.Index(rawStr, "[core]")
		profilesIdx := strings.Index(rawStr, "profiles")
		standardsIdx := strings.Index(rawStr, "[standards")
		if profilesIdx > coreIdx && (standardsIdx == -1 || profilesIdx < standardsIdx) {
			t.Errorf("profiles must be under [standards], not [core]:\n%s", rawStr)
		}
	}

	got, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Core.Home != "/opt/grimoire" {
		t.Errorf("Core.Home = %q", got.Core.Home)
	}
	if len(got.Core.Profiles) != 1 || got.Core.Profiles[0] != "oop" {
		t.Errorf("Core.Profiles = %v", got.Core.Profiles)
	}
	ds := got.Sections["engineering"]
	if ds.ComplianceThreshold != 80 {
		t.Errorf("threshold = %v", ds.ComplianceThreshold)
	}
}
