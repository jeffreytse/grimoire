// Package manifest loads, merges, and writes grimoire.toml files.
// grimoire.toml is the single config format used at all scopes:
//
//	./grimoire.toml                       project scope
//	~/.config/grimoire/grimoire.toml      global/user scope
//	/etc/grimoire/grimoire.toml           system scope
//
// Sections:
//
//	[package]        project identity metadata (project scope only)
//	[dependencies]   skill deps with semver constraints
//	[dev-dependencies]
//	[[package]]     additional packages (official package is always implicit)
//	[core]           install-mode, home dir, agents (global/system scope)
//	[standards]      BPDD compliance config (profiles, report-path, etc.)
//	[standards.*]    per-domain compliance sections
package manifest

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"github.com/jeffreytse/grimoire/internal/config"
)

// PackageMeta holds the [package] section — project identity.
type PackageMeta struct {
	Name        string   `toml:"name"`
	Version     string   `toml:"version"`
	Description string   `toml:"description"`
	Authors     []string `toml:"authors"`
	License     string   `toml:"license"`
}

// DepSpec holds the version constraint for a dependency.
// The dep key (map key in Deps/DevDeps) is a PackageRef string that encodes the package:
//
//	apply-solid = "^1.0.0"                      // official package (no prefix)
//	"acmecorp/practices:apply-tdd" = "~2.0.0"   // package in key; no [[package]] config needed
//
// ParsePackageRef(key) extracts the package URL and skill path directly.
type DepSpec struct {
	Version string // semver range or "*"
}

// StandardsSection holds the scalar fields under [standards].
type StandardsSection struct {
	Profiles      []string
	ReportPath    string
	StalenessDays int
	Extends       []string
}

// ManifestFile is one parsed grimoire.toml file.
// DomainSections uses dotted keys: "engineering", "engineering.testing".
// Dep map keys are PackageRef strings — the package is embedded in the key, not stored separately.
type ManifestFile struct {
	Package        PackageMeta
	Deps           map[string]DepSpec // key = PackageRef string
	DevDeps        map[string]DepSpec // key = PackageRef string
	Standards      StandardsSection
	DomainSections map[string]config.DomainSection
	InlineProfiles map[string]config.InlineProfileDef // from [profiles.*]
	Core           config.CoreSection
}

// Resolved holds the effective manifest after merging all file layers.
// Higher-priority layers win per key (project > global > system).
type Resolved struct {
	Package        PackageMeta
	Deps           map[string]DepSpec // merged; project-layer first wins per key
	DevDeps        map[string]DepSpec
	Standards      StandardsSection
	DomainSections map[string]config.DomainSection
	InlineProfiles map[string]config.InlineProfileDef // from [profiles.*]
	Core           config.CoreSection
	Sources        map[string]string // dotted key → source file path
}

// ParseFile reads one grimoire.toml. Returns zero-value ManifestFile when absent.
func ParseFile(path string) (ManifestFile, error) {
	mf := ManifestFile{
		Deps:           make(map[string]DepSpec),
		DevDeps:        make(map[string]DepSpec),
		DomainSections: make(map[string]config.DomainSection),
		InlineProfiles: make(map[string]config.InlineProfileDef),
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return mf, nil
	}
	if err != nil {
		return mf, err
	}
	var raw map[string]any
	if err := toml.Unmarshal(data, &raw); err != nil {
		return mf, err
	}
	return parseRaw(raw), nil
}

// Load reads all three grimoire.toml layers for projectDir and merges them.
// Layers (highest → lowest): ./grimoire.toml → global → system.
func Load(projectDir string) (Resolved, error) {
	layerPaths := []string{
		ProjectPath(projectDir),
		GlobalPath(),
		SystemPath(),
	}
	var layers []ManifestFile
	var paths []string
	var errs []error

	for _, p := range layerPaths {
		mf, err := ParseFile(p)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		layers = append(layers, mf)
		paths = append(paths, p)
	}
	if len(layers) == 0 && len(errs) > 0 {
		return Resolved{}, errs[0]
	}
	// [core].home is machine-level; clear from the project layer so each machine
	// uses its own home. [[package]] stays — project packages are team-shared
	// and committed to VCS so teammates get them automatically.
	if len(layers) > 0 {
		layers[0].Core.Home = ""
		layers[0].Core.InstallMode = ""
	}
	return merge(layers, paths), nil
}

