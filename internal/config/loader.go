package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// ParseFile reads one settings.toml file.
// Returns a zero-value FileConfig when the file is absent — callers treat missing as defaults.
func ParseFile(path string) (FileConfig, error) {
	fs := FileConfig{Sections: make(map[string]DomainSection)}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return fs, nil
	}
	if err != nil {
		return fs, err
	}
	var raw map[string]any
	if err := toml.Unmarshal(data, &raw); err != nil {
		return fs, err
	}
	return parseRaw(raw), nil
}

// Load reads all three file layers for projectDir in priority order and merges them.
// Layers (highest → lowest): grimoire.toml → ~/.config/grimoire/grimoire.toml → /etc/grimoire/grimoire.toml.
func Load(projectDir string) (Config, error) {
	layerPaths := []string{
		ProjectPath(projectDir),
		GlobalPath(),
		SystemPath(),
	}

	layers := make([]FileConfig, 0, len(layerPaths))
	paths := make([]string, 0, len(layerPaths))
	var layerErrors []error

	for _, p := range layerPaths {
		fs, err := ParseFile(p)
		if err != nil {
			layerErrors = append(layerErrors, err)
			continue
		}
		layers = append(layers, fs)
		paths = append(paths, p)
	}

	// Fail only if every layer failed (nothing at all could be loaded).
	if len(layers) == 0 && len(layerErrors) > 0 {
		return Config{}, layerErrors[0]
	}

	// [core] and [[package]] are user-level; strip from the local project layer.
	// Core.Profiles is NOT cleared — it comes from [standards], not [core].
	if len(layers) > 0 {
		layers[0].Core.Home = ""
		layers[0].Packages = nil
	}

	// Expand [dependencies] skills into a synthetic lowest-priority PackageDef layer.
	// Each package ref is parsed into a PackageDef for the package it references.
	// Official/local refs and duplicates are skipped.
	// Using a synthetic appended layer ensures [dependencies] entries never override
	// user's [[package]] entries (Merge dedupes by name, first-seen wins).
	{
		seen := make(map[string]bool)
		var syntheticPkgs []PackageDef
		for i := range layers {
			for _, dep := range layers[i].Dependencies.Skills {
				ref := ParsePackageRef(dep)
				if ref.IsLocal() || ref.IsOfficialRepoPath() || ref.PackageName == "" {
					continue
				}
				if seen[ref.PackageName] {
					continue
				}
				seen[ref.PackageName] = true
				syntheticPkgs = append(syntheticPkgs, PackageDef{
					Name:     ref.PackageName,
					URL:      ref.PackageURL,
					Enabled:  true,
					Priority: 50,
				})
			}
		}
		if len(syntheticPkgs) > 0 {
			layers = append(layers, FileConfig{
				Sections:       make(map[string]DomainSection),
				InlineProfiles: make(map[string]InlineProfileDef),
				Packages:       syntheticPkgs,
			})
			paths = append(paths, "<dependencies>")
		}
	}

	// Load grimoire.toml from non-official [[package]] entries as base layers.
	// Each package can ship a grimoire.toml with practices/profiles for automatic inheritance.
	seenExt := make(map[string]bool)
	for i := range layers {
		layer := &layers[i]
		for _, rd := range layer.Packages {
			if rd.Official || !rd.Enabled || rd.Name == "" {
				continue
			}
			if seenExt[rd.Name] {
				continue
			}
			seenExt[rd.Name] = true
			regHome := packageDefHome(rd)
			rf, err := ParseFile(filepath.Join(regHome, "grimoire.toml"))
			if err != nil {
				continue
			}
			rf.Core = CoreSection{}
			rf.Packages = nil
			rf.InlineProfiles = nil
			layers = append(layers, rf)
			paths = append(paths, filepath.Join(regHome, "grimoire.toml"))
		}
	}

	// Build package name → home dir map for extends resolution.
	regHomes := make(map[string]string)
	for i := range layers {
		for _, rd := range layers[i].Packages {
			if rd.Name != "" {
				if _, seen := regHomes[rd.Name]; !seen {
					regHomes[rd.Name] = packageDefHome(rd)
				}
			}
		}
	}

	// Load standards.extends targets as additional base layers (lowest priority).
	seenExtends := make(map[string]bool)
	var missingExtends []string
	seenMissing := make(map[string]bool)
	for i := range layers {
		for _, ref := range layers[i].Extends {
			p, err := resolveExtendsRef(ref, regHomes)
			if err != nil {
				if !seenMissing[ref] {
					seenMissing[ref] = true
					missingExtends = append(missingExtends, ref)
				}
				continue
			}
			if seenExtends[p] {
				continue
			}
			seenExtends[p] = true
			ef, err := ParseFile(p)
			if err != nil {
				continue
			}
			ef.Core = CoreSection{}
			ef.Packages = nil
			ef.Extends = nil // no recursive extends in v1
			ef.InlineProfiles = nil
			layers = append(layers, ef)
			paths = append(paths, p)
		}
	}

	r := Merge(layers, paths)
	r.MissingExtends = missingExtends

	// GRIMOIRE_HOME env var override — highest priority.
	if v := os.Getenv("GRIMOIRE_HOME"); v != "" {
		r.Core.Home = v
		r.Sources["core.home"] = "$GRIMOIRE_HOME"
	}

	return r, nil
}

