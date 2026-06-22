package profiles

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"github.com/jeffreytse/grimoire/internal/skills"
)

// ResolveOptions controls optional behaviour for profile resolution.
type ResolveOptions struct {
	// Registries enables tag-query fallback and tag-field resolution.
	// When nil, tag queries are skipped.
	Registries []skills.SkillsRegistry
	// InlineProfiles holds profiles defined inline in settings.toml [profiles.*].
	// Checked after file lookup fails and before tag/domain fallback.
	InlineProfiles map[string]Profile
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
	Source      string     // absolute path of the file; "" if not found; "(tag query)" if tag-only; "(settings.toml)" if inline
	// Compliance recommendations carried by the profile (0/-1 = not set).
	ComplianceThreshold      float64
	ComplianceThresholdError int // -1 = not set
}

// ParseProfileRef resolves a qualified profile ref against known registry names.
// Qualified form: "<registry>/<profile>" where registry is any known registry name.
// Longest match wins (handles overlapping prefixes across hosts).
// Returns ("", ref) for unqualified refs — caller uses first-wins resolution.
func ParseProfileRef(ref string, registryNames []string) (registry, name string) {
	sorted := append([]string(nil), registryNames...)
	sort.Slice(sorted, func(i, j int) bool { return len(sorted[i]) > len(sorted[j]) })
	for _, reg := range sorted {
		prefix := reg + "/"
		if strings.HasPrefix(ref, prefix) {
			return reg, ref[len(prefix):]
		}
	}
	return "", ref
}

// allRegistryNames returns the Name field of every configured registry.
func allRegistryNames() []string {
	regs := skills.AllRegistries()
	names := make([]string, len(regs))
	for i, r := range regs {
		names[i] = r.Name
	}
	return names
}

// SearchPaths returns the ordered list of paths searched for a named profile ref.
// Qualified ref (e.g. "acmecorp/standards/engineering", "official/engineering",
// "gitlab.com/org/repo/name") — resolved against known registry names, returns a
// single path in that registry only.
// Unqualified ref — covers project-local, then each registry's profiles dir.
func SearchPaths(ref, projectDir string) []string {
	reg, name := ParseProfileRef(ref, allRegistryNames())
	if reg != "" {
		return []string{filepath.Join(skills.RegistryHome(reg), "profiles", name+".toml")}
	}
	// Unqualified: project-local first, then all registries, then default fallback.
	paths := []string{
		filepath.Join(projectDir, ".grimoire", "profiles", name+".toml"),
	}
	for _, r := range skills.AllRegistries() {
		paths = append(paths, filepath.Join(r.Home, "profiles", name+".toml"))
	}
	paths = append(paths, filepath.Join(projectDir, ".grimoire", "profiles", "default.toml"))
	for _, r := range skills.AllRegistries() {
		paths = append(paths, filepath.Join(r.Home, "profiles", "default.toml"))
	}
	return paths
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

// resolveByDomain returns all skills from sources whose domain path matches domain.
func resolveByDomain(domain string, regs []skills.SkillsRegistry) []SkillRef {
	var refs []SkillRef
	seen := make(map[string]struct{})
	for _, reg := range regs {
		all, err := skills.ListAllSkills(reg.Root)
		if err != nil {
			continue
		}
		for _, sk := range all {
			if sk.Domain == domain {
				if _, ok := seen[sk.Name]; !ok {
					seen[sk.Name] = struct{}{}
					refs = append(refs, SkillRef{Name: sk.Name})
				}
			}
		}
	}
	return refs
}

// ResolveByTags returns all skills from sources whose tags contain profileName.
func ResolveByTags(profileName string, regs []skills.SkillsRegistry) []SkillRef {
	var refs []SkillRef
	for _, reg := range regs {
		all, err := skills.ListAllSkills(reg.Root)
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
func ResolveSkills(p *Profile, projectDir string, regs []skills.SkillsRegistry, visited map[string]bool) []SkillRef {
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
			for _, ref := range ResolveSkills(&parent, projectDir, regs, visited) {
				insert(ref)
			}
		}
	}

	// Layer 2: tags — try frontmatter tag match, fall back to domain-path match
	for _, tag := range p.Tags {
		tagRefs := ResolveByTags(tag, regs)
		if len(tagRefs) == 0 {
			tagRefs = resolveByDomain(tag, regs)
		}
		for _, ref := range tagRefs {
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

	// If no profile file, check inline settings.toml definition first.
	if p.Source == "" {
		if inline, ok := opts.InlineProfiles[name]; ok {
			p = inline
		}
	}

	// If still no source, try tag query then domain-path fallback.
	if p.Source == "" && len(opts.Registries) > 0 {
		refs := ResolveByTags(name, opts.Registries)
		src := "(tag query)"
		if len(refs) == 0 {
			refs = resolveByDomain(name, opts.Registries)
			src = "(domain)"
		}
		if len(refs) > 0 {
			p.Skills = refs
			p.Source = src
			return p, nil
		}
	}

	// Full resolution via ResolveSkills (handles extends, tags, explicit skills, exclude).
	if p.Source != "" || len(p.Extends) > 0 || len(p.Tags) > 0 {
		p.Skills = ResolveSkills(&p, projectDir, opts.Registries, visited)
	}

	return p, nil
}

// ProfileMeta holds lightweight metadata for ranking profiles without full resolution.
type ProfileMeta struct {
	Extends []string
	Tags    []string
}

// ReadMeta parses only extends and tags from a profile TOML file — cheap, no skill resolution.
func ReadMeta(path string) ProfileMeta {
	data, err := os.ReadFile(path)
	if err != nil {
		return ProfileMeta{}
	}
	var raw struct {
		Extends []string `toml:"extends"`
		Tags    []string `toml:"tags"`
	}
	_ = toml.Unmarshal(data, &raw)
	return ProfileMeta{Extends: raw.Extends, Tags: raw.Tags}
}

func parseFile(path, profileName string) (Profile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Profile{}, err
	}

	var raw struct {
		Name                     string     `toml:"name"`
		Description              string     `toml:"description"`
		Extends                  []string   `toml:"extends"`
		Tags                     []string   `toml:"tags"`
		Skills                   []SkillRef `toml:"skills"`
		Exclude                  []string   `toml:"exclude"`
		ComplianceThreshold      float64    `toml:"compliance-threshold"`
		ComplianceThresholdError *int       `toml:"compliance-threshold-error"` // pointer: nil = absent
	}
	if err := toml.Unmarshal(data, &raw); err != nil {
		return Profile{}, err
	}

	name := raw.Name
	if name == "" {
		name = profileName
	}
	maxErr := -1 // sentinel: not set
	if raw.ComplianceThresholdError != nil {
		maxErr = *raw.ComplianceThresholdError
	}
	return Profile{
		Name:                     name,
		Description:              raw.Description,
		Extends:                  raw.Extends,
		Tags:                     raw.Tags,
		Skills:                   raw.Skills,
		Exclude:                  raw.Exclude,
		Source:                   path,
		ComplianceThreshold:      raw.ComplianceThreshold,
		ComplianceThresholdError: maxErr,
	}, nil
}
