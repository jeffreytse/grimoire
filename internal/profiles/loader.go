package profiles

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"github.com/jeffreytse/grimoire/internal/config"
	"github.com/jeffreytse/grimoire/internal/skills"
)

// ResolveOptions controls optional behaviour for profile resolution.
type ResolveOptions struct {
	// Packages enables tag-query fallback and tag-field resolution.
	// When nil, tag queries are skipped.
	Packages []skills.SkillsPackage
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
	Source      string     // absolute path of the file; "" if not found; "(tag query)" if tag-only; "(grimoire.toml)" if inline
	// Compliance recommendations carried by the profile (0/-1 = not set).
	ComplianceThreshold      float64
	ComplianceThresholdError int // -1 = not set
}

// ParseProfileRef resolves a qualified profile ref against known package names.
// Qualified form: "<package>/<profile>" where package is any known package name.
// Longest match wins (handles overlapping prefixes across hosts).
// Returns ("", ref) for unqualified refs — caller uses first-wins resolution.
func ParseProfileRef(ref string, packageNames []string) (pkgName, name string) {
	sorted := append([]string(nil), packageNames...)
	sort.Slice(sorted, func(i, j int) bool { return len(sorted[i]) > len(sorted[j]) })
	for _, reg := range sorted {
		prefix := reg + "/"
		if strings.HasPrefix(ref, prefix) {
			return reg, ref[len(prefix):]
		}
	}
	return "", ref
}

// allPackageNames returns the Name field of every configured package.
func allPackageNames() []string {
	regs := skills.AllPackages()
	names := make([]string, len(regs))
	for i, r := range regs {
		names[i] = r.Name
	}
	return names
}

// ResolveByGlob finds all profiles whose path matches name using depth-anywhere glob semantics.
// Profile names use implicit ** prefix: "engineering" matches **/engineering.toml anywhere in
// the package tree and all .toml files under any **/engineering/ directory.
// For qualified refs ("acmecorp/engineering"), only the named package is searched.
// For unqualified refs, project .grimoire is checked first, then all packages in priority order.
// Deduplicates by Source path across packages.
func ResolveByGlob(name, projectDir string) ([]Profile, error) {
	// Package URL form: has ':' that is not a Windows drive letter (index 1).
	// e.g. "acmecorp/standards@v0.1.0:engineering", "github.com/acmecorp/standards:engineering"
	if strings.Contains(name, ":") && (len(name) < 2 || name[1] != ':') {
		ref := config.ParsePackageRef(name)
		var root string
		switch {
		case ref.IsLocal():
			root = filepath.FromSlash(ref.LocalPath)
		case ref.IsOfficialRepoPath():
			root = skills.OfficialPackageHome()
		default:
			root = skills.PackageHome(ref.PackageName)
		}
		if root == "" {
			return nil, nil
		}
		return ProfilesMatchingGlob(root, ref.Path)
	}

	reg, localName := ParseProfileRef(name, allPackageNames())
	if reg != "" {
		return ProfilesMatchingGlob(skills.PackageHome(reg), localName)
	}
	seen := map[string]bool{}
	seenRel := map[string]bool{}
	var result []Profile
	// Project-local profiles: no rel-path dedup (local always wins, path collisions are intentional).
	if ps, err := ProfilesMatchingGlob(filepath.Join(projectDir, ".grimoire"), name); err == nil {
		for i := range ps {
			p := ps[i]
			if !seen[p.Source] {
				seen[p.Source] = true
				result = append(result, p)
			}
		}
	}
	// Package profiles: dedup by both absolute Source and root-relative path.
	// First (highest-priority) package wins when two packages have the same rel-path.
	for _, r := range skills.AllPackages() {
		ps, err := ProfilesMatchingGlob(r.Home, name)
		if err != nil {
			continue
		}
		for i := range ps {
			p := ps[i]
			rel, _ := filepath.Rel(r.Home, p.Source)
			rel = filepath.ToSlash(rel)
			if seen[p.Source] || seenRel[rel] {
				continue
			}
			seenRel[rel] = true
			seen[p.Source] = true
			result = append(result, p)
		}
	}
	return result, nil
}

// Resolve finds the first profile matching name using depth-anywhere glob semantics.
// Returns a zero-value Profile (Source == "") with no error when no match is found.
// Call ResolveWithOptions to get the fully resolved skill list (extends + tags + exclude applied).
func Resolve(name, projectDir string) (Profile, error) {
	ps, err := ResolveByGlob(name, projectDir)
	if err != nil || len(ps) == 0 {
		return Profile{Name: name}, err
	}
	return ps[0], nil
}

// ResolveAll expands each name via ResolveByGlob and returns the flat deduplicated list.
// A single name may match multiple profiles (e.g. "engineering" activates all .toml under
// any engineering/ directory). Deduplicates by Source path across all names.
func ResolveAll(names []string, projectDir string) ([]Profile, error) {
	seen := map[string]bool{}
	var out []Profile
	for _, name := range names {
		ps, err := ResolveByGlob(name, projectDir)
		if err != nil {
			return nil, err
		}
		for i := range ps {
			p := ps[i]
			if !seen[p.Source] {
				seen[p.Source] = true
				out = append(out, p)
			}
		}
	}
	return out, nil
}

