package settings

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// CoreSection holds machine-level and top-level runtime settings.
// These live under the [core] TOML section.
type CoreSection struct {
	Home              string
	Profiles          []string // from [standards] profiles — stored here for convenience
	Registry          string   // [core] registry — single source ref ("owner/repo[@version]" or full URL)
	Agents            []string // [core] agents — pinned agent targets (empty = auto-detect)
	InstallMode       string   // [core] install-mode — "symlink" (default) | "copy"
	UpdateConcurrency *int     // [core] update-concurrency — nil=default(8), 0=unlimited, N=cap at N
}

// DomainSection holds skill practice settings for one domain or subdomain.
// ComplianceThresholdError uses -1 as "unset" sentinel because 0 is meaningful (allow no errors).
type DomainSection struct {
	Practices                []string
	Disabled                 []string // skill names to suppress from spec regardless of profiles
	Fallback                 string   // "ask" | "skip" | ""
	ComplianceThreshold      float64  // 0 = unset
	ComplianceThresholdError int      // -1 = unset; 0 = allow none; N = allow up to N
}

// InlineSkillRef is a skill entry inside an inline profile definition.
type InlineSkillRef struct {
	Name     string
	Priority int
}

// InlineProfileDef is a profile definition embedded inside settings.toml under [profiles.*].
// It mirrors the profile TOML file format and may also carry compliance settings.
type InlineProfileDef struct {
	Name                     string         // optional — overrides map key; defaults to map key
	Description              string
	Tags                     []string
	Extends                  []string
	Skills                   []InlineSkillRef
	Exclude                  []string
	ComplianceThreshold      float64
	ComplianceThresholdError int // -1 = unset
}

// FileSettings is one parsed settings.toml file.
type FileSettings struct {
	Core             CoreSection
	StandardsExtends []string                    // from [standards] extends = [...]
	ReportPath       string                      // from [standards] report-path
	StalenessDays    int                         // from [standards] staleness-days (0 = unset)
	Sections         map[string]DomainSection    // dotted keys: "engineering", "engineering.architecture"
	InlineProfiles   map[string]InlineProfileDef // from [profiles.*]
}

// Resolved holds the effective settings after merging all file layers.
type Resolved struct {
	Core             CoreSection
	StandardsExtends []string                    // deduped union from all file layers
	ReportPath       string                      // first non-empty across layers
	StalenessDays    int                         // first nonzero across layers (default 7 when 0)
	sections         map[string]DomainSection
	InlineProfiles   map[string]InlineProfileDef // merged, higher-priority layers win per name
	// Sources maps dotted key paths to the file that provided them.
	// E.g. "core.home" → "/path/to/settings.toml"
	Sources map[string]string
}

// GlobalPath returns the path to the user-global settings file, respecting XDG_CONFIG_HOME.
func GlobalPath() string {
	if h := os.Getenv("XDG_CONFIG_HOME"); h != "" {
		return filepath.Join(h, "grimoire", "settings.toml")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "grimoire", "settings.toml")
}

// SystemPath returns the system-wide settings file path.
// Linux/macOS: /etc/grimoire/settings.toml
// Windows: %PROGRAMDATA%\grimoire\settings.toml.
func SystemPath() string {
	if runtime.GOOS == "windows" {
		pd := os.Getenv("PROGRAMDATA")
		if pd == "" {
			pd = `C:\ProgramData`
		}
		return filepath.Join(pd, "grimoire", "settings.toml")
	}
	return "/etc/grimoire/settings.toml"
}

var validStandardsFields = map[string]bool{
	"extends":                    true,
	"profiles":                   true,
	"report-path":                true,
	"staleness-days":             true,
	"practices":                  true,
	"disabled":                   true,
	"fallback":                   true,
	"compliance-threshold":       true,
	"compliance-threshold-error": true,
}

// ParseStandardsKey splits a dotted standards key into domain and field.
// "standards.profiles"                    → domain="",                  field="profiles"
// "standards.engineering.practices"       → domain="engineering",       field="practices"
// "standards.engineering.testing.fallback"→ domain="engineering.testing", field="fallback".
func ParseStandardsKey(dotted string) (domain, field string, err error) {
	parts := strings.SplitN(dotted, ".", 2)
	if len(parts) < 2 || parts[0] != "standards" {
		return "", "", fmt.Errorf("key must start with \"standards.\"")
	}
	rest := parts[1]
	// check if rest is a top-level field (no domain)
	if validStandardsFields[rest] {
		return "", rest, nil
	}
	// last segment is field, everything before is domain
	idx := strings.LastIndex(rest, ".")
	if idx < 0 {
		return "", "", fmt.Errorf("unknown standards key %q", dotted)
	}
	field = rest[idx+1:]
	if !validStandardsFields[field] {
		return "", "", fmt.Errorf("unknown standards field %q in key %q", field, dotted)
	}
	return rest[:idx], field, nil
}

