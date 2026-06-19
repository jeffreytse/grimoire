package profiles

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"

	"github.com/jeffreytse/grimoire/internal/skills"
)

// ResolveOptions controls optional behaviour for profile resolution.
type ResolveOptions struct {
	// Sources enables tag-query fallback (step 5). When nil, tag query is skipped.
	Sources []skills.SkillsSource
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

// ResolveWithOptions resolves a profile with optional tag-query fallback.
func ResolveWithOptions(name, projectDir string, opts ResolveOptions) (Profile, error) {
	p, err := Resolve(name, projectDir)
	if err != nil {
		return Profile{}, err
	}
	// file found — return as-is
	if p.Source != "" {
		return p, nil
	}
	// tag query fallback
	if len(opts.Sources) > 0 {
		refs := ResolveByTags(name, opts.Sources)
		if len(refs) > 0 {
			p.Skills = refs
			p.Source = "(tag query)"
		}
	}
	return p, nil
}

// SkillRef is a skill entry inside a profile file.
type SkillRef struct {
	Name string `toml:"name"`
}

// Profile is a resolved profile — either from a file or an empty value when no file exists.
type Profile struct {
	Name        string
	Description string
	Skills      []SkillRef
	Source      string // absolute path of the file that was loaded; "" if not found
}

// SearchPaths returns the ordered list of paths searched for a named profile.
// Callers can use this for diagnostics without triggering file I/O.
func SearchPaths(name, projectDir string) []string {
	grimoireHome := skills.GrimoireHome()
	return []string{
		filepath.Join(projectDir, ".grimoire", "profiles", name+".toml"),
		filepath.Join(grimoireHome, "profiles", name+".toml"),
		filepath.Join(projectDir, ".grimoire", "profiles", "default.toml"),
		filepath.Join(grimoireHome, "profiles", "default.toml"),
	}
}

// Resolve finds and loads a named profile following the resolution order from docs/profiles.md.
// Returns a zero-value Profile (Source == "") with no error when no file is found —
// tag-based resolution is the LLM skill layer's responsibility.
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
// Skills are NOT deduplicated here — callers merge as needed.
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

func parseFile(path, profileName string) (Profile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Profile{}, err
	}

	var raw struct {
		Name        string     `toml:"name"`
		Description string     `toml:"description"`
		Skills      []SkillRef `toml:"skills"`
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
		Skills:      raw.Skills,
		Source:      path,
	}, nil
}