// LoadGlobal reads only the global grimoire.toml.
func LoadGlobal() (ManifestFile, error) {
	return ParseFile(GlobalPath())
}

// WriteFile serializes mf to a grimoire.toml file, creating parent dirs as needed.
func WriteFile(path string, mf *ManifestFile) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	m := toRawMap(mf)
	data, err := toml.Marshal(m)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// AppendDep adds key = 'version' to the [dependencies] section of the TOML file
// at path without rewriting the rest of the file (preserves comments and formatting).
func AppendDep(path, key, version string) error {
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return os.WriteFile(path, []byte("[dependencies]\n"+depLine(key, version)), 0o644)
	}
	if err != nil {
		return err
	}

	lines := strings.Split(string(raw), "\n")

	// Find [dependencies] header (exact match, not a subtable like [dependencies.foo]).
	depStart := -1
	for i, l := range lines {
		if strings.TrimSpace(l) == "[dependencies]" {
			depStart = i
			break
		}
	}

	if depStart < 0 {
		// No [dependencies] section — append it.
		text := strings.TrimRight(string(raw), "\n")
		return os.WriteFile(path, []byte(text+"\n\n[dependencies]\n"+depLine(key, version)), 0o644)
	}

	// Find next top-level section header (line starting with `[` but not `[[`).
	nextSection := len(lines)
	for i := depStart + 1; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "[") && !strings.HasPrefix(trimmed, "[[") {
			nextSection = i
			break
		}
	}

	// Insert before trailing blank lines that precede the next section.
	insertAt := nextSection
	for insertAt > depStart+1 && strings.TrimSpace(lines[insertAt-1]) == "" {
		insertAt--
	}

	newLine := strings.TrimRight(depLine(key, version), "\n")
	out := make([]string, 0, len(lines)+1)
	out = append(out, lines[:insertAt]...)
	out = append(out, newLine)
	out = append(out, lines[insertAt:]...)
	return os.WriteFile(path, []byte(strings.Join(out, "\n")), 0o644)
}

// RemoveDep removes key from the [dependencies] section of the TOML file at path.
// Idempotent: returns nil if the key is not found.
func RemoveDep(path, key string) error {
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	// Build candidate line prefixes for both bare and quoted forms of the key.
	quotedKey := `"` + strings.ReplaceAll(key, `"`, `\"`) + `"`
	candidates := []string{
		key + " =",
		key + "=",
		quotedKey + " =",
		quotedKey + "=",
	}

	lines := strings.Split(string(raw), "\n")
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		matched := false
		for _, c := range candidates {
			if strings.HasPrefix(trimmed, c) {
				matched = true
				break
			}
		}
		if !matched {
			out = append(out, l)
		}
	}

	if len(out) == len(lines) {
		return nil // key not found — idempotent
	}
	return os.WriteFile(path, []byte(strings.Join(out, "\n")), 0o644)
}

