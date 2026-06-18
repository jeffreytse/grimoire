package skills

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jeffreytse/grimoire/internal/config"
)

const GrimoireRepo = "https://github.com/jeffreytse/grimoire-skills.git"

func GrimoireHome() string {
	if h := os.Getenv("GRIMOIRE_HOME"); h != "" {
		return h
	}
	cfg, _ := config.Load()
	if cfg.Home != "" {
		return cfg.Home
	}
	if cfg.Source != "" && !IsGitURL(cfg.Source) {
		return cfg.Source
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".grimoire")
	}
	return filepath.Join(home, ".grimoire")
}

// GrimoireRepoURL returns the git URL to clone/pull skills from.
// Falls back to the hardcoded GrimoireRepo constant when no custom URL is configured.
func GrimoireRepoURL() string {
	cfg, _ := config.Load()
	if cfg.Source != "" && IsGitURL(cfg.Source) {
		return cfg.Source
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
