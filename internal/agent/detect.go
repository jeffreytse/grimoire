package agent

import (
	"os"
	"os/exec"
	"path/filepath"
)

// Binary returns the CLI binary name for an agent key (may differ from the key itself).
func Binary(ag string) string { return agentBinary(ag) }

// agentBinary returns the CLI binary name for an agent key.
// Most agents share their binary name with their key; exceptions are listed here.
func agentBinary(ag string) string {
	switch ag {
	case "antigravity":
		return "agy"
	default:
		return ag
	}
}

// Detected returns the names of agents found in PATH.
func Detected() []string {
	var found []string
	for _, ag := range All {
		if _, err := exec.LookPath(agentBinary(ag)); err == nil {
			found = append(found, ag)
		}
	}
	return found
}

// DetectedOrInstalled returns agents that are either detectable via PATH or have a
// non-empty skills directory. Used by uninstall so cleanup covers agents whose binary
// is no longer in PATH but whose skills were installed in a prior run.
func DetectedOrInstalled() []string {
	seen := map[string]bool{}
	var found []string
	for _, ag := range All {
		if _, err := exec.LookPath(agentBinary(ag)); err == nil {
			found = append(found, ag)
			seen[ag] = true
		}
	}
	for _, ag := range All {
		if seen[ag] {
			continue
		}
		dir := SkillsDir(ag)
		if dir == "" {
			continue
		}
		if entries, err := os.ReadDir(dir); err == nil && len(entries) > 0 {
			found = append(found, ag)
			seen[ag] = true
		}
	}
	return found
}

// DetectedCheckAgents returns check-capable agent names found on the system,
// in CheckAgents order. "copilot" requires the "gh" binary (gh copilot extension).
func DetectedCheckAgents() []string {
	var found []string
	for _, ag := range CheckAgents {
		binary := agentBinary(ag)
		if ag == "copilot" {
			binary = "gh"
		}
		if _, err := exec.LookPath(binary); err == nil {
			found = append(found, ag)
		}
	}
	return found
}

// Version returns the version string of an agent binary, or "".
func Version(ag string) string {
	out, err := exec.Command(agentBinary(ag), "--version").Output()
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
