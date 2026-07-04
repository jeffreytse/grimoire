package manifest

import (
	"os"
	"path/filepath"
	"runtime"
)

// ProjectPath returns ./grimoire.toml relative to the given project directory.
func ProjectPath(projectDir string) string {
	return filepath.Join(projectDir, "grimoire.toml")
}

// GlobalPath returns the user-global grimoire.toml, respecting XDG_CONFIG_HOME.
func GlobalPath() string {
	if h := os.Getenv("XDG_CONFIG_HOME"); h != "" {
		return filepath.Join(h, "grimoire", "grimoire.toml")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "grimoire", "grimoire.toml")
}

// SystemPath returns the system-wide grimoire.toml.
// Linux/macOS: /etc/grimoire/grimoire.toml
// Windows: %PROGRAMDATA%\grimoire\grimoire.toml.
func SystemPath() string {
	if runtime.GOOS == "windows" {
		pd := os.Getenv("PROGRAMDATA")
		if pd == "" {
			pd = `C:\ProgramData`
		}
		return filepath.Join(pd, "grimoire", "grimoire.toml")
	}
	return "/etc/grimoire/grimoire.toml"
}

// LockPath returns the grimoire.lock path alongside the given grimoire.toml path.
func LockPath(manifestPath string) string {
	return filepath.Join(filepath.Dir(manifestPath), "grimoire.lock")
}
