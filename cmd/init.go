package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/detect"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize .grimoire/ in the current project",
	RunE:  runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	dir := filepath.Join(cwd, ".grimoire")
	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf(".grimoire/ already exists — edit .grimoire/settings.toml directly or run /configure-grimoire")
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating .grimoire/: %w", err)
	}

	profile := detect.Profile(cwd)
	if err := writeSettings(dir, profile); err != nil {
		return err
	}
	if err := writeGitignore(dir); err != nil {
		return err
	}

	fmt.Println("✓ Grimoire initialized. Run /suggest-best-practice to get started.")
	return nil
}

func writeSettings(dir, profile string) error {
	var profileLine string
	if profile != "" {
		profileLine = fmt.Sprintf("profiles = [%q]", profile)
	} else {
		profileLine = `# profiles = ["engineering"]   # uncomment and set your profile`
	}

	content := fmt.Sprintf(`# Grimoire settings
# Docs: https://github.com/jeffreytse/grimoire/blob/main/docs/settings.md

[core]
%s

# Domain-level settings example:
# [engineering]
# practices = ["apply-solid-principles", "apply-kiss-principle"]
# compliance-threshold = 80
# compliance-threshold-error = 0
`, profileLine)

	path := filepath.Join(dir, "settings.toml")
	return os.WriteFile(path, []byte(content), 0o644)
}

func writeGitignore(dir string) error {
	content := "settings.local.toml\n"
	path := filepath.Join(dir, ".gitignore")
	return os.WriteFile(path, []byte(content), 0o644)
}
