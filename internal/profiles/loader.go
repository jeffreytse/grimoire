package profiles

import (
	"errors"
	"os"
	"path/filepath"
	"sort"

	"github.com/pelletier/go-toml/v2"

	"github.com/jeffreytse/grimoire/internal/skills"
)

// ResolveOptions controls optional behaviour for profile resolution.
type ResolveOptions struct {
	// Sources enables tag-query fallback and tag-field resolution.
	// When nil, tag queries are skipped.
	Sources []skills.SkillsSource
}

// SkillRef is a skill entry inside a profile file or resolved list.
type SkillRef struct {
	Name     string `toml:"name"`
	Priority int    `toml:"priority"` // 0 = unset → treated as 50 during resolution
}

// Profile is a resolved profile — either from a file or an empty value when no file exists.
type Profile struct {
	Name        string
	Description string
	Extends     []string   // profile names to inherit from
	Tags        []string   // skill tags to bulk-activate
	Skills      []SkillRef // after Resolve: explicit [[skills]]; after ResolveWithOptions: fully resolved
	Exclude     []string   // skill names to remove after all inclusions
	Source      string     // absolute path of the file; "" if not found; "(tag query)" if tag-only
}

// SearchPaths returns the ordered list of paths searched for a named profile.
func SearchPaths(name, projectDir string) []string {
	grimoireHome := skills.GrimoireHome()
	return []string{
		filepath.Join(projectDir, ".grimoire", "profiles", name+".toml"),
		filepath.Join(grimoireHome, "profiles", name+".toml"),
		filepath.Join(projectDir, ".grimoire", "profiles", "default.toml"),
		filepath.Join(grimoireHome, "profiles", "default.toml"),
	}
}

// Resolve finds and parses a named profile file following the resolution order from docs/profiles.md.
// The returned Profile.Skills contains only the explicit [[skills]] from the file.
// Call ResolveWithOptions to get the fully resolved skill list (extends + tags + exclude applied).
// Returns a zero-value Profile (Source == "") with no error when no file is found.
func Resolve(name, projectDir string) (Profile, error) {
	for _, path := range SearchPaths(name, projectDir) {
		p, err := parseFile(path, name)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return Profile{}, err
		}
		return p, nil
	}
	return Profile{Name: name}, nil
}

// ResolveAll loads all named profiles and returns them in declaration order.
// Each profile's Skills field contains only explicit [[skills]] entries.
// Use ResolveWithOptions per profile for full resolution.
func ResolveAll(names []string, projectDir string) ([]Profile, error) {
	out := make([]Profile, 0, len(names))
	for _, name := range names {
		p, err := Resolve(name, projectDir)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

// ResolveByTags returns all skills from sources whose tags contain profileName.
func ResolveByTags(profileName string, sources []skills.SkillsSource) []SkillRef {
	var refs []SkillRef
	for _, src := range sources {
		all, err := skills.ListAllSkills(src.Root)
		if err != nil {
			continue
		}
		for _, sk := range all {
			for _, tag := range sk.Tags {
				if tag == profileName {
					refs = append(refs, SkillRef{Name: sk.Name})
					break
				}
			}
		}
	}
	return refs
}

// ResolveSkills resolves the full skill list for p by applying:
//  1. extends — inherit skills from named profiles (recursive, cycle-safe)
//  2. tags — bulk-activate all skills matching any tag in p.Tags
//  3. [[skills]] — explicit additions; override priority of already-collected skills
//  4. exclude — remove named skills from the final list
//
// Results are sorted by priority ascending (lower = higher priority), then insertion order.
// visited tracks the current resolution stack to prevent infinite recursion.
func ResolveSkills(p *Profile, projectDir string, sources []skills.SkillsSource, visited map[string]bool) []SkillRef {
	if visited == nil {
		visited = make(map[string]bool)
	}

	type indexed struct {
		ref   SkillRef
		order int
	}
	byName := make(map[string]indexed)
	counter := 0

	insert := func(ref SkillRef) {
		eff := ref
		if eff.Priority == 0 {
			eff.Priority = 50
		}
		if _, exists := byName[eff.Name]; !exists {
			byName[eff.Name] = indexed{ref: eff, order: counter}
			counter++
		}
	}

	insertOrOverride := func(ref SkillRef) {
		eff := ref
		if eff.Priority == 0 {
			eff.Priority = 50
		}
		if existing, exists := byName[eff.Name]; exists {
			byName[eff.Name] = indexed{ref: eff, order: existing.order}
		} else {
			byName[eff.Name] = indexed{ref: eff, order: counter}
			counter++
		}
	}

	// Layer 1: extends
	for _, parentName := range p.Extends {
		if visited[parentName] {
			continue
		}
		visited[parentName] = true
		parent, err := Resolve(parentName, projectDir)
		if err == nil {
			for _, ref := range ResolveSkills(&parent, projectDir, sources, visited) {
				insert(ref)
			}
		}
	}

	// Layer 2: tags
	for _, tag := range p.Tags {
		for _, ref := range ResolveByTags(tag, sources) {
			insert(ref)
		}
	}

	// Layer 3: explicit [[skills]] — override priority if already present
	for _, ref := range p.Skills {
		insertOrOverride(ref)
	}

	// Layer 4: exclude
	excludeSet := make(map[string]bool, len(p.Exclude))
	for _, name := range p.Exclude {
		excludeSet[name] = true
	}

	var collected []indexed
	for _, item := range byName {
		if !excludeSet[item.ref.Name] {
			collected = append(collected, item)
		}
	}
	sort.SliceStable(collected, func(i, j int) bool {
		pi, pj := collected[i].ref.Priority, collected[j].ref.Priority
		if pi != pj {
			return pi < pj
		}
		return collected[i].order < collected[j].order
	})

	result := make([]SkillRef, len(collected))
	for i, item := range collected {
		result[i] = item.ref
	}
	return result
}

// ResolveWithOptions resolves a profile with full composition (extends + tags + exclude).
// The returned Profile.Skills contains the fully resolved list.
func ResolveWithOptions(name, projectDir string, opts ResolveOptions) (Profile, error) {
	p, err := Resolve(name, projectDir)
	if err != nil {
		return Profile{}, err
	}

	visited := map[string]bool{name: true}

	// If no profile file, treat the profile name itself as a tag query (backward compat).
	if p.Source == "" && len(opts.Sources) > 0 {
		refs := ResolveByTags(name, opts.Sources)
		if len(refs) > 0 {
			p.Skills = refs
			p.Source = "(tag query)"
			return p, nil
		}
	}

	// Full resolution via ResolveSkills (handles extends, tags, explicit skills, exclude).
	if p.Source != "" || len(p.Extends) > 0 || len(p.Tags) > 0 {
		p.Skills = ResolveSkills(&p, projectDir, opts.Sources, visited)
	}

	return p, nil
}

func parseFile(path, profileName string) (Profile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Profile{}, err
	}

	var raw struct {
		Name        string     `toml:"name"`
		Description string     `toml:"description"`
		Extends     []string   `toml:"extends"`
		Tags        []string   `toml:"tags"`
		Skills      []SkillRef `toml:"skills"`
		Exclude     []string   `toml:"exclude"`
	}
	if err := toml.Unmarshal(data, &raw); err != nil {
		return Profile{}, err
	}

	name := raw.Name
	if name == "" {
		name = profileName
	}
	return Profile{
		Name:        name,
		Description: raw.Description,
		Extends:     raw.Extends,
		Tags:        raw.Tags,
		Skills:      raw.Skills,
		Exclude:     raw.Exclude,
		Source:      path,
	}, nil
}
