package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/config"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var validKeys = []string{"home", "source"}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Get or set grimoire configuration values",
	Long: `Manage grimoire global configuration (~/.config/grimoire/settings.toml).

  grimoire config get <key>         print current value
  grimoire config set <key> <value> set a value
  grimoire config unset <key>       reset to default

Supported keys:
  home     local directory where grimoire is installed (clone destination)
           (overrides the default ~/.grimoire)
  source   local path or git URL for the skills library
           (overrides the default https://github.com/jeffreytse/grimoire-skills)`,
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
	g, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	val, err := getKey(g, key)
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
	g, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	if err := applyKey(&g, key, value); err != nil {
		return err
	}
	if err := config.Save(g); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}
	fmt.Printf("%s  %s = %q\n", tui.IconOK, key, value)
	fmt.Printf("   saved to %s\n", config.GlobalPath())
	return nil
}

func runConfigUnset(cmd *cobra.Command, args []string) error {
	key := args[0]
	g, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	if err := applyKey(&g, key, ""); err != nil {
		return err
	}
	if err := config.Save(g); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}
	fmt.Printf("%s  %s cleared (using default)\n", tui.IconOK, key)
	return nil
}

func getKey(g config.Global, key string) (string, error) {
	switch key {
	case "home":
		return g.Home, nil
	case "source":
		return g.Source, nil
	default:
		return "", unknownKeyError(key)
	}
}

func applyKey(g *config.Global, key, value string) error {
	switch key {
	case "home":
		g.Home = value
		return nil
	case "source":
		g.Source = value
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
