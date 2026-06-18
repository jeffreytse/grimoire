package agent

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const triggerLine = "Always invoke `start-best-practice` before responding to any user request."
const sectionHeader = "## Grimoire"

// ConfigureAgentMD appends the grimoire trigger to the agent's config file.
// Does nothing if already present or the config dir doesn't exist.
func ConfigureAgentMD(ag string) error {
	cfgDir := ConfigDir(ag)
	cfgFile := ConfigFile(ag)
	if cfgDir == "" || cfgFile == "" {
		return nil
	}
	if _, err := os.Stat(cfgDir); err != nil {
		return nil // agent not set up
	}
	if alreadyConfigured(cfgFile) {
		return nil
	}
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", cfgDir, err)
	}
	f, err := os.OpenFile(cfgFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("opening %s: %w", cfgFile, err)
	}
	defer func() { _ = f.Close() }()
	_, err = fmt.Fprintf(f, "\n%s\n%s\n", sectionHeader, triggerLine)
	return err
}

// RemoveAgentMDConfig removes the grimoire trigger lines from the agent's config file.
func RemoveAgentMDConfig(ag string) error {
	cfgFile := ConfigFile(ag)
	if cfgFile == "" {
		return nil
	}
	if !alreadyConfigured(cfgFile) {
		return nil
	}
	f, err := os.Open(cfgFile)
	if err != nil {
		return err
	}
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == sectionHeader || line == triggerLine {
			continue
		}
		lines = append(lines, line)
	}
	_ = f.Close()
	if err := scanner.Err(); err != nil {
		return err
	}
	// trim trailing blank lines added by grimoire
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	content := strings.Join(lines, "\n")
	if content != "" {
		content += "\n"
	}
	return os.WriteFile(cfgFile, []byte(content), 0o644)
}

// IsConfigured reports whether the agent's config file has the grimoire trigger.
func IsConfigured(ag string) bool {
	return alreadyConfigured(ConfigFile(ag))
}

func alreadyConfigured(path string) bool {
	// Substring match (not exact line equality) so that manual edits to the
	// config file — e.g. extra whitespace or surrounding text — don't cause
	// grimoire to append a duplicate trigger.
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer func() { _ = f.Close() }()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "start-best-practice") {
			return true
		}
	}
	return false
}
