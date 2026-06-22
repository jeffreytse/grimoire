package skills

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ListPresets returns preset names found under <registryHome>/presets/.
// Each preset is a subdirectory containing at least a settings.toml.
func ListPresets(registryHome string) []string {
	dir := filepath.Join(registryHome, "presets")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
			names = append(names, e.Name())
		}
	}
	return names
}

// ApplyPreset copies a preset's settings.toml and optional profiles/ into destDir.
// destDir must already exist. The preset must contain a settings.toml at minimum.
func ApplyPreset(registryHome, presetName, destDir string) error {
	presetDir := filepath.Join(registryHome, "presets", presetName)

	src := filepath.Join(presetDir, "settings.toml")
	dst := filepath.Join(destDir, "settings.toml")
	if err := copyFile(src, dst); err != nil {
		return fmt.Errorf("copying settings.toml: %w", err)
	}

	profilesSrc := filepath.Join(presetDir, "profiles")
	if _, err := os.Stat(profilesSrc); err == nil {
		profilesDst := filepath.Join(destDir, "profiles")
		if err := os.MkdirAll(profilesDst, 0o755); err != nil {
			return fmt.Errorf("creating profiles dir: %w", err)
		}
		fentries, _ := os.ReadDir(profilesSrc)
		for _, fe := range fentries {
			if !fe.IsDir() && strings.HasSuffix(fe.Name(), ".toml") {
				if err := copyFile(
					filepath.Join(profilesSrc, fe.Name()),
					filepath.Join(profilesDst, fe.Name()),
				); err != nil {
					return fmt.Errorf("copying profile %s: %w", fe.Name(), err)
				}
			}
		}
	}

	return nil
}
