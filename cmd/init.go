package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/compliance"
	"github.com/jeffreytse/grimoire/internal/config"
	grimctx "github.com/jeffreytse/grimoire/internal/context"
	"github.com/jeffreytse/grimoire/internal/detect"
	gitops "github.com/jeffreytse/grimoire/internal/git"
	"github.com/jeffreytse/grimoire/internal/manifest"
	"github.com/jeffreytse/grimoire/internal/profiles"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var (
	flagInitYes       bool
	flagInitAuto      bool
	flagInitGlobal    bool
	flagInitProfile   string
	flagInitThreshold int
	flagInitMaxErrors int
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize .grimoire/ in the current project",
	RunE:  runInit,
}

func init() {
	initCmd.Flags().BoolVarP(&flagInitYes, "yes", "y", false, "accept all defaults without prompting")
	initCmd.Flags().BoolVar(&flagInitAuto, "auto", false, "detect project context and pre-populate grimoire.toml automatically")
	initCmd.Flags().BoolVar(&flagInitGlobal, "global", false, "initialize global grimoire.toml (~/.config/grimoire/grimoire.toml)")
	initCmd.Flags().StringVar(&flagInitProfile, "profile", "", "profile to activate (skips prompt)")
	initCmd.Flags().IntVar(&flagInitThreshold, "threshold", 0, "compliance threshold % (skips prompt)")
	initCmd.Flags().IntVar(&flagInitMaxErrors, "max-errors", -1, "max allowed errors (skips prompt)")
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

	// --global flag: initialize user-global grimoire.toml
	if flagInitGlobal {
		return runInitGlobal()
	}

	// --auto flag: detect context, populate grimoire.toml, run install
	if flagInitAuto {
		return runInitAuto(cwd, dir, reinit)
	}

	hasExplicit := flagInitProfile != "" || flagInitThreshold != 0 || flagInitMaxErrors != -1

	// Skip wizard hint when user already has profiles configured in global/extends settings.
	if !hasExplicit && tui.IsTTY() && !flagInitYes {
		resolved, _ := config.Load(cwd)
		if !reinit && len(resolved.Core.Profiles) > 0 {
			fmt.Printf("  Using profiles from your settings: %s\n", strings.Join(resolved.Core.Profiles, ", "))
			// fall through to wizard below — loadExistingInitConfig will pick these up
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
	manifestPath := manifest.ProjectPath(cwd)
	pkgName := filepath.Base(cwd)
	if created, err := scaffoldManifest(manifestPath, pkgName, cfg.Profile, cfg.Threshold, cfg.MaxErrors); err != nil {
		fmt.Fprintf(os.Stderr, "  warn: could not create grimoire.toml: %v\n", err)
	} else if created {
		fmt.Printf("  Created grimoire.toml (skill manifest and config)\n")
	}
	printInitSuccess(reinit)
	printProfilePreview(cfg.Profile, cwd)
	printInitReport(cwd)
	return nil
}

func runInitGlobal() error {
	cwd := getProjectDir()
	manifestPath := manifest.GlobalPath()
	dir := filepath.Dir(manifestPath)
	_, statErr := os.Stat(manifestPath)
	reinit := statErr == nil
	detected := detect.Profile(cwd)

	existing := loadExistingInitConfig(cwd, detected)
	cfg := runInitWizard(existing, detected, cwd)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}
	if created, err := scaffoldManifest(manifestPath, "global", cfg.Profile, cfg.Threshold, cfg.MaxErrors); err != nil {
		fmt.Fprintf(os.Stderr, "  warn: could not create grimoire.toml: %v\n", err)
	} else if created {
		fmt.Printf("  Created grimoire.toml (%s)\n", manifestPath)
	}
	printInitSuccess(reinit)
	printProfilePreview(cfg.Profile, cwd)
	printInitReport(cwd)
	return nil
}

func printInitSuccess(reinit bool) {
	if reinit {
		fmt.Printf("✓ Grimoire reinitialized.\n")
	} else {
		fmt.Printf("✓ Grimoire initialized.\n")
	}
}

// loadExistingInitConfig reads current grimoire.toml (if any) and returns an
// initConfig pre-populated with those values — used as wizard defaults on re-init.
func loadExistingInitConfig(cwd, detected string) initConfig {
	cfg := initConfig{Profile: detected, Threshold: 0, MaxErrors: -1}

	r, err := config.Load(cwd)
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

	// Exit cleanly on Ctrl+C (SIGINT), even if bubbletea left the terminal in raw mode.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		if _, ok := <-sigCh; ok {
			fmt.Fprintln(os.Stderr)
			os.Exit(0) //nolint:revive // intentional: cleanly exit on Ctrl+C inside blocking wizard goroutine
		}
	}()
	defer signal.Stop(sigCh)

	fmt.Println()

	// Profile — show pick list when available
	cfg.Profile = promptProfile(r, defaults.Profile, detected)

	// Compliance defaults: profile recommendation → package settings.toml → hardcoded fallback
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
		// Fall back to package grimoire.toml if profile didn't carry compliance fields
		if thresholdDefault == 0 || maxDefault < 0 {
			if regHome := findProfilePackage(cfg.Profile); regHome != "" {
				regPath := filepath.Join(regHome, "grimoire.toml")
				if fs, err := config.ParseFile(regPath); err == nil {
					r2 := config.Merge([]config.FileConfig{fs}, []string{regPath})
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

	cfg.Threshold = thresholdDefault
	cfg.MaxErrors = maxDefault

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

	hasProfiles := len(items) > 2 // more than just ProfileSelectNone + ProfileSelectOther
	if !hasProfiles {
		fmt.Fprintf(os.Stderr, "  %s  no profiles found — run: grimoire package update\n", tui.IconWarn)
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
	name string
	pkg  string
	tier int    // 1=exact, 2=extends detected, 3=tags detected, 4=other
	ann  string // picker annotation
}

// listRankedProfileItems scans all packages and returns a ranked picker list.
// items may contain ProfileSectionPrefix headers, ProfileSelectNone, ProfileSelectOther.
// annotations maps name → annotation string.
func listRankedProfileItems(detected string) (items []string, annotations map[string]string) {
	annotations = make(map[string]string)

	regs := skills.AllPackages()
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

		// Named profile TOML files: <package-home>/profiles/*.toml
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
	line, err := r.ReadString('\n')
	// Terminal in raw mode sends 0x03 for Ctrl+C instead of SIGINT.
	if err != nil || strings.ContainsRune(line, '\x03') {
		fmt.Fprintln(os.Stderr)
		return ""
	}
	return strings.TrimSpace(line)
}

func printProfilePreview(profileName, projectDir string) {
	if profileName == "" {
		return
	}
	if len(skills.AllSkillsPackages()) == 0 {
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

// runInitAuto implements `grimoire init --auto`:
// detect project context → load profile's skill list → scaffold grimoire.toml with deps → run install.
func runInitAuto(cwd, grimDir string, reinit bool) error {
	// Bootstrap: auto-install the official package if nothing is installed yet.
	if len(skills.AllPackages()) == 0 {
		home := skills.OfficialPackageHome()
		url := skills.GrimoireRepoURL()
		if !filepath.IsAbs(url) {
			fmt.Printf("  No packages installed — cloning official package…\n")
			if mkErr := os.MkdirAll(filepath.Dir(home), 0o755); mkErr != nil {
				return fmt.Errorf("creating package dir: %w", mkErr)
			}
			if cloneErr := gitops.Clone(url, home); cloneErr != nil {
				return fmt.Errorf("cloning official package: %w", cloneErr)
			}
			fmt.Printf("  %s  Package installed to %s\n", tui.IconOK, home)
		}
	}

	ctx := grimctx.Detect(cwd)

	if ctx.Profile == "" {
		fmt.Println("  No profile detected — falling back to interactive init.")
		fmt.Println("  Run: grimoire init (interactive wizard)")
		return nil
	}

	fmt.Printf("  Detected profile: %s\n", ctx.Profile)

	// Scaffold grimoire.toml if absent, then populate with profile skills.
	manifestPath := manifest.ProjectPath(cwd)
	mf := manifest.ManifestFile{}
	if _, err := os.Stat(manifestPath); err == nil {
		mf, _ = manifest.ParseFile(manifestPath)
	} else {
		// new file — fill in package name
		mf.Package.Name = filepath.Base(cwd)
		mf.Package.Version = "0.1.0"
	}

	// Resolve profile to get its skill list.
	opts := resolveOpts(cwd)
	p, err := profiles.ResolveWithOptions(ctx.Profile, cwd, opts)
	if err != nil {
		return fmt.Errorf("resolving profile %q: %w", ctx.Profile, err)
	}
	if len(p.Skills) == 0 {
		fmt.Printf("  Profile %q has no skills — nothing to add.\n", ctx.Profile)
	} else {
		if mf.Deps == nil {
			mf.Deps = make(map[string]manifest.DepSpec)
		}
		added := 0
		for _, sk := range p.Skills {
			if _, exists := mf.Deps[sk.Name]; !exists {
				mf.Deps[sk.Name] = manifest.DepSpec{Version: "*"}
				added++
			}
		}
		fmt.Printf("  Added %d skill(s) from profile %q to grimoire.toml\n", added, ctx.Profile)
	}

	// Write standards section with detected profile.
	if len(mf.Standards.Profiles) == 0 {
		mf.Standards.Profiles = []string{ctx.Profile}
	}

	if err := manifest.WriteFile(manifestPath, &mf); err != nil {
		return fmt.Errorf("writing grimoire.toml: %w", err)
	}

	if err := os.MkdirAll(grimDir, 0o755); err != nil {
		return fmt.Errorf("creating .grimoire/: %w", err)
	}

	printInitSuccess(reinit)
	printProfilePreview(ctx.Profile, cwd)

	// Run install to materialize the deps we just wrote.
	fmt.Println()
	fmt.Println("  Installing skills from grimoire.toml…")
	return runInstallFromManifest(cwd)
}

// scaffoldManifest writes a starter grimoire.toml at manifestPath if one does not already exist.
// Returns (true, nil) when a new file is created; (false, nil) when it already exists.
func scaffoldManifest(manifestPath, pkgName, profile string, threshold, maxErrors int) (bool, error) {
	if _, err := os.Stat(manifestPath); err == nil {
		return false, nil // already exists
	}

	if pkgName == "" || pkgName == "." || pkgName == "/" {
		pkgName = "my-project"
	}

	profileLine := `# profiles = ["engineering"]   # uncomment and set your profile`
	if profile != "" {
		profileLine = fmt.Sprintf("profiles = [%q]", profile)
	}

	var thresholdBlock string
	if profile != "" && (threshold > 0 || maxErrors >= 0) {
		thresholdBlock = fmt.Sprintf("\n[standards.%s]\n", profile)
		if threshold > 0 {
			thresholdBlock += fmt.Sprintf("compliance-threshold = %d\n", threshold)
		}
		if maxErrors >= 0 {
			thresholdBlock += fmt.Sprintf("compliance-threshold-error = %d\n", maxErrors)
		}
	}

	content := fmt.Sprintf(`# grimoire.toml — skill manifest and config
# Docs: https://grimoire.jeffreytse.net/docs/manifest

[package]
name = %q
version = "0.1.0"

[dependencies]
# Install skills with: grimoire install <skill-name>
# apply-solid-principles = "*"
# apply-dry-principle = "*"

[standards]
%s
%s`, pkgName, profileLine, thresholdBlock)

	return true, os.WriteFile(manifestPath, []byte(content), 0o644)
}
