package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/skills"
)

var cliVersion = "dev"

func SetVersion(v string) { cliVersion = v }

var flagProjectDir string

// getProjectDir returns the project directory to use for all project-scoped
// operations. Priority: --project-dir flag > GRIMOIRE_PROJECT_DIR env var > cwd.
func getProjectDir() string {
	p := flagProjectDir
	if p == "" {
		p = os.Getenv("GRIMOIRE_PROJECT_DIR")
	}
	if p == "" {
		cwd, _ := os.Getwd()
		return cwd
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return p
	}
	return abs
}

var rootCmd = &cobra.Command{
	SilenceUsage: true,
	Use:          "grimoire",
	Short:        "Grimoire — best practice enforcement for AI assistants",
	Long: `Grimoire — skills package manager for AI agents.

  grimoire init          Initialize grimoire.toml in the current project
  grimoire install       Install skills — from grimoire.toml, or add+install a ref in one step
  grimoire uninstall     Remove skills — all, or remove+unlink a ref in one step
  grimoire update        Update all packages to latest
  grimoire list          List available domains, sub-domains, and skills
  grimoire check         Run BPDD compliance check against declared practices
  grimoire watch         Re-run compliance check whenever files change
  grimoire status        Show project compliance health at a glance
  grimoire wizard        Open the interactive TUI wizard
  grimoire doctor        Run a health check on the grimoire installation
  grimoire clean         Remove broken skill symlinks
  grimoire config        Get or set grimoire configuration values
  grimoire package       Manage skill packages (add, remove, list, update)
  grimoire profile       Manage profiles (list, show, init)
  grimoire self-update   Update the grimoire CLI binary to the latest release`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(skills.OfficialPackageHome()); err != nil {
			fmt.Println("New to grimoire? Run: grimoire wizard")
			fmt.Println()
		}
		return cmd.Help()
	},
}

func Execute() {
	rootCmd.Version = cliVersion
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1) //nolint:revive // intentional: propagate cobra error as exit code 1
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagProjectDir, "project-dir", "", "project directory (default: current working directory)")
	rootCmd.AddCommand(wizardCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(watchCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(selfUpdateCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(settingsCmd)
	rootCmd.AddCommand(packageCmd)
	rootCmd.AddCommand(profileCmd)
	rootCmd.AddCommand(contextCmd)
	rootCmd.AddCommand(mcpCmd)
	rootCmd.AddCommand(lspCmd)
}
