package cmd

import (
	"os"
	"path/filepath"

	"github.com/jeffreytse/grimoire/internal/config"
	"github.com/jeffreytse/grimoire/internal/profiles"
	"github.com/jeffreytse/grimoire/internal/skills"
)

// resolveOpts builds ResolveOptions for the given project directory.
// Inline profiles are loaded from all package grimoire.toml files (lower priority)
// then from the project's grimoire.toml (higher priority, overwrites package).
func resolveOpts(projectDir string) profiles.ResolveOptions {
	opts := profiles.ResolveOptions{Packages: skills.AllSkillsPackages()}

	inlines := make(map[string]profiles.Profile)
	// Iterate packages in reverse so higher-priority packages overwrite lower-priority ones,
	// matching the first-match priority semantics of file-based SearchPaths.
	regs := skills.AllPackages()
	for i := len(regs) - 1; i >= 0; i-- {
		regPath := filepath.Join(regs[i].Home, "grimoire.toml")
		if fs, err := config.ParseFile(regPath); err == nil {
			r := config.Merge([]config.FileConfig{fs}, []string{regPath})
			for name, p := range inlineProfilesFromSettings(&r) { //nolint:gocritic // map range copy is unavoidable
				inlines[name] = p
			}
		}
	}
	// Project always wins last
	if r, err := config.Load(projectDir); err == nil {
		for name, p := range inlineProfilesFromSettings(&r) { //nolint:gocritic // map range copy is unavoidable
			inlines[name] = p
		}
	}
	if len(inlines) > 0 {
		opts.InlineProfiles = inlines
	}
	return opts
}

// inlineProfilesFromSettings converts config.Config InlineProfiles to profiles.Profile map.
func inlineProfilesFromSettings(r *config.Config) map[string]profiles.Profile {
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
			Source:                   "(grimoire.toml)",
			ComplianceThreshold:      def.ComplianceThreshold,
			ComplianceThresholdError: def.ComplianceThresholdError,
		}
	}
	return result
}

// findProfilePackage returns the package home dir that ships profiles/<name>.toml.
// Uses AllPackages() so profile-only bundles (no skills/) are included.
func findProfilePackage(profileName string) string {
	for _, reg := range skills.AllPackages() {
		if _, err := os.Stat(filepath.Join(reg.Home, "profiles", profileName+".toml")); err == nil {
			return reg.Home
		}
	}
	return ""
}
