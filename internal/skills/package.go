package skills

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/jeffreytse/grimoire/internal/config"
)

const OfficialPackageName = "official"

// SkillsPackage is a named package with its resolved skills directory path.
// Narrower than PackageEntry: only packages that have a skills/ directory.
type SkillsPackage struct {
	Name string // derived package name, e.g. "github.com/jeffreytse/grimoire-core"
	Root string // absolute path to the skills/ directory
}

// SkillConflict records a skill that was shadowed by a higher-priority package.
type SkillConflict struct {
	CanonicalPath string // e.g. "engineering/development/apply-solid"
	WinnerPackage string // package whose version will be installed
	LoserPackage  string // package whose version was suppressed
}

// PackageEntry is a configured package with its home directory.
// Includes packages that have no skills/ directory (e.g. profiles-only bundles).
type PackageEntry struct {
	Name     string // package name: user-given or URL-derived
	Home     string // absolute path to package root dir (parent of skills/)
	Priority int    // higher = searched first; 100 for official, 50 for user, 0 = old-model default
	Official bool   // true for the STANDARD.md-compliant package
}

// PackageHome returns the local clone directory for a named package.
// All git-hosted packages live under PackagesRoot().
// Local packages (absolute paths) are returned as-is.
func PackageHome(name string) string {
	if filepath.IsAbs(name) {
		return name
	}
	return filepath.Join(PackagesRoot(), name)
}

// AllPackages returns all configured packages, priority-ordered.
//
// Uses [[package]] entries from global settings, sorted by priority descending.
// When no [[package]] is configured, returns an implicit official package entry
// so fresh installs (no settings file) still work.
func AllPackages() []PackageEntry {
	cfg, _ := config.LoadGlobal()

	if len(cfg.Packages) == 0 {
		return []PackageEntry{{
			Name:     OfficialPackageDerivedName(),
			Home:     OfficialPackageHome(),
			Priority: 100,
			Official: true,
		}}
	}

	sorted := make([]config.PackageDef, len(cfg.Packages))
	copy(sorted, cfg.Packages)
	sort.Slice(sorted, func(i, j int) bool {
		return effectivePriority(sorted[i]) > effectivePriority(sorted[j])
	})

	var all []PackageEntry
	officialSeen := false
	for _, rd := range sorted {
		if !rd.Enabled {
			continue
		}
		isOfficial := rd.Official && !officialSeen
		if isOfficial {
			officialSeen = true
		}
		all = append(all, PackageEntry{
			Name:     rd.Name,
			Home:     packageDefHome(rd),
			Priority: effectivePriority(rd),
			Official: isOfficial,
		})
	}
	return all
}

// packageDefHome resolves the local clone directory for a PackageDef.
// Path uses the versioned derived name: packages/<host>/<owner>/<repo>@<version>.
func packageDefHome(rd config.PackageDef) string {
	u, ver := config.ParseRef(rd.URL)
	if u == "" {
		u = rd.URL
	}
	if filepath.IsAbs(u) {
		return u
	}
	versionedName := filepath.FromSlash(config.DeriveVersionedName(u, ver))
	return filepath.Join(PackagesRoot(), versionedName)
}

// effectivePriority returns rd.Priority, defaulting to 100 for official and 50 for user packages.
func effectivePriority(rd config.PackageDef) int {
	if rd.Priority > 0 {
		return rd.Priority
	}
	if rd.Official {
		return 100
	}
	return 50
}

// AllSkillsPackages returns all configured packages' skills roots in priority order.
// Derived from AllPackages(); skips entries with no skills/ dir.
func AllSkillsPackages() []SkillsPackage {
	var pkgs []SkillsPackage
	for _, pkg := range AllPackages() {
		skillsRoot := filepath.Join(pkg.Home, "skills")
		if _, err := os.Stat(skillsRoot); err == nil {
			pkgs = append(pkgs, SkillsPackage{Name: pkg.Name, Root: skillsRoot})
		}
	}
	return pkgs
}

// listAllSkillsFromPackages is the unexported implementation with configurable body loading.
func listAllSkillsFromPackages(pkgs []SkillsPackage, withBody bool) ([]Skill, []SkillConflict, error) {
	seen := make(map[string]string) // canonical path → package name that claimed it
	var all []Skill
	var conflicts []SkillConflict

	for _, pkg := range pkgs {
		var list []Skill
		var err error
		if !withBody {
			// Install path: check on-disk cache before parsing SKILL.md files.
			// pkg.Root = <pkgHome>/skills/, so parent = pkgHome.
			pkgHome := filepath.Dir(pkg.Root)
			if cached, ok := readSkillCache(pkgHome); ok {
				list = cached
			} else {
				list, err = listAllSkills(pkg.Root, false)
				if err == nil {
					cp := append([]Skill(nil), list...)
					go writeSkillCache(pkgHome, cp)
				}
			}
		} else {
			list, err = listAllSkills(pkg.Root, true)
		}
		if err != nil {
			continue
		}
		for i := range list {
			sk := list[i]
			dirName := filepath.Base(sk.Path)
			key := sk.Domain + "/" + dirName
			if sk.Subdomain != "" {
				key = sk.Domain + "/" + sk.Subdomain + "/" + dirName
			}
			if winner, exists := seen[key]; exists {
				conflicts = append(conflicts, SkillConflict{
					CanonicalPath: key,
					WinnerPackage: winner,
					LoserPackage:  pkg.Name,
				})
				continue
			}
			seen[key] = pkg.Name
			sk.Package = pkg.Name
			all = append(all, sk)
		}
	}

	return all, conflicts, nil
}

// ListAllSkillsFromPackages lists skills from all packages, tagging each with its package.
// First (highest-priority) package to provide a given canonical path wins.
// Conflicts contains skills from lower-priority packages that were shadowed.
func ListAllSkillsFromPackages(pkgs []SkillsPackage) ([]Skill, []SkillConflict, error) {
	return listAllSkillsFromPackages(pkgs, true)
}

// ListAllSkillsFromPackagesMeta is identical to ListAllSkillsFromPackages but skips
// loading sk.Body — use for install paths that never access sk.Body.
func ListAllSkillsFromPackagesMeta(pkgs []SkillsPackage) ([]Skill, []SkillConflict, error) {
	return listAllSkillsFromPackages(pkgs, false)
}
