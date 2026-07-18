package skills

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/bmatcuk/doublestar/v4"

	"github.com/jeffreytse/grimoire/internal/config"
)

// GlobMatch reports whether path matches pattern using Standard Glob with globstar.
// Supports *, ?, **, {a,b}, [...], [^...]/[!...], and POSIX character classes
// such as [[:alpha:]], [[:digit:]], [[:alnum:]], [[:lower:]], [[:upper:]], [[:xdigit:]].
// All patterns are root-anchored; use ** prefix to match at any depth.
// pattern="" matches all paths.
func GlobMatch(pattern, path string) bool {
	if pattern == "" {
		return true
	}
	// .gitignore style: pattern without '/' matches at any path depth.
	if !strings.Contains(pattern, "/") {
		pattern = "**/" + pattern
	}
	expanded := expandPOSIXClasses(pattern)
	matched, _ := doublestar.Match(expanded, path)
	return matched
}

// walkSkills is the unexported implementation with configurable body loading.
// Phase 1: walk directory tree collecting skill paths, skipping hidden dirs
// (.git/, dot-prefixed) and known build artifact dirs. Phase 2: parse
// SKILL.md files concurrently (8-slot semaphore) to parallelise I/O.
func walkSkills(root string, withBody bool) ([]Skill, error) {
	// Phase 1: collect skill dirs — no file reads, skip hidden noise.
	var paths []string
	walkErr := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			return nil
		}
		name := d.Name()
		if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" {
			return filepath.SkipDir
		}
		if _, statErr := os.Stat(filepath.Join(p, "SKILL.md")); statErr == nil {
			paths = append(paths, p)
		}
		return nil
	})
	if walkErr != nil {
		return nil, walkErr
	}

	if !withBody {
		// Install path: build minimal Skills from paths — zero file reads.
		out := make([]Skill, 0, len(paths))
		for _, p := range paths {
			rel, relErr := filepath.Rel(root, p)
			if relErr != nil {
				continue
			}
			out = append(out, Skill{
				Name:   filepath.Base(p),
				Path:   p,
				Domain: filepath.ToSlash(rel),
			})
		}
		return out, nil
	}

	// Phase 2 (withBody=true): parse SKILL.md files concurrently.
	concurrency := runtime.GOMAXPROCS(0) * 2
	if concurrency > 32 {
		concurrency = 32
	}
	result := make([]Skill, len(paths))
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)
	for i, p := range paths {
		wg.Add(1)
		go func(i int, p string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			rel, relErr := filepath.Rel(root, p)
			if relErr != nil {
				return
			}
			rel = filepath.ToSlash(rel)
			meta, body := parseSkillMeta(p, true)
			name := filepath.Base(p)
			if meta.Name != "" {
				name = meta.Name
			}
			result[i] = Skill{
				Name:              name,
				Path:              p,
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
				// Domain carries the root-relative path for glob matching.
				Domain: rel,
			}
		}(i, p)
	}
	wg.Wait()

	var out []Skill
	for i := range result {
		sk := result[i]
		if sk.Path != "" {
			out = append(out, sk)
		}
	}
	return out, nil
}

// WalkSkills recursively finds all skills (directories containing SKILL.md) under root.
// root is the package (repo) root — no assumed skills/ prefix.
// Returns skills whose Domain field is relative to root.
func WalkSkills(root string) ([]Skill, error) {
	return walkSkills(root, true)
}

// skillsMatchingGlob is the unexported implementation with configurable body loading.
func skillsMatchingGlob(root, glob string, withBody bool) ([]Skill, error) {
	all, err := walkSkills(root, withBody)
	if err != nil {
		return nil, err
	}
	if glob == "" {
		return all, nil
	}
	var matched []Skill
	for i := range all {
		sk := all[i]
		if GlobMatch(glob, sk.Domain) {
			matched = append(matched, sk)
		}
	}
	return matched, nil
}

// SkillsMatchingGlob returns all skills under root whose root-relative path matches glob.
// root is the package (repo) root. glob="" returns all skills.
// Follows Standard Glob with globstar semantics (always root-anchored).
func SkillsMatchingGlob(root, glob string) ([]Skill, error) {
	return skillsMatchingGlob(root, glob, true)
}

// skillsMatchingAny is the unexported implementation with configurable body loading.
// Each unique package root is walked exactly once regardless of how many refs point to it.
func skillsMatchingAny(refs []config.PackageRef, allRegs []PackageEntry, withBody bool) ([]Skill, []SkillConflict, error) {
	seen := make(map[string]string) // rel-path → ref.Raw that claimed it
	var result []Skill
	var conflicts []SkillConflict

	type refResolved struct{ root, glob, pkgName, raw string }
	rr := make([]refResolved, 0, len(refs))
	for i := range refs {
		ref := refs[i]
		var root, glob, pkgName string
		switch {
		case ref.IsLocal():
			root = filepath.FromSlash(ref.LocalPath)
			glob = ref.Path
			pkgName = ref.LocalPath
		case ref.IsOfficialRepoPath():
			root = officialPackageHome(allRegs)
			pkgName = officialName(allRegs)
			glob = ref.Path
		default:
			reg := findPackage(ref.PackageName, allRegs)
			if reg == nil {
				continue
			}
			root = reg.Home
			pkgName = reg.Name
			glob = ref.Path
		}
		if root != "" {
			rr = append(rr, refResolved{root, glob, pkgName, ref.Raw})
		}
	}

	// Walk each unique root exactly once.
	type walkResult struct {
		skills []Skill
		err    error
	}
	walkCache := make(map[string]walkResult)
	for _, r := range rr {
		if _, exists := walkCache[r.root]; !exists {
			ss, err := walkSkills(r.root, withBody)
			walkCache[r.root] = walkResult{ss, err}
		}
	}

	// Filter cached walk by glob per ref, dedup.
	for _, r := range rr {
		cached := walkCache[r.root]
		if cached.err != nil {
			continue
		}
		for i := range cached.skills {
			sk := cached.skills[i]
			if r.glob != "" && !GlobMatch(r.glob, sk.Domain) {
				continue
			}
			key := sk.Domain
			if winner, exists := seen[key]; exists {
				conflicts = append(conflicts, SkillConflict{
					CanonicalPath: key,
					WinnerPackage: winner,
					LoserPackage:  r.raw,
				})
				continue
			}
			seen[key] = r.raw
			if sk.Package == "" {
				sk.Package = r.pkgName
			}
			result = append(result, sk)
		}
	}
	return result, conflicts, nil
}

