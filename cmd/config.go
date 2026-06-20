package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/settings"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var (
	flagConfigLocal  bool
	flagConfigGlobal bool
	flagConfigSystem bool
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Get or set grimoire configuration values",
	Long: `Manage grimoire configuration across three levels (highest priority first):

  --local    .grimoire/settings.toml       project, committed — shared with team
  --global   ~/.config/grimoire/settings.toml  per-user (default for core.* keys)
  --system   /etc/grimoire/settings.toml   machine-wide (requires admin)

  grimoire config get <key>          print current resolved value and source
  grimoire config set <key> <value>  write to target level
  grimoire config unset <key>        clear from target level

Core keys (default level: --global):
  core.home     local directory where grimoire is installed
  core.source   local path or git URL for the skills library

Standards keys (default level: --local):
  standards.profiles                       active profiles (comma-separated)
  standards.<domain>.practices             practice names (comma-separated)
  standards.<domain>.disabled              skill names to suppress (comma-separated)
  standards.<domain>.fallback              "ask" | "skip"
  standards.<domain>.compliance-threshold  number 0–100
  standards.<domain>.compliance-threshold-error  max allowed errors (-1 = unset)

Examples:
  grimoire config set standards.profiles "clean-architecture,tdd"
  grimoire config set standards.engineering.compliance-threshold 75
  grimoire config set standards.engineering.testing.practices "apply-tdd"
  grimoire config set core.home ~/.grimoire --global`,
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Print a config value and its source",
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
	Short: "Clear a config value",
	Args:  cobra.ExactArgs(1),
	RunE:  runConfigUnset,
}

func init() {
	for _, cmd := range []*cobra.Command{configGetCmd, configSetCmd, configUnsetCmd} {
		cmd.Flags().BoolVar(&flagConfigLocal, "local", false, "use project settings (.grimoire/settings.toml)")
		cmd.Flags().BoolVar(&flagConfigGlobal, "global", false, "use user settings (~/.config/grimoire/settings.toml)")
		cmd.Flags().BoolVar(&flagConfigSystem, "system", false, "use system settings (/etc/grimoire/settings.toml)")
	}
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configUnsetCmd)
}

func runConfigGet(_ *cobra.Command, args []string) error {
	key := args[0]
	r, err := settings.Load(".")
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}
	val, src, err := getKeyResolved(r, key)
	if err != nil {
		return err
	}
	if val == "" {
		fmt.Printf("(unset)\n")
	} else {
		fmt.Printf("%s%s\n", val, sourceTag(src))
	}
	return nil
}

func runConfigSet(_ *cobra.Command, args []string) error {
	key, value := args[0], args[1]
	path, err := targetFilePath(key)
	if err != nil {
		return err
	}
	fs, err := settings.LoadFile(path)
	if err != nil {
		return fmt.Errorf("loading %s: %w", path, err)
	}
	if err := applyKey(&fs, key, value); err != nil {
		return err
	}
	if err := settings.WriteFile(path, fs); err != nil {
		return fmt.Errorf("saving %s: %w", path, err)
	}
	fmt.Printf("%s  %s = %s\n", tui.IconOK, key, value)
	fmt.Printf("   saved to %s\n", path)
	return nil
}

func runConfigUnset(_ *cobra.Command, args []string) error {
	key := args[0]
	path, err := targetFilePath(key)
	if err != nil {
		return err
	}
	fs, err := settings.LoadFile(path)
	if err != nil {
		return fmt.Errorf("loading %s: %w", path, err)
	}
	if err := applyKey(&fs, key, ""); err != nil {
		return err
	}
	if err := settings.WriteFile(path, fs); err != nil {
		return fmt.Errorf("saving %s: %w", path, err)
	}
	fmt.Printf("%s  %s cleared\n", tui.IconOK, key)
	return nil
}

