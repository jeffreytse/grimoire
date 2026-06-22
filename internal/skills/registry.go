package skills

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/jeffreytse/grimoire/internal/settings"
)

const OfficialRegistryName = "official"

// SkillsSource is a named registry with its resolved skills directory path.
type SkillsSource struct {
	Name string // derived registry name, e.g. "jeffreytse/grimoire-hub", "acmecorp/standards"
	Root string // absolute path to the skills/ directory
}

// RegistryEntry is a configured registry with its home directory.
// Unlike SkillsSource, this includes registries that have no skills/ directory
// (e.g. profiles-only or settings-only bundles).
type RegistryEntry struct {
	Name     string // registry name: user-given (new model) or URL-derived (old model)
	Home     string // absolute path to registry root dir (parent of skills/)
	Priority int    // higher = searched first; 100 for official, 50 for user, 0 = old-model default
	Official bool   // true for the STANDARD.md-compliant registry
}

// RegistryHome returns the local clone directory for a named registry.
// All git-hosted registries live under RegistriesRoot().
// Local registries (absolute paths) are returned as-is.
func RegistryHome(name string) string {
	if filepath.IsAbs(name) {
		return name
	}
	return filepath.Join(RegistriesRoot(), name)
}

// AllRegistries returns all configured registries, priority-ordered.
//
// Uses [[registry]] entries from global settings, sorted by priority descending.
// When no [[registry]] is configured, returns an implicit official registry entry
// so fresh installs (no settings file) still work.
func AllRegistries() []RegistryEntry {
	cfg, _ := settings.LoadGlobal()

	if len(cfg.Registries) == 0 {
		return []RegistryEntry{{
			Name:     OfficialRegistryName,
			Home:     OfficialRegistryHome(),
			Priority: 100,
			Official: true,
		}}
	}

	sorted := make([]settings.RegistryDef, len(cfg.Registries))
	copy(sorted, cfg.Registries)
	sort.Slice(sorted, func(i, j int) bool {
		return effectivePriority(sorted[i]) > effectivePriority(sorted[j])
	})

	var all []RegistryEntry
	for _, rd := range sorted {
		if !rd.Enabled {
			continue
		}
		all = append(all, RegistryEntry{
			Name:     rd.Name,
			Home:     registryDefHome(rd),
			Priority: effectivePriority(rd),
			Official: rd.Official,
		})
	}
	return all
}

// registryDefHome resolves the local clone directory for a RegistryDef.
// For absolute-path URLs, returns the path directly.
// For git URLs and shorthands, uses registries/<name>/.
func registryDefHome(rd settings.RegistryDef) string {
	u, _ := settings.ParseRef(rd.URL)
	if u == "" {
		u = rd.URL
	}
	if filepath.IsAbs(u) {
		return u
	}
	return filepath.Join(RegistriesRoot(), rd.Name)
}

// effectivePriority returns rd.Priority, defaulting to 100 for official and 50 for user registries.
func effectivePriority(rd settings.RegistryDef) int {
	if rd.Priority > 0 {
		return rd.Priority
	}
	if rd.Official {
		return 100
	}
	return 50
}


// AllSkillsSources returns all configured registries' skills roots in priority order.
// Derived from AllRegistries(); skips entries with no skills/ dir.
func AllSkillsSources() []SkillsSource {
	var sources []SkillsSource
	for _, reg := range AllRegistries() {
		skillsRoot := filepath.Join(reg.Home, "skills")
		if _, err := os.Stat(skillsRoot); err == nil {
			sources = append(sources, SkillsSource{Name: reg.Name, Root: skillsRoot})
		}
	}
	return sources
}

// ListAllSkillsFromSources lists skills from all sources, tagging each with its registry.
// First source to provide a given skill name wins (official has highest priority).
func ListAllSkillsFromSources(sources []SkillsSource) ([]Skill, error) {
	seen := make(map[string]struct{}) // skill name → already included
	var all []Skill

	for _, src := range sources {
		skills, err := ListAllSkills(src.Root)
		if err != nil {
			continue
		}
		for _, sk := range skills {
			if _, exists := seen[sk.Name]; exists {
				continue // higher-priority registry already provided this skill
			}
			seen[sk.Name] = struct{}{}
			sk.Registry = src.Name
			all = append(all, sk)
		}
	}

	return all, nil
}