// SkillsMatchingAny resolves the install pool from a list of PackageRefs.
// For each ref:
//   - Local: walks ref.LocalPath
//   - IsOfficialRepoPath: walks the official package home with ref.Path as glob
//   - Otherwise: locates package by PackageName in allRegs, walks its Home
//
// Deduplicates by canonical rel-path (first ref wins).
// Sets sk.Package on every returned skill so callers can derive source/commit.
func SkillsMatchingAny(refs []config.PackageRef, allRegs []PackageEntry) ([]Skill, []SkillConflict, error) {
	return skillsMatchingAny(refs, allRegs, true)
}

// SkillsMatchingAnyMeta is identical to SkillsMatchingAny but skips loading
// sk.Body — use for install paths that never access sk.Body.
func SkillsMatchingAnyMeta(refs []config.PackageRef, allRegs []PackageEntry) ([]Skill, []SkillConflict, error) {
	return skillsMatchingAny(refs, allRegs, false)
}

// officialName returns the Name of the official PackageEntry, falling back to
// OfficialPackageDerivedName when none is marked official in allRegs.
func officialName(allRegs []PackageEntry) string {
	for _, r := range allRegs {
		if r.Official {
			return r.Name
		}
	}
	return OfficialPackageDerivedName()
}

// officialPackageHome returns the Home of the official package in allRegs,
// falling back to OfficialPackageHome() if none is marked official.
func officialPackageHome(allRegs []PackageEntry) string {
	for _, r := range allRegs {
		if r.Official {
			return r.Home
		}
	}
	return OfficialPackageHome()
}

// findPackage returns the PackageEntry matching name, or nil if not found.
func findPackage(name string, allRegs []PackageEntry) *PackageEntry {
	for i := range allRegs {
		if allRegs[i].Name == name {
			return &allRegs[i]
		}
	}
	return nil
}

// expandPOSIXClasses expands POSIX named character classes inside bracket expressions.
// e.g. [[:alpha:]] → [A-Za-z], [[:digit:][:upper:]] → [0-9A-Z].
func expandPOSIXClasses(pattern string) string {
	if !strings.Contains(pattern, "[:") {
		return pattern
	}
	posixMap := map[string]string{
		"[:alnum:]":  "0-9A-Za-z",
		"[:alpha:]":  "A-Za-z",
		"[:digit:]":  "0-9",
		"[:lower:]":  "a-z",
		"[:upper:]":  "A-Z",
		"[:xdigit:]": "0-9A-Fa-f",
		"[:word:]":   "0-9A-Za-z_",
		"[:space:]":  " \t\n\r",
		"[:blank:]":  " \t",
		"[:graph:]":  "!-~",
		"[:print:]":  " -~",
		"[:punct:]":  "!-/:-@\\[-`{-~",
	}
	var b strings.Builder
	i := 0
	for i < len(pattern) {
		if pattern[i] != '[' {
			_ = b.WriteByte(pattern[i])
			i++
			continue
		}
		// Scan to the matching ']' of this bracket expression.
		// POSIX classes [: ... :] inside the brackets contain their own ']',
		// so we must skip over [: ... :] sequences instead of stopping at them.
		j := i + 1
		if j < len(pattern) && (pattern[j] == '!' || pattern[j] == '^') {
			j++
		}
		if j < len(pattern) && pattern[j] == ']' {
			j++ // literal ']' as first char of bracket expression
		}
		for j < len(pattern) && pattern[j] != ']' {
			// Skip over embedded POSIX class [: ... :]
			if j+1 < len(pattern) && pattern[j] == '[' && pattern[j+1] == ':' {
				k := j + 2
				for k+1 < len(pattern) && (pattern[k] != ':' || pattern[k+1] != ']') {
					k++
				}
				if k+1 < len(pattern) {
					j = k + 2 // advance past the closing :]
				} else {
					j++ // malformed — advance conservatively
				}
			} else {
				j++
			}
		}
		if j >= len(pattern) {
			_ = b.WriteByte(pattern[i])
			i++
			continue
		}
		inner := pattern[i+1 : j]
		for cls, exp := range posixMap {
			inner = strings.ReplaceAll(inner, cls, exp)
		}
		_ = b.WriteByte('[')
		b.WriteString(inner)
		_ = b.WriteByte(']')
		i = j + 1
	}
	return b.String()
}
