package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/compliance"
	"github.com/jeffreytse/grimoire/internal/detect"
	"github.com/jeffreytse/grimoire/internal/profiles"
	"github.com/jeffreytse/grimoire/internal/settings"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var (
	flagInitYes       bool
	flagInitProfile   string
	flagInitThreshold int
	flagInitMaxErrors int
	flagInitPreset    string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize .grimoire/ in the current project",
	RunE:  runInit,
}

func init() {
	initCmd.Flags().BoolVarP(&flagInitYes, "yes", "y", false, "accept all defaults without prompting")
	initCmd.Flags().StringVar(&flagInitProfile, "profile", "", "profile to activate (skips prompt)")
	initCmd.Flags().IntVar(&flagInitThreshold, "threshold", 0, "compliance threshold % (skips prompt)")
	initCmd.Flags().IntVar(&flagInitMaxErrors, "max-errors", -1, "max allowed errors (skips prompt)")
	initCmd.Flags().StringVar(&flagInitPreset, "preset", "", "apply a named preset from an installed registry (skips wizard)")
}

// initConfig holds the answers collected by the wizard (or defaults).
type initConfig struct {
	Profile   string
	Threshold int // 0 = omit from settings
	MaxErrors int // -1 = omit, ≥0 = write
}

func runInit(cmd *cobra.Command, args []string) error {
	cwd := getProjectDir()
	dir := filepath.Join(cwd, ".grimoire")
	reinit := dirExists(dir)
	detected := detect.Profile(cwd)

	// --preset flag: apply named preset directly, skip wizard
	if flagInitPreset != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("creating .grimoire/: %w", err)
		}
		return applyNamedPreset(flagInitPreset, dir, cwd, reinit)
	}

	hasExplicit := flagInitProfile != "" || flagInitThreshold != 0 || flagInitMaxErrors != -1

	// TTY interactive path: offer preset picker before wizard.
	// Skip preset picker when the user already has profiles configured in global/extends
	// settings — they've already committed to a workflow, so we pre-populate the wizard
	// with those choices instead of interrupting with a preset menu.
	if !hasExplicit && tui.IsTTY() && !flagInitYes {
		resolved, _ := settings.Load(cwd)
		if !reinit && len(resolved.Core.Profiles) > 0 {
			fmt.Printf("  Using profiles from your settings: %s\n", strings.Join(resolved.Core.Profiles, ", "))
			// fall through to wizard below — loadExistingInitConfig will pick these up
		} else {
			cand, mode := promptPreset(detected)
			switch mode {
			case "skip":
				fmt.Println("  Skipped — run `grimoire init` to set up later.")
				return nil
			case "preset":
				if err := os.MkdirAll(dir, 0o755); err != nil {
					return fmt.Errorf("creating .grimoire/: %w", err)
				}
				if err := skills.ApplyPreset(cand.regHome, cand.presetName, dir); err != nil {
					return fmt.Errorf("applying preset: %w", err)
				}
				printInitSuccess(reinit, cand.presetName)
				if r, err := settings.Load(cwd); err == nil && len(r.Core.Profiles) > 0 {
					printProfilePreview(r.Core.Profiles[0], cwd)
				}
				printInitReport(cwd)
				return nil
				// case "custom": fall through to wizard below
			}
		}
	}

	// Wizard / explicit-flags path (existing behaviour)
	existing := loadExistingInitConfig(cwd, detected)
	var cfg initConfig
	if !hasExplicit && tui.IsTTY() && !flagInitYes {
		cfg = runInitWizard(existing, detected, cwd)
	} else {
		cfg = existing
	}
	if flagInitProfile != "" {
		cfg.Profile = flagInitProfile
	}
	if flagInitThreshold != 0 {
		cfg.Threshold = flagInitThreshold
	}
	if flagInitMaxErrors != -1 {
		cfg.MaxErrors = flagInitMaxErrors
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating .grimoire/: %w", err)
	}
	if err := writeSettings(dir, cfg); err != nil {
		return err
	}
	printInitSuccess(reinit, "")
	printProfilePreview(cfg.Profile, cwd)
	printInitReport(cwd)
	return nil
}

