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
// Layers: .grimoire/settings.local.toml → .grimoire/settings.toml → ~/.config/grimoire/settings.toml
func Load(projectDir string) (Resolved, error) {
	type entry struct {
		path string
		fs   FileSettings
	}

	layerPaths := []string{
		filepath.Join(projectDir, ".grimoire", "settings.local.toml"),
		filepath.Join(projectDir, ".grimoire", "settings.toml"),
		GlobalPath(),
	}

	layers := make([]FileSettings, 0, len(layerPaths))
	paths := make([]string, 0, len(layerPaths))

	for _, p := range layerPaths {
		fs, err := ParseFile(p)
		if err != nil {
			return Resolved{}, err
		}
		layers = append(layers, fs)
		paths = append(paths, p)
	}

	return Merge(layers, paths), nil
}

// LoadGlobal reads only the global settings file (no project layers).
// Used by grimoire config get/set/unset which operate on the global file only.
func LoadGlobal() (FileSettings, error) {
	return ParseFile(GlobalPath())
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
		default:
			parseDomainInto(key, m, &fs)
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
	if arr, ok := m["profiles"].([]any); ok {
		for _, p := range arr {
			if s, ok := p.(string); ok {
				cs.Profiles = append(cs.Profiles, s)
			}
		}
	}
	return cs
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
	if len(fs.Core.Profiles) > 0 {
		core["profiles"] = fs.Core.Profiles
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

	for _, key := range keys {
		ds := fs.Sections[key]
		dm := domainToMap(ds)
		if len(dm) == 0 {
			continue
		}
		parts := strings.SplitN(key, ".", 2)
		if len(parts) == 1 {
			existing, ok := m[key].(map[string]any)
			if !ok {
				existing = map[string]any{}
			}
			for k, v := range dm {
				existing[k] = v
			}
			m[key] = existing
		} else {
			parent, ok := m[parts[0]].(map[string]any)
			if !ok {
				parent = map[string]any{}
				m[parts[0]] = parent
			}
			parent[parts[1]] = dm
		}
	}

	return m
}

func domainToMap(ds DomainSection) map[string]any {
	m := map[string]any{}
	if len(ds.Practices) > 0 {
		m["practices"] = ds.Practices
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
