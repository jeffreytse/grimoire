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

// GrimoireRepoURL returns the git URL or absolute local path for the official skills source.
// GRIMOIRE_SOURCE env var overrides all other configuration.
// Otherwise reads core.registry from global settings.
// Returns a git URL or an absolute path — callers must handle both.
func GrimoireRepoURL() string {
	if s := os.Getenv("GRIMOIRE_SOURCE"); s != "" {
		if IsGitURL(s) || filepath.IsAbs(s) {
			return s
		}
	}
	cfg, _ := settings.LoadGlobal()
	if cfg.Core.Registry != "" {
		u, _ := settings.ParseRef(cfg.Core.Registry)
		if u != "" && (IsGitURL(u) || filepath.IsAbs(u)) {
			return u
		}
	}
	return GrimoireRepo
}

// OfficialRegistryHome returns the local directory for the official registry.
// For git-hosted registries the path is derived from the URL and lives under RegistriesRoot.
// For local registries (absolute paths) the path itself is returned directly.
// Default: <RegistriesRoot>/jeffreytse/grimoire-hub
func OfficialRegistryHome() string {
	url := GrimoireRepoURL()
	if filepath.IsAbs(url) {
		return url
	}
	name := settings.DeriveRegistryName(url)
	return filepath.Join(RegistriesRoot(), name)
}

// ExtendsHome returns the local clone directory for a standards extends target.
// For absolute paths (local registries), the path itself is returned directly.
func ExtendsHome(name string) string {
	if filepath.IsAbs(name) {
		return name
	}
	return filepath.Join(RegistriesRoot(), name)
}

// IsGitURL reports whether s looks like a git remote URL.
func IsGitURL(s string) bool {
	return strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://") ||
		strings.HasPrefix(s, "git://") ||
		strings.HasPrefix(s, "git@")
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
