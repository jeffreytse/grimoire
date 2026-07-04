package skills

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jeffreytse/grimoire/internal/config"
)

const GrimoireRepo = "https://github.com/jeffreytse/grimoire-core.git"

func GrimoireHome() string {
	if h := os.Getenv("GRIMOIRE_HOME"); h != "" {
		return h
	}
	cfg, _ := config.LoadGlobal()
	if cfg.Core.Home != "" {
		return cfg.Core.Home
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".grimoire")
	}
	return filepath.Join(home, ".grimoire")
}

// PackagesRoot returns the directory that contains all cloned skill packages.
// Neither GrimoireHome() nor PackagesRoot() is itself a git repo.
func PackagesRoot() string {
	return filepath.Join(GrimoireHome(), "packages")
}

// highestPriorityOfficialDef returns the official=true PackageDef with the highest priority.
// Priorities are already normalized to 100/50 by Merge(), so comparison is valid.
// Cannot use AllPackages() here — circular: AllPackages() calls OfficialPackageHome() in
// its no-packages fallback path.
func highestPriorityOfficialDef(regs []config.PackageDef) (config.PackageDef, bool) {
	var best config.PackageDef
	found := false
	for _, rd := range regs {
		if rd.Official && (!found || rd.Priority > best.Priority) {
			best = rd
			found = true
		}
	}
	return best, found
}

// GrimoireRepoURL returns the git URL or absolute local path for the official skills source.
// Resolution order:
//  1. GRIMOIRE_SOURCE env var
//  2. highest-priority official=true entry in [[package]]
//  3. Built-in GrimoireRepo constant
func GrimoireRepoURL() string {
	if s := os.Getenv("GRIMOIRE_SOURCE"); s != "" {
		if IsGitURL(s) || filepath.IsAbs(s) {
			return s
		}
	}
	cfg, _ := config.LoadGlobal()
	if rd, ok := highestPriorityOfficialDef(cfg.Packages); ok && rd.URL != "" {
		u, _ := config.ParseRef(rd.URL)
		if u == "" {
			u = rd.URL
		}
		if IsGitURL(u) || filepath.IsAbs(u) {
			return u
		}
	}
	return GrimoireRepo
}

// OfficialPackageHome returns the local directory for the official skill package.
// Uses the name from the highest-priority official=true [[package]] entry as the subdir.
// When no [[package]] is configured, derives the name from the GrimoireRepo constant.
// For local packages (absolute paths) the path itself is returned directly.
func OfficialPackageHome() string {
	url := GrimoireRepoURL()
	if filepath.IsAbs(url) {
		return url
	}
	cfg, _ := config.LoadGlobal()
	if rd, ok := highestPriorityOfficialDef(cfg.Packages); ok && rd.URL != "" {
		return filepath.Join(PackagesRoot(), rd.Name)
	}
	return filepath.Join(PackagesRoot(), config.DerivePackageName(GrimoireRepo))
}

// IsGitURL reports whether s looks like a git remote URL.
func IsGitURL(s string) bool {
	return strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://") ||
		strings.HasPrefix(s, "git://") ||
		strings.HasPrefix(s, "git@")
}

// OfficialPackageDerivedName returns the path-derived name for the official skill package.
// Matches the directory name under PackagesRoot().
func OfficialPackageDerivedName() string {
	return config.DerivePackageName(GrimoireRepoURL())
}

func SkillsRoot() string {
	return filepath.Join(OfficialPackageHome(), "skills")
}

func GrimoireVersion() string {
	data, err := os.ReadFile(filepath.Join(OfficialPackageHome(), "VERSION"))
	if err != nil {
		return "unknown"
	}
	v := string(data)
	for i, c := range v {
		if c == '\n' || c == '\r' || c == ' ' {
			return v[:i]
		}
	}
	return v
}
