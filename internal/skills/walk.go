package skills

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// skillFrontmatter holds all recognized SKILL.md YAML frontmatter fields.
// Compatible with OpenCode's recognized fields (name, description, license,
// compatibility, metadata). Grimoire extends with version, authors, tags,
// dependencies. OpenCode silently ignores unknown fields.
type skillFrontmatter struct {
	Name          string            `yaml:"name"`
	Version       string            `yaml:"version"`
	Description   string            `yaml:"description"`
	Authors       []string          `yaml:"authors"`
	License       string            `yaml:"license"`
	Tags          []string          `yaml:"tags"`
	Compatibility []string          `yaml:"compatibility"`
	Metadata      map[string]string `yaml:"metadata"`
	Dependencies  map[string]string `yaml:"dependencies"`
	Criteria      []string          `yaml:"criteria"`
	// Lifecycle and citation fields (STANDARD.md)
	Source            string   `yaml:"source"`
	Emerging          bool     `yaml:"emerging"`
	Stable            bool     `yaml:"stable"`
	Deprecated        bool     `yaml:"deprecated"`
	DeprecatedBy      string   `yaml:"deprecated_by"`
	Practitioner      bool     `yaml:"practitioner"`
	Verified          bool     `yaml:"verified"`
	Related           []string `yaml:"related"`
	DuplicateReviewed bool     `yaml:"duplicate-reviewed"`
}

// parseSkillMeta reads SKILL.md and returns its frontmatter and raw body.
// When withBody is false, only the first 4 KB are read (enough for any
// frontmatter) and the returned body is always empty — skips 80-95% of I/O
// for callers that never use sk.Body (e.g. install).
// Returns zero values when the file is absent or unparseable.
func parseSkillMeta(skillPath string, withBody bool) (meta skillFrontmatter, body string) {
	fullPath := filepath.Join(skillPath, "SKILL.md")
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return skillFrontmatter{}, ""
	}
	if !withBody {
		// Frontmatter-only fast path: parse without constructing body string.
		content := string(data)
		if strings.HasPrefix(content, "---") {
			rest := content[3:]
			end := strings.Index(rest, "\n---")
			if end != -1 {
				var m skillFrontmatter
				_ = yaml.Unmarshal([]byte(rest[:end]), &m)
				return m, ""
			}
		}
		return skillFrontmatter{}, ""
	}
	meta, body, _ = parseSkillMetaFromContent(string(data))
	return meta, body
}

// parseSkillMetaFromContent parses frontmatter and body from already-read
// SKILL.md content. Used by parseSkillMeta (withBody=true) and
// ParseSkillFromContent to avoid duplicate file reads.
// Returns the YAML unmarshal error so callers can distinguish missing-field
// errors from invalid-YAML errors without a second parse.
func parseSkillMetaFromContent(content string) (meta skillFrontmatter, body string, yamlErr error) {
	if strings.HasPrefix(content, "---") {
		rest := content[3:]
		end := strings.Index(rest, "\n---")
		if end != -1 {
			yamlErr = yaml.Unmarshal([]byte(rest[:end]), &meta)
			return meta, rest[end+4:], yamlErr // skip past the closing ---\n
		}
	}
	return skillFrontmatter{}, content, nil
}

type Domain struct {
	Name   string
	Path   string
	Nested bool
}

type Subdomain struct {
	Domain string
	Name   string
	Path   string
}

type Skill struct {
	// Identity fields (existing)
	Package   string `json:"package,omitempty"`
	Domain    string `json:"domain"`
	Subdomain string `json:"subdomain,omitempty"`
	Name      string `json:"name"`
	Path      string `json:"path"`

	// Metadata from SKILL.md YAML frontmatter (primary source).
	// SKILL.toml overrides these when present. All empty when neither declares them.
	Tags          []string          `json:"tags,omitempty"`
	Version       string            `json:"version,omitempty"`
	Description   string            `json:"description,omitempty"`
	Authors       []string          `json:"authors,omitempty"`
	License       string            `json:"license,omitempty"`
	Compatibility []string          `json:"compatibility,omitempty"` // e.g. ["opencode","claude"]
	Metadata      map[string]string `json:"metadata,omitempty"`      // opencode-compatible string map
	Dependencies  map[string]string `json:"dependencies,omitempty"`  // skill name → semver constraint

	// Criteria is the explicit list of compliance criteria parsed from the SKILL.md
	// frontmatter. When present, the AI is instructed to evaluate exactly these
	// criteria (using these names verbatim in criteria_matrix), making d.Total precise.
	Criteria []string `json:"criteria,omitempty"`

	// Body is the raw SKILL.md content after the frontmatter block.
	// Used as the compliance rubric — the AI infers criteria from the full text.
	Body string `json:"body,omitempty"`

	// Lifecycle and citation fields (STANDARD.md).
	Source            string   `json:"source,omitempty"`
	Emerging          bool     `json:"emerging,omitempty"`
	Stable            bool     `json:"stable,omitempty"`
	Deprecated        bool     `json:"deprecated,omitempty"`
	DeprecatedBy      string   `json:"deprecated_by,omitempty"`
	Practitioner      bool     `json:"practitioner,omitempty"`
	Verified          bool     `json:"verified,omitempty"`
	Related           []string `json:"related,omitempty"`
	DuplicateReviewed bool     `json:"duplicate_reviewed,omitempty"`
}

