package settings

import (
	"net/url"
	"path/filepath"
	"strings"
)

// ParseRef parses a registry ref string and returns the full git clone URL and version.
//
// Supported formats:
//   - "owner/repo"            → https://github.com/owner/repo.git, version ""
//   - "owner/repo@v1.0.0"    → https://github.com/owner/repo.git, version "v1.0.0"
//   - "https://host/r.git"   → https://host/r.git, version ""
//   - "https://host/r.git@main" → https://host/r.git, version "main"
//   - "git@github.com:org/r.git@v1" → git@github.com:org/r.git, version "v1"
func ParseRef(s string) (gitURL, version string) {
	gitURL, version = splitRef(s)
	if filepath.IsAbs(gitURL) {
		return // local path — no git URL expansion
	}
	// expand shorthand owner/repo → https://github.com/owner/repo.git
	if !strings.Contains(gitURL, "://") && !strings.HasPrefix(gitURL, "git@") {
		gitURL = "https://github.com/" + strings.TrimSuffix(gitURL, ".git") + ".git"
	}
	return
}

// DeriveRegistryName derives a stable local directory name from a git URL.
// Used as the path component under ~/.grimoire/registries/.
//
// GitHub URLs use the shorthand "owner/repo" (no host prefix) for brevity.
// All other hosts include the hostname to prevent collisions between identically
// named repos on different services:
//
//	https://github.com/acmecorp/standards.git  → "acmecorp/standards"
//	https://gitlab.com/acmecorp/standards.git  → "gitlab.com/acmecorp/standards"
//	git@github.com:acmecorp/standards.git      → "acmecorp/standards"
//	git@gitlab.com:acmecorp/standards.git      → "gitlab.com/acmecorp/standards"
func DeriveRegistryName(rawURL string) string {
	if filepath.IsAbs(rawURL) {
		return filepath.ToSlash(rawURL) // local registry: absolute path is the name
	}
	// shorthand "owner/repo" — no host present
	if !strings.Contains(rawURL, "://") && !strings.HasPrefix(rawURL, "git@") {
		return strings.TrimSuffix(rawURL, ".git")
	}
	// SSH: git@<host>:<path>.git
	if strings.HasPrefix(rawURL, "git@") {
		raw := strings.TrimPrefix(rawURL, "git@")
		host, path, _ := strings.Cut(raw, ":")
		path = strings.TrimSuffix(path, ".git")
		if host == "github.com" {
			return path
		}
		return host + "/" + path
	}
	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		return rawURL
	}
	path := strings.TrimSuffix(strings.Trim(u.Path, "/"), ".git")
	if u.Host == "github.com" {
		return path
	}
	return u.Host + "/" + path
}

// splitRef splits a ref string at the last '@' that appears after the last '/'.
// This correctly handles git@ SSH URLs (their '@' appears before the last '/').
func splitRef(s string) (rawURL, version string) {
	lastSlash := strings.LastIndex(s, "/")
	atIdx := strings.LastIndex(s, "@")
	if atIdx > lastSlash && atIdx != -1 {
		return s[:atIdx], s[atIdx+1:]
	}
	return s, ""
}
