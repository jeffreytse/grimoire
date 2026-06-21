package detect

import (
	"os"
	"path/filepath"
	"strings"
)

// Profile returns the grimoire profile best matching the current directory.
// Returns "" when no confident signal is found.
func Profile(dir string) string {
	if dir == "" {
		dir = "."
	}

	// Named manifest files — unambiguous signals.
	manifests := []struct {
		file    string
		profile string
	}{
		// engineering
		{"package.json", "engineering"},
		{"pyproject.toml", "engineering"},
		{"Cargo.toml", "engineering"},
		{"go.mod", "engineering"},
		{"pom.xml", "engineering"},
		{"build.gradle", "engineering"},
		{"Gemfile", "engineering"},
		{"requirements.txt", "engineering"},
		{"Makefile", "engineering"},
		{"Dockerfile", "engineering"},
		// design (UI/UX/graphic)
		{"*.fig", "design"},
		{"*.sketch", "design"},
		{"*.xd", "design"},
		// art (digital painting / illustration)
		{"*.procreate", "art"},
		{"*.kra", "art"}, // Krita
		{"*.xcf", "art"}, // GIMP
		{"*.csp", "art"}, // Clip Studio Paint
		// film / video production
		{"*.prproj", "film"},    // Adobe Premiere Pro
		{"*.aep", "film"},       // Adobe After Effects
		{"*.drp", "film"},       // DaVinci Resolve
		{"*.fcpbundle", "film"}, // Final Cut Pro (bundle dir)
		// music production
		{"*.als", "music"},   // Ableton Live
		{"*.flp", "music"},   // FL Studio
		{"*.logic", "music"}, // Logic Pro (bundle dir)
	}

	for _, c := range manifests {
		if strings.ContainsRune(c.file, '*') {
			matches, _ := filepath.Glob(filepath.Join(dir, c.file))
			if len(matches) > 0 {
				return c.profile
			}
		} else if fileExists(filepath.Join(dir, c.file)) {
			return c.profile
		}
	}

	// Source code globs → engineering
	for _, pattern := range []string{"*.py", "*.js", "*.ts", "*.go", "*.rs", "*.java", "*.rb", "*.cs", "*.cpp", "*.c"} {
		matches, _ := filepath.Glob(filepath.Join(dir, pattern))
		if len(matches) > 0 {
			return "engineering"
		}
	}

	// Jupyter notebooks → science/data analysis
	if nb, _ := filepath.Glob(filepath.Join(dir, "*.ipynb")); len(nb) > 0 {
		return "science"
	}

	// R scripts/markdown → science
	for _, pattern := range []string{"*.r", "*.R", "*.rmd", "*.Rmd"} {
		if m, _ := filepath.Glob(filepath.Join(dir, pattern)); len(m) > 0 {
			return "science"
		}
	}

	// Raw photo formats → photography
	for _, pattern := range []string{"*.raw", "*.RAW", "*.cr2", "*.CR2", "*.nef", "*.NEF", "*.arw", "*.ARW", "*.dng", "*.DNG"} {
		if m, _ := filepath.Glob(filepath.Join(dir, pattern)); len(m) > 0 {
			return "photography"
		}
	}

	// Mostly markdown → writing (3+ files)
	mdMatches, _ := filepath.Glob(filepath.Join(dir, "*.md"))
	if len(mdMatches) > 2 {
		return "writing"
	}

	return ""
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