func printInitSuccess(reinit bool, preset string) {
	suffix := ""
	if preset != "" {
		suffix = fmt.Sprintf(" with preset %q", preset)
	}
	if reinit {
		fmt.Printf("✓ Grimoire reinitialized%s.\n", suffix)
	} else {
		fmt.Printf("✓ Grimoire initialized%s.\n", suffix)
	}
}

// applyNamedPreset locates preset by name across all installed registries and applies it.
func applyNamedPreset(name, dir, cwd string, reinit bool) error {
	for _, reg := range skills.AllRegistries() {
		for _, p := range skills.ListPresets(reg.Home) {
			if p != name {
				continue
			}
			if err := skills.ApplyPreset(reg.Home, name, dir); err != nil {
				return fmt.Errorf("applying preset: %w", err)
			}
			printInitSuccess(reinit, name)
			if r, err := settings.Load(cwd); err == nil && len(r.Core.Profiles) > 0 {
				printProfilePreview(r.Core.Profiles[0], cwd)
			}
			printInitReport(cwd)
			return nil
		}
	}
	return fmt.Errorf("preset %q not found — run `grimoire registry update` to refresh installed registries", name)
}

// presetCandidate carries the registry context for a preset shown in the picker.
type presetCandidate struct {
	presetName string
	regName    string
	regHome    string
}

// promptPreset shows a TUI preset picker and returns the chosen candidate plus
// a mode string: "preset" (apply it), "custom" (fall through to wizard), or "skip".
func promptPreset(detected string) (cand *presetCandidate, mode string) {
	items, annotations, candMap := listRankedPresetItems(detected)
	if len(items) == 0 {
		return nil, "custom"
	}

	fmt.Printf("  Initialize with:")
	chosen, ok := tui.RunProfileSelect(items, annotations, detected)
	if !ok {
		return nil, "custom"
	}
	switch chosen {
	case tui.PresetSelectCustom:
		return nil, "custom"
	case tui.PresetSelectSkip:
		return nil, "skip"
	default:
		c := candMap[chosen]
		return &c, "preset"
	}
}

// listRankedPresetItems builds the TUI item list for the preset picker.
// Returns items (for RunProfileSelect), annotations, and a lookup map from item key → presetCandidate.
func listRankedPresetItems(detected string) ([]string, map[string]string, map[string]presetCandidate) { //nolint:gocritic // three-tuple return; naming doesn't improve clarity
	annotations := make(map[string]string)
	candMap := make(map[string]presetCandidate)
	seen := make(map[string]struct{})

	var suggested, others []string
	for _, reg := range skills.AllRegistries() {
		for _, pname := range skills.ListPresets(reg.Home) {
			if _, ok := seen[pname]; ok {
				continue
			}
			seen[pname] = struct{}{}
			annotations[pname] = "[" + reg.Name + "]"
			candMap[pname] = presetCandidate{pname, reg.Name, reg.Home}
			if detected != "" && pname == detected {
				suggested = append(suggested, pname)
			} else {
				others = append(others, pname)
			}
		}
	}

	if len(suggested) == 0 && len(others) == 0 {
		return nil, annotations, candMap
	}

	var items []string
	if len(suggested) > 0 {
		items = append(items, tui.ProfileSectionPrefix+"Suggested")
		items = append(items, suggested...)
	}
	if len(others) > 0 {
		if len(suggested) > 0 {
			items = append(items, tui.ProfileSectionPrefix+"All presets")
		}
		items = append(items, others...)
	}
	items = append(items, tui.PresetSelectCustom, tui.PresetSelectSkip)
	return items, annotations, candMap
}

// loadExistingInitConfig reads current settings.toml (if any) and returns an
// initConfig pre-populated with those values — used as wizard defaults on re-init.
func loadExistingInitConfig(cwd, detected string) initConfig {
	cfg := initConfig{Profile: detected, Threshold: 0, MaxErrors: -1}

	r, err := settings.Load(cwd)
	if err != nil {
		return cfg
	}
	if len(r.Core.Profiles) > 0 {
		cfg.Profile = r.Core.Profiles[0]
	}
	if cfg.Profile != "" {
		sec := r.ResolveSection(cfg.Profile)
		if sec.ComplianceThreshold > 0 {
			cfg.Threshold = int(sec.ComplianceThreshold)
		}
		cfg.MaxErrors = sec.ComplianceThresholdError // -1 = unset sentinel preserved
	}
	return cfg
}

