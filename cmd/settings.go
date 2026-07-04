package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/config"
	"github.com/jeffreytse/grimoire/internal/profiles"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var (
	flagSettingsDomain string
	flagSettingsJSON   bool
)

var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Show resolved grimoire settings for the current project",
	Long: `Show the effective config after merging all layers (highest priority first):

  1. grimoire.toml                       (project — committed, --local)
  2. ~/.config/grimoire/grimoire.toml    (user global, --global)
  3. /etc/grimoire/grimoire.toml         (system-wide, --system)

Each key shows the source file that provided it.
Use grimoire config get/set/unset to manage all keys.`,
	RunE: runSettings,
}

func init() {
	settingsCmd.Flags().StringVar(&flagSettingsDomain, "domain", "", "show only sections matching this domain prefix")
	settingsCmd.Flags().BoolVar(&flagSettingsJSON, "json", false, "output resolved settings as JSON")
}

func runSettings(cmd *cobra.Command, args []string) error {
	resolved, err := config.Load(getProjectDir())
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}

	if flagSettingsJSON {
		return printSettingsJSON(resolved)
	}
	printSettingsHuman(resolved)
	return nil
}

func printSettingsHuman(r config.Config) { //nolint:gocritic // value semantics intentional; callers pass local Resolved vars
	printed := false

	// [core] section — machine-only keys
	core := r.Core
	if core.Home != "" {
		fmt.Println()
		fmt.Printf("    %s\n", tui.StyleDim.Render("[core]"))
		fmt.Printf("    home: %s%s\n", core.Home, sourceTag(r.Sources["core.home"]))
		printed = true
	}

	// [standards] section — profiles + domain sections
	hasProfiles := len(core.Profiles) > 0
	keys := r.SectionKeys()
	sort.Strings(keys)
	filteredKeys := keys[:0]
	for _, k := range keys {
		if flagSettingsDomain == "" || strings.HasPrefix(k, flagSettingsDomain) {
			filteredKeys = append(filteredKeys, k)
		}
	}

	if hasProfiles || len(filteredKeys) > 0 {
		fmt.Println()
		fmt.Printf("    %s\n", tui.StyleDim.Render("[standards]"))
		if hasProfiles {
			fmt.Printf("    profiles: %s%s\n", strings.Join(core.Profiles, ", "), sourceTag(r.Sources["standards.profiles"]))
			printExpandedProfiles(core.Profiles)
		}
		printed = true
	}

	for _, key := range filteredKeys {
		ds := r.ResolveSection(key)
		lines := domainSectionLines(key, &ds, r.Sources)
		if len(lines) == 0 {
			continue
		}
		fmt.Println()
		fmt.Printf("    %s\n", tui.StyleDim.Render("[standards."+key+"]"))
		for _, l := range lines {
			fmt.Println(l)
		}
	}

	if !printed {
		fmt.Printf("  %s  no settings configured\n", tui.IconWarn)
		fmt.Printf("  run: grimoire init\n")
		fmt.Printf("  or edit: %s\n", config.GlobalPath())
	} else {
		fmt.Println()
	}
}

func printExpandedProfiles(profileNames []string) {
	cwd := getProjectDir()
	opts := resolveOpts(cwd)
	for _, name := range profileNames {
		p, err := profiles.ResolveWithOptions(name, cwd, opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "    warn: loading profile %q: %v\n", name, err)
			continue
		}
		if p.Source == "" {
			fmt.Printf("      %s\n", tui.StyleDim.Render(p.Name+":"))
			fmt.Printf("        %s\n", tui.StyleDim.Render("(no installed skills match — AI applies semantically)"))
			continue
		}
		src := sourceTag(p.Source)
		fmt.Printf("      %s%s\n", tui.StyleDim.Render(p.Name+":"), src)
		for _, sk := range p.Skills {
			fmt.Printf("        %s %s\n", tui.StyleCyan.Render("→"), sk.Name)
		}
		if len(p.Skills) == 0 {
			fmt.Printf("        %s\n", tui.StyleDim.Render("(no installed skills match — AI applies semantically)"))
		}
	}
}

func domainSectionLines(key string, ds *config.DomainSection, sources map[string]string) []string {
	var lines []string
	if len(ds.Practices) > 0 {
		lines = append(lines, fmt.Sprintf("    practices: %s%s",
			strings.Join(ds.Practices, ", "), sourceTag(sources[key+".practices"])))
	}
	if len(ds.Disabled) > 0 {
		lines = append(lines, fmt.Sprintf("    disabled: %s%s",
			strings.Join(ds.Disabled, ", "), sourceTag(sources[key+".disabled"])))
	}
	if ds.Fallback != "" {
		lines = append(lines, fmt.Sprintf("    fallback: %s%s",
			ds.Fallback, sourceTag(sources[key+".fallback"])))
	}
	if ds.ComplianceThreshold > 0 {
		lines = append(lines, fmt.Sprintf("    compliance-threshold: %.0f%s",
			ds.ComplianceThreshold, sourceTag(sources[key+".compliance-threshold"])))
	}
	if ds.ComplianceThresholdError >= 0 {
		lines = append(lines, fmt.Sprintf("    compliance-threshold-error: %d%s",
			ds.ComplianceThresholdError, sourceTag(sources[key+".compliance-threshold-error"])))
	}
	return lines
}

func sourceTag(path string) string {
	if path == "" {
		return ""
	}
	// shorten home dir for readability
	if home, err := os.UserHomeDir(); err == nil {
		path = strings.Replace(path, home, "~", 1)
	}
	return tui.StyleDim.Render("   (" + path + ")")
}

func settingsToMap(r config.Config) map[string]any { //nolint:gocritic // value semantics intentional
	out := map[string]any{}
	core := map[string]any{}
	if r.Core.Home != "" {
		core["home"] = r.Core.Home
	}
	if len(r.Core.Profiles) > 0 {
		core["profiles"] = r.Core.Profiles
	}
	if len(core) > 0 {
		out["core"] = core
	}
	for _, key := range r.SectionKeys() {
		ds := r.ResolveSection(key)
		d := map[string]any{}
		if len(ds.Practices) > 0 {
			d["practices"] = ds.Practices
		}
		if ds.Fallback != "" {
			d["fallback"] = ds.Fallback
		}
		if ds.ComplianceThreshold > 0 {
			d["compliance-threshold"] = ds.ComplianceThreshold
		}
		if ds.ComplianceThresholdError >= 0 {
			d["compliance-threshold-error"] = ds.ComplianceThresholdError
		}
		out[key] = d
	}
	return out
}

func printSettingsJSON(r config.Config) error { //nolint:gocritic // value semantics intentional
	m := settingsToMap(r)
	if flagSettingsDomain != "" {
		for k := range m {
			if k != "core" && !strings.HasPrefix(k, flagSettingsDomain) {
				delete(m, k)
			}
		}
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(m)
}