// packageDefHome returns the local clone directory for a package definition.
// Path uses the versioned derived name: packages/<host>/<owner>/<repo>@<version>.
func packageDefHome(rd PackageDef) string {
	u, ver := ParseRef(rd.URL)
	if u == "" {
		u = rd.URL
	}
	if filepath.IsAbs(u) {
		return u
	}
	versionedName := filepath.FromSlash(DeriveVersionedName(u, ver))
	if h := os.Getenv("GRIMOIRE_HOME"); h != "" {
		return filepath.Join(h, "packages", versionedName)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".grimoire", "packages", versionedName)
}

// resolveExtendsRef maps a standards.extends ref to a grimoire.toml path.
// Format: "<package-name>" → package root grimoire.toml
//
//	"<package-name>/<preset-name>" → preset grimoire.toml
//
// Longest installed package name wins when names share a prefix.
func resolveExtendsRef(ref string, regHomes map[string]string) (string, error) {
	best := ""
	for name := range regHomes {
		if (ref == name || strings.HasPrefix(ref, name+"/")) && len(name) > len(best) {
			best = name
		}
	}
	if best == "" {
		return "", fmt.Errorf("no installed package matches extends ref %q", ref)
	}
	remainder := strings.TrimPrefix(strings.TrimPrefix(ref, best), "/")
	if remainder == "" {
		return filepath.Join(regHomes[best], "grimoire.toml"), nil
	}
	return filepath.Join(regHomes[best], "presets", remainder, "grimoire.toml"), nil
}

// LoadGlobal reads only the global settings file (no project layers).
func LoadGlobal() (FileConfig, error) {
	return ParseFile(GlobalPath())
}

// LoadFile reads a single settings file by path.
// Used by grimoire config get/set/unset to target a specific level.
func LoadFile(path string) (FileConfig, error) {
	return ParseFile(path)
}

// SaveGlobal writes fs to the global settings file, creating parent dirs as needed.
func SaveGlobal(fs FileConfig) error { //nolint:gocritic // value semantics intentional for config snapshot
	return WriteFile(GlobalPath(), fs)
}

