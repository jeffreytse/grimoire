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
	Name string // "official", "my-org", ...
	Root string // absolute path to the skills/ directory
}

// RegistriesDir returns the directory where custom registry clones live.
func RegistriesDir() string {
	return filepath.Join(GrimoireHome(), "registries")
}

// RegistryHome returns the local clone directory for a named registry.
// The official registry uses GrimoireHome() for backward compatibility.
// Custom registries live at RegistriesDir()/<name>/.
func RegistryHome(name string) string {
	if name == OfficialRegistryName || name == "" {
		return GrimoireHome()
	}
	return filepath.Join(RegistriesDir(), name)
}

// AllSkillsSources returns all configured registries' skills roots in priority order.
// Official is always first. Registries whose skills/ dir doesn't exist yet are skipped.
func AllSkillsSources() []SkillsSource {
	var sources []SkillsSource

	// official always first
	officialRoot := SkillsRoot()
	if _, err := os.Stat(officialRoot); err == nil {
		sources = append(sources, SkillsSource{Name: OfficialRegistryName, Root: officialRoot})
	}

	// custom registries from global settings, sorted for deterministic order
	fs, _ := settings.LoadGlobal()
	names := make([]string, 0, len(fs.Registries))
	for name := range fs.Registries {
		if name != OfficialRegistryName {
			names = append(names, name)
		}
	}
	sort.Strings(names)

	for _, name := range names {
		regRoot := filepath.Join(RegistryHome(name), "skills")
		if _, err := os.Stat(regRoot); err == nil {
			sources = append(sources, SkillsSource{Name: name, Root: regRoot})
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