// depLine formats one TOML dependency line. Keys with non-barekey characters are quoted.
func depLine(key, version string) string {
	k := key
	for _, ch := range key {
		if (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') &&
			(ch < '0' || ch > '9') && ch != '-' && ch != '_' {
			k = `"` + strings.ReplaceAll(key, `"`, `\"`) + `"`
			break
		}
	}
	v := version
	if v == "" {
		v = "*"
	}
	return fmt.Sprintf("%s = '%s'\n", k, v)
}

// merge combines multiple ManifestFile layers. Earlier layers (index 0) win.
func merge(layers []ManifestFile, paths []string) Resolved {
	r := Resolved{
		Deps:           make(map[string]DepSpec),
		DevDeps:        make(map[string]DepSpec),
		DomainSections: make(map[string]config.DomainSection),
		InlineProfiles: make(map[string]config.InlineProfileDef),
		Sources:        make(map[string]string),
	}
	seenDeps := make(map[string]bool)
	seenDevDeps := make(map[string]bool)

	for i := range layers {
		mf := &layers[i]
		src := ""
		if i < len(paths) {
			src = paths[i]
		}

		// [package] — project layer (index 0) wins entirely
		if i == 0 && mf.Package.Name != "" {
			r.Package = mf.Package
			r.Sources["package"] = src
		}

		// [dependencies] — first occurrence per skill name wins
		for name, spec := range mf.Deps {
			if !seenDeps[name] {
				seenDeps[name] = true
				r.Deps[name] = spec
				r.Sources["dependencies."+name] = src
			}
		}

		// [dev-dependencies]
		for name, spec := range mf.DevDeps {
			if !seenDevDeps[name] {
				seenDevDeps[name] = true
				r.DevDeps[name] = spec
				r.Sources["dev-dependencies."+name] = src
			}
		}

		// [standards] scalar fields — first non-empty wins
		if r.Standards.ReportPath == "" && mf.Standards.ReportPath != "" {
			r.Standards.ReportPath = mf.Standards.ReportPath
			r.Sources["standards.report-path"] = src
		}
		if r.Standards.StalenessDays == 0 && mf.Standards.StalenessDays > 0 {
			r.Standards.StalenessDays = mf.Standards.StalenessDays
			r.Sources["standards.staleness-days"] = src
		}
		if len(r.Standards.Extends) == 0 && len(mf.Standards.Extends) > 0 {
			r.Standards.Extends = mf.Standards.Extends
			r.Sources["standards.extends"] = src
		}
		// profiles: union across all layers
		r.Standards.Profiles = appendUnique(r.Standards.Profiles, mf.Standards.Profiles...)

		// [standards.*] domain sections — first occurrence per key wins
		for key, ds := range mf.DomainSections {
			if _, exists := r.DomainSections[key]; !exists {
				r.DomainSections[key] = ds
				r.Sources["standards."+key] = src
			}
		}

		// [profiles.*] — first occurrence per name wins
		for name := range mf.InlineProfiles {
			if _, exists := r.InlineProfiles[name]; !exists {
				r.InlineProfiles[name] = mf.InlineProfiles[name]
				r.Sources["profiles."+name] = src
			}
		}

		// [core] scalar fields — first non-empty wins
		if r.Core.Home == "" && mf.Core.Home != "" {
			r.Core.Home = mf.Core.Home
			r.Sources["core.home"] = src
		}
		if r.Core.InstallMode == "" && mf.Core.InstallMode != "" {
			r.Core.InstallMode = mf.Core.InstallMode
			r.Sources["core.install-mode"] = src
		}
		if len(r.Core.Agents) == 0 && len(mf.Core.Agents) > 0 {
			r.Core.Agents = mf.Core.Agents
			r.Sources["core.agents"] = src
		}
		if r.Core.UpdateConcurrency == nil && mf.Core.UpdateConcurrency != nil {
			r.Core.UpdateConcurrency = mf.Core.UpdateConcurrency
			r.Sources["core.update-concurrency"] = src
		}
		// search-packages: union across all layers (additive)
		r.Core.SearchPackages = appendUnique(r.Core.SearchPackages, mf.Core.SearchPackages...)
	}

	// GRIMOIRE_HOME env override — highest priority.
	if v := os.Getenv("GRIMOIRE_HOME"); v != "" {
		r.Core.Home = v
		r.Sources["core.home"] = "$GRIMOIRE_HOME"
	}

	return r
}

// appendUnique appends items to dst, skipping duplicates.
func appendUnique(dst []string, items ...string) []string {
	seen := make(map[string]bool, len(dst))
	for _, v := range dst {
		seen[v] = true
	}
	for _, v := range items {
		if !seen[v] {
			seen[v] = true
			dst = append(dst, v)
		}
	}
	return dst
}

// parseRaw converts a map[string]any (from toml.Unmarshal) into a ManifestFile.
func parseRaw(raw map[string]any) ManifestFile {
	mf := ManifestFile{
		Deps:           make(map[string]DepSpec),
		DevDeps:        make(map[string]DepSpec),
		DomainSections: make(map[string]config.DomainSection),
		InlineProfiles: make(map[string]config.InlineProfileDef),
	}
	for key, val := range raw {
		switch key {
		case "package":
			if m, ok := val.(map[string]any); ok {
				mf.Package = parsePackageMeta(m)
			}
		case "dependencies":
			if m, ok := val.(map[string]any); ok {
				mf.Deps = parseDepMap(m)
			}
		case "dev-dependencies":
			if m, ok := val.(map[string]any); ok {
				mf.DevDeps = parseDepMap(m)
			}
		case "core":
			if m, ok := val.(map[string]any); ok {
				mf.Core = parseCoreSection(m)
			}
		case "standards":
			if m, ok := val.(map[string]any); ok {
				parseStandardsInto(m, &mf)
			}
		case "profiles":
			if m, ok := val.(map[string]any); ok {
				mf.InlineProfiles = parseInlineProfiles(m)
			}
		}
	}
	return mf
}

func parsePackageMeta(m map[string]any) PackageMeta {
	var p PackageMeta
	if v, ok := m["name"].(string); ok {
		p.Name = v
	}
	if v, ok := m["version"].(string); ok {
		p.Version = v
	}
	if v, ok := m["description"].(string); ok {
		p.Description = v
	}
	if v, ok := m["license"].(string); ok {
		p.License = v
	}
	p.Authors = parseStringSlice(m["authors"])
	return p
}

// parseDepMap parses a TOML [dependencies] or [dev-dependencies] map.
// Keys are PackageRef strings; values are semver constraint strings.
// Table form is not supported — package info is encoded in the key, not the value.
func parseDepMap(m map[string]any) map[string]DepSpec {
	result := make(map[string]DepSpec, len(m))
	for name, val := range m {
		if v, ok := val.(string); ok {
			result[name] = DepSpec{Version: v}
		}
	}
	return result
}

func parseStandardsInto(m map[string]any, mf *ManifestFile) {
	mf.Standards.Extends = parseStringSlice(m["extends"])
	mf.Standards.Profiles = parseStringSlice(m["profiles"])
	if v, ok := m["report-path"].(string); ok {
		mf.Standards.ReportPath = v
	}
	if v, ok := anyToInt(m["staleness-days"]); ok && v > 0 {
		mf.Standards.StalenessDays = v
	}
	for domainName, val := range m {
		switch domainName {
		case "profiles", "extends", "report-path", "staleness-days":
			continue
		}
		if sub, ok := val.(map[string]any); ok {
			parseDomainInto(domainName, sub, mf)
		}
	}
}

// parseDomainInto extracts a DomainSection from m and recurses into sub-maps.
func parseDomainInto(prefix string, m map[string]any, mf *ManifestFile) {
	ds := config.DomainSection{ComplianceThresholdError: -1}
	for k, v := range m {
		switch k {
		case "practices":
			ds.Practices = parseStringSlice(v)
		case "disabled":
			ds.Disabled = parseStringSlice(v)
		case "fallback":
			if s, ok := v.(string); ok {
				ds.Fallback = s
			}
		case "compliance-threshold":
			ds.ComplianceThreshold = anyToFloat64(v)
		case "compliance-threshold-error":
			if n, ok := anyToInt(v); ok {
				ds.ComplianceThresholdError = n
			}
		default:
			if sub, ok := v.(map[string]any); ok {
				parseDomainInto(prefix+"."+k, sub, mf)
			}
		}
	}
	mf.DomainSections[prefix] = ds
}

func parseCoreSection(m map[string]any) config.CoreSection {
	var cs config.CoreSection
	if v, ok := m["home"].(string); ok {
		cs.Home = v
	}
	if v, ok := m["install-mode"].(string); ok {
		cs.InstallMode = v
	}
	cs.Agents = parseStringSlice(m["agents"])
	cs.SearchPackages = parseStringSlice(m["search-packages"])
	if v, ok := anyToInt(m["update-concurrency"]); ok && v >= 0 {
		cs.UpdateConcurrency = &v
	}
	return cs
}

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

// toRawMap converts ManifestFile to a nested map[string]any for TOML marshaling.
func toRawMap(mf *ManifestFile) map[string]any {
	m := map[string]any{}

	// [package]
	if mf.Package.Name != "" {
		pm := map[string]any{}
		if mf.Package.Name != "" {
			pm["name"] = mf.Package.Name
		}
		if mf.Package.Version != "" {
			pm["version"] = mf.Package.Version
		}
		if mf.Package.Description != "" {
			pm["description"] = mf.Package.Description
		}
		if len(mf.Package.Authors) > 0 {
			pm["authors"] = mf.Package.Authors
		}
		if mf.Package.License != "" {
			pm["license"] = mf.Package.License
		}
		m["package"] = pm
	}

	// [dependencies] — keys are PackageRef strings; values are semver constraint strings
	if len(mf.Deps) > 0 {
		dm := make(map[string]any, len(mf.Deps))
		for name, spec := range mf.Deps {
			dm[name] = spec.Version
		}
		m["dependencies"] = dm
	}

	// [dev-dependencies]
	if len(mf.DevDeps) > 0 {
		dm := make(map[string]any, len(mf.DevDeps))
		for name, spec := range mf.DevDeps {
			dm[name] = spec.Version
		}
		m["dev-dependencies"] = dm
	}

	// [core]
	core := map[string]any{}
	if mf.Core.Home != "" {
		core["home"] = mf.Core.Home
	}
	if len(mf.Core.Agents) > 0 {
		core["agents"] = mf.Core.Agents
	}
	if mf.Core.InstallMode != "" {
		core["install-mode"] = mf.Core.InstallMode
	}
	if mf.Core.UpdateConcurrency != nil {
		core["update-concurrency"] = *mf.Core.UpdateConcurrency
	}
	if len(mf.Core.SearchPackages) > 0 {
		core["search-packages"] = mf.Core.SearchPackages
	}
	if len(core) > 0 {
		m["core"] = core
	}

	// [standards] and [standards.*]
	standardsMap := map[string]any{}
	if len(mf.Standards.Profiles) > 0 {
		standardsMap["profiles"] = mf.Standards.Profiles
	}
	if mf.Standards.ReportPath != "" {
		standardsMap["report-path"] = mf.Standards.ReportPath
	}
	if mf.Standards.StalenessDays > 0 {
		standardsMap["staleness-days"] = mf.Standards.StalenessDays
	}
	if len(mf.Standards.Extends) > 0 {
		standardsMap["extends"] = mf.Standards.Extends
	}

	// sort domain keys so parents appear before children
	keys := make([]string, 0, len(mf.DomainSections))
	for k := range mf.DomainSections {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		ds := mf.DomainSections[key]
		dm := domainToMap(&ds)
		if len(dm) == 0 {
			continue
		}
		parts := strings.SplitN(key, ".", 2)
		if len(parts) == 1 {
			existing, ok := standardsMap[key].(map[string]any)
			if !ok {
				existing = map[string]any{}
			}
			for k, v := range dm {
				existing[k] = v
			}
			standardsMap[key] = existing
		} else {
			parent, ok := standardsMap[parts[0]].(map[string]any)
			if !ok {
				parent = map[string]any{}
				standardsMap[parts[0]] = parent
			}
			parent[parts[1]] = dm
		}
	}
	if len(standardsMap) > 0 {
		m["standards"] = standardsMap
	}

	// [profiles.*]
	if len(mf.InlineProfiles) > 0 {
		profilesMap := make(map[string]any, len(mf.InlineProfiles))
		for name := range mf.InlineProfiles {
			def := mf.InlineProfiles[name]
			pm := map[string]any{}
			if def.Description != "" {
				pm["description"] = def.Description
			}
			if len(def.Tags) > 0 {
				pm["tags"] = def.Tags
			}
			if len(def.Extends) > 0 {
				pm["extends"] = def.Extends
			}
			if len(def.Exclude) > 0 {
				pm["exclude"] = def.Exclude
			}
			if def.ComplianceThreshold > 0 {
				pm["compliance-threshold"] = def.ComplianceThreshold
			}
			if def.ComplianceThresholdError >= 0 {
				pm["compliance-threshold-error"] = def.ComplianceThresholdError
			}
			if len(def.Skills) > 0 {
				var skillsArr []any
				for _, ref := range def.Skills {
					sm := map[string]any{"name": ref.Name}
					if ref.Priority != 0 {
						sm["priority"] = ref.Priority
					}
					skillsArr = append(skillsArr, sm)
				}
				pm["skills"] = skillsArr
			}
			profilesMap[name] = pm
		}
		m["profiles"] = profilesMap
	}

	return m
}

// parseInlineProfiles parses the [profiles.*] top-level section.
func parseInlineProfiles(m map[string]any) map[string]config.InlineProfileDef {
	result := make(map[string]config.InlineProfileDef)
	for name, val := range m {
		sub, ok := val.(map[string]any)
		if !ok {
			continue
		}
		def := config.InlineProfileDef{ComplianceThresholdError: -1}
		if v, ok := sub["name"].(string); ok {
			def.Name = v
		}
		if v, ok := sub["description"].(string); ok {
			def.Description = v
		}
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
					ref := config.InlineSkillRef{}
					if v, ok := sm["name"].(string); ok {
						ref.Name = v
					}
					if n, ok := anyToInt(sm["priority"]); ok {
						ref.Priority = n
					}
					if ref.Name != "" {
						def.Skills = append(def.Skills, ref)
					}
				}
			}
		}
		result[name] = def
	}
	return result
}

func domainToMap(ds *config.DomainSection) map[string]any {
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
