package cmd

import (
	"os"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

var cliVersion = "dev"

func SetVersion(v string) { cliVersion = v }

var rootCmd = &cobra.Command{
	SilenceUsage: true,
	Use:          "grimoire",
	Short: "Grimoire — best practice enforcement for AI assistants",
	Long: `Grimoire skills enforce best practices in AI-assisted development.

  grimoire               Interactive TUI — install, update, doctor
  grimoire install       Install skills to AI agent directories
  grimoire uninstall     Remove skills from AI agent directories
  grimoire update        Pull the latest grimoire skills and relink
  grimoire list          List available domains, sub-domains, and skills
  grimoire doctor        Run a health check on the grimoire installation
  grimoire clean         Remove broken skill symlinks
  grimoire init          Initialize .grimoire/ in the current project
  grimoire check         Evaluate a compliance report
  grimoire config        Get or set grimoire configuration values
  grimoire registry      Manage skill registries (add, remove, list, update)
  grimoire self-update   Update the grimoire CLI binary to the latest release`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !isatty.IsTerminal(os.Stdin.Fd()) {
			return cmd.Help()
		}
		return runInteractive()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1) //nolint:revive // intentional: propagate cobra error as exit code 1
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(checkCmd)
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
	rootCmd.AddCommand(registryCmd)
}
