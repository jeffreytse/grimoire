package skills

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jeffreytse/grimoire/internal/settings"
)

const GrimoireRepo = "https://github.com/jeffreytse/grimoire-skills.git"

func GrimoireHome() string {
	if h := os.Getenv("GRIMOIRE_HOME"); h != "" {
		return h
	}
	cfg, _ := settings.LoadGlobal()
	if cfg.Core.Home != "" {
		return cfg.Core.Home
	}
	if cfg.Core.Source != "" && !IsGitURL(cfg.Core.Source) {
		return cfg.Core.Source
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".grimoire")
	}
	return filepath.Join(home, ".grimoire")
}

// GrimoireRepoURL returns the git URL to clone/pull skills from.
// source is global-only: a project-level override would affect the shared GrimoireHome()
// for all projects on this machine, causing cross-project contamination.
func GrimoireRepoURL() string {
	cfg, _ := settings.LoadGlobal()
	if cfg.Core.Source != "" && IsGitURL(cfg.Core.Source) {
		return cfg.Core.Source
	}
	return GrimoireRepo
}

// IsGitURL reports whether s looks like a git remote URL.
func IsGitURL(s string) bool {
	return strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://") ||
		strings.HasPrefix(s, "git://") ||
		strings.HasPrefix(s, "git@")
}

func SkillsRoot() string {
	return filepath.Join(GrimoireHome(), "skills")
}

func GrimoireVersion() string {
	data, err := os.ReadFile(filepath.Join(GrimoireHome(), "VERSION"))
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
