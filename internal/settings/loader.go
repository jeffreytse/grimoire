package settings

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// ParseFile reads one settings.toml file.
// Returns a zero-value FileSettings when the file is absent — callers treat missing as defaults.
func ParseFile(path string) (FileSettings, error) {
	fs := FileSettings{Sections: make(map[string]DomainSection)}
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
// Layers (highest → lowest): .grimoire/settings.toml → ~/.config/grimoire/settings.toml → /etc/grimoire/settings.toml.
func Load(projectDir string) (Resolved, error) {
	layerPaths := []string{
		filepath.Join(projectDir, ".grimoire", "settings.toml"),
		GlobalPath(),
		SystemPath(),
	}

	layers := make([]FileSettings, 0, len(layerPaths))
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
		return Resolved{}, layerErrors[0]
	}

	// [core] section is user-level; strip it from the local project layer.
	// Core.Profiles is NOT cleared — it comes from [standards], not [core].
	if len(layers) > 0 {
		layers[0].Core.Home = ""
		layers[0].Core.Registry = ""
	}

	// Resolve [standards] extends: load each target's settings.toml as a base layer.
	// Dedup by derived registry name across all layers.
	seenExt := make(map[string]bool)
	for _, layer := range layers {
		for _, ref := range layer.StandardsExtends {
			u, _ := ParseRef(ref)
			name := DeriveRegistryName(u)
			if seenExt[name] {
				continue
			}
			seenExt[name] = true
			regHome := extendsHomePath(name)
			rf, err := ParseFile(filepath.Join(regHome, "settings.toml"))
			if err != nil {
				continue // not cloned yet — skip silently
			}
			// inherit only [standards.*] sections
			rf.Core = CoreSection{}
			rf.StandardsExtends = nil
			rf.InlineProfiles = nil
			layers = append(layers, rf)
			paths = append(paths, filepath.Join(regHome, "settings.toml"))
		}
	}

	r := Merge(layers, paths)

	// GRIMOIRE_HOME env var override — highest priority.
	if v := os.Getenv("GRIMOIRE_HOME"); v != "" {
		r.Core.Home = v
		r.Sources["core.home"] = "$GRIMOIRE_HOME"
	}

	return r, nil
}

// LoadGlobal reads only the global settings file (no project layers).
func LoadGlobal() (FileSettings, error) {
	return ParseFile(GlobalPath())
}

// LoadFile reads a single settings file by path.
// Used by grimoire config get/set/unset to target a specific level.
func LoadFile(path string) (FileSettings, error) {
	return ParseFile(path)
}

// SaveGlobal writes fs to the global settings file, creating parent dirs as needed.
func SaveGlobal(fs FileSettings) error {
	return WriteFile(GlobalPath(), fs)
}

// WriteFile serializes fs to a settings.toml file, creating parent dirs as needed.
func WriteFile(path string, fs FileSettings) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	m := toMap(fs)
	data, err := toml.Marshal(m)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// parseRaw converts a map[string]any (from toml.Unmarshal) into FileSettings.
func parseRaw(raw map[string]any) FileSettings {
	fs := FileSettings{
		Sections:       make(map[string]DomainSection),
		InlineProfiles: make(map[string]InlineProfileDef),
	}
	for key, val := range raw {
		switch key {
		case "core":
			if m, ok := val.(map[string]any); ok {
				fs.Core = parseCoreSection(m)
			}
		case "standards":
			if m, ok := val.(map[string]any); ok {
				fs.StandardsExtends = parseStringSlice(m["extends"])
				fs.Core.Profiles = append(fs.Core.Profiles, parseProfilesFromMap(m)...)
				fs.ReportPath, _ = m["report-path"].(string)
				if v, ok := anyToInt(m["staleness-days"]); ok && v > 0 {
					fs.StalenessDays = v
				}
				for domainName, val := range m {
					if domainName == "profiles" || domainName == "extends" {
						continue
					}
					if sub, ok := val.(map[string]any); ok {
						parseDomainInto(domainName, sub, &fs)
					}
				}
			}
		case "profiles":
			if m, ok := val.(map[string]any); ok {
				fs.InlineProfiles = parseInlineProfiles(m)
			}
		}
	}
	return fs
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
		def.Name, _ = sub["name"].(string)
		def.Description, _ = sub["description"].(string)
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
					ref.Name, _ = sm["name"].(string)
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
	cs.Home, _ = m["home"].(string)               //nolint:revive // wrong type silently skipped
	cs.Registry, _ = m["registry"].(string)        //nolint:revive
	cs.InstallMode, _ = m["install-mode"].(string) //nolint:revive
	cs.Agents = parseStringSlice(m["agents"])
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
func parseDomainInto(prefix string, m map[string]any, fs *FileSettings) {
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

// toMap converts FileSettings to a nested map[string]any for TOML marshaling.
func toMap(fs FileSettings) map[string]any {
	m := map[string]any{}

	core := map[string]any{}
	if fs.Core.Home != "" {
		core["home"] = fs.Core.Home
	}
	if fs.Core.Registry != "" {
		core["registry"] = fs.Core.Registry
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
	if len(fs.StandardsExtends) > 0 {
		standardsMap["extends"] = fs.StandardsExtends
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

// extendsHomePath returns the local clone directory for a standards extends target.
// Inline equivalent of skills.ExtendsHome — avoids circular import (skills imports settings).
func extendsHomePath(name string) string {
	if filepath.IsAbs(name) {
		return name // local registry: path is already the home
	}
	if h := os.Getenv("GRIMOIRE_HOME"); h != "" {
		return filepath.Join(h, "registries", name)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".grimoire", "registries", name)
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
