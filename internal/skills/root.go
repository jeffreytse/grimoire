package skills

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jeffreytse/grimoire/internal/settings"
)

const GrimoireRepo = "https://github.com/jeffreytse/grimoire-hub.git"

func GrimoireHome() string {
	if h := os.Getenv("GRIMOIRE_HOME"); h != "" {
		return h
	}
	cfg, _ := settings.LoadGlobal()
	if cfg.Core.Home != "" {
		return cfg.Core.Home
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".grimoire")
	}
	return filepath.Join(home, ".grimoire")
}

// RegistriesRoot returns the directory that contains all cloned registries.
// Neither GrimoireHome() nor RegistriesRoot() is itself a git repo.
func RegistriesRoot() string {
	return filepath.Join(GrimoireHome(), "registries")
}

// highestPriorityOfficialDef returns the official=true RegistryDef with the highest priority.
// Priorities are already normalized to 100/50 by Merge(), so comparison is valid.
// Cannot use AllRegistries() here — circular: AllRegistries() calls OfficialRegistryHome() in
// its no-registries fallback path.
func highestPriorityOfficialDef(regs []settings.RegistryDef) (settings.RegistryDef, bool) {
	var best settings.RegistryDef
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
//  2. highest-priority official=true entry in [[registry]]
//  3. Built-in GrimoireRepo constant
func GrimoireRepoURL() string {
	if s := os.Getenv("GRIMOIRE_SOURCE"); s != "" {
		if IsGitURL(s) || filepath.IsAbs(s) {
			return s
		}
	}
	cfg, _ := settings.LoadGlobal()
	if rd, ok := highestPriorityOfficialDef(cfg.Registries); ok && rd.URL != "" {
		u, _ := settings.ParseRef(rd.URL)
		if u == "" {
			u = rd.URL
		}
		if IsGitURL(u) || filepath.IsAbs(u) {
			return u
		}
	}
	return GrimoireRepo
}

// OfficialRegistryHome returns the local directory for the official registry.
// Uses the name from the highest-priority official=true [[registry]] entry as the subdir.
// When no [[registry]] is configured, derives the name from the GrimoireRepo constant.
// For local registries (absolute paths) the path itself is returned directly.
func OfficialRegistryHome() string {
	url := GrimoireRepoURL()
	if filepath.IsAbs(url) {
		return url
	}
	cfg, _ := settings.LoadGlobal()
	if rd, ok := highestPriorityOfficialDef(cfg.Registries); ok && rd.URL != "" {
		return filepath.Join(RegistriesRoot(), rd.Name)
	}
	return filepath.Join(RegistriesRoot(), settings.DeriveRegistryName(GrimoireRepo))
}

// IsGitURL reports whether s looks like a git remote URL.
func IsGitURL(s string) bool {
	return strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://") ||
		strings.HasPrefix(s, "git://") ||
		strings.HasPrefix(s, "git@")
}

// OfficialRegistryDerivedName returns the path-derived name for the official registry.
// Matches the directory name under RegistriesRoot().
func OfficialRegistryDerivedName() string {
	return settings.DeriveRegistryName(GrimoireRepoURL())
}

func SkillsRoot() string {
	return filepath.Join(OfficialRegistryHome(), "skills")
}

func GrimoireVersion() string {
	data, err := os.ReadFile(filepath.Join(OfficialRegistryHome(), "VERSION"))
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
