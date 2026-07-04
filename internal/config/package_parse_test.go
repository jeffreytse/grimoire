package config

import (
	"testing"
)

func TestParseFile_PackageTableArray(t *testing.T) {
	dir := t.TempDir()
	path := writeSettingsFile(t, dir, "settings.toml", `
[[registry]]
name = "official"
url = "https://github.com/jeffreytse/grimoire-core.git"
official = true
priority = 100

[[registry]]
name = "my-team"
url = "https://github.com/acme/grimoire.git"
priority = 50

[[registry]]
name = "local-draft"
url = "/opt/my-skills"
enabled = false
`)
	fs, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(fs.Packages) != 3 {
		t.Fatalf("Packages = %d, want 3", len(fs.Packages))
	}

	official := fs.Packages[0]
	if official.Name != "official" {
		t.Errorf("Packages[0].Name = %q, want official", official.Name)
	}
	if !official.Official {
		t.Error("Packages[0].Official should be true")
	}
	if official.Priority != 100 {
		t.Errorf("Packages[0].Priority = %d, want 100", official.Priority)
	}
	if !official.Enabled {
		t.Error("Packages[0].Enabled should default to true")
	}

	user := fs.Packages[1]
	if user.Name != "my-team" {
		t.Errorf("Packages[1].Name = %q, want my-team", user.Name)
	}
	if user.Official {
		t.Error("Packages[1].Official should be false")
	}

	draft := fs.Packages[2]
	if draft.Enabled {
		t.Error("Packages[2].Enabled should be false")
	}
}

func TestMerge_RegistriesDedup(t *testing.T) {
	layers := []FileConfig{
		{
			Sections: make(map[string]DomainSection),
			Packages: []PackageDef{
				{Name: "official", URL: "https://github.com/org/hub.git", Official: true, Priority: 100, Enabled: true},
				{Name: "team-a", URL: "https://github.com/team-a/hub.git", Priority: 50, Enabled: true},
			},
		},
		{
			Sections: make(map[string]DomainSection),
			Packages: []PackageDef{
				// "official" appears again in lower-priority layer — should be ignored.
				{Name: "official", URL: "https://github.com/other/hub.git", Official: true, Priority: 80, Enabled: true},
				{Name: "team-b", URL: "https://github.com/team-b/hub.git", Priority: 30, Enabled: true},
			},
		},
	}
	r := Merge(layers, []string{"high.toml", "low.toml"})

	if len(r.Packages) != 3 {
		t.Fatalf("Packages = %d, want 3 (dedup of official)", len(r.Packages))
	}
	// Higher-priority layer's "official" URL should win.
	if r.Packages[0].URL != "https://github.com/org/hub.git" {
		t.Errorf("official URL = %q, want org/hub (higher layer)", r.Packages[0].URL)
	}
}

func TestParseFile_PackageRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := writeSettingsFile(t, dir, "settings.toml", "")

	fs, _ := ParseFile(path)
	fs.Packages = []PackageDef{
		{Name: "official", URL: "https://github.com/jeffreytse/grimoire-core.git", Official: true, Priority: 100, Enabled: true},
		{Name: "my-team", URL: "https://github.com/acme/hub.git", Priority: 50, Enabled: false},
	}
	if err := WriteFile(path, fs); err != nil {
		t.Fatal(err)
	}

	fs2, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(fs2.Packages) != 2 {
		t.Fatalf("round-trip Packages = %d, want 2", len(fs2.Packages))
	}
	if fs2.Packages[0].Name != "official" || !fs2.Packages[0].Official {
		t.Errorf("round-trip Packages[0] = %+v", fs2.Packages[0])
	}
	if fs2.Packages[1].Enabled {
		t.Error("round-trip Packages[1].Enabled should be false")
	}
}
