package config

import (
	"net/url"
	"path/filepath"
	"strings"
)

// ParseRef parses a package ref string and returns the full git clone URL and version.
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

// DerivePackageName derives a versioned local directory name from a git URL.
// Used as the path component under ~/.grimoire/packages/.
// Always includes host prefix and version: <host>/<owner>/<repo>@<version>.
// When no version is present in rawURL, "@latest" is appended (floating HEAD).
//
//	https://github.com/acmecorp/standards.git        → "github.com/acmecorp/standards@latest"
//	https://github.com/acmecorp/standards.git@v1.0.0 → "github.com/acmecorp/standards@v1.0.0"
//	git@github.com:acmecorp/standards.git            → "github.com/acmecorp/standards@latest"
//	acmecorp/standards@v1.0.0 (shorthand)            → "github.com/acmecorp/standards@v1.0.0"
func DerivePackageName(rawURL string) string {
	if filepath.IsAbs(rawURL) {
		return filepath.ToSlash(rawURL) // local package: absolute path is the name
	}
	base, ver := splitRef(rawURL)
	if ver == "" {
		ver = "latest"
	}
	return deriveBaseName(base) + "@" + ver
}

// DeriveVersionedName builds a versioned package name from a base URL (no @tag)
// and an explicit version string. Use this when the version is known separately
// (e.g. from PackageRef.Tag). If version is empty, "@latest" is used.
func DeriveVersionedName(baseURL, version string) string {
	if filepath.IsAbs(baseURL) {
		return filepath.ToSlash(baseURL)
	}
	if version == "" {
		version = "latest"
	}
	return deriveBaseName(baseURL) + "@" + version
}

// deriveBaseName returns "<host>/<owner>/<repo>" from a git URL (no version).
func deriveBaseName(rawURL string) string {
	// shorthand "owner/repo" — expand to github.com
	if !strings.Contains(rawURL, "://") && !strings.HasPrefix(rawURL, "git@") {
		return "github.com/" + strings.TrimSuffix(rawURL, ".git")
	}
	// SSH: git@<host>:<path>.git
	if strings.HasPrefix(rawURL, "git@") {
		raw := strings.TrimPrefix(rawURL, "git@")
		host, path, _ := strings.Cut(raw, ":")
		return host + "/" + strings.TrimSuffix(path, ".git")
	}
	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		return rawURL
	}
	return u.Host + "/" + strings.TrimSuffix(strings.Trim(u.Path, "/"), ".git")
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
