package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/agent"
	"github.com/jeffreytse/grimoire/internal/compliance"
	"github.com/jeffreytse/grimoire/internal/git"
	"github.com/jeffreytse/grimoire/internal/rules"
	"github.com/jeffreytse/grimoire/internal/settings"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var flagContextJSON bool

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Show grimoire context for AI assistants",
	Long: `Output the full grimoire environment context in one shot.

AI assistants can call this command at session start to get deterministic
ground truth: installation state, resolved settings, compliance status,
skill inventory, and structural rule findings.

Use --json for machine-readable output.`,
	RunE: runContext,
}

func init() {
	contextCmd.Flags().BoolVar(&flagContextJSON, "json", false, "output as JSON")
}

type contextAgentInfo struct {
	Name           string `json:"name"`
	DisplayName    string `json:"display_name"`
	Detected       bool   `json:"detected"`
	Version        string `json:"version,omitempty"`
	SkillsDir      string `json:"skills_dir,omitempty"`
	SkillsCount    int    `json:"skills_count"`
	BrokenSymlinks int    `json:"broken_symlinks"`
	Configured     bool   `json:"configured"`
}

type contextRegistryInfo struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	SkillsCount int    `json:"skills_count"`
	Cloned      bool   `json:"cloned"`
}

type contextOutput struct {
	CLIVersion      string                   `json:"cli_version"`
	GrimoireVersion string                   `json:"grimoire_version,omitempty"`
	GrimoireHome    string                   `json:"grimoire_home"`
	Agents          []contextAgentInfo       `json:"agents"`
	Settings        map[string]any           `json:"settings,omitempty"`
	Compliance      *compliance.Report       `json:"compliance,omitempty"`
	Registries      []contextRegistryInfo    `json:"registries"`
	RuleFindings    []compliance.Diagnostic  `json:"rule_findings,omitempty"`
}

func runContext(cmd *cobra.Command, args []string) error {
	home := skills.GrimoireHome()

	// grimoire version
	grimoireVer := ""
	if state, err := git.CurrentState(home); err == nil {
		grimoireVer = state.Version
	} else {
		grimoireVer = skills.GrimoireVersion()
	}

	// agents
	agentInfos := buildAgentInfos()

	// settings — reuse the same JSON shape as `grimoire settings --json`
	settingsMap := buildSettingsMap()

	// compliance — nil if no report
	var complianceReport *compliance.Report
	if r, err := compliance.Load(""); err == nil {
		complianceReport = r
	}

	// registries
	registries := buildRegistryInfos()

	// rule findings
	eng := &rules.Engine{
		SkillsSources: skills.AllSkillsSources(),
		ProjectDir:    ".",
	}
	ruleFindings := eng.Run()

	out := contextOutput{
		CLIVersion:      strings.TrimPrefix(cliVersion, "v"),
		GrimoireVersion: grimoireVer,
		GrimoireHome:    home,
		Agents:          agentInfos,
		Settings:        settingsMap,
		Compliance:      complianceReport,
		Registries:      registries,
		RuleFindings:    ruleFindings,
	}

	if flagContextJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}

	printContextHuman(out)
	return nil
}

func buildAgentInfos() []contextAgentInfo {
	var infos []contextAgentInfo
	for _, ag := range agent.All {
		info := contextAgentInfo{
			Name:        ag,
			DisplayName: agent.DisplayName(ag),
		}
		if _, err := exec.LookPath(ag); err == nil {
			info.Detected = true
			info.Version = agent.Version(ag)
			info.SkillsDir = agent.SkillsDir(ag)
			info.SkillsCount = agent.SkillCount(ag)
			info.BrokenSymlinks = agent.BrokenSymlinkCount(ag)
			info.Configured = agent.IsConfigured(ag)
		}
		infos = append(infos, info)
	}
	return infos
}

func buildSettingsMap() map[string]any {
	r, err := settings.Load(".")
	if err != nil {
		return nil
	}
	m := map[string]any{}
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
		m["core"] = core
	}
	for _, key := range r.SectionKeys() {
		ds := r.ResolveSection(key)
		entry := map[string]any{}
		if len(ds.Practices) > 0 {
			entry["practices"] = ds.Practices
		}
		if len(ds.Disabled) > 0 {
			entry["disabled"] = ds.Disabled
		}
		if ds.Fallback != "" {
			entry["fallback"] = ds.Fallback
		}
		if ds.ComplianceThreshold > 0 {
			entry["compliance-threshold"] = ds.ComplianceThreshold
		}
		if ds.ComplianceThresholdError >= 0 {
			entry["compliance-threshold-error"] = ds.ComplianceThresholdError
		}
		if len(entry) > 0 {
			m[key] = entry
		}
	}
	return m
}

