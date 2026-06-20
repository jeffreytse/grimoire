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
// Layers (highest → lowest): .grimoire/settings.toml → ~/.config/grimoire/settings.toml → /etc/grimoire/settings.toml
func Load(projectDir string) (Resolved, error) {
	type entry struct {
		path string
		fs   FileSettings
	}

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

	r := Merge(layers, paths)

	// Env var overrides — highest priority, above all file layers.
	// Useful in CI/CD where writing config files is impractical.
	if v := os.Getenv("GRIMOIRE_HOME"); v != "" {
		r.Core.Home = v
		r.Sources["core.home"] = "$GRIMOIRE_HOME"
	}
	if v := os.Getenv("GRIMOIRE_SOURCE"); v != "" {
		r.Core.Source = v
		r.Sources["core.source"] = "$GRIMOIRE_SOURCE"
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
		Sections:   make(map[string]DomainSection),
		Registries: make(map[string]RegistryConfig),
	}
	for key, val := range raw {
		m, ok := val.(map[string]any)
		if !ok {
			continue
		}
		switch key {
		case "core":
			fs.Core = parseCoreSection(m)
		case "registry":
			fs.Registries = parseRegistries(m)
		case "standards":
			fs.Core.Profiles = append(fs.Core.Profiles, parseProfilesFromMap(m)...)
			for domainName, val := range m {
				if domainName == "profiles" {
					continue
				}
				if sub, ok := val.(map[string]any); ok {
					parseDomainInto(domainName, sub, &fs)
				}
			}
		// unknown top-level keys silently ignored
		}
	}
	return fs
}

func parseRegistries(m map[string]any) map[string]RegistryConfig {
	result := make(map[string]RegistryConfig)
	for name, val := range m {
		sub, ok := val.(map[string]any)
		if !ok {
			continue
		}
		rc := RegistryConfig{}
		rc.URL, _ = sub["url"].(string)
		if rc.URL != "" {
			result[name] = rc
		}
	}
	return result
}

func parseCoreSection(m map[string]any) CoreSection {
	var cs CoreSection
	cs.Home, _ = m["home"].(string)
	cs.Source, _ = m["source"].(string)
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
			ds.Fallback, _ = v.(string)
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

	// registry sections
	if len(fs.Registries) > 0 {
		regMap := map[string]any{}
		for name, rc := range fs.Registries {
			if rc.URL != "" {
				regMap[name] = map[string]any{"url": rc.URL}
			}
		}
		if len(regMap) > 0 {
			m["registry"] = regMap
		}
	}

	core := map[string]any{}
	if fs.Core.Home != "" {
		core["home"] = fs.Core.Home
	}
	if fs.Core.Source != "" {
		core["source"] = fs.Core.Source
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
	if len(fs.Core.Profiles) > 0 {
		standardsMap["profiles"] = fs.Core.Profiles
	}
	for _, key := range keys {
		dm := domainToMap(fs.Sections[key])
		if len(dm) == 0 {
			continue
		}
		parts := strings.SplitN(key, ".", 2)
		if len(parts) == 1 {
			existing, _ := standardsMap[key].(map[string]any)
			if existing == nil {
				existing = map[string]any{}
			}
			for k, v := range dm {
				existing[k] = v
			}
			standardsMap[key] = existing
		} else {
			parent, _ := standardsMap[parts[0]].(map[string]any)
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

func domainToMap(ds DomainSection) map[string]any {
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
