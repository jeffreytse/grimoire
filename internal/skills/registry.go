package skills

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/jeffreytse/grimoire/internal/settings"
)

const OfficialRegistryName = "official"

// SkillsRegistry is a named registry with its resolved skills directory path.
// Narrower than RegistryEntry: only registries that have a skills/ directory.
type SkillsRegistry struct {
	Name string // derived registry name, e.g. "jeffreytse/grimoire-hub", "acmecorp/standards"
	Root string // absolute path to the skills/ directory
}

// SkillConflict records a skill that was shadowed by a higher-priority registry.
type SkillConflict struct {
	CanonicalPath  string // e.g. "engineering/development/apply-solid"
	WinnerRegistry string // registry whose version will be installed
	LoserRegistry  string // registry whose version was suppressed
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

// AllSkillsRegistries returns all configured registries' skills roots in priority order.
// Derived from AllRegistries(); skips entries with no skills/ dir.
func AllSkillsRegistries() []SkillsRegistry {
	var regs []SkillsRegistry
	for _, reg := range AllRegistries() {
		skillsRoot := filepath.Join(reg.Home, "skills")
		if _, err := os.Stat(skillsRoot); err == nil {
			regs = append(regs, SkillsRegistry{Name: reg.Name, Root: skillsRoot})
		}
	}
	return regs
}

// ListAllSkillsFromRegistries lists skills from all registries, tagging each with its registry.
// First (highest-priority) registry to provide a given canonical path wins.
// Conflicts contains skills from lower-priority registries that were shadowed.
func ListAllSkillsFromRegistries(regs []SkillsRegistry) ([]Skill, []SkillConflict, error) {
	seen := make(map[string]string) // canonical path → registry name that claimed it
	var all []Skill
	var conflicts []SkillConflict

	for _, reg := range regs {
		list, err := ListAllSkills(reg.Root)
		if err != nil {
			continue
		}
		for _, sk := range list {
			dirName := filepath.Base(sk.Path)
			key := sk.Domain + "/" + dirName
			if sk.Subdomain != "" {
				key = sk.Domain + "/" + sk.Subdomain + "/" + dirName
			}
			if winner, exists := seen[key]; exists {
				conflicts = append(conflicts, SkillConflict{
					CanonicalPath:  key,
					WinnerRegistry: winner,
					LoserRegistry:  reg.Name,
				})
				continue
			}
			seen[key] = reg.Name
			sk.Registry = reg.Name
			all = append(all, sk)
		}
	}

	return all, conflicts, nil
}
