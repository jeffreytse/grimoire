package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/manifest"
)

var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Regenerate grimoire.lock from the current grimoire.toml without installing",
	RunE:  runLock,
}

func init() {
	rootCmd.AddCommand(lockCmd)
}

func runLock(_ *cobra.Command, _ []string) error {
	projectDir := getProjectDir()
	manifestPath := manifest.ProjectPath(projectDir)

	if _, err := os.Stat(manifestPath); err != nil {
		return fmt.Errorf("grimoire.toml not found in %s — run: grimoire init", projectDir)
	}

	r, err := manifest.Load(projectDir)
	if err != nil {
		return fmt.Errorf("loading grimoire.toml: %w", err)
	}

	if len(r.Deps) == 0 && len(r.DevDeps) == 0 {
		fmt.Println("grimoire.toml has no [dependencies] — nothing to lock")
		return nil
	}

	// Build skill list from packages without installing
	allDeps := make(map[string]manifest.DepSpec, len(r.Deps)+len(r.DevDeps))
	for k, v := range r.Deps {
		allDeps[k] = v
	}
	for k, v := range r.DevDeps {
		allDeps[k] = v
	}

	if err := updateLockFile(projectDir, allDeps, nil); err != nil {
		return fmt.Errorf("regenerating grimoire.lock: %w", err)
	}

	lockPath := manifest.LockPath(manifestPath)
	fmt.Printf("grimoire.lock written to %s\n", lockPath)
	return nil
}