// IsNested returns true when the domain uses subdomain directories
// (no direct skills/ folder, or skills/ is empty).
// Defaults to true when the skills/ dir is absent so that callers iterate
// subdomains rather than assuming a flat layout for an unrecognised domain.
func IsNested(domainDir string) bool {
	skillsDir := filepath.Join(domainDir, "skills")
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return true
	}
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), ".") {
			return false
		}
	}
	return true
}

func ListDomains(skillsRoot string) ([]Domain, error) {
	entries, err := os.ReadDir(skillsRoot)
	if err != nil {
		return nil, err
	}
	var domains []Domain
	for _, e := range entries {
		// .claude-plugin is a Claude plugin metadata directory that lives
		// alongside skill domains but is not itself a domain.
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") || e.Name() == ".claude-plugin" {
			continue
		}
		path := filepath.Join(skillsRoot, e.Name())
		domains = append(domains, Domain{
			Name:   e.Name(),
			Path:   path,
			Nested: IsNested(path),
		})
	}
	return domains, nil
}

func ListSubdomains(domainDir string) ([]Subdomain, error) {
	domain := filepath.Base(domainDir)
	entries, err := os.ReadDir(domainDir)
	if err != nil {
		return nil, err
	}
	var subs []Subdomain
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") || e.Name() == ".claude-plugin" {
			continue
		}
		subPath := filepath.Join(domainDir, e.Name())
		if _, err := os.Stat(filepath.Join(subPath, "skills")); err != nil {
			continue
		}
		subs = append(subs, Subdomain{
			Domain: domain,
			Name:   e.Name(),
			Path:   subPath,
		})
	}
	return subs, nil
}

// skillEntry holds the minimal path information needed to build a Skill without
// reading any files. Used to separate directory enumeration from YAML parsing.
type skillEntry struct {
	path, domain, subdomain string
}

// collectSkillEntries enumerates all skill directories under a domain using
// only ReadDir + Stat — no file reads. Safe to call from goroutines.
func collectSkillEntries(d Domain) []skillEntry {
	var entries []skillEntry
	if d.Nested {
		subs, err := ListSubdomains(d.Path)
		if err != nil {
			return nil
		}
		for _, s := range subs {
			skillsDir := filepath.Join(s.Path, "skills")
			es, err := os.ReadDir(skillsDir)
			if err != nil {
				continue
			}
			for _, e := range es {
				if strings.HasPrefix(e.Name(), ".") {
					continue
				}
				sp := filepath.Join(skillsDir, e.Name())
				if _, statErr := os.Stat(filepath.Join(sp, "SKILL.md")); statErr == nil {
					entries = append(entries, skillEntry{sp, d.Name, s.Name})
				}
			}
		}
	} else {
		skillsDir := filepath.Join(d.Path, "skills")
		es, err := os.ReadDir(skillsDir)
		if err != nil {
			return nil
		}
		for _, e := range es {
			if strings.HasPrefix(e.Name(), ".") {
				continue
			}
			sp := filepath.Join(skillsDir, e.Name())
			if _, statErr := os.Stat(filepath.Join(sp, "SKILL.md")); statErr == nil {
				entries = append(entries, skillEntry{sp, d.Name, ""})
			}
		}
	}
	return entries
}