// targetFilePath resolves which settings file to write based on flags and key prefix.
func targetFilePath(key string) (string, error) {
	count := 0
	if flagConfigLocal {
		count++
	}
	if flagConfigGlobal {
		count++
	}
	if flagConfigSystem {
		count++
	}
	if count > 1 {
		return "", fmt.Errorf("only one of --local, --global, --system may be specified")
	}

	if flagConfigLocal {
		cwd, _ := os.Getwd()
		return filepath.Join(cwd, ".grimoire", "settings.toml"), nil
	}
	if flagConfigGlobal {
		return settings.GlobalPath(), nil
	}
	if flagConfigSystem {
		return settings.SystemPath(), nil
	}
	// defaults
	if strings.HasPrefix(key, "standards.") {
		cwd, _ := os.Getwd()
		return filepath.Join(cwd, ".grimoire", "settings.toml"), nil
	}
	return settings.GlobalPath(), nil
}

// getKeyResolved returns the resolved value and source for a key.
func getKeyResolved(r settings.Resolved, key string) (val, src string, err error) {
	switch key {
	case "core.home":
		return r.Core.Home, r.Sources["core.home"], nil
	case "core.source":
		return r.Core.Source, r.Sources["core.source"], nil
	}
	if strings.HasPrefix(key, "standards.") {
		domain, field, err := settings.ParseStandardsKey(key)
		if err != nil {
			return "", "", err
		}
		switch field {
		case "profiles":
			return strings.Join(r.Core.Profiles, ", "), r.Sources["standards.profiles"], nil
		default:
			if domain == "" {
				return "", "", fmt.Errorf("field %q requires a domain (e.g. standards.engineering.%s)", field, field)
			}
			ds := r.ResolveSection(domain)
			return domainFieldString(ds, field, r.Sources, domain+"."+field), r.Sources[domain+"."+field], nil
		}
	}
	return "", "", fmt.Errorf("unknown config key %q", key)
}

func domainFieldString(ds settings.DomainSection, field string, sources map[string]string, srcKey string) string {
	switch field {
	case "practices":
		return strings.Join(ds.Practices, ", ")
	case "disabled":
		return strings.Join(ds.Disabled, ", ")
	case "fallback":
		return ds.Fallback
	case "compliance-threshold":
		if ds.ComplianceThreshold == 0 {
			return ""
		}
		return fmt.Sprintf("%.0f", ds.ComplianceThreshold)
	case "compliance-threshold-error":
		if ds.ComplianceThresholdError == -1 {
			return ""
		}
		return strconv.Itoa(ds.ComplianceThresholdError)
	}
	return ""
}

// applyKey mutates fs with the given key=value. value="" means unset.
func applyKey(fs *settings.FileSettings, key, value string) error {
	switch key {
	case "core.home":
		fs.Core.Home = value
		return nil
	case "core.source":
		fs.Core.Source = value
		return nil
	}
	if strings.HasPrefix(key, "standards.") {
		domain, field, err := settings.ParseStandardsKey(key)
		if err != nil {
			return err
		}
		if domain == "" {
			// top-level standards field
			switch field {
			case "profiles":
				fs.Core.Profiles = splitCSV(value)
				return nil
			}
		}
		// domain field
		if fs.Sections == nil {
			fs.Sections = make(map[string]settings.DomainSection)
		}
		ds := fs.Sections[domain]
		if ds.ComplianceThresholdError == 0 {
			ds.ComplianceThresholdError = -1 // preserve sentinel on new entries
		}
		switch field {
		case "practices":
			ds.Practices = splitCSV(value)
		case "disabled":
			ds.Disabled = splitCSV(value)
		case "fallback":
			ds.Fallback = value
		case "compliance-threshold":
			if value == "" {
				ds.ComplianceThreshold = 0
			} else {
				f, err := strconv.ParseFloat(strings.TrimSuffix(value, "%"), 64)
				if err != nil {
					return fmt.Errorf("compliance-threshold must be a number: %w", err)
				}
				ds.ComplianceThreshold = f
			}
		case "compliance-threshold-error":
			if value == "" {
				ds.ComplianceThresholdError = -1
			} else {
				n, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("compliance-threshold-error must be an integer: %w", err)
				}
				ds.ComplianceThresholdError = n
			}
		}
		fs.Sections[domain] = ds
		return nil
	}
	return fmt.Errorf("unknown config key %q", key)
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
