package resolver

import (
	"testing"

	"github.com/jeffreytse/grimoire/internal/manifest"
)

func deps(kv ...string) map[string]manifest.DepSpec {
	m := make(map[string]manifest.DepSpec, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i]] = manifest.DepSpec{Version: kv[i+1]}
	}
	return m
}

func meta(kv ...string) map[string]SkillMeta {
	m := make(map[string]SkillMeta, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i]] = SkillMeta{Name: kv[i], Version: kv[i+1], Source: "jeffreytse/grimoire-core"}
	}
	return m
}

func TestSkillNameFromKey(t *testing.T) {
	cases := []struct {
		key  string
		want string
	}{
		{"apply-solid", "apply-solid"},
		{"acmecorp/practices:apply-tdd", "apply-tdd"},
		{"github.com/myteam/skills:domain/skill", "skill"},
		{"apply-dry-principle", "apply-dry-principle"},
	}
	for _, tc := range cases {
		if got := skillNameFromKey(tc.key); got != tc.want {
			t.Errorf("skillNameFromKey(%q) = %q, want %q", tc.key, got, tc.want)
		}
	}
}

func TestResolve_Simple(t *testing.T) {
	r := New(meta("apply-solid", "1.2.3", "apply-dry", "2.0.0"))
	entries, err := r.Resolve(deps("apply-solid", "^1.0.0", "apply-dry", "*"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("want 2 entries, got %d", len(entries))
	}
	// entries sorted by name
	if entries[0].Name != "apply-dry" {
		t.Errorf("entries[0].Name = %q, want apply-dry", entries[0].Name)
	}
	if entries[1].Name != "apply-solid" {
		t.Errorf("entries[1].Name = %q, want apply-solid", entries[1].Name)
	}
	if entries[1].Version != "1.2.3" {
		t.Errorf("apply-solid version = %q, want 1.2.3", entries[1].Version)
	}
}

func TestResolve_PackageInKey(t *testing.T) {
	metaMap := map[string]SkillMeta{
		"apply-tdd": {Name: "apply-tdd", Version: "2.0.5", Source: "acmecorp/practices",
			Resolved: "https://github.com/acmecorp/practices.git"},
	}
	r := New(metaMap)
	entries, err := r.Resolve(deps("acmecorp/practices:apply-tdd", "~2.0.0"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("want 1 entry, got %d", len(entries))
	}
	if entries[0].Name != "apply-tdd" {
		t.Errorf("name = %q", entries[0].Name)
	}
	if entries[0].Source != "acmecorp/practices" {
		t.Errorf("source = %q", entries[0].Source)
	}
}

func TestResolve_Transitive(t *testing.T) {
	metaMap := map[string]SkillMeta{
		"apply-solid": {
			Name: "apply-solid", Version: "1.0.0", Source: "jeffreytse/grimoire-core",
			Dependencies: map[string]string{"apply-dry": "^2.0.0"},
		},
		"apply-dry": {Name: "apply-dry", Version: "2.3.0", Source: "jeffreytse/grimoire-core"},
	}
	r := New(metaMap)
	entries, err := r.Resolve(deps("apply-solid", "^1.0.0"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("want 2 entries (including transitive), got %d", len(entries))
	}
	names := make(map[string]bool)
	for _, e := range entries {
		names[e.Name] = true
	}
	if !names["apply-solid"] || !names["apply-dry"] {
		t.Errorf("missing transitive dep; got names: %v", names)
	}
}

func TestResolve_Conflict(t *testing.T) {
	// apply-solid 1.2.3 doesn't satisfy ^2.0.0
	r := New(meta("apply-solid", "1.2.3"))
	_, err := r.Resolve(deps("apply-solid", "^2.0.0"))
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
	ce, ok := err.(*ErrConflict)
	if !ok {
		t.Fatalf("expected *ErrConflict, got %T", err)
	}
	if len(ce.Conflicts) != 1 {
		t.Fatalf("want 1 conflict, got %d", len(ce.Conflicts))
	}
	if ce.Conflicts[0].Skill != "apply-solid" {
		t.Errorf("conflict skill = %q", ce.Conflicts[0].Skill)
	}
}

func TestResolve_TransitiveConflict(t *testing.T) {
	// root requires apply-solid ^1.0.0 AND apply-dry ^3.0.0
	// apply-solid transitively requires apply-dry ^2.0.0
	// apply-dry installed as 2.5.0 — satisfies ^2.0.0 but NOT ^3.0.0
	metaMap := map[string]SkillMeta{
		"apply-solid": {
			Name: "apply-solid", Version: "1.0.0", Source: "jeffreytse/grimoire-core",
			Dependencies: map[string]string{"apply-dry": "^2.0.0"},
		},
		"apply-dry": {Name: "apply-dry", Version: "2.5.0", Source: "jeffreytse/grimoire-core"},
	}
	r := New(metaMap)
	_, err := r.Resolve(deps("apply-solid", "^1.0.0", "apply-dry", "^3.0.0"))
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
	ce, ok := err.(*ErrConflict)
	if !ok {
		t.Fatalf("expected *ErrConflict, got %T", err)
	}
	if ce.Conflicts[0].Skill != "apply-dry" {
		t.Errorf("conflict skill = %q", ce.Conflicts[0].Skill)
	}
}

func TestResolve_Versionless(t *testing.T) {
	// skill with no metadata — any constraint should pass (no version to conflict with)
	r := New(map[string]SkillMeta{})
	entries, err := r.Resolve(deps("apply-solid", "^1.0.0"))
	if err != nil {
		t.Fatalf("versionless skill should not conflict: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("want 1 entry, got %d", len(entries))
	}
	if entries[0].Version != "" {
		t.Errorf("versionless entry should have empty version, got %q", entries[0].Version)
	}
}

func TestResolve_Wildcard(t *testing.T) {
	r := New(meta("apply-dry", "3.0.0"))
	entries, err := r.Resolve(deps("apply-dry", "*"))
	if err != nil {
		t.Fatalf("wildcard should always satisfy: %v", err)
	}
	if len(entries) != 1 || entries[0].Version != "3.0.0" {
		t.Errorf("unexpected entries: %v", entries)
	}
}

func TestResolve_Empty(t *testing.T) {
	r := New(map[string]SkillMeta{})
	entries, err := r.Resolve(map[string]manifest.DepSpec{})
	if err != nil {
		t.Fatalf("empty deps should not error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("want 0 entries, got %d", len(entries))
	}
}
