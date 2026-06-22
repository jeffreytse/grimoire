package settings

import (
	"testing"
)

func TestParseFile_RegistryTableArray(t *testing.T) {
	dir := t.TempDir()
	path := writeSettingsFile(t, dir, "settings.toml", `
[[registry]]
name = "official"
url = "https://github.com/jeffreytse/grimoire-hub.git"
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
	if len(fs.Registries) != 3 {
		t.Fatalf("Registries = %d, want 3", len(fs.Registries))
	}

	official := fs.Registries[0]
	if official.Name != "official" {
		t.Errorf("Registries[0].Name = %q, want official", official.Name)
	}
	if !official.Official {
		t.Error("Registries[0].Official should be true")
	}
	if official.Priority != 100 {
		t.Errorf("Registries[0].Priority = %d, want 100", official.Priority)
	}
	if !official.Enabled {
		t.Error("Registries[0].Enabled should default to true")
	}

	user := fs.Registries[1]
	if user.Name != "my-team" {
		t.Errorf("Registries[1].Name = %q, want my-team", user.Name)
	}
	if user.Official {
		t.Error("Registries[1].Official should be false")
	}

	draft := fs.Registries[2]
	if draft.Enabled {
		t.Error("Registries[2].Enabled should be false")
	}
}

func TestMerge_RegistriesDedup(t *testing.T) {
	layers := []FileSettings{
		{
			Sections: make(map[string]DomainSection),
			Registries: []RegistryDef{
				{Name: "official", URL: "https://github.com/org/hub.git", Official: true, Priority: 100, Enabled: true},
				{Name: "team-a", URL: "https://github.com/team-a/hub.git", Priority: 50, Enabled: true},
			},
		},
		{
			Sections: make(map[string]DomainSection),
			Registries: []RegistryDef{
				// "official" appears again in lower-priority layer — should be ignored.
				{Name: "official", URL: "https://github.com/other/hub.git", Official: true, Priority: 80, Enabled: true},
				{Name: "team-b", URL: "https://github.com/team-b/hub.git", Priority: 30, Enabled: true},
			},
		},
	}
	r := Merge(layers, []string{"high.toml", "low.toml"})

	if len(r.Registries) != 3 {
		t.Fatalf("Registries = %d, want 3 (dedup of official)", len(r.Registries))
	}
	// Higher-priority layer's "official" URL should win.
	if r.Registries[0].URL != "https://github.com/org/hub.git" {
		t.Errorf("official URL = %q, want org/hub (higher layer)", r.Registries[0].URL)
	}
}

func TestParseFile_RegistryRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := writeSettingsFile(t, dir, "settings.toml", "")

	fs, _ := ParseFile(path)
	fs.Registries = []RegistryDef{
		{Name: "official", URL: "https://github.com/jeffreytse/grimoire-hub.git", Official: true, Priority: 100, Enabled: true},
		{Name: "my-team", URL: "https://github.com/acme/hub.git", Priority: 50, Enabled: false},
	}
	if err := WriteFile(path, fs); err != nil {
		t.Fatal(err)
	}

	fs2, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(fs2.Registries) != 2 {
		t.Fatalf("round-trip Registries = %d, want 2", len(fs2.Registries))
	}
	if fs2.Registries[0].Name != "official" || !fs2.Registries[0].Official {
		t.Errorf("round-trip Registries[0] = %+v", fs2.Registries[0])
	}
	if fs2.Registries[1].Enabled {
		t.Error("round-trip Registries[1].Enabled should be false")
	}
}
