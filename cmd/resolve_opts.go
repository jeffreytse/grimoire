package cmd

import (
	"os"
	"path/filepath"

	"github.com/jeffreytse/grimoire/internal/profiles"
	"github.com/jeffreytse/grimoire/internal/settings"
	"github.com/jeffreytse/grimoire/internal/skills"
)

// resolveOpts builds ResolveOptions for the given project directory.
// Inline profiles are loaded from all registry settings.toml files (lower priority)
// then from the project's .grimoire/settings.toml (higher priority, overwrites registry).
func resolveOpts(projectDir string) profiles.ResolveOptions {
	opts := profiles.ResolveOptions{Sources: skills.AllSkillsSources()}

	inlines := make(map[string]profiles.Profile)
	// Iterate registries in reverse so higher-priority registries overwrite lower-priority ones,
	// matching the first-match priority semantics of file-based SearchPaths.
	regs := skills.AllRegistries()
	for i := len(regs) - 1; i >= 0; i-- {
		regPath := filepath.Join(regs[i].Home, "settings.toml")
		if fs, err := settings.ParseFile(regPath); err == nil {
			r := settings.Merge([]settings.FileSettings{fs}, []string{regPath})
			for name, p := range inlineProfilesFromSettings(&r) { //nolint:gocritic // map range copy is unavoidable
				inlines[name] = p
			}
		}
	}
	// Project always wins last
	if r, err := settings.Load(projectDir); err == nil {
		for name, p := range inlineProfilesFromSettings(&r) { //nolint:gocritic // map range copy is unavoidable
			inlines[name] = p
		}
	}
	if len(inlines) > 0 {
		opts.InlineProfiles = inlines
	}
	return opts
}

// inlineProfilesFromSettings converts settings.Resolved InlineProfiles to profiles.Profile map.
func inlineProfilesFromSettings(r *settings.Resolved) map[string]profiles.Profile {
	if len(r.InlineProfiles) == 0 {
		return nil
	}
	result := make(map[string]profiles.Profile, len(r.InlineProfiles))
	for key, def := range r.InlineProfiles { //nolint:gocritic // map range copy unavoidable
		refs := make([]profiles.SkillRef, len(def.Skills))
		for i := range def.Skills {
			refs[i] = profiles.SkillRef{Name: def.Skills[i].Name, Priority: def.Skills[i].Priority}
		}
		profileName := key
		if def.Name != "" {
			profileName = def.Name
		}
		result[key] = profiles.Profile{
			Name:                     profileName,
			Description:              def.Description,
			Tags:                     def.Tags,
			Extends:                  def.Extends,
			Exclude:                  def.Exclude,
			Skills:                   refs,
			Source:                   "(settings.toml)",
			ComplianceThreshold:      def.ComplianceThreshold,
			ComplianceThresholdError: def.ComplianceThresholdError,
		}
	}
	return result
}

// findProfileRegistry returns the registry home dir that ships profiles/<name>.toml.
// Uses AllRegistries() so profile-only bundles (no skills/) are included.
func findProfileRegistry(profileName string) string {
	for _, reg := range skills.AllRegistries() {
		if _, err := os.Stat(filepath.Join(reg.Home, "profiles", profileName+".toml")); err == nil {
			return reg.Home
		}
	}
	return ""
}