func runInitWizard(defaults initConfig, detected, projectDir string) initConfig {
	cfg := defaults
	r := bufio.NewReader(os.Stdin)

	fmt.Println()

	// Profile — show pick list when available
	cfg.Profile = promptProfile(r, defaults.Profile, detected)

	// Compliance defaults: profile recommendation → registry settings.toml → hardcoded fallback
	thresholdDefault := defaults.Threshold
	maxDefault := defaults.MaxErrors
	if cfg.Profile != "" {
		opts := resolveOpts(projectDir)
		if p, err := profiles.ResolveWithOptions(cfg.Profile, projectDir, opts); err == nil {
			if thresholdDefault == 0 && p.ComplianceThreshold > 0 {
				thresholdDefault = int(p.ComplianceThreshold)
			}
			if maxDefault < 0 && p.ComplianceThresholdError >= 0 {
				maxDefault = p.ComplianceThresholdError
			}
		}
		// Fall back to registry settings.toml if profile didn't carry compliance fields
		if thresholdDefault == 0 || maxDefault < 0 {
			if regHome := findProfileRegistry(cfg.Profile); regHome != "" {
				regPath := filepath.Join(regHome, "settings.toml")
				if fs, err := settings.ParseFile(regPath); err == nil {
					r2 := settings.Merge([]settings.FileSettings{fs}, []string{regPath})
					sec := r2.ResolveSection(cfg.Profile)
					if thresholdDefault == 0 && sec.ComplianceThreshold > 0 {
						thresholdDefault = int(sec.ComplianceThreshold)
					}
					if maxDefault < 0 && sec.ComplianceThresholdError >= 0 {
						maxDefault = sec.ComplianceThresholdError
					}
				}
			}
		}
	}
	if thresholdDefault <= 0 {
		thresholdDefault = 80
	}
	if maxDefault < 0 {
		maxDefault = 0
	}

	// Compliance threshold
	fmt.Printf("  Compliance threshold %% [%d]: ", thresholdDefault)
	line := readLine(r)
	if line == "" {
		cfg.Threshold = thresholdDefault
	} else if n, err := strconv.Atoi(line); err == nil && n >= 0 && n <= 100 {
		cfg.Threshold = n
	} else {
		fmt.Fprintf(os.Stderr, "  invalid threshold %q — using %d\n", line, thresholdDefault)
		cfg.Threshold = thresholdDefault
	}

	// Max allowed errors
	fmt.Printf("  Max allowed errors [%d]: ", maxDefault)
	line = readLine(r)
	if line == "" {
		cfg.MaxErrors = maxDefault
	} else if n, err := strconv.Atoi(line); err == nil && n >= 0 {
		cfg.MaxErrors = n
	} else {
		fmt.Fprintf(os.Stderr, "  invalid value %q — using %d\n", line, maxDefault)
		cfg.MaxErrors = maxDefault
	}

	fmt.Println()
	return cfg
}

// promptProfile shows a scrollable TUI profile picker with ranked suggestions.
// Falls back to free-text when skills are not installed or TUI fails.
func promptProfile(r *bufio.Reader, current, detected string) string {
	items, annotations := listRankedProfileItems(detected)

	// fallback: returned on cancel — prefer saved setting, then detected
	fallback := current
	if fallback == "" {
		fallback = detected
	}

	// cursor: detected wins; fallback to current/saved
	cursorItem := detected
	if cursorItem == "" {
		cursorItem = fallback
	}

	if len(items) == 0 {
		// Skills not installed yet — plain text input
		hint := fallback
		if detected != "" {
			hint = detected + " (detected)"
		}
		fmt.Printf("  Profile [%s]: ", hint)
		if line := readLine(r); line != "" {
			return line
		}
		return fallback
	}

	fmt.Printf("  Profile [%s]:", cursorItem)

	chosen, ok := tui.RunProfileSelect(items, annotations, cursorItem)
	if !ok {
		return fallback
	}
	switch chosen {
	case tui.ProfileSelectNone:
		return ""
	case tui.ProfileSelectOther:
		fmt.Printf("  Profile name: ")
		if line := readLine(r); line != "" {
			return line
		}
		return fallback
	}
	return chosen
}

