package config

import (
	"path/filepath"
	"strings"
)

// PackageRef is a parsed package URL of the form:
//
//	[host/][owner/repo][@version][:globbing_path]
//
// User-facing name: grimoire-ref. Format: [host/][owner/repo][@version][:glob_path]
// Single segment with no '/' and no ':' = path within the official package.
// Local paths (Unix absolute/relative, Windows absolute) are detected first.
type PackageRef struct {
	Raw         string
	Local       bool
	LocalPath   string // normalized filesystem path (ToSlash); empty if not local
	Host        string // e.g. "github.com"; empty if local or official shorthand
	Owner       string
	Repo        string
	Tag         string
	Path        string // glob pattern within package/local root; empty = all
	PackageName string // e.g. "github.com/acmecorp/standards"; derived from URL
	PackageURL  string // e.g. "https://github.com/acmecorp/standards.git"
}

// IsOfficialRepoPath reports whether this ref is a path shorthand into the
// official package — no host/owner/repo specified.
func (r PackageRef) IsOfficialRepoPath() bool { //nolint:gocritic
	return !r.Local && r.Owner == "" && r.Path != ""
}

// IsLocal reports whether this ref is a local filesystem path.
func (r PackageRef) IsLocal() bool { return r.Local } //nolint:gocritic

// ParsePackageRef parses a package URL string into a PackageRef.
func ParsePackageRef(s string) PackageRef {
	ref := PackageRef{Raw: s}

	// 1. Local path detection — checked before any ':' split.
	if isLocalPath(s) {
		ref.Local = true
		// Find glob delimiter: first ':' not at index 1 (Windows drive letter).
		start := 0
		if isWindowsAbsPath(s) {
			start = 2 // skip the drive-letter ':'
		}
		if idx := strings.Index(s[start:], ":"); idx >= 0 {
			colonIdx := start + idx
			ref.LocalPath = filepath.ToSlash(s[:colonIdx])
			ref.Path = s[colonIdx+1:]
		} else {
			ref.LocalPath = filepath.ToSlash(s)
		}
		return ref
	}

	// 2. Split on first ':' → locator : glob path.
	locator := s
	if colonIdx := strings.Index(s, ":"); colonIdx >= 0 {
		locator = s[:colonIdx]
		ref.Path = s[colonIdx+1:]
	} else if !strings.Contains(s, "/") {
		// 3. Single segment, no ':', no '/': official package path shorthand.
		ref.Path = s
		return ref
	}
	// 4. If no ':' and has '/', locator = full string; Path = "" (all skills).

	// 5. Parse locator — extract optional @version suffix.
	lastSlash := strings.LastIndex(locator, "/")
	if atIdx := strings.LastIndex(locator, "@"); atIdx > lastSlash && atIdx > 0 {
		ref.Tag = locator[atIdx+1:]
		locator = locator[:atIdx]
	}

	if locator == "" {
		return ref
	}

	segments := strings.Split(locator, "/")
	if strings.Contains(segments[0], ".") {
		// Explicit host (contains '.')
		ref.Host = segments[0]
		if len(segments) >= 3 {
			ref.Owner = segments[1]
			ref.Repo = segments[2]
		}
	} else {
		// Implicit github.com
		ref.Host = "github.com"
		if len(segments) >= 2 {
			ref.Owner = segments[0]
			ref.Repo = segments[1]
		}
	}

	if ref.Owner != "" && ref.Repo != "" {
		ref.PackageURL = "https://" + ref.Host + "/" + ref.Owner + "/" + ref.Repo + ".git"
		ref.PackageName = DeriveVersionedName(ref.PackageURL, ref.Tag)
	}

	return ref
}

func isWindowsAbsPath(s string) bool {
	return len(s) >= 3 && s[1] == ':' && (s[2] == '\\' || s[2] == '/')
}

func isLocalPath(s string) bool {
	return strings.HasPrefix(s, "/") ||
		strings.HasPrefix(s, "./") ||
		strings.HasPrefix(s, "../") ||
		isWindowsAbsPath(s)
}
