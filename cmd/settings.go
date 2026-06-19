package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/profiles"
	"github.com/jeffreytse/grimoire/internal/settings"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var (
	flagSettingsDomain        string
	flagSettingsJSON          bool
	flagSettingsExpandProfiles bool
)

var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Show resolved grimoire settings for the current project",
	Long: `Show the effective settings after merging all layers (highest priority first):

  1. .grimoire/settings.local.toml   (project personal — gitignored)
  2. .grimoire/settings.toml         (project shared — committed)
  3. ~/.config/grimoire/settings.toml (global)

Each key shows the source file that provided it.
Use grimoire config get/set/unset to manage [core] keys (home, source).`,
	RunE: runSettings,
}

func init() {
	settingsCmd.Flags().StringVar(&flagSettingsDomain, "domain", "", "show only sections matching this domain prefix")
	settingsCmd.Flags().BoolVar(&flagSettingsJSON, "json", false, "output resolved settings as JSON")
	settingsCmd.Flags().BoolVar(&flagSettingsExpandProfiles, "expand-profiles", false, "show skills defined in each profile file")
}

func runSettings(cmd *cobra.Command, args []string) error {
	resolved, err := settings.Load(".")
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}

	if flagSettingsJSON {
		return printSettingsJSON(resolved)
	}
	printSettingsHuman(resolved)
	return nil
}

func printSettingsHuman(r settings.Resolved) {
	printed := false

	// [core] section — machine-only keys
	core := r.Core
	if core.Home != "" || core.Source != "" {
		fmt.Println()
		fmt.Printf("  %s\n", tui.StyleDim.Render("[core]"))
		if core.Home != "" {
			fmt.Printf("    home: %s%s\n", core.Home, sourceTag(r.Sources["core.home"]))
		}
		if core.Source != "" {
			fmt.Printf("    source: %s%s\n", core.Source, sourceTag(r.Sources["core.source"]))
		}
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
		fmt.Printf("  %s\n", tui.StyleDim.Render("[standards]"))
		if hasProfiles {
			fmt.Printf("    profiles: %s%s\n", strings.Join(core.Profiles, ", "), sourceTag(r.Sources["standards.profiles"]))
			if flagSettingsExpandProfiles {
				printExpandedProfiles(core.Profiles)
			}
		}
		printed = true
	}

	for _, key := range filteredKeys {
		ds := r.ResolveSection(key)
		lines := domainSectionLines(key, ds, r.Sources)
		if len(lines) == 0 {
			continue
		}
		fmt.Println()
		fmt.Printf("  %s\n", tui.StyleDim.Render("[standards."+key+"]"))
		for _, l := range lines {
			fmt.Println(l)
		}
	}

	if !printed {
		fmt.Printf("  %s  no settings configured\n", tui.IconWarn)
		fmt.Printf("  run: grimoire init\n")
		fmt.Printf("  or edit: %s\n", settings.GlobalPath())
	} else {
		fmt.Println()
	}
}

func printExpandedProfiles(profileNames []string) {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	opts := profiles.ResolveOptions{Sources: skills.AllSkillsSources()}
	for _, name := range profileNames {
		p, err := profiles.ResolveWithOptions(name, cwd, opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "    warn: loading profile %q: %v\n", name, err)
			continue
		}
		if p.Source == "" {
			fmt.Printf("    %s\n", tui.StyleDim.Render(p.Name+":"))
			fmt.Printf("      %s\n", tui.StyleDim.Render("(no profile file found — resolved by LLM tag query at runtime)"))
			continue
		}
		src := sourceTag(p.Source)
		fmt.Printf("    %s%s\n", tui.StyleDim.Render(p.Name+":"), src)
		for _, sk := range p.Skills {
			fmt.Printf("      %s %s\n", tui.StyleCyan.Render("→"), sk.Name)
		}
		if len(p.Skills) == 0 {
			fmt.Printf("      %s\n", tui.StyleDim.Render("(no skills defined in profile file)"))
		}
	}
}

func domainSectionLines(key string, ds settings.DomainSection, sources map[string]string) []string {
	var lines []string
	if len(ds.Practices) > 0 {
		lines = append(lines, fmt.Sprintf("    practices: %s%s",
			strings.Join(ds.Practices, ", "), sourceTag(sources[key+".practices"])))
	}
	if ds.Fallback != "" {
		lines = append(lines, fmt.Sprintf("    fallback: %s%s",
			ds.Fallback, sourceTag(sources[key+".fallback"])))
	}
	if ds.ComplianceThreshold > 0 {
		lines = append(lines, fmt.Sprintf("    compliance-threshold: %.0f%%%s",
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

func printSettingsJSON(r settings.Resolved) error {
	type domainOut struct {
		Practices                []string `json:"practices,omitempty"`
		Fallback                 string   `json:"fallback,omitempty"`
		ComplianceThreshold      float64  `json:"compliance-threshold,omitempty"`
		ComplianceThresholdError *int     `json:"compliance-threshold-error,omitempty"`
	}

	out := map[string]any{}

	core := map[string]any{}
	if r.Core.Home != "" {
		core["home"] = r.Core.Home
	}
	if r.Core.Source != "" {
		core["source"] = r.Core.Source
	}
	if len(r.Core.Profiles) > 0 {
		core["profiles"] = r.Core.Profiles
	}
	if len(core) > 0 {
		out["core"] = core
	}

	for _, key := range r.SectionKeys() {
		if flagSettingsDomain != "" && !strings.HasPrefix(key, flagSettingsDomain) {
			continue
		}
		ds := r.ResolveSection(key)
		d := domainOut{
			Practices:           ds.Practices,
			Fallback:            ds.Fallback,
			ComplianceThreshold: ds.ComplianceThreshold,
		}
		if ds.ComplianceThresholdError >= 0 {
			v := ds.ComplianceThresholdError
			d.ComplianceThresholdError = &v
		}
		out[key] = d
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