// profileEntry is used internally for ranking before building the picker list.
type profileEntry struct {
	name     string
	registry string
	tier     int    // 1=exact, 2=extends detected, 3=tags detected, 4=other
	ann      string // picker annotation
}

// listRankedProfileItems scans all registries and returns a ranked picker list.
// items may contain ProfileSectionPrefix headers, ProfileSelectNone, ProfileSelectOther.
// annotations maps name → annotation string.
func listRankedProfileItems(detected string) (items []string, annotations map[string]string) {
	annotations = make(map[string]string)

	regs := skills.AllRegistries()
	if len(regs) == 0 {
		return nil, annotations
	}
	seen := make(map[string]struct{})
	var entries []profileEntry

	addEntry := func(name, regName string, tier int, ann string) {
		key := name + "\x00" + regName
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		regLabel := "[" + regName + "]"
		if ann != "" {
			ann += " " + regLabel
		} else {
			ann = regLabel
		}
		entries = append(entries, profileEntry{name, regName, tier, ann})
	}

	for _, reg := range regs {
		regName := reg.Name

		// Named profile TOML files: <registry-home>/profiles/*.toml
		profDir := filepath.Join(reg.Home, "profiles")
		if files, err := os.ReadDir(profDir); err == nil {
			for _, f := range files {
				if f.IsDir() || !strings.HasSuffix(f.Name(), ".toml") {
					continue
				}
				name := strings.TrimSuffix(f.Name(), ".toml")
				tier, ann := rankProfileFile(name, filepath.Join(profDir, f.Name()), detected)
				addEntry(name, regName, tier, ann)
			}
		}

	}

	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].tier != entries[j].tier {
			return entries[i].tier < entries[j].tier
		}
		return entries[i].name < entries[j].name
	})

	var suggested, others []profileEntry
	for _, e := range entries {
		if e.tier <= 3 {
			suggested = append(suggested, e)
		} else {
			others = append(others, e)
		}
	}

	if len(suggested) > 0 && detected != "" {
		items = append(items, tui.ProfileSectionPrefix+"Suggested — "+detected)
		for _, e := range suggested {
			items = append(items, e.name)
			if e.ann != "" {
				annotations[e.name] = e.ann
			}
		}
	}
	if len(others) > 0 {
		if len(suggested) > 0 {
			items = append(items, tui.ProfileSectionPrefix+"All profiles")
		}
		for _, e := range others {
			items = append(items, e.name)
			if e.ann != "" {
				annotations[e.name] = e.ann
			}
		}
	}

	items = append(items, tui.ProfileSelectNone, tui.ProfileSelectOther)
	return items, annotations
}

func rankProfileFile(name, path, detected string) (tier int, ann string) {
	if detected == "" {
		return 4, ""
	}
	if name == detected {
		return 1, "(detected)"
	}
	meta := profiles.ReadMeta(path)
	for _, e := range meta.Extends {
		if e == detected || strings.HasSuffix(e, "/"+detected) {
			return 2, "extends " + detected
		}
	}
	for _, t := range meta.Tags {
		if t == detected {
			return 3, "tagged " + detected
		}
	}
	return 4, ""
}

func readLine(r *bufio.Reader) string {
	line, _ := r.ReadString('\n')
	return strings.TrimSpace(line)
}

func printProfilePreview(profileName, projectDir string) {
	if profileName == "" {
		return
	}
	if len(skills.AllSkillsRegistries()) == 0 {
		return // grimoire not yet installed — skip silently
	}
	p, err := profiles.ResolveWithOptions(profileName, projectDir, resolveOpts(projectDir))
	if err != nil || len(p.Skills) == 0 {
		return
	}
	fmt.Println()
	fmt.Printf("  Active profile: %s (%d skills)\n", profileName, len(p.Skills))
	const preview = 6
	for i, sk := range p.Skills {
		if i >= preview {
			fmt.Printf("    %s\n", colorize(ansiGray, fmt.Sprintf("… and %d more — run `grimoire list` to see all", len(p.Skills)-preview)))
			break
		}
		fmt.Printf("    %s %s\n", colorize(ansiGreen, "·"), sk.Name)
	}
}

