package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/skills"
)

var cliVersion = "dev"

func SetVersion(v string) { cliVersion = v }

var flagProjectDir string

// errNoSkillLibrary is the cold-start error for commands that need a
// downloaded skill library. `grimoire doctor` shows where grimoire looked.
var errNoSkillLibrary = errors.New("no skill library downloaded yet — run: grimoire update")

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
	Short:        "Grimoire — makes AI agents follow best practices",
	Long: `Grimoire makes AI agents follow best practices — a package manager for best-practice skills.

Start here:  grimoire mcp install   (connect grimoire to your AI agent)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(skills.OfficialPackageHome()); err != nil {
			fmt.Println("New to grimoire? Run: grimoire mcp install  (or grimoire wizard to choose step by step)")
			fmt.Println()
		}
		return cmd.Help()
	},
}

func Execute() {
	rootCmd.Version = cliVersion
	groupCommands(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1) //nolint:revive // intentional: propagate cobra error as exit code 1
	}
}

func init() {
	// Config layers (and the packages they list) follow --project-dir too.
	skills.ProjectDir = getProjectDir
	rootCmd.PersistentFlags().StringVar(&flagProjectDir, "project-dir", "", "project directory (default: current working directory)")
	rootCmd.AddCommand(wizardCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(watchCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(selfUpdateCmd)
	rootCmd.AddCommand(configCmd)

	rootCmd.AddCommand(packageCmd)
	rootCmd.AddCommand(profileCmd)
	rootCmd.AddCommand(contextCmd)
	rootCmd.AddCommand(mcpCmd)
	rootCmd.AddCommand(lspCmd)
}
