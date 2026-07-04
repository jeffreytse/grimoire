package skills

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type skillCacheFile struct {
	HeadCommit string  `json:"head_commit"` // resolved HEAD commit SHA
	Skills     []Skill `json:"skills"`
}

// readSkillCache returns cached skills for pkgHome if .grimoire-cache.json exists
// and its HEAD commit SHA matches the current HEAD. Returns nil, false on any miss.
// Only valid for git-backed packages (non-git packages always return false).
func readSkillCache(pkgHome string) ([]Skill, bool) {
	commit := headCommit(pkgHome)
	if commit == "" {
		return nil, false
	}

	data, err := os.ReadFile(filepath.Join(pkgHome, ".grimoire-cache.json"))
	if err != nil {
		return nil, false
	}

	var cache skillCacheFile
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, false
	}

	if cache.HeadCommit != commit {
		return nil, false
	}

	return cache.Skills, true
}

// writeSkillCache persists skills to <pkgHome>/.grimoire-cache.json, keyed by
// the current HEAD commit SHA. Silently no-ops for non-git packages or on error.
// Intended to be called in a goroutine — does not block the caller.
// Uses an atomic temp-file-then-rename write to prevent partial reads.
func writeSkillCache(pkgHome string, skills []Skill) {
	commit := headCommit(pkgHome)
	if commit == "" {
		return
	}

	cache := skillCacheFile{
		HeadCommit: commit,
		Skills:     skills,
	}

	data, err := json.Marshal(cache)
	if err != nil {
		return
	}

	tmp := filepath.Join(pkgHome, ".grimoire-cache.json.tmp")
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return
	}
	_ = os.Rename(tmp, filepath.Join(pkgHome, ".grimoire-cache.json"))
}

// InvalidateSkillCache deletes the on-disk skill cache for pkgHome.
// Call after any git operation that changes package content (pull, checkout, clone).
// Safe to call when no cache exists — silently no-ops.
func InvalidateSkillCache(pkgHome string) {
	_ = os.Remove(filepath.Join(pkgHome, ".grimoire-cache.json"))
}

// headCommit returns the resolved HEAD commit SHA for a git-backed package,
// or "" for non-git packages or on any read error.
// Handles both symbolic refs (branches) and detached HEAD (SHA directly).
func headCommit(pkgHome string) string {
	data, err := os.ReadFile(filepath.Join(pkgHome, ".git", "HEAD"))
	if err != nil {
		return ""
	}
	line := strings.TrimSpace(string(data))
	if strings.HasPrefix(line, "ref: ") {
		ref := strings.TrimPrefix(line, "ref: ")
		sha, err := os.ReadFile(filepath.Join(pkgHome, ".git", ref))
		if err != nil {
			return resolvePackedRef(pkgHome, ref)
		}
		return strings.TrimSpace(string(sha))
	}
	return line // detached HEAD — line is the SHA
}

// resolvePackedRef looks up ref in .git/packed-refs, used as fallback when
// the loose ref file is absent (e.g. after git gc or shallow clones).
func resolvePackedRef(pkgHome, ref string) string {
	data, err := os.ReadFile(filepath.Join(pkgHome, ".git", "packed-refs"))
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasSuffix(line, " "+ref) {
			fields := strings.Fields(line)
			if len(fields) >= 1 {
				return fields[0]
			}
		}
	}
	return ""
}