// WriteFile serializes fs to a settings.toml file, creating parent dirs as needed.
func WriteFile(path string, fs FileConfig) error { //nolint:gocritic // value semantics intentional for config snapshot
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	m := toMap(&fs)
	data, err := toml.Marshal(m)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// parseRaw converts a map[string]any (from toml.Unmarshal) into FileConfig.
func parseRaw(raw map[string]any) FileConfig {
	fs := FileConfig{
		Sections:       make(map[string]DomainSection),
		InlineProfiles: make(map[string]InlineProfileDef),
	}
	for key, val := range raw {
		switch key {
		case "core":
			if m, ok := val.(map[string]any); ok {
				fs.Core = parseCoreSection(m)
			}
		case "package":
			// [[package]] is a TOML table array → []any of map[string]any
			if arr, ok := val.([]any); ok {
				fs.Packages = parsePackageDefs(arr)
			}
		case "registry":
			// [[package]] — backward-compat alias for [[package]]
			if arr, ok := val.([]any); ok && len(fs.Packages) == 0 {
				fs.Packages = parsePackageDefs(arr)
			}
		case "standards":
			if m, ok := val.(map[string]any); ok {
				fs.Extends = parseStringSlice(m["extends"])
				fs.CheckExclude = parseStringSlice(m["exclude"])
				fs.Core.Profiles = append(fs.Core.Profiles, parseProfilesFromMap(m)...)
				fs.ReportPath, _ = m["report-path"].(string) //nolint:revive // wrong type silently skipped
				if v, ok := anyToInt(m["staleness-days"]); ok && v > 0 {
					fs.StalenessDays = v
				}
				for domainName, val := range m {
					if domainName == "profiles" || domainName == "extends" || domainName == "exclude" {
						continue
					}
					if sub, ok := val.(map[string]any); ok {
						parseDomainInto(domainName, sub, &fs)
					}
				}
			}
		case "dependencies":
			if m, ok := val.(map[string]any); ok {
				fs.Dependencies.Skills = parseStringSlice(m["skills"])
			}
		case "profiles":
			if m, ok := val.(map[string]any); ok {
				fs.InlineProfiles = parseInlineProfiles(m)
			}
		}
	}
	return fs
}

// parsePackageDefs parses a [[package]] TOML table array into PackageDef entries.
func parsePackageDefs(arr []any) []PackageDef {
	var result []PackageDef
	for _, item := range arr {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rd := PackageDef{Enabled: true}
		rd.Name, _ = m["name"].(string)       //nolint:revive // wrong type silently skipped
		rd.URL, _ = m["url"].(string)         //nolint:revive // wrong type silently skipped
		rd.Official, _ = m["official"].(bool) //nolint:revive // wrong type silently skipped
		if n, ok := anyToInt(m["priority"]); ok {
			rd.Priority = n
		}
		if v, ok := m["enabled"].(bool); ok {
			rd.Enabled = v
		}
		if rd.Name != "" && rd.URL != "" {
			result = append(result, rd)
		}
	}
	return result
}

// parseInlineProfiles parses the top-level [profiles.*] section into InlineProfileDef entries.
func parseInlineProfiles(m map[string]any) map[string]InlineProfileDef {
	result := make(map[string]InlineProfileDef)
	for key, val := range m {
		sub, ok := val.(map[string]any)
		if !ok {
			continue
		}
		def := InlineProfileDef{ComplianceThresholdError: -1}
		def.Name, _ = sub["name"].(string)               //nolint:revive // wrong type silently skipped
		def.Description, _ = sub["description"].(string) //nolint:revive // wrong type silently skipped
		def.Tags = parseStringSlice(sub["tags"])
		def.Extends = parseStringSlice(sub["extends"])
		def.Exclude = parseStringSlice(sub["exclude"])
		if ct := anyToFloat64(sub["compliance-threshold"]); ct > 0 {
			def.ComplianceThreshold = ct
		}
		if sub["compliance-threshold-error"] != nil {
			if n, ok := anyToInt(sub["compliance-threshold-error"]); ok {
				def.ComplianceThresholdError = n
			}
		}
		if arr, ok := sub["skills"].([]any); ok {
			for _, item := range arr {
				if sm, ok := item.(map[string]any); ok {
					ref := InlineSkillRef{}
					ref.Name, _ = sm["name"].(string) //nolint:revive // wrong type silently skipped
					if n, ok := anyToInt(sm["priority"]); ok {
						ref.Priority = n
					}
					if ref.Name != "" {
						def.Skills = append(def.Skills, ref)
					}
				}
			}
		}
		result[key] = def
	}
	return result
}

// parseStringSlice converts a []any (from TOML unmarshal) to []string.
func parseStringSlice(v any) []string {
	arr, ok := v.([]any)
	if !ok {
		return nil
	}
	var result []string
	for _, item := range arr {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

func parseCoreSection(m map[string]any) CoreSection {
	var cs CoreSection
	cs.Home, _ = m["home"].(string)                //nolint:revive // wrong type silently skipped
	cs.InstallMode, _ = m["install-mode"].(string) //nolint:revive // wrong type silently skipped
	cs.Agents = parseStringSlice(m["agents"])
	cs.CheckExclude = parseStringSlice(m["check-exclude"])
	if v, ok := anyToInt(m["update-concurrency"]); ok && v >= 0 {
		cs.UpdateConcurrency = &v
	}
	return cs
}

func parseProfilesFromMap(m map[string]any) []string {
	arr, ok := m["profiles"].([]any)
	if !ok {
		return nil
	}
	var profiles []string
	for _, p := range arr {
		if s, ok := p.(string); ok {
			profiles = append(profiles, s)
		}
	}
	return profiles
}

// parseDomainInto extracts a DomainSection from m (at prefix) and recurses into sub-maps.
func parseDomainInto(prefix string, m map[string]any, fs *FileConfig) {
	ds := DomainSection{ComplianceThresholdError: -1}

	for k, v := range m {
		switch k {
		case "practices":
			if arr, ok := v.([]any); ok {
				for _, p := range arr {
					if s, ok := p.(string); ok {
						ds.Practices = append(ds.Practices, s)
					}
				}
			}
		case "disabled":
			if arr, ok := v.([]any); ok {
				for _, p := range arr {
					if s, ok := p.(string); ok {
						ds.Disabled = append(ds.Disabled, s)
					}
				}
			}
		case "fallback":
			ds.Fallback, _ = v.(string) //nolint:revive // wrong type silently skipped
		case "compliance-threshold":
			ds.ComplianceThreshold = anyToFloat64(v)
		case "compliance-threshold-error":
			if n, ok := anyToInt(v); ok {
				ds.ComplianceThresholdError = n
			}
		default:
			// nested map = subdomain
			if sub, ok := v.(map[string]any); ok {
				parseDomainInto(prefix+"."+k, sub, fs)
			}
		}
	}

	fs.Sections[prefix] = ds
}

// toMap converts FileConfig to a nested map[string]any for TOML marshaling.
func toMap(fs *FileConfig) map[string]any {
	m := map[string]any{}

	// [dependencies] — written before [[package]] for clarity
	if len(fs.Dependencies.Skills) > 0 {
		m["dependencies"] = map[string]any{
			"skills": fs.Dependencies.Skills,
		}
	}

	// [[package]] table array
	if len(fs.Packages) > 0 {
		var regArr []any
		for _, rd := range fs.Packages {
			rm := map[string]any{
				"name": rd.Name,
				"url":  rd.URL,
			}
			if rd.Official {
				rm["official"] = true
			}
			if rd.Priority > 0 {
				rm["priority"] = rd.Priority
			}
			if !rd.Enabled {
				rm["enabled"] = false
			}
			regArr = append(regArr, rm)
		}
		m["package"] = regArr
	}

	core := map[string]any{}
	if fs.Core.Home != "" {
		core["home"] = fs.Core.Home
	}
	if len(fs.Core.Agents) > 0 {
		core["agents"] = fs.Core.Agents
	}
	if fs.Core.InstallMode != "" {
		core["install-mode"] = fs.Core.InstallMode
	}
	if fs.Core.UpdateConcurrency != nil {
		core["update-concurrency"] = *fs.Core.UpdateConcurrency
	}
	if len(core) > 0 {
		m["core"] = core
	}

	// sort keys so parents are written before children
	keys := make([]string, 0, len(fs.Sections))
	for k := range fs.Sections {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// profiles and domain sections all nest under [standards.*]
	standardsMap := map[string]any{}
	if len(fs.Extends) > 0 {
		standardsMap["extends"] = fs.Extends
	}
	if len(fs.Core.Profiles) > 0 {
		standardsMap["profiles"] = fs.Core.Profiles
	}
	if fs.ReportPath != "" {
		standardsMap["report-path"] = fs.ReportPath
	}
	if fs.StalenessDays > 0 {
		standardsMap["staleness-days"] = fs.StalenessDays
	}
	for _, key := range keys {
		s := fs.Sections[key]
		dm := domainToMap(&s)
		if len(dm) == 0 {
			continue
		}
		parts := strings.SplitN(key, ".", 2)
		if len(parts) == 1 {
			existing, _ := standardsMap[key].(map[string]any) //nolint:revive // nil if absent, handled below
			if existing == nil {
				existing = map[string]any{}
			}
			for k, v := range dm {
				existing[k] = v
			}
			standardsMap[key] = existing
		} else {
			parent, _ := standardsMap[parts[0]].(map[string]any) //nolint:revive // nil if absent, handled below
			if parent == nil {
				parent = map[string]any{}
				standardsMap[parts[0]] = parent
			}
			parent[parts[1]] = dm
		}
	}
	if len(standardsMap) > 0 {
		m["standards"] = standardsMap
	}

	return m
}

func domainToMap(ds *DomainSection) map[string]any {
	m := map[string]any{}
	if len(ds.Practices) > 0 {
		m["practices"] = ds.Practices
	}
	if len(ds.Disabled) > 0 {
		m["disabled"] = ds.Disabled
	}
	if ds.Fallback != "" {
		m["fallback"] = ds.Fallback
	}
	if ds.ComplianceThreshold > 0 {
		m["compliance-threshold"] = ds.ComplianceThreshold
	}
	if ds.ComplianceThresholdError >= 0 {
		m["compliance-threshold-error"] = ds.ComplianceThresholdError
	}
	return m
}

func anyToFloat64(v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int64:
		return float64(x)
	case int:
		return float64(x)
	}
	return 0
}

func anyToInt(v any) (int, bool) {
	switch x := v.(type) {
	case int64:
		return int(x), true
	case int:
		return x, true
	case float64:
		return int(x), true
	}
	return 0, false
}