func buildRegistryInfos() []contextRegistryInfo {
	officialURL := skills.GrimoireRepoURL()
	officialRoot := skills.SkillsRoot()
	officialCloned := false
	if _, err := os.Stat(officialRoot); err == nil {
		officialCloned = true
	}
	infos := []contextRegistryInfo{{
		Name:        "official",
		URL:         officialURL,
		SkillsCount: countSkills(officialRoot),
		Cloned:      officialCloned,
	}}

	fs, err := settings.LoadGlobal()
	if err != nil {
		return infos
	}
	for _, name := range sortedRegistryNames(fs.Registries) {
		rc := fs.Registries[name]
		regHome := skills.RegistryHome(name)
		regRoot := regHome + "/skills"
		cloned := false
		if _, err := os.Stat(regHome); err == nil {
			cloned = true
		}
		infos = append(infos, contextRegistryInfo{
			Name:        name,
			URL:         rc.URL,
			SkillsCount: countSkills(regRoot),
			Cloned:      cloned,
		})
	}
	return infos
}

func printContextHuman(out contextOutput) {
	fmt.Printf("\nGrimoire context\n")
	fmt.Printf("  cli:      %s\n", out.CLIVersion)
	if out.GrimoireVersion != "" {
		fmt.Printf("  grimoire: %s\n", out.GrimoireVersion)
	} else {
		fmt.Printf("  grimoire: %s\n", tui.StyleDim.Render("not installed"))
	}
	fmt.Printf("  home:     %s\n", out.GrimoireHome)

	// Agents
	fmt.Println()
	fmt.Println("  Agents")
	for _, ag := range out.Agents {
		if !ag.Detected {
			fmt.Printf("    %s  %-16s not found\n", tui.IconSkip, ag.DisplayName)
			continue
		}
		vs := ag.Version
		if vs == "" {
			vs = "detected"
		}
		cfgMark := tui.IconOK
		if !ag.Configured {
			cfgMark = tui.IconWarn
		}
		fmt.Printf("    %s  %-16s %-10s %d skills", cfgMark, ag.DisplayName, vs, ag.SkillsCount)
		if ag.BrokenSymlinks > 0 {
			fmt.Printf(", %d broken", ag.BrokenSymlinks)
		}
		fmt.Println()
	}

	// Registries
	fmt.Println()
	fmt.Println("  Registries")
	for _, reg := range out.Registries {
		icon := tui.IconOK
		if !reg.Cloned {
			icon = tui.IconWarn
		}
		fmt.Printf("    %s  %-12s %d skills  %s\n", icon, reg.Name, reg.SkillsCount, tui.StyleDim.Render(reg.URL))
	}

	// Settings summary
	if out.Settings != nil {
		if core, ok := out.Settings["core"].(map[string]any); ok {
			if profiles, ok := core["profiles"]; ok {
				fmt.Println()
				fmt.Printf("  Profiles: %v\n", profiles)
			}
		}
	}

	// Compliance
	fmt.Println()
	if out.Compliance != nil {
		r := out.Compliance
		statusIcon := tui.IconOK
		if r.Threshold.Status != "pass" {
			statusIcon = tui.IconFail
		}
		fmt.Printf("  Compliance: %.1f%% %s  %s\n", r.Coverage.OverallPct, statusIcon,
			tui.StyleDim.Render("("+r.Timestamp+")"))
	} else {
		fmt.Printf("  Compliance: %s\n", tui.StyleDim.Render("no report — run /check-best-practice-compliance"))
	}

	// Rule findings
	if len(out.RuleFindings) > 0 {
		fmt.Println()
		fmt.Println("  Structural issues:")
		for i := range out.RuleFindings {
			f := &out.RuleFindings[i]
			icon := tui.IconWarn
			if f.Severity == 1 {
				icon = tui.IconFail
			}
			fmt.Printf("    %s  [%s] %s\n", icon, f.Code, f.Message)
		}
	}

	fmt.Println()
}
