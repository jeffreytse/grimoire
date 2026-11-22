package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Plugin agents load skills only from a plugin bundle, not from a plain skills
// directory. grimoire installs skills into a grimoire-managed bundle
// (<bundle>/skills/<name>/), then FinalizePlugin writes the agent's manifest
// and registers the bundle through the agent's own plugin CLI.
//
//   - minimax: MiniMax Code loads every directory under
//     ${MINIMAX_DATA_DIR:-~/.minimax}/plugins/ as a local plugin through its
//     Claude Code reader, so the bundle lives there and needs no CLI step.
//   - grok:    `grok plugin install <bundle> --trust` copies the bundle, so a
//     refresh is uninstall (--keep-data) + install + enable.
//   - muse:    `muse plugins install <bundle>`; Muse requires every skill to be
//     listed in .muse-plugin/plugin.json.
//
// Plugin agents are global-only: they have no project skills directory.

const pluginName = "grimoire"
const pluginVersion = "1.0.0"
const pluginDesc = "Every field's professional best practices, installed by the grimoire CLI."

// IsPluginAgent reports whether ag installs through a plugin bundle.
func IsPluginAgent(ag string) bool {
	switch ag {
	case "grok", "muse", "minimax":
		return true
	}
	return false
}

// PluginBundleDir returns the bundle root grimoire manages for a plugin agent.
func PluginBundleDir(ag string) string {
	home := mustHome()
	switch ag {
	case "minimax":
		return filepath.Join(envDir("MINIMAX_DATA_DIR", filepath.Join(home, ".minimax")), "plugins", pluginName)
	case "grok", "muse":
		return filepath.Join(home, ".grimoire", "agent-plugins", ag, pluginName)
	}
	return ""
}

// FinalizePlugin writes the agent-specific manifest for the bundle and
// registers it with the agent's CLI. It is a no-op for non-plugin agents.
// run executes an external command; pass nil to use os/exec.
func FinalizePlugin(ag string, run func(name string, args ...string) error) error {
	if !IsPluginAgent(ag) {
		return nil
	}
	if run == nil {
		run = runCmd
	}
	bundle := PluginBundleDir(ag)
	names, err := bundleSkillNames(filepath.Join(bundle, "skills"))
	if err != nil {
		return err
	}
	if len(names) == 0 {
		return nil
	}
	if err := writePluginManifest(ag, bundle, names); err != nil {
		return err
	}

	switch ag {
	case "grok":
		if _, err := exec.LookPath("grok"); err != nil {
			return fmt.Errorf("grok not in PATH; register manually: grok plugin install %q --trust && grok plugin enable %s", bundle, pluginName)
		}
		_ = run("grok", "plugin", "uninstall", pluginName, "--keep-data") // absent on first install
		if err := run("grok", "plugin", "install", bundle, "--trust"); err != nil {
			return fmt.Errorf("grok plugin install: %w", err)
		}
		if err := run("grok", "plugin", "enable", pluginName); err != nil {
			return fmt.Errorf("grok plugin enable: %w", err)
		}
	case "muse":
		if _, err := exec.LookPath("muse"); err != nil {
			return fmt.Errorf("muse not in PATH; register manually: muse plugins install %q", bundle)
		}
		if err := run("muse", "plugins", "install", bundle); err != nil {
			return fmt.Errorf("muse plugins install: %w", err)
		}
	}
	return nil
}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func bundleSkillNames(skillsDir string) ([]string, error) {
	entries, err := os.ReadDir(skillsDir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if hasHiddenPrefix(e.Name()) {
			continue
		}
		if _, err := os.Stat(filepath.Join(skillsDir, e.Name(), "SKILL.md")); err == nil {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

type museSkill struct {
	ID             string `json:"id"`
	Path           string `json:"path"`
	EnabledDefault bool   `json:"enabledDefault"`
}

func writePluginManifest(ag, bundle string, names []string) error {
	base := map[string]any{
		"name":        pluginName,
		"version":     pluginVersion,
		"description": pluginDesc,
	}
	var rel string
	var data any
	switch ag {
	case "minimax":
		rel, data = filepath.Join(".claude-plugin", "plugin.json"), base
	case "grok":
		base["skills"] = []string{"./skills"}
		rel, data = filepath.Join(".grok-plugin", "plugin.json"), base
	case "muse":
		skills := make([]museSkill, len(names))
		for i, n := range names {
			skills[i] = museSkill{ID: n, Path: "skills/" + n + "/SKILL.md", EnabledDefault: true}
		}
		rel = filepath.Join(".muse-plugin", "plugin.json")
		data = map[string]any{
			"schemaVersion": 1,
			"name":          pluginName,
			"version":       pluginVersion,
			"description":   pluginDesc,
			"compat":        map[string]string{"source": "native", "manifestDir": ".muse-plugin"},
			"capabilities":  map[string]any{"skills": skills},
		}
	default:
		return nil
	}
	path := filepath.Join(bundle, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}
