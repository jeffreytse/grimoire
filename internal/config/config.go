package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// LLMProviderConfig configures the API provider used by `grimoire check --independent`.
// Named providers (anthropic, openai, openrouter, grok, ollama, groq) have built-in
// defaults; only override fields that differ. Use name = "custom" for arbitrary endpoints.
type LLMProviderConfig struct {
	Name      string `toml:"name"`        // built-in name ("anthropic","openai","openrouter","grok","ollama","groq") or "custom"
	BaseURL   string `toml:"base-url"`    // empty = use built-in default for named provider
	APIKeyEnv string `toml:"api-key-env"` // env var holding the API key; empty = built-in default or none
	Model     string `toml:"model"`       // empty = use built-in default for named provider
	Format    string `toml:"format"`      // "openai" or "anthropic"; empty = derived from Name
	MaxTokens int    `toml:"max-tokens"`  // 0 = use default (8192)
}

// CoreSection holds machine-level and top-level runtime settings.
// These live under the [core] TOML section.
type CoreSection struct {
	Home              string
	Profiles          []string          // from [standards] profiles — stored here for convenience
	Agents            []string          // [core] agents — pinned agent targets (empty = auto-detect)
	InstallMode       string            // [core] install-mode — "symlink" (default) | "copy"
	UpdateConcurrency *int              // [core] update-concurrency — nil=default(8), 0=unlimited, N=cap at N
	SearchPackages    []string          `toml:"search-packages"` // additional packages for `grimoire search` (global only)
	CheckAgents       []string          `toml:"check-agents"`    // local CLI resolution order for --independent; empty = default
	CheckProvider     LLMProviderConfig `toml:"check-provider"`  // API provider for --independent; zero = auto-detect
	CheckExclude      []string          `toml:"check-exclude"`   // glob patterns for files to exclude from grimoire check
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

// CriterionRef links a stable compliance ID to a body bullet by substring match.
// ID is what appears in the compliance report (e.g. "srp").
// Matches is a substring matched against the extracted ## Criteria or ## Anti-patterns bullet.
type CriterionRef struct {
	ID      string `toml:"id"`
	Matches string `toml:"matches"`
}

// InlineSkillRef is a skill entry inside an inline profile definition.
type InlineSkillRef struct {
	Name         string
	Priority     int
	Criteria     []CriterionRef // optional: stable IDs mapped to ## Criteria bullets
	AntiPatterns []CriterionRef // optional: stable IDs mapped to ## Anti-patterns bullets
}

// PackageDef is one entry in the [[package]] table array.
// It represents a named, priority-ordered skill package source (git URL or local path).
// Multiple packages are searched in priority order (highest first) for skill resolution.
type PackageDef struct {
	Name     string // user-chosen identifier, e.g. "official", "my-team"
	URL      string // git URL, owner/repo[@version] shorthand, or absolute path
	Official bool   // exactly one entry should be official (STANDARD.md-compliant package)
	Priority int    // 0 = unset; normalized in Merge: 100 if Official, else 50
	Enabled  bool   // false = skip in resolution without removing the entry
}

// InlineProfileDef is a profile definition embedded inside grimoire.toml under [profiles.*].
// It mirrors the profile TOML file format and may also carry compliance settings.
type InlineProfileDef struct {
	Name                     string // optional — overrides map key; defaults to map key
	Description              string
	Tags                     []string
	Extends                  []string
	Skills                   []InlineSkillRef
	Exclude                  []string
	ComplianceThreshold      float64
	ComplianceThresholdError int // -1 = unset
}

// DependenciesSection holds the project/global dependency declaration.
// Each entry is a package URL: [host/][owner/repo][@version][:glob_path].
// List order encodes priority: first entry = highest priority.
// Commenting out a line disables it (TOML omits it from parsing).
type DependenciesSection struct {
	Skills []string // ordered package refs: first=highest priority
}

// FileConfig is one parsed grimoire.toml file.
type FileConfig struct {
	Core           CoreSection
	Packages       []PackageDef                // from [[package]] table array
	Dependencies   DependenciesSection         // from [dependencies] (preferred)
	ReportPath     string                      // from [standards] report-path
	StalenessDays  int                         // from [standards] staleness-days (0 = unset)
	Extends        []string                    // from [standards] extends
	CheckExclude   []string                    // from [standards] exclude
	Sections       map[string]DomainSection    // dotted keys: "engineering", "engineering.architecture"
	InlineProfiles map[string]InlineProfileDef // from [profiles.*]
}

// Config holds the effective settings after merging all file layers.
type Config struct {
	Core           CoreSection
	Packages       []PackageDef // merged [[package]] entries, deduped by name
	ReportPath     string       // first non-empty across layers
	StalenessDays  int          // first nonzero across layers (default 7 when 0)
	sections       map[string]DomainSection
	InlineProfiles map[string]InlineProfileDef // merged, higher-priority layers win per name
	// Sources maps dotted key paths to the file that provided them.
	// E.g. "core.home" → "/path/to/grimoire.toml"
	Sources map[string]string
	// MissingExtends holds standards.extends refs that could not be resolved
	// because the referenced package is not installed.
	MissingExtends []string
	// DepSkills is the union of all layers' [dependencies] skills, deduped,
	// in declaration order (project layer first, then global, then system).
	DepSkills []string
	// CheckExclude is the union of all layers' [standards] exclude patterns.
	CheckExclude []string
}

// ProjectPath returns the path to the project-level config file.
func ProjectPath(dir string) string {
	return filepath.Join(dir, "grimoire.toml")
}

// GlobalPath returns the path to the user-global config file, respecting XDG_CONFIG_HOME.
func GlobalPath() string {
	if h := os.Getenv("XDG_CONFIG_HOME"); h != "" {
		return filepath.Join(h, "grimoire", "grimoire.toml")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "grimoire", "grimoire.toml")
}

// SystemPath returns the system-wide config file path.
// Linux/macOS: /etc/grimoire/grimoire.toml
// Windows: %PROGRAMDATA%\grimoire\grimoire.toml.
func SystemPath() string {
	if runtime.GOOS == "windows" {
		pd := os.Getenv("PROGRAMDATA")
		if pd == "" {
			pd = `C:\ProgramData`
		}
		return filepath.Join(pd, "grimoire", "grimoire.toml")
	}
	return "/etc/grimoire/grimoire.toml"
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
func Merge(layers []FileConfig, paths []string) Config {
	r := Config{
		sections:       make(map[string]DomainSection),
		InlineProfiles: make(map[string]InlineProfileDef),
		Sources:        make(map[string]string),
	}

	// Core scalar fields: first non-empty wins (layers[0] = highest priority).
	for i := range layers {
		fs := &layers[i]
		src := paths[i]
		if r.Core.Home == "" && fs.Core.Home != "" {
			r.Core.Home = fs.Core.Home
			r.Sources["core.home"] = src
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
	}

	// Packages: union by name; first occurrence (highest-priority layer) wins per name.
	// Normalize unset priorities: 100 for official, 50 for user packages.
	seenReg := make(map[string]bool)
	for i := range layers {
		for _, rd := range layers[i].Packages {
			if rd.Name == "" || seenReg[rd.Name] {
				continue
			}
			seenReg[rd.Name] = true
			if rd.Priority == 0 {
				if rd.Official {
					rd.Priority = 100
				} else {
					rd.Priority = 50
				}
			}
			r.Packages = append(r.Packages, rd)
		}
	}

	// collect all section keys across all layers
	allKeys := make(map[string]struct{})
	for i := range layers {
		for k := range layers[i].Sections {
			allKeys[k] = struct{}{}
		}
	}

	// InlineProfiles: higher-priority layers win per profile name (project > global > system)
	for i := len(layers) - 1; i >= 0; i-- {
		for name, def := range layers[i].InlineProfiles { //nolint:gocritic // map range copy unavoidable
			r.InlineProfiles[name] = def
		}
	}

	for key := range allKeys {
		merged := DomainSection{ComplianceThresholdError: -1}
		for i := range layers {
			src := paths[i]
			ds, ok := layers[i].Sections[key]
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

	// DepSkills: union across all layers in declaration order, deduped.
	seenDep := make(map[string]bool)
	for i := range layers {
		for _, dep := range layers[i].Dependencies.Skills {
			if !seenDep[dep] {
				seenDep[dep] = true
				r.DepSkills = append(r.DepSkills, dep)
			}
		}
	}

	// CheckExclude: union across all layers in declaration order, deduped.
	seenExclude := make(map[string]bool)
	for i := range layers {
		for _, pat := range layers[i].CheckExclude {
			if !seenExclude[pat] {
				seenExclude[pat] = true
				r.CheckExclude = append(r.CheckExclude, pat)
			}
		}
	}

	return r
}

// SectionKeys returns all domain/subdomain keys present in the resolved settings.
func (r Config) SectionKeys() []string { //nolint:gocritic // value receiver is intentional for immutable Config
	keys := make([]string, 0, len(r.sections))
	for k := range r.sections {
		keys = append(keys, k)
	}
	return keys
}

// ResolveSection returns the effective DomainSection for scope (e.g. "engineering.testing").
// Subdomain keys overlay domain keys; unset subdomain keys inherit from the domain.
// Only one level of nesting is supported (domain.subdomain).
func (r Config) ResolveSection(scope string) DomainSection { //nolint:gocritic // value receiver is intentional for immutable Config
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