// listSkillsInDir is the unexported implementation with configurable body loading.
// Called from within domain goroutines in listAllSkills — kept sequential to avoid
// nested semaphore contention.
func listSkillsInDir(dir, domainName, subdomainName string, withBody bool) ([]Skill, error) {
	skillsDir := filepath.Join(dir, "skills")
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil, err
	}
	var skills []Skill
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		skillPath := filepath.Join(skillsDir, e.Name())
		if _, err := os.Stat(filepath.Join(skillPath, "SKILL.md")); err != nil {
			continue
		}
		meta, body := parseSkillMeta(skillPath, withBody)
		name := e.Name()
		if meta.Name != "" {
			name = meta.Name
		}
		skills = append(skills, Skill{
			Domain:            domainName,
			Subdomain:         subdomainName,
			Name:              name,
			Path:              skillPath,
			Tags:              meta.Tags,
			Version:           meta.Version,
			Description:       meta.Description,
			Authors:           meta.Authors,
			License:           meta.License,
			Compatibility:     meta.Compatibility,
			Metadata:          meta.Metadata,
			Dependencies:      meta.Dependencies,
			Criteria:          meta.Criteria,
			Body:              body,
			Source:            meta.Source,
			Emerging:          meta.Emerging,
			Stable:            meta.Stable,
			Deprecated:        meta.Deprecated,
			DeprecatedBy:      meta.DeprecatedBy,
			Practitioner:      meta.Practitioner,
			Verified:          meta.Verified,
			Related:           meta.Related,
			DuplicateReviewed: meta.DuplicateReviewed,
		})
	}
	return skills, nil
}

// ListSkillsInDir returns all skills in dir/skills/ with full body loading.
func ListSkillsInDir(dir, domainName, subdomainName string) ([]Skill, error) {
	return listSkillsInDir(dir, domainName, subdomainName, true)
}

// listAllSkills is the unexported implementation with configurable body loading.
//
// !withBody (install path): enumerates domains in parallel using collectSkillEntries
// (ReadDir + Stat only — zero file reads), then builds Skills from directory names.
// Handles 1000+ skills in <50ms regardless of YAML parse cost.
//
// withBody (check path): sequential enumeration then concurrent YAML parse.
func listAllSkills(skillsRoot string, withBody bool) ([]Skill, error) {
	domains, err := ListDomains(skillsRoot)
	if err != nil {
		return nil, err
	}

	concurrency := runtime.GOMAXPROCS(0) * 2
	if concurrency > 32 {
		concurrency = 32
	}

	if !withBody {
		// Install path: enumerate skill directories in parallel, zero file reads.
		entrySets := make([][]skillEntry, len(domains))
		var wg sync.WaitGroup
		sem := make(chan struct{}, concurrency)
		for i, d := range domains {
			wg.Add(1)
			go func(i int, d Domain) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()
				entrySets[i] = collectSkillEntries(d)
			}(i, d)
		}
		wg.Wait()

		var all []Skill
		for _, es := range entrySets {
			for _, e := range es {
				all = append(all, Skill{
					Domain:    e.domain,
					Subdomain: e.subdomain,
					Name:      filepath.Base(e.path),
					Path:      e.path,
				})
			}
		}
		return all, nil
	}

	// Check path (withBody=true): collect entries sequentially, then parse concurrently.
	var entries []skillEntry
	for _, d := range domains {
		entries = append(entries, collectSkillEntries(d)...)
	}

	result := make([]Skill, len(entries))
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)
	for i, e := range entries {
		wg.Add(1)
		go func(i int, e skillEntry) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			meta, body := parseSkillMeta(e.path, true)
			name := filepath.Base(e.path)
			if meta.Name != "" {
				name = meta.Name
			}
			result[i] = Skill{
				Domain:            e.domain,
				Subdomain:         e.subdomain,
				Name:              name,
				Path:              e.path,
				Tags:              meta.Tags,
				Version:           meta.Version,
				Description:       meta.Description,
				Authors:           meta.Authors,
				License:           meta.License,
				Compatibility:     meta.Compatibility,
				Metadata:          meta.Metadata,
				Dependencies:      meta.Dependencies,
				Criteria:          meta.Criteria,
				Body:              body,
				Source:            meta.Source,
				Emerging:          meta.Emerging,
				Stable:            meta.Stable,
				Deprecated:        meta.Deprecated,
				DeprecatedBy:      meta.DeprecatedBy,
				Practitioner:      meta.Practitioner,
				Verified:          meta.Verified,
				Related:           meta.Related,
				DuplicateReviewed: meta.DuplicateReviewed,
			}
		}(i, e)
	}
	wg.Wait()

	var all []Skill
	for i := range result {
		sk := result[i]
		if sk.Path != "" {
			all = append(all, sk)
		}
	}
	return all, nil
}

