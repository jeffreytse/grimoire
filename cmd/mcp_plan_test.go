package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// planFiles returns the files a plan says it would edit or write.
func planFiles(steps []string) []string {
	var files []string
	for _, s := range steps {
		for _, verb := range []string{"edit ", "write "} {
			if rest, ok := strings.CutPrefix(s, verb); ok {
				files = append(files, strings.SplitN(rest, " (", 2)[0])
			}
		}
	}
	return files
}

// The dry-run plan must match what a real install does: every file it names
// exists after a real install, and the dry run itself changes nothing.
func TestMCPDryRun_MatchesRealInstall(t *testing.T) {
	for _, key := range []string{"claude", "codex", "gemini", "opencode", "pi", "antigravity", "copilot", "kimi", "omp"} {
		t.Run(key, func(t *testing.T) {
			f, home := setupMCPSetupEnv(t, true)
			flagMCPSetupTarget = key
			t.Cleanup(func() { flagMCPSetupDryRun = false })

			targets, err := resolveMCPAgents()
			if err != nil {
				t.Fatal(err)
			}
			files := planFiles(mcpPlan(targets[0], flagMCPSetupScope, "grimoire", true))
			if len(files) == 0 && targets[0].configPath != nil {
				t.Fatal("plan names no files")
			}

			flagMCPSetupDryRun = true
			if err := runMCPInstall(nil, nil); err != nil {
				t.Fatal(err)
			}
			if len(f.calls) != 0 {
				t.Errorf("dry run ran commands: %v", f.calls)
			}
			for _, p := range files {
				if _, err := os.Stat(p); err == nil {
					t.Errorf("dry run created %s", p)
				}
			}

			flagMCPSetupDryRun = false
			if err := runMCPInstall(nil, nil); err != nil {
				t.Fatal(err)
			}
			for _, p := range files {
				if _, err := os.Stat(p); err != nil {
					t.Errorf("plan names %s, but a real install didn't write it", p)
				}
				if !strings.HasPrefix(p, home) && !strings.HasPrefix(p, getProjectDir()) {
					t.Errorf("plan path %s is outside the sandbox", filepath.Clean(p))
				}
			}
		})
	}
}
