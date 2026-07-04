package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/config"
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
	Long: `Manage grimoire configuration (like git config).

Levels (highest priority first):
  default    grimoire.toml                         project — committed, shared with team
  -g/--global  ~/.config/grimoire/grimoire.toml   per-user
  --system   /etc/grimoire/grimoire.toml           machine-wide (requires admin)

  grimoire config get <key>               print current resolved value and source
  grimoire config set <key> <value>       write to target level (overwrites)
  grimoire config add <key> <value>       append a value to a list key (idempotent)
  grimoire config remove <key> <value>    remove a value from a list key
  grimoire config unset <key>             clear from target level

Core keys (always global — no flag needed, --local rejected):
  core.home      directory where grimoire is installed
  core.package   official package ref (owner/repo[@version] or URL)

Standards keys (default: project; use -g for user level):
  standards.extends                        additional packages (list)
  standards.profiles                       active profiles (list)
  standards.<domain>.practices             practice names (list)
  standards.<domain>.disabled              skill names to suppress (list)
  standards.<domain>.fallback              "ask" | "skip"
  standards.<domain>.compliance-threshold  number 0–100
  standards.<domain>.compliance-threshold-error  max allowed errors (-1 = unset)

Examples:
  grimoire config set standards.profiles "clean-architecture,tdd"
  grimoire config add standards.extends acmecorp/standards
  grimoire config add standards.extends acmecorp/standards -g
  grimoire config add standards.engineering.practices apply-solid-principles
  grimoire config set standards.engineering.compliance-threshold 75
  grimoire config set core.home ~/.grimoire
  grimoire config set core.package acmecorp/fork`,
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

var configAddCmd = &cobra.Command{
	Use:   "add <key> <value>",
	Short: "Append a value to a list config key (idempotent)",
	Args:  cobra.ExactArgs(2),
	RunE:  runConfigAdd,
}

var configRemoveCmd = &cobra.Command{
	Use:   "remove <key> <value>",
	Short: "Remove a value from a list config key",
	Args:  cobra.ExactArgs(2),
	RunE:  runConfigRemove,
}

func init() {
	for _, cmd := range []*cobra.Command{configGetCmd, configSetCmd, configUnsetCmd, configAddCmd, configRemoveCmd} {
		cmd.Flags().BoolVar(&flagConfigLocal, "local", false, "use project config (grimoire.toml) — same as default")
		cmd.Flags().BoolVarP(&flagConfigGlobal, "global", "g", false, "use user config (~/.config/grimoire/grimoire.toml)")
		cmd.Flags().BoolVar(&flagConfigSystem, "system", false, "use system config (/etc/grimoire/grimoire.toml)")
	}
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configUnsetCmd)
	configCmd.AddCommand(configAddCmd)
	configCmd.AddCommand(configRemoveCmd)
}

func runConfigGet(_ *cobra.Command, args []string) error {
	key := args[0]
	r, err := config.Load(getProjectDir())
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}
	val, src, err := getKeyResolved(&r, key)
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
	fs, err := config.LoadFile(path)
	if err != nil {
		return fmt.Errorf("loading %s: %w", path, err)
	}
	if err := applyKey(&fs, key, value); err != nil {
		return err
	}
	if err := config.WriteFile(path, fs); err != nil {
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
	fs, err := config.LoadFile(path)
	if err != nil {
		return fmt.Errorf("loading %s: %w", path, err)
	}
	if err := applyKey(&fs, key, ""); err != nil {
		return err
	}
	if err := config.WriteFile(path, fs); err != nil {
		return fmt.Errorf("saving %s: %w", path, err)
	}
	fmt.Printf("%s  %s cleared\n", tui.IconOK, key)
	return nil
}

