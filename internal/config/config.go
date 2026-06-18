package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// Global holds user-level grimoire settings persisted to
// ~/.config/grimoire/settings.toml.
type Global struct {
	// Home overrides the local directory where grimoire is installed (clone destination).
	// Takes precedence over Source when both are set to local paths.
	Home string `toml:"home,omitempty"`
	// Source overrides where grimoire pulls skills from.
	// A local path uses that directory as GrimoireHome;
	// a git URL is used in place of the default repo constant.
	Source string `toml:"source,omitempty"`
}

// GlobalPath returns the path to the global settings file.
func GlobalPath() string {
	if h := os.Getenv("XDG_CONFIG_HOME"); h != "" {
		return filepath.Join(h, "grimoire", "settings.toml")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "grimoire", "settings.toml")
}

// Load reads the global config. Returns a zero-value Global when the
// file is absent — callers should treat a missing file as defaults, not an error.
func Load() (Global, error) {
	var g Global
	data, err := os.ReadFile(GlobalPath())
	if errors.Is(err, os.ErrNotExist) {
		return g, nil
	}
	if err != nil {
		return g, err
	}
	return g, toml.Unmarshal(data, &g)
}

// Save writes g to the global config file, creating parent directories as needed.
func Save(g Global) error {
	path := GlobalPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := toml.Marshal(g)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