// MergeProfiles merges the skill lists of multiple profiles into one deduplicated list.
// Profiles must have their Skills field already populated via ResolveSkills or ResolveWithOptions.
// First listed profile wins when the same skill name appears in more than one profile —
// reflecting declaration order priority in standards.profiles.
func MergeProfiles(ps []Profile) []SkillRef {
	seen := map[string]bool{}
	var result []SkillRef
	for i := range ps {
		p := ps[i]
		for _, sk := range p.Skills {
			if !seen[sk.Name] {
				seen[sk.Name] = true
				result = append(result, sk)
			}
		}
	}
	return result
}

// ResolveEffectiveSkills resolves the effective ordered skill list for a set of profile refs.
// Pipeline:
//  1. Resolve all profiles (depth-anywhere glob, deduplicated by Source and rel-path).
//  2. Fully resolve each profile (extends + tags + exclude) via ResolveSkills.
//  3. Merge: first listed profile wins for same skill name.
//  4. Practices override: re-priorities matching skills in practices order; adds missing ones.
//  5. Disabled: removes skills unconditionally, even if listed in practices.
func ResolveEffectiveSkills(names []string, projectDir string, practices, disabled []string) ([]SkillRef, error) {
	rawProfiles, err := ResolveAll(names, projectDir)
	if err != nil {
		return nil, err
	}

	for i := range rawProfiles {
		rawProfiles[i].Skills = ResolveSkills(&rawProfiles[i], projectDir, nil, nil)
	}

	merged := MergeProfiles(rawProfiles)

	if len(practices) > 0 {
		byName := make(map[string]*SkillRef, len(merged))
		for i := range merged {
			byName[merged[i].Name] = &merged[i]
		}
		for i, name := range practices {
			pri := i + 1
			if sk, ok := byName[name]; ok {
				sk.Priority = pri
			} else {
				merged = append(merged, SkillRef{Name: name, Priority: pri})
				byName[name] = &merged[len(merged)-1]
			}
		}
		sort.SliceStable(merged, func(i, j int) bool {
			return merged[i].Priority < merged[j].Priority
		})
	}

	if len(disabled) > 0 {
		disSet := make(map[string]bool, len(disabled))
		for _, name := range disabled {
			disSet[name] = true
		}
		out := merged[:0]
		for _, sk := range merged {
			if !disSet[sk.Name] {
				out = append(out, sk)
			}
		}
		merged = out
	}

	return merged, nil
}

// resolveByDomain returns all skills from sources whose domain path matches domain.
func resolveByDomain(domain string, regs []skills.SkillsPackage) []SkillRef {
	var refs []SkillRef
	seen := make(map[string]struct{})
	for _, reg := range regs {
		all, err := skills.ListAllSkills(reg.Root)
		if err != nil {
			continue
		}
		for i := range all {
			sk := all[i]
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
func ResolveByTags(profileName string, regs []skills.SkillsPackage) []SkillRef {
	var refs []SkillRef
	for _, reg := range regs {
		all, err := skills.ListAllSkills(reg.Root)
		if err != nil {
			continue
		}
		for i := range all {
			sk := all[i]
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
func ResolveSkills(p *Profile, projectDir string, regs []skills.SkillsPackage, visited map[string]bool) []SkillRef {
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

	// If no profile file, check inline grimoire.toml definition first.
	if p.Source == "" {
		if inline, ok := opts.InlineProfiles[name]; ok {
			p = inline
		}
	}

	// If still no source, try tag query then domain-path fallback.
	if p.Source == "" && len(opts.Packages) > 0 {
		refs := ResolveByTags(name, opts.Packages)
		src := "(tag query)"
		if len(refs) == 0 {
			refs = resolveByDomain(name, opts.Packages)
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
		p.Skills = ResolveSkills(&p, projectDir, opts.Packages, visited)
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

// toProfileGlob auto-prepends "**/" to a pattern that has no "**", enabling depth-anywhere
// matching. Patterns that already contain "**" are used as-is. Empty pattern is unchanged.
func toProfileGlob(p string) string {
	if p == "" || strings.Contains(p, "**") {
		return p
	}
	return "**/" + p
}

// ProfilesMatchingGlob finds all profiles under root whose root-relative path matches glob.
// Profile files are identified by the .toml extension (inferred — pattern must omit it).
// A pattern matches either the file itself or any file under a matched directory.
// Patterns without "**" are automatically treated as depth-anywhere ("engineering" → "**/engineering").
// glob="" returns all profiles.
func ProfilesMatchingGlob(root, glob string) ([]Profile, error) {
	effectiveGlob := toProfileGlob(glob)
	var result []Profile
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if filepath.Ext(p) != ".toml" {
			return nil
		}
		rel, relErr := filepath.Rel(root, p)
		if relErr != nil {
			return nil
		}
		rel = filepath.ToSlash(rel)
		relNoExt := strings.TrimSuffix(rel, ".toml")
		dirPart := filepath.ToSlash(filepath.Dir(rel))
		if dirPart == "." {
			dirPart = ""
		}

		if glob == "" || skills.GlobMatch(effectiveGlob, relNoExt) || (dirPart != "" && skills.GlobMatch(effectiveGlob, dirPart)) {
			name := filepath.Base(relNoExt)
			prof, parseErr := parseFile(p, name)
			if parseErr == nil {
				result = append(result, prof)
			}
		}
		return nil
	})
	return result, err
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
