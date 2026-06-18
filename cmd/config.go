package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/settings"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var validKeys = []string{"home", "source"}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Get or set grimoire core configuration values",
	Long: `Manage grimoire core configuration in the global settings file.
These are machine-level keys stored under the [core] section.

  grimoire config get <key>         print current value
  grimoire config set <key> <value> set a value
  grimoire config unset <key>       reset to default

Supported keys:
  home     local directory where grimoire is installed (clone destination)
           (overrides the default ~/.grimoire)
  source   local path or git URL for the skills library
           (overrides the default https://github.com/jeffreytse/grimoire-skills)

To manage skill practice settings (profiles, practices, thresholds), use:
  grimoire settings`,
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Print a config value",
	Args:  cobra.ExactArgs(1),
	RunE:  runConfigGet,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a config value",
	Args:  cobra.ExactArgs(2),
	RunE:  runConfigSet,
}

var configUnsetCmd = &cobra.Command{
	Use:   "unset <key>",
	Short: "Clear a config value (reset to default)",
	Args:  cobra.ExactArgs(1),
	RunE:  runConfigUnset,
}

func init() {
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configUnsetCmd)
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	key := args[0]
	fs, err := settings.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	val, err := getCoreKey(fs, key)
	if err != nil {
		return err
	}
	if val == "" {
		fmt.Printf("(default: %s)\n", defaultFor(key))
	} else {
		fmt.Println(val)
	}
	return nil
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	key, value := args[0], args[1]
	fs, err := settings.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	if err := applyCoreKey(&fs, key, value); err != nil {
		return err
	}
	if err := settings.SaveGlobal(fs); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}
	fmt.Printf("%s  %s = %q\n", tui.IconOK, key, value)
	fmt.Printf("   saved to %s\n", settings.GlobalPath())
	return nil
}

func runConfigUnset(cmd *cobra.Command, args []string) error {
	key := args[0]
	fs, err := settings.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	if err := applyCoreKey(&fs, key, ""); err != nil {
		return err
	}
	if err := settings.SaveGlobal(fs); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}
	fmt.Printf("%s  %s cleared (using default)\n", tui.IconOK, key)
	return nil
}

func getCoreKey(fs settings.FileSettings, key string) (string, error) {
	switch key {
	case "home":
		return fs.Core.Home, nil
	case "source":
		return fs.Core.Source, nil
	default:
		return "", unknownKeyError(key)
	}
}

func applyCoreKey(fs *settings.FileSettings, key, value string) error {
	switch key {
	case "home":
		fs.Core.Home = value
		return nil
	case "source":
		fs.Core.Source = value
		return nil
	default:
		return unknownKeyError(key)
	}
}

func defaultFor(key string) string {
	switch key {
	case "home":
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".grimoire")
	case "source":
		return "https://github.com/jeffreytse/grimoire-skills"
	default:
		return "(none)"
	}
}

func unknownKeyError(key string) error {
	return fmt.Errorf("unknown config key %q — valid keys: %v", key, validKeys)
}
