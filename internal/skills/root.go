package skills

import (
	"os"
	"path/filepath"
)

const GrimoireRepo = "https://github.com/jeffreytse/grimoire-skills.git"

func GrimoireHome() string {
	if h := os.Getenv("GRIMOIRE_HOME"); h != "" {
		return h
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".grimoire")
	}
	return filepath.Join(home, ".grimoire")
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