// targetFilePath resolves which settings file to write based on flags and key prefix.
// core.* keys are always written to global (they have no effect at project level).
// All other keys default to project level; -g/--global overrides to user level.
func targetFilePath(key string) (string, error) {
	if strings.HasPrefix(key, "core.") {
		if flagConfigLocal {
			return "", fmt.Errorf("[core] keys are always global — omit --local or use -g")
		}
		return config.GlobalPath(), nil
	}

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
		return "", fmt.Errorf("only one of --local, -g/--global, --system may be specified")
	}

	if flagConfigGlobal {
		return config.GlobalPath(), nil
	}
	if flagConfigSystem {
		return config.SystemPath(), nil
	}
	// default: project level (like git config)
	return config.ProjectPath(getProjectDir()), nil
}

// levelToFilePath maps an explicit level string to a settings file path.
// Used by MCP tools where cobra flags are not available.
func levelToFilePath(key, level string) (string, error) {
	switch level {
	case "local":
		if strings.HasPrefix(key, "core.") {
			return "", fmt.Errorf("[core] keys are always global — omit level or use global")
		}
		return config.ProjectPath(getProjectDir()), nil
	case "global":
		return config.GlobalPath(), nil
	case "system":
		return config.SystemPath(), nil
	default: // "" — same auto-defaults as targetFilePath
		if strings.HasPrefix(key, "core.") {
			return config.GlobalPath(), nil
		}
		return config.ProjectPath(getProjectDir()), nil
	}
}

// getKeyResolved returns the resolved value and source for a key.
func getKeyResolved(r *config.Config, key string) (val, src string, err error) {
	switch key {
	case "core.home":
		return r.Core.Home, r.Sources["core.home"], nil
	case "core.agents":
		return strings.Join(r.Core.Agents, ", "), r.Sources["core.agents"], nil
	case "core.install-mode":
		return r.Core.InstallMode, r.Sources["core.install-mode"], nil
	}
	if strings.HasPrefix(key, "standards.") {
		domain, field, err := config.ParseStandardsKey(key)
		if err != nil {
			return "", "", err
		}
		switch field {
		case "profiles":
			return strings.Join(r.Core.Profiles, ", "), r.Sources["standards.profiles"], nil
		case "report-path":
			return r.ReportPath, r.Sources["standards.report-path"], nil
		case "staleness-days":
			if r.StalenessDays == 0 {
				return "", r.Sources["standards.staleness-days"], nil
			}
			return strconv.Itoa(r.StalenessDays), r.Sources["standards.staleness-days"], nil
		default:
			if domain == "" {
				return "", "", fmt.Errorf("field %q requires a domain (e.g. standards.engineering.%s)", field, field)
			}
			ds := r.ResolveSection(domain)
			return domainFieldString(&ds, field, r.Sources, domain+"."+field), r.Sources[domain+"."+field], nil
		}
	}
	return "", "", fmt.Errorf("unknown config key %q", key)
}

