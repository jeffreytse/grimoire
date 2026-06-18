package detect

import (
	"os"
	"path/filepath"
)

// Profile returns the grimoire profile best matching the current directory.
func Profile(dir string) string {
	if dir == "" {
		dir = "."
	}

	checks := []struct {
		file    string
		profile string
	}{
		{"package.json", "engineering"},
		{"pyproject.toml", "engineering"},
		{"Cargo.toml", "engineering"},
		{"go.mod", "engineering"},
		{"pom.xml", "engineering"},
		{"build.gradle", "engineering"},
		{"Gemfile", "engineering"},
		{"requirements.txt", "engineering"},
	}

	for _, c := range checks {
		if fileExists(filepath.Join(dir, c.file)) {
			return c.profile
		}
	}

	// Glob for common source extensions
	for _, pattern := range []string{"*.py", "*.js", "*.ts", "*.go", "*.rs", "*.java", "*.rb"} {
		matches, _ := filepath.Glob(filepath.Join(dir, pattern))
		if len(matches) > 0 {
			return "engineering"
		}
	}

	// Mostly markdown → writing
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
