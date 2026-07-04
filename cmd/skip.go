package cmd

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"

	"github.com/jeffreytse/grimoire/internal/config"
)

var flagNoGitignore bool

// gitignoreSkip is a nil-safe skip predicate set by initGitignoreMatcher.
// Nil when --no-gitignore is set or no .gitignore patterns exist.
var gitignoreSkip func(rel string, isDir bool) bool

// globalCheckExclude holds merged --exclude + [core] check-exclude + [standards] exclude.
// Set by initExcludePatterns; used by watchers and file walkers.
var globalCheckExclude []string

// initGitignoreMatcher builds the gitignore matcher for projectDir and stores
// it in gitignoreSkip. Call once at command start before any goroutines launch.
func initGitignoreMatcher(projectDir string) {
	if flagNoGitignore {
		gitignoreSkip = nil
		return
	}
	m := buildGitignoreMatcher(projectDir)
	if m == nil {
		gitignoreSkip = nil
		return
	}
	gitignoreSkip = func(rel string, isDir bool) bool {
		return m.Match(strings.Split(filepath.ToSlash(rel), "/"), isDir)
	}
}

// initExcludePatterns merges --exclude flag, [core] check-exclude, and [standards] exclude
// into globalCheckExclude. Call once at command start after initGitignoreMatcher.
func initExcludePatterns(projectDir string, flagExcludes []string) {
	cfg, err := config.Load(projectDir)
	if err == nil {
		globalCheckExclude = append(globalCheckExclude, flagExcludes...)
		globalCheckExclude = append(globalCheckExclude, cfg.Core.CheckExclude...)
		globalCheckExclude = append(globalCheckExclude, cfg.CheckExclude...)
	} else {
		globalCheckExclude = flagExcludes
	}
}

// shouldSkip returns true if a relative path should be excluded from watching and checking.
func shouldSkip(rel string, isDir bool) bool {
	return (gitignoreSkip != nil && gitignoreSkip(rel, isDir)) ||
		matchesExcludePatterns(rel, globalCheckExclude)
}

// buildGitignoreMatcher walks projectDir, reads all .gitignore files with
// correct domain prefixes, and returns a Matcher. Returns nil if no patterns found.
func buildGitignoreMatcher(projectDir string) gitignore.Matcher {
	var patterns []gitignore.Pattern
	_ = filepath.WalkDir(projectDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return nil
		}
		name := d.Name()
		if path != projectDir {
			if strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules" ||
				name == "target" || name == "dist" || name == "build" || name == "__pycache__" {
				return filepath.SkipDir
			}
		}
		data, err := os.ReadFile(filepath.Join(path, ".gitignore"))
		if err != nil {
			return nil
		}
		var domain []string
		if path != projectDir {
			rel, _ := filepath.Rel(projectDir, path)
			domain = strings.Split(filepath.ToSlash(rel), "/")
		}
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			patterns = append(patterns, gitignore.ParsePattern(line, domain))
		}
		return nil
	})
	if len(patterns) == 0 {
		return nil
	}
	return gitignore.NewMatcher(patterns)
}