func domainFieldString(ds *config.DomainSection, field string, sources map[string]string, srcKey string) string {
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
func applyKey(fs *config.FileConfig, key, value string) error {
	switch key {
	case "core.home":
		fs.Core.Home = value
		return nil
	case "core.agents":
		fs.Core.Agents = splitCSV(value)
		return nil
	case "core.install-mode":
		if value != "" && value != "symlink" && value != "copy" {
			return fmt.Errorf("core.install-mode must be \"symlink\" or \"copy\"")
		}
		fs.Core.InstallMode = value
		return nil
	}
	if strings.HasPrefix(key, "standards.") {
		domain, field, err := config.ParseStandardsKey(key)
		if err != nil {
			return err
		}
		if domain == "" {
			// top-level standards field
			switch field {
			case "profiles":
				fs.Core.Profiles = splitCSV(value)
				return nil
			case "report-path":
				fs.ReportPath = value
				return nil
			case "staleness-days":
				if value == "" {
					fs.StalenessDays = 0
					return nil
				}
				n, err := strconv.Atoi(value)
				if err != nil || n < 0 {
					return fmt.Errorf("staleness-days must be a non-negative integer")
				}
				fs.StalenessDays = n
				return nil
			}
		}
		// domain field
		if fs.Sections == nil {
			fs.Sections = make(map[string]config.DomainSection)
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

func runConfigAdd(_ *cobra.Command, args []string) error {
	key, value := args[0], args[1]
	path, err := targetFilePath(key)
	if err != nil {
		return err
	}
	fs, err := config.LoadFile(path)
	if err != nil {
		return fmt.Errorf("loading %s: %w", path, err)
	}
	if err := appendToKey(&fs, key, value); err != nil {
		return err
	}
	if err := config.WriteFile(path, fs); err != nil {
		return fmt.Errorf("saving %s: %w", path, err)
	}
	fmt.Printf("%s  %s += %s\n", tui.IconOK, key, value)
	fmt.Printf("   saved to %s\n", path)
	return nil
}

func runConfigRemove(_ *cobra.Command, args []string) error {
	key, value := args[0], args[1]
	path, err := targetFilePath(key)
	if err != nil {
		return err
	}
	fs, err := config.LoadFile(path)
	if err != nil {
		return fmt.Errorf("loading %s: %w", path, err)
	}
	if err := removeFromKey(&fs, key, value); err != nil {
		return err
	}
	if err := config.WriteFile(path, fs); err != nil {
		return fmt.Errorf("saving %s: %w", path, err)
	}
	fmt.Printf("%s  %s -= %s\n", tui.IconOK, key, value)
	fmt.Printf("   saved to %s\n", path)
	return nil
}

// appendToKey appends value to a list-type key in fs (idempotent).
func appendToKey(fs *config.FileConfig, key, value string) error {
	switch key {
	case "core.agents":
		for _, v := range fs.Core.Agents {
			if v == value {
				return nil
			}
		}
		fs.Core.Agents = append(fs.Core.Agents, value)
		return nil
	case "standards.profiles":
		for _, v := range fs.Core.Profiles {
			if v == value {
				return nil
			}
		}
		fs.Core.Profiles = append(fs.Core.Profiles, value)
		return nil
	}
	if strings.HasPrefix(key, "standards.") {
		domain, field, err := config.ParseStandardsKey(key)
		if err != nil {
			return err
		}
		if domain != "" && (field == "practices" || field == "disabled") {
			if fs.Sections == nil {
				fs.Sections = make(map[string]config.DomainSection)
			}
			ds := fs.Sections[domain]
			if ds.ComplianceThresholdError == 0 {
				ds.ComplianceThresholdError = -1
			}
			ptr := &ds.Practices
			if field == "disabled" {
				ptr = &ds.Disabled
			}
			for _, v := range *ptr {
				if v == value {
					fs.Sections[domain] = ds
					return nil
				}
			}
			*ptr = append(*ptr, value)
			fs.Sections[domain] = ds
			return nil
		}
	}
	return fmt.Errorf("key %q is not a list — use 'config set' instead", key)
}

// removeFromKey removes value from a list-type key in fs (idempotent — no error if not present).
func removeFromKey(fs *config.FileConfig, key, value string) error {
	switch key {
	case "core.agents":
		fs.Core.Agents = filterOut(fs.Core.Agents, value)
		return nil
	case "standards.profiles":
		fs.Core.Profiles = filterOut(fs.Core.Profiles, value)
		return nil
	}
	if strings.HasPrefix(key, "standards.") {
		domain, field, err := config.ParseStandardsKey(key)
		if err != nil {
			return err
		}
		if domain != "" && (field == "practices" || field == "disabled") {
			if fs.Sections == nil {
				return nil
			}
			ds := fs.Sections[domain]
			if field == "practices" {
				ds.Practices = filterOut(ds.Practices, value)
			} else {
				ds.Disabled = filterOut(ds.Disabled, value)
			}
			fs.Sections[domain] = ds
			return nil
		}
	}
	return fmt.Errorf("key %q is not a list — use 'config unset' instead", key)
}

// filterOut returns a new slice with all occurrences of value removed.
func filterOut(slice []string, value string) []string {
	out := slice[:0:0]
	for _, v := range slice {
		if v != value {
			out = append(out, v)
		}
	}
	return out
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