// Merge combines layers in priority order — layers[0] wins over layers[1], etc.
// paths must be the same length as layers and contains the source file path for each.
func Merge(layers []FileSettings, paths []string) Resolved {
	r := Resolved{
		sections:       make(map[string]DomainSection),
		InlineProfiles: make(map[string]InlineProfileDef),
		Sources:        make(map[string]string),
	}

	// Core scalar fields: first non-empty wins (layers[0] = highest priority).
	// StandardsExtends: union across all layers, deduped by derived name.
	seenExt := make(map[string]bool)
	for i, fs := range layers {
		src := paths[i]
		if r.Core.Home == "" && fs.Core.Home != "" {
			r.Core.Home = fs.Core.Home
			r.Sources["core.home"] = src
		}
		if r.Core.Registry == "" && fs.Core.Registry != "" {
			r.Core.Registry = fs.Core.Registry
			r.Sources["core.registry"] = src
		}
		if len(r.Core.Agents) == 0 && len(fs.Core.Agents) > 0 {
			r.Core.Agents = fs.Core.Agents
			r.Sources["core.agents"] = src
		}
		if r.Core.InstallMode == "" && fs.Core.InstallMode != "" {
			r.Core.InstallMode = fs.Core.InstallMode
			r.Sources["core.install-mode"] = src
		}
		if r.Core.UpdateConcurrency == nil && fs.Core.UpdateConcurrency != nil {
			r.Core.UpdateConcurrency = fs.Core.UpdateConcurrency
			r.Sources["core.update-concurrency"] = src
		}
		if len(r.Core.Profiles) == 0 && len(fs.Core.Profiles) > 0 {
			r.Core.Profiles = fs.Core.Profiles
			r.Sources["standards.profiles"] = src
		}
		if r.ReportPath == "" && fs.ReportPath != "" {
			r.ReportPath = fs.ReportPath
			r.Sources["standards.report-path"] = src
		}
		if r.StalenessDays == 0 && fs.StalenessDays > 0 {
			r.StalenessDays = fs.StalenessDays
			r.Sources["standards.staleness-days"] = src
		}
		for _, ref := range fs.StandardsExtends {
			u, _ := ParseRef(ref)
			name := DeriveRegistryName(u)
			if !seenExt[name] {
				seenExt[name] = true
				r.StandardsExtends = append(r.StandardsExtends, ref)
				if r.Sources["standards.extends"] == "" {
					r.Sources["standards.extends"] = src
				}
			}
		}
	}

	// collect all section keys across all layers
	allKeys := make(map[string]struct{})
	for _, fs := range layers {
		for k := range fs.Sections {
			allKeys[k] = struct{}{}
		}
	}

	// InlineProfiles: higher-priority layers win per profile name (project > global > system)
	for i := len(layers) - 1; i >= 0; i-- {
		for name, def := range layers[i].InlineProfiles {
			r.InlineProfiles[name] = def
		}
	}

	for key := range allKeys {
		merged := DomainSection{ComplianceThresholdError: -1}
		for i, fs := range layers {
			src := paths[i]
			ds, ok := fs.Sections[key]
			if !ok {
				continue
			}
			if len(merged.Practices) == 0 && len(ds.Practices) > 0 {
				merged.Practices = ds.Practices
				r.Sources[key+".practices"] = src
			}
			if len(merged.Disabled) == 0 && len(ds.Disabled) > 0 {
				merged.Disabled = ds.Disabled
				r.Sources[key+".disabled"] = src
			}
			if merged.Fallback == "" && ds.Fallback != "" {
				merged.Fallback = ds.Fallback
				r.Sources[key+".fallback"] = src
			}
			if merged.ComplianceThreshold == 0 && ds.ComplianceThreshold > 0 {
				merged.ComplianceThreshold = ds.ComplianceThreshold
				r.Sources[key+".compliance-threshold"] = src
			}
			if merged.ComplianceThresholdError == -1 && ds.ComplianceThresholdError != -1 {
				merged.ComplianceThresholdError = ds.ComplianceThresholdError
				r.Sources[key+".compliance-threshold-error"] = src
			}
		}
		r.sections[key] = merged
	}

	return r
}

// SectionKeys returns all domain/subdomain keys present in the resolved settings.
func (r Resolved) SectionKeys() []string {
	keys := make([]string, 0, len(r.sections))
	for k := range r.sections {
		keys = append(keys, k)
	}
	return keys
}

// ResolveSection returns the effective DomainSection for scope (e.g. "engineering.testing").
// Subdomain keys overlay domain keys; unset subdomain keys inherit from the domain.
// Only one level of nesting is supported (domain.subdomain).
func (r Resolved) ResolveSection(scope string) DomainSection {
	parts := strings.SplitN(scope, ".", 2)
	domain := r.sections[parts[0]]
	if len(parts) == 1 {
		return domain
	}

	sub, ok := r.sections[scope]
	if !ok {
		return domain
	}

	// overlay: subdomain wins per-key; domain fills gaps
	result := DomainSection{ComplianceThresholdError: -1}
	if len(sub.Practices) > 0 {
		result.Practices = sub.Practices
	} else {
		result.Practices = domain.Practices
	}
	if len(sub.Disabled) > 0 {
		result.Disabled = sub.Disabled
	} else {
		result.Disabled = domain.Disabled
	}
	if sub.Fallback != "" {
		result.Fallback = sub.Fallback
	} else {
		result.Fallback = domain.Fallback
	}
	if sub.ComplianceThreshold > 0 {
		result.ComplianceThreshold = sub.ComplianceThreshold
	} else {
		result.ComplianceThreshold = domain.ComplianceThreshold
	}
	if sub.ComplianceThresholdError != -1 {
		result.ComplianceThresholdError = sub.ComplianceThresholdError
	} else {
		result.ComplianceThresholdError = domain.ComplianceThresholdError
	}
	return result
}
