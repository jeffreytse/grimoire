package agent

import (
	"os"
	"os/exec"
	"path/filepath"
)

// Detected returns the names of agents found in PATH.
func Detected() []string {
	var found []string
	for _, ag := range All {
		if _, err := exec.LookPath(ag); err == nil {
			found = append(found, ag)
		}
	}
	return found
}

// Version returns the version string of an agent binary, or "".
func Version(ag string) string {
	out, err := exec.Command(ag, "--version").Output()
	if err != nil {
		return ""
	}
	// extract first version-like token (digits and dots)
	s := string(out)
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			j := i
			for j < len(s) && (s[j] >= '0' && s[j] <= '9' || s[j] == '.') {
				j++
			}
			return s[i:j]
		}
	}
	return ""
}

// SkillCount returns how many skills are installed for an agent.
func SkillCount(ag string) int {
	dir := SkillsDir(ag)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if !hasHiddenPrefix(e.Name()) {
			count++
		}
	}
	return count
}

// BrokenSymlinkCount counts broken symlinks in an agent's skills dir.
func BrokenSymlinkCount(ag string) int {
	dir := SkillsDir(ag)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if hasHiddenPrefix(e.Name()) {
			continue
		}
		full := filepath.Join(dir, e.Name())
		if isBrokenSymlink(full) {
			count++
		}
	}
	return count
}

func isBrokenSymlink(path string) bool {
	// Lstat succeeds for symlinks even when the target is gone; Stat follows
	// the link and fails if the target is missing. Together they identify a
	// symlink whose target no longer exists.
	if _, err := os.Lstat(path); err != nil {
		return false
	}
	_, err := os.Stat(path)
	return err != nil
}

func hasHiddenPrefix(name string) bool {
	return name != "" && name[0] == '.'
}
