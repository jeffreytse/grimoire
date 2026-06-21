package skills

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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
	Name string // derived registry name, e.g. "jeffreytse/grimoire-hub"
	Home string // absolute path to registry root dir (parent of skills/)
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

// AllRegistries returns all cloned registries found under RegistriesRoot().
// Discovery is filesystem-driven: any directory containing at least one of
// skills/, profiles/, presets/, or settings.toml is treated as a registry root.
// The recursive walk handles 2-level (GitHub) and 3-level (non-GitHub) paths.
// The official registry is always first so its skills and profiles take priority
// over community registries in first-wins resolution.
func AllRegistries() []RegistryEntry {
	root := RegistriesRoot()
	var all []RegistryEntry
	walkRegistryDirs(root, root, &all)

	// Include local registries declared in global settings (absolute paths outside RegistriesRoot).
	seen := make(map[string]bool)
	for _, e := range all {
		seen[e.Home] = true
	}
	cfg, _ := settings.LoadGlobal()
	for _, ref := range cfg.StandardsExtends {
		u, _ := settings.ParseRef(ref)
		if !filepath.IsAbs(u) || seen[u] {
			continue
		}
		seen[u] = true
		if isRegistryDir(u) {
			all = append(all, RegistryEntry{Name: u, Home: u})
		}
	}

	officialHome := OfficialRegistryHome()
	var ordered, rest []RegistryEntry
	for _, e := range all {
		if e.Home == officialHome {
			ordered = append(ordered, e)
		} else {
			rest = append(rest, e)
		}
	}
	return append(ordered, rest...)
}

// walkRegistryDirs recursively finds registry roots under cur.
// A directory that looks like a registry is added to out; otherwise its
// children are searched. Symlinks are skipped.
func walkRegistryDirs(root, cur string, out *[]RegistryEntry) {
	dirs, err := os.ReadDir(cur)
	if err != nil {
		return
	}
	for _, d := range dirs {
		if d.Type()&fs.ModeSymlink != 0 {
			continue
		}
		if !d.IsDir() || strings.HasPrefix(d.Name(), ".") {
			continue
		}
		path := filepath.Join(cur, d.Name())
		if isRegistryDir(path) {
			name, _ := filepath.Rel(root, path)
			*out = append(*out, RegistryEntry{Name: filepath.ToSlash(name), Home: path})
		} else {
			walkRegistryDirs(root, path, out)
		}
	}
}

// isRegistryDir reports whether path looks like a registry root.
func isRegistryDir(path string) bool {
	for _, marker := range []string{"skills", "profiles", "presets", "settings.toml"} {
		if _, err := os.Stat(filepath.Join(path, marker)); err == nil {
			return true
		}
	}
	return false
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