func printInitReport(projectDir string) {
	reportPath := resolvedReportPath(projectDir)
	report, err := compliance.Load(reportPath)
	if err != nil {
		// No report yet — show next-step guidance
		fmt.Println()
		fmt.Println("  Next steps:")
		fmt.Printf("    1. Ask your AI to run %s\n", colorize(ansiGreen, "/check-best-practice-compliance"))
		fmt.Printf("    2. Run %s to see your compliance score\n", colorize(ansiGreen, "grimoire check"))
		return
	}

	pass := report.Threshold.Status == "pass"
	statusLabel := colorize(ansiGreen, "PASS")
	if !pass {
		statusLabel = colorize(ansiRed, "FAIL")
	}

	fmt.Println()
	fmt.Printf("  Existing report: %.1f%% — %s\n", report.Coverage.OverallPct, statusLabel)

	errors := filterBySeverity(report.Diagnostics, 1)
	warnings := filterBySeverity(report.Diagnostics, 2)
	shown := 0
	for i := range errors {
		if shown >= 3 {
			break
		}
		loc := formatLoc(&errors[i])
		fmt.Printf("    %s %s%s\n", colorize(ansiRed, "✗"), errors[i].Message, loc)
		shown++
	}
	for i := range warnings {
		if shown >= 3 {
			break
		}
		loc := formatLoc(&warnings[i])
		fmt.Printf("    %s %s%s\n", colorize(ansiYellow, "⚠"), warnings[i].Message, loc)
		shown++
	}
	total := len(errors) + len(warnings)
	if total > 3 {
		fmt.Printf("    %s\n", colorize(ansiGray, fmt.Sprintf("… and %d more — run `grimoire check` for full report", total-3)))
	}

	fmt.Println()
	fmt.Printf("  Fix first finding: ask your AI to run %s\n", colorize(ansiGreen, "/fix-best-practice-finding"))
}

func writeSettings(dir string, cfg initConfig) error {
	var profileLine string
	if cfg.Profile != "" {
		profileLine = fmt.Sprintf("profiles = [%q]", cfg.Profile)
	} else {
		profileLine = `# profiles = ["engineering"]   # uncomment and set your profile`
	}

	var thresholdBlock string
	if cfg.Profile != "" && (cfg.Threshold > 0 || cfg.MaxErrors >= 0) {
		thresholdBlock = fmt.Sprintf("\n[standards.%s]\n", cfg.Profile)
		if cfg.Threshold > 0 {
			thresholdBlock += fmt.Sprintf("compliance-threshold = %d\n", cfg.Threshold)
		}
		if cfg.MaxErrors >= 0 {
			thresholdBlock += fmt.Sprintf("compliance-threshold-error = %d\n", cfg.MaxErrors)
		}
	}

	content := fmt.Sprintf(`# Grimoire settings
# Docs: https://github.com/jeffreytse/grimoire/blob/main/docs/settings.md

[core]
# home = "~/.grimoire"                  # override grimoire home
# agents = ["claude", "codex"]          # pinned targets; empty = auto-detect
# install-mode = "symlink"              # "symlink" (default) | "copy"
# update-concurrency = 8               # max concurrent registry pulls (default 8); 0 = unlimited

[standards]
%s
%s
# report-path = ".grimoire/report.json" # custom compliance report path
# staleness-days = 14                   # warn if registry not pulled in N days
#
# Domain standards example:
# [standards.engineering]
# practices = ["apply-solid-principles", "apply-kiss-principle"]
# compliance-threshold = 80
# compliance-threshold-error = 0

# Inline profile (alternative to profiles/<name>.toml file):
# [profiles.myteam]
# description = "My team standards"
# tags = ["engineering"]
# compliance-threshold = 80
# compliance-threshold-error = 0
#
# [[profiles.myteam.skills]]
# name = "apply-solid-principles"
# priority = 1
`, profileLine, thresholdBlock)

	path := filepath.Join(dir, "settings.toml")
	return os.WriteFile(path, []byte(content), 0o644)
}