// ListAllSkills lists all skills under skillsRoot with full body loading.
func ListAllSkills(skillsRoot string) ([]Skill, error) {
	return listAllSkills(skillsRoot, true)
}

// ResolveSkillPath resolves a skill reference like "engineering/development/propose-commit"
// to its absolute path in the skills root.
func ResolveSkillPath(skillsRoot, ref string) (string, error) {
	parts := strings.Split(ref, "/")
	var skillPath string
	switch len(parts) {
	case 3:
		skillPath = filepath.Join(skillsRoot, parts[0], parts[1], "skills", parts[2])
	case 2:
		skillPath = filepath.Join(skillsRoot, parts[0], "skills", parts[1])
	default:
		return "", &ErrInvalidSkillRef{Ref: ref}
	}
	if _, err := os.Stat(filepath.Join(skillPath, "SKILL.md")); err != nil {
		return "", &ErrSkillNotFound{Path: skillPath}
	}
	return skillPath, nil
}

type ErrInvalidSkillRef struct{ Ref string }
type ErrSkillNotFound struct{ Path string }

func (e *ErrInvalidSkillRef) Error() string {
	return "invalid skill reference: " + e.Ref + " (use domain/subdomain/skill or domain/skill)"
}
func (e *ErrSkillNotFound) Error() string { return "skill not found: " + e.Path }

// ParseSkillFile reads and parses a single SKILL.md file. path may be the
// SKILL.md file itself or the directory containing it.
func ParseSkillFile(path string) (Skill, error) {
	skillMD := path
	if filepath.Base(path) != "SKILL.md" {
		skillMD = filepath.Join(path, "SKILL.md")
	}
	if _, err := os.Stat(skillMD); err != nil {
		return Skill{}, err
	}
	dir := filepath.Dir(skillMD)
	meta, body := parseSkillMeta(dir, true)
	name := filepath.Base(dir)
	if meta.Name != "" {
		name = meta.Name
	}
	return Skill{
		Name:              name,
		Path:              dir,
		Tags:              meta.Tags,
		Version:           meta.Version,
		Description:       meta.Description,
		Authors:           meta.Authors,
		License:           meta.License,
		Compatibility:     meta.Compatibility,
		Metadata:          meta.Metadata,
		Dependencies:      meta.Dependencies,
		Criteria:          meta.Criteria,
		Body:              body,
		Source:            meta.Source,
		Emerging:          meta.Emerging,
		Stable:            meta.Stable,
		Deprecated:        meta.Deprecated,
		DeprecatedBy:      meta.DeprecatedBy,
		Practitioner:      meta.Practitioner,
		Verified:          meta.Verified,
		Related:           meta.Related,
		DuplicateReviewed: meta.DuplicateReviewed,
	}, nil
}

// ParseSkillFromContent builds a Skill from already-read SKILL.md content.
// dir must be the directory that contains the SKILL.md file.
// The returned error is the YAML unmarshal error (if any); the Skill is still
// populated with whatever could be parsed. Use instead of ParseSkillFile when
// the content is already in memory to avoid a redundant file read.
func ParseSkillFromContent(content, dir string) (Skill, error) {
	meta, body, yamlErr := parseSkillMetaFromContent(content)
	name := filepath.Base(dir)
	if meta.Name != "" {
		name = meta.Name
	}
	return Skill{
		Name:              name,
		Path:              dir,
		Tags:              meta.Tags,
		Version:           meta.Version,
		Description:       meta.Description,
		Authors:           meta.Authors,
		License:           meta.License,
		Compatibility:     meta.Compatibility,
		Metadata:          meta.Metadata,
		Dependencies:      meta.Dependencies,
		Criteria:          meta.Criteria,
		Body:              body,
		Source:            meta.Source,
		Emerging:          meta.Emerging,
		Stable:            meta.Stable,
		Deprecated:        meta.Deprecated,
		DeprecatedBy:      meta.DeprecatedBy,
		Practitioner:      meta.Practitioner,
		Verified:          meta.Verified,
		Related:           meta.Related,
		DuplicateReviewed: meta.DuplicateReviewed,
	}, yamlErr
}

// WalkSkillFiles returns the absolute paths of all SKILL.md files found
// recursively under root, skipping hidden and vendor directories.
func WalkSkillFiles(root string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if (len(name) > 1 && name[0] == '.') || name == "node_modules" || name == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() == "SKILL.md" {
			paths = append(paths, p)
		}
		return nil
	})
	return paths, err
}
