package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/config"
	gitops "github.com/jeffreytse/grimoire/internal/git"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Manage grimoire skill packages",
	Long: `Manage grimoire skill packages — official and user-defined.

Multiple packages are searched in priority order (highest first).
User packages do not need to follow the official STANDARD.md.

  grimoire package list                        list all packages
  grimoire package add <name> <grimoire-ref>    add a package and clone it
  grimoire package remove <name>               remove a package
  grimoire package enable <name>               enable a disabled package
  grimoire package disable <name>              disable a package without removing it
  grimoire package set <grimoire-ref>           set the official package URL
  grimoire package reset                       revert official package to built-in default
  grimoire package update [<name>]             pull latest from all (or named) packages
  grimoire package validate [<path-or-name>]   validate a package's structure
  grimoire package publish                      guide for publishing a package to the community`,
}

var (
	flagPackageListJSON    bool
	flagPackageUpdateRef   string
	flagPackageAddPriority int
)

var packageListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured packages",
	RunE:  runPackageList,
}

var packageSetCmd = &cobra.Command{
	Use:   "set <grimoire-ref>",
	Short: "Set the official package URL (owner/repo[@version] or full URL)",
	Args:  cobra.ExactArgs(1),
	RunE:  runPackageSet,
}

var packageResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Clear core.package (revert to built-in default)",
	Args:  cobra.NoArgs,
	RunE:  runPackageReset,
}

var packageAddCmd = &cobra.Command{
	Use:   "add <name> <grimoire-ref>",
	Short: "Add a named package and clone it",
	Long: `Add a named package to the [[package]] list and clone it locally.

<name>         short identifier for this package (e.g. "my-team", "plugins")
<grimoire-ref>  git URL, owner/repo[@version] shorthand, or absolute local path

User packages do not need to follow the official STANDARD.md — any
directory with a skills/ folder (or flat .md files) is accepted.`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runPackageAdd,
}

var packageRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a package by name",
	Args:  cobra.ExactArgs(1),
	RunE:  runPackageRemove,
}

var packageEnableCmd = &cobra.Command{
	Use:   "enable <name>",
	Short: "Enable a previously disabled package",
	Args:  cobra.ExactArgs(1),
	RunE:  runPackageEnable,
}

var packageDisableCmd = &cobra.Command{
	Use:   "disable <name>",
	Short: "Disable a package without removing it",
	Args:  cobra.ExactArgs(1),
	RunE:  runPackageDisable,
}

var packageUpdateCmd = &cobra.Command{
	Use:   "update [<name>]",
	Short: "Pull latest from all (or one named) package",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runPackageUpdate,
}

var packageValidateCmd = &cobra.Command{
	Use:   "validate [<path-or-name>]",
	Short: "Validate a package's structure (for package authors)",
	Long: `Validate a package's directory structure before publishing.

Without arguments, validates the current directory (useful when authoring a package).
With an argument, validates an installed package by name or an absolute path.

Checks:
  - Package markers present (skills/, profiles/, grimoire.toml)
  - skills/: domain directories with SKILL.md files
  - profiles/: .toml files are valid TOML
  - grimoire.toml: valid TOML if present`,
	Args: cobra.MaximumNArgs(1),
	RunE: runPackageValidate,
}

func init() {
	packageListCmd.Flags().BoolVar(&flagPackageListJSON, "json", false, "output as JSON")
	packageUpdateCmd.Flags().StringVar(&flagPackageUpdateRef, "ref", "", "change pinned version/branch/commit")
	packageAddCmd.Flags().IntVar(&flagPackageAddPriority, "priority", 0, "skill resolution priority (higher wins; default 50 for user packages)")
	packageCmd.AddCommand(packageListCmd)
	packageCmd.AddCommand(packageAddCmd)
	packageCmd.AddCommand(packageRemoveCmd)
	packageCmd.AddCommand(packageEnableCmd)
	packageCmd.AddCommand(packageDisableCmd)
	packageCmd.AddCommand(packageSetCmd)
	packageCmd.AddCommand(packageResetCmd)
	packageCmd.AddCommand(packageUpdateCmd)
	packageCmd.AddCommand(packageValidateCmd)
	packageCmd.AddCommand(packagePublishCmd)
}

type packageListEntry struct {
	Name          string `json:"name"`
	URL           string `json:"url"`
	Version       string `json:"version,omitempty"`
	Priority      int    `json:"priority"`
	SkillsCount   int    `json:"skills_count"`
	ProfilesCount int    `json:"profiles_count"`
	Cloned        bool   `json:"cloned"`
	Enabled       bool   `json:"enabled"`
	Kind          string `json:"kind"` // "official" | "user" | "local"
}

func runPackageList(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	regs := skills.AllPackages()
	entries := make([]packageListEntry, 0, len(regs))
	officialSeen := false

	for _, reg := range regs {
		var url, ver string
		enabled := true
		for _, rd := range cfg.Packages {
			if rd.Name == reg.Name {
				url, ver = config.ParseRef(rd.URL)
				if url == "" {
					url = rd.URL
				}
				enabled = rd.Enabled
				break
			}
		}
		if url == "" {
			url = skills.GrimoireRepoURL()
		}
		kind := "user"
		if reg.Official && !officialSeen {
			kind = "official"
			officialSeen = true
		}
		if filepath.IsAbs(url) {
			kind = "local"
		}
		if ver == "" && reg.Official {
			ver = skills.GrimoireVersion()
		}
		entries = append(entries, packageListEntry{
			Name:          reg.Name,
			URL:           url,
			Version:       ver,
			Priority:      reg.Priority,
			SkillsCount:   countSkills(filepath.Join(reg.Home, "skills")),
			ProfilesCount: countProfiles(filepath.Join(reg.Home, "profiles")),
			Cloned:        dirExists(reg.Home),
			Enabled:       enabled,
			Kind:          kind,
		})
	}

	if flagPackageListJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(entries)
	}

	for _, e := range entries {
		icon := tui.IconOK
		if !e.Cloned {
			icon = tui.IconWarn
		}
		if !e.Enabled {
			icon = tui.IconSkip
		}
		ref := e.URL
		if e.Version != "" {
			ref += "@" + e.Version
		}
		tag := tui.StyleDim.Render("[" + e.Kind + "]")
		priStr := fmt.Sprintf("p%d", e.Priority)
		if !e.Enabled {
			priStr = "disabled"
		}
		fmt.Printf("  %s  %-30s %s %-10s %s\n", icon, e.Name, tag, priStr, tui.StyleDim.Render(ref))
		switch {
		case !e.Enabled:
			fmt.Printf("         disabled — run: grimoire package enable %s\n\n", e.Name)
		case e.Cloned:
			fmt.Printf("         %d skills · %d profiles\n\n", e.SkillsCount, e.ProfilesCount)
		default:
			fmt.Printf("         not cloned — run: grimoire package update %s\n\n", e.Name)
		}
	}
	return nil
}

func runPackageSet(cmd *cobra.Command, args []string) error {
	ref := args[0]
	u, _ := config.ParseRef(ref)
	if !skills.IsGitURL(u) && !filepath.IsAbs(u) {
		return fmt.Errorf("invalid grimoire-ref %q — expected owner/repo[@version], git URL, or absolute path", ref)
	}
	if filepath.IsAbs(u) {
		if _, err := os.Stat(u); err != nil {
			return fmt.Errorf("local path %q not found", u)
		}
	}

	cfg, err := config.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	// Update or add the official=true [[package]] entry.
	for i, rd := range cfg.Packages {
		if !rd.Official {
			continue
		}
		cfg.Packages[i].URL = ref
		if err := config.SaveGlobal(cfg); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}
		fmt.Printf("%s  official package = %s\n", tui.IconOK, ref)
		if filepath.IsAbs(u) {
			fmt.Printf("   local package set as official\n")
		} else {
			fmt.Printf("   run: grimoire package update official  to apply\n")
		}
		return nil
	}
	cfg.Packages = append(cfg.Packages, config.PackageDef{
		Name:     "official",
		URL:      ref,
		Official: true,
		Priority: 100,
		Enabled:  true,
	})
	if err := config.SaveGlobal(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}
	fmt.Printf("%s  official package = %s\n", tui.IconOK, ref)
	if filepath.IsAbs(u) {
		fmt.Printf("   local package set as official\n")
	} else {
		fmt.Printf("   run: grimoire package update official  to apply\n")
	}
	return nil
}

func runPackageReset(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	// Remove the official=true [[package]] entry so GrimoireRepo constant takes effect.
	var kept []config.PackageDef
	for _, rd := range cfg.Packages {
		if !rd.Official {
			kept = append(kept, rd)
		}
	}
	cfg.Packages = kept
	if err := config.SaveGlobal(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}
	fmt.Printf("%s  official package reset — using built-in default (%s)\n", tui.IconOK, skills.GrimoireRepo)
	return nil
}

func runPackageEnable(cmd *cobra.Command, args []string) error {
	return setPackageEnabled(args[0], true)
}

func runPackageDisable(cmd *cobra.Command, args []string) error {
	return setPackageEnabled(args[0], false)
}

func setPackageEnabled(name string, enabled bool) error {
	cfg, err := config.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	for i, rd := range cfg.Packages {
		if rd.Name != name {
			continue
		}
		cfg.Packages[i].Enabled = enabled
		if err := config.SaveGlobal(cfg); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}
		verb := "enabled"
		if !enabled {
			verb = "disabled"
		}
		fmt.Printf("%s  %s %s\n", tui.IconOK, verb, name)
		return nil
	}
	return fmt.Errorf("package %q not found in [[package]] — check: grimoire package list", name)
}

// updateNamedPackage clones or pulls a named [[package]] entry.
func updateNamedPackage(name, ref, forceVer string, w io.Writer) error {
	u, ver := config.ParseRef(ref)
	if u == "" {
		u = ref
	}
	if forceVer != "" {
		ver = forceVer
	}
	if flagPackageUpdateRef != "" {
		ver = flagPackageUpdateRef
	}

	if filepath.IsAbs(u) {
		if _, err := os.Stat(u); err != nil {
			return fmt.Errorf("local package %q not found", u)
		}
		fmt.Fprintf(w, "  %s  %s is a local path — no update needed\n", tui.IconOK, name)
		return nil
	}

	// Derive path from URL+version so each version gets its own directory.
	dest := skills.PackageHome(config.DeriveVersionedName(u, ver))
	if !dirExists(dest) {
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return fmt.Errorf("creating dir: %w", err)
		}
		if err := gitops.Clone(u, dest); err != nil {
			return fmt.Errorf("cloning %s: %w", name, err)
		}
		skills.InvalidateSkillCache(dest)
		if ver != "" {
			if err := gitops.CheckoutVersion(dest, ver); err != nil {
				return fmt.Errorf("checkout %s: %w", name, err)
			}
			skills.InvalidateSkillCache(dest)
		}
		fmt.Fprintf(w, "  %s  %s cloned\n", tui.IconOK, name)
		return nil
	}

	oldState, _ := gitops.CurrentState(dest)
	if ver != "" {
		if err := gitops.CheckoutVersion(dest, ver); err != nil {
			return fmt.Errorf("checking out version for %s: %w", name, err)
		}
		skills.InvalidateSkillCache(dest)
	} else {
		if err := gitops.PullWithForceFallback(dest); err != nil {
			return fmt.Errorf("updating %s: %w", name, err)
		}
		skills.InvalidateSkillCache(dest)
	}
	fmt.Fprintf(w, "  %s  %s up to date\n", tui.IconOK, name)
	if changes, err := gitops.PackageChangesSince(dest, oldState.Commit); err == nil {
		printPackageChanges(changes, dest, oldState.Commit, w)
	}
	return nil
}

func runPackageUpdate(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	r, _ := config.Load(getProjectDir())

	if len(args) == 1 {
		name := args[0]
		board := tui.NewStatusBoard([]string{name})
		stopSpinner := board.StartSpinner()
		board.SetUpdating(0)
		buf := &bytes.Buffer{}
		var updateErr error
		if name == skills.OfficialPackageName || isOfficialByName(name) {
			updateErr = updateCorePackage(buf)
		} else if ref := findPackageRef(name, &cfg); ref != "" {
			updateErr = updateNamedPackage(name, ref, "", buf)
		} else {
			updateErr = fmt.Errorf("package %q not found — check: grimoire package list", name)
		}
		stopSpinner()
		if updateErr != nil {
			board.SetDone(0, tui.IconError)
			board.Finish()
			return updateErr
		}
		board.SetDone(0, tui.IconDone)
		board.Finish()
		os.Stdout.Write(buf.Bytes())
		return nil
	}

	// Build work list: new-model packages + old-model extends targets.
	type workItem struct {
		name   string
		ref    string
		isCore bool
	}
	var items []workItem
	var pinnedSkipped []string

	if len(cfg.Packages) > 0 {
		for _, rd := range cfg.Packages {
			if !rd.Enabled {
				continue
			}
			u, ver := config.ParseRef(rd.URL)
			if u == "" {
				u = rd.URL
			}
			if filepath.IsAbs(u) {
				continue // local package — user manages checkout
			}
			if ver != "" {
				pinnedSkipped = append(pinnedSkipped, rd.Name)
				continue // pinned package — immutable by intent; skip bulk update
			}
			items = append(items, workItem{
				name:   rd.Name,
				ref:    rd.URL,
				isCore: rd.Official,
			})
		}
	} else {
		// No [[package]] configured — update the implicit official package.
		items = append(items, workItem{name: skills.OfficialPackageName, isCore: true})
	}

	names := make([]string, len(items))
	for i, item := range items {
		names[i] = item.name
	}
	board := tui.NewStatusBoard(names)
	stopSpinner := board.StartSpinner()

	limit := 8
	if r.Core.UpdateConcurrency != nil {
		if *r.Core.UpdateConcurrency == 0 {
			limit = len(items)
		} else {
			limit = *r.Core.UpdateConcurrency
		}
	}
	if limit > len(items) {
		limit = len(items)
	}
	if limit == 0 {
		limit = 1
	}
	sem := make(chan struct{}, limit)

	type result struct {
		buf *bytes.Buffer
		err error
	}
	results := make([]result, len(items))
	var wg sync.WaitGroup
	for i, item := range items {
		wg.Add(1)
		sem <- struct{}{}
		board.SetUpdating(i)
		go func(i int, item workItem) {
			defer wg.Done()
			defer func() { <-sem }()
			buf := &bytes.Buffer{}
			var err error
			switch {
			case item.isCore:
				err = updateCorePackage(buf)
			default:
				err = updateNamedPackage(item.name, item.ref, "", buf)
			}
			if err != nil {
				board.SetDone(i, tui.IconError)
			} else {
				board.SetDone(i, tui.IconDone)
			}
			results[i] = result{buf: buf, err: err}
		}(i, item)
	}
	wg.Wait()
	stopSpinner()
	board.Finish()

	if len(pinnedSkipped) > 0 {
		fmt.Printf("  %s  pinned (skipped): %s\n", tui.IconSkip, strings.Join(pinnedSkipped, ", "))
		fmt.Printf("       to force-update: grimoire package update <name>\n")
	}

	for i, res := range results {
		if res.err != nil {
			fmt.Fprintf(os.Stderr, "  warn: %s: %v\n", items[i].name, res.err)
		}
		os.Stdout.Write(res.buf.Bytes())
	}
	return nil
}

// isOfficialByName checks if a named package is the (demote-resolved) official package.
// Uses AllPackages() so demoted lower-priority official entries return false.
func isOfficialByName(name string) bool {
	for _, reg := range skills.AllPackages() {
		if reg.Official && reg.Name == name {
			return true
		}
	}
	return false
}

// findPackageRef returns the URL ref for a named package in [[package]], or "".
func findPackageRef(name string, cfg *config.FileConfig) string {
	for _, rd := range cfg.Packages {
		if rd.Name == name {
			return rd.URL
		}
	}
	return ""
}

func updateCorePackage(w io.Writer) error {
	dest := skills.OfficialPackageHome()
	url := skills.GrimoireRepoURL()

	// Local path: verify it exists, no git ops
	if filepath.IsAbs(url) {
		if _, err := os.Stat(dest); err != nil {
			return fmt.Errorf("local package %q not found", dest)
		}
		fmt.Fprintf(w, "  %s  official package is a local path — no update needed\n", tui.IconOK)
		return nil
	}

	if !dirExists(dest) {
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return fmt.Errorf("creating dir: %w", err)
		}
		if err := gitops.Clone(url, dest); err != nil {
			return fmt.Errorf("cloning: %w", err)
		}
		skills.InvalidateSkillCache(dest)
		fmt.Fprintf(w, "  %s  official cloned\n", tui.IconOK)
		return nil
	}

	ver := ""
	cfg, _ := config.LoadGlobal()
	for _, rd := range cfg.Packages {
		if rd.Official {
			_, ver = config.ParseRef(rd.URL)
			break
		}
	}
	if flagPackageUpdateRef != "" {
		ver = flagPackageUpdateRef
		for i, rd := range cfg.Packages {
			if rd.Official {
				cfg.Packages[i].URL = url + "@" + ver
				_ = config.SaveGlobal(cfg)
				break
			}
		}
	}

	oldState, _ := gitops.CurrentState(dest)
	if ver != "" {
		if err := gitops.CheckoutVersion(dest, ver); err != nil {
			return fmt.Errorf("checking out version: %w", err)
		}
		skills.InvalidateSkillCache(dest)
	} else {
		if err := gitops.PullWithForceFallback(dest); err != nil {
			return fmt.Errorf("updating: %w", err)
		}
		skills.InvalidateSkillCache(dest)
	}
	fmt.Fprintf(w, "  %s  official up to date\n", tui.IconOK)
	if changes, err := gitops.PackageChangesSince(dest, oldState.Commit); err == nil {
		printPackageChanges(changes, dest, oldState.Commit, w)
	}
	return nil
}

func runPackageAdd(cmd *cobra.Command, args []string) error {
	if len(args) == 1 {
		return fmt.Errorf("usage: grimoire package add <name> <grimoire-ref>\n\nExample:\n  grimoire package add my-team https://github.com/acme/grimoire.git")
	}

	name := args[0]
	ref := args[1]

	if strings.ContainsAny(name, "/\\") {
		return fmt.Errorf("package name %q must not contain path separators — use a short identifier like %q", name, strings.ReplaceAll(name, "/", "-"))
	}

	u, _ := config.ParseRef(ref)
	if u == "" {
		u = ref
	}
	if !skills.IsGitURL(u) && !filepath.IsAbs(u) {
		return fmt.Errorf("invalid grimoire-ref %q — expected owner/repo[@version], git URL, or absolute path", ref)
	}
	if filepath.IsAbs(u) {
		if _, err := os.Stat(u); err != nil {
			return fmt.Errorf("local path %q not found", u)
		}
	}

	cfg, err := config.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Idempotent: if name already exists, update URL/priority.
	for i, existing := range cfg.Packages {
		if existing.Name != name {
			continue
		}
		cfg.Packages[i].URL = ref
		if flagPackageAddPriority > 0 {
			cfg.Packages[i].Priority = flagPackageAddPriority
		}
		if err := config.SaveGlobal(cfg); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}
		fmt.Printf("%s  updated package %s → %s\n", tui.IconOK, name, u)
		if filepath.IsAbs(u) {
			return nil
		}
		return updateNamedPackage(name, ref, "", os.Stdout)
	}

	rd := config.PackageDef{
		Name:    name,
		URL:     ref,
		Enabled: true,
	}
	if flagPackageAddPriority > 0 {
		rd.Priority = flagPackageAddPriority
	}
	cfg.Packages = append(cfg.Packages, rd)
	if err := config.SaveGlobal(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}
	fmt.Printf("%s  added package %s\n", tui.IconOK, name)

	if filepath.IsAbs(u) {
		sc := countSkills(filepath.Join(u, "skills"))
		pc := countProfiles(filepath.Join(u, "profiles"))
		if sc > 0 {
			fmt.Printf("   %d skills available from %s\n", sc, name)
		}
		if pc > 0 {
			fmt.Printf("   %d profiles available from %s\n", pc, name)
		}
		return nil
	}

	if err := updateNamedPackage(name, ref, "", os.Stdout); err != nil {
		return err
	}
	home := skills.PackageHome(name)
	if sc := countSkills(filepath.Join(home, "skills")); sc > 0 {
		fmt.Printf("   %d skills available from %s\n", sc, name)
	}
	if pc := countProfiles(filepath.Join(home, "profiles")); pc > 0 {
		fmt.Printf("   %d profiles available from %s\n", pc, name)
	}
	return nil
}

func runPackageRemove(cmd *cobra.Command, args []string) error {
	target := args[0]

	cfg, err := config.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	var kept []config.PackageDef
	removed := false
	for _, rd := range cfg.Packages {
		if rd.Name == target {
			removed = true
			continue
		}
		kept = append(kept, rd)
	}
	if !removed {
		return fmt.Errorf("package %q not found — check: grimoire package list", target)
	}
	cfg.Packages = kept
	if err := config.SaveGlobal(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}
	home := skills.PackageHome(target)
	fmt.Printf("%s  removed %s from [[package]]\n", tui.IconOK, target)
	fmt.Printf("   local clone at %s preserved — delete manually if no longer needed\n", home)
	return nil
}

var packagePublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Guide for publishing a grimoire package to the community",
	Long: `Print the publishing checklist and open GitHub to create a new package repo.

Community packages follow the naming convention:  grimoire-package-<name>
GitHub topic:  grimoire-package  (makes your package discoverable via search)`,
	Args: cobra.NoArgs,
	RunE: runPackagePublish,
}

func runPackagePublish(_ *cobra.Command, _ []string) error {
	fmt.Printf("\n%s\n\n", tui.StyleBold.Render("Publishing a Grimoire Package"))

	fmt.Printf("  %s  Step 1 — validate your package structure\n\n", tui.StyleDim.Render("1."))
	fmt.Printf("    grimoire package validate\n\n")

	fmt.Printf("  %s  Step 2 — create a GitHub repo\n\n", tui.StyleDim.Render("2."))
	fmt.Printf("    Naming convention:  grimoire-package-<yourname>\n")
	fmt.Printf("    Add GitHub topic:   grimoire-package\n")
	fmt.Printf("    (topic makes your package discoverable by other grimoire users)\n\n")
	fmt.Printf("    %s\n\n", "https://github.com/new")

	fmt.Printf("  %s  Step 3 — push and share\n\n", tui.StyleDim.Render("3."))
	fmt.Printf("    git remote add origin https://github.com/<you>/grimoire-package-<name>\n")
	fmt.Printf("    git push -u origin main\n\n")

	fmt.Printf("  %s  Step 4 — install from anywhere\n\n", tui.StyleDim.Render("4."))
	fmt.Printf("    grimoire package add <name> <owner>/<repo>\n\n")

	fmt.Printf("  %s  Tip: include a README with the install command above\n", tui.StyleDim.Render("★"))
	fmt.Printf("       so users can add your package with one line.\n\n")

	openBrowser("https://github.com/new")
	return nil
}

// openBrowser tries to open url in the default browser. Fails silently.
func openBrowser(url string) {
	var args []string
	switch {
	case commandExists("open"):
		args = []string{"open", url}
	case commandExists("xdg-open"):
		args = []string{"xdg-open", url}
	case commandExists("start"):
		args = []string{"cmd", "/c", "start", url}
	default:
		return
	}
	_ = exec.Command(args[0], args[1:]...).Start() //nolint:gosec // url is a hard-coded constant
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func dirExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func countSkills(skillsRoot string) int {
	all, _ := skills.ListAllSkills(skillsRoot)
	return len(all)
}

func countProfiles(profilesDir string) int {
	entries, err := os.ReadDir(profilesDir)
	if err != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".toml") {
			count++
		}
	}
	return count
}

func runPackageValidate(cmd *cobra.Command, args []string) error {
	var target string
	if len(args) == 1 {
		ref := args[0]
		// Treat as a filesystem path when absolute or starts with . / ..
		isPath := filepath.IsAbs(ref) || ref == "." || ref == ".." ||
			strings.HasPrefix(ref, "./") || strings.HasPrefix(ref, "../")
		if isPath {
			abs, err := filepath.Abs(ref)
			if err != nil {
				return err
			}
			target = abs
		} else {
			// Treat as installed package name or shorthand ref
			u, _ := config.ParseRef(ref)
			name := config.DerivePackageName(u)
			home := skills.PackageHome(name)
			if !dirExists(home) {
				return fmt.Errorf("package %q not installed — run: grimoire package update %s", ref, ref)
			}
			target = home
		}
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		target = cwd
	}

	type vcheck struct {
		status string
		detail string
	}
	var checks []vcheck
	allOK := true

	fail := func(detail string) { checks = append(checks, vcheck{"error", detail}); allOK = false }
	warn := func(detail string) { checks = append(checks, vcheck{"warn", detail}); allOK = false }
	pass := func(detail string) { checks = append(checks, vcheck{"ok", detail}) }
	skip := func(detail string) { checks = append(checks, vcheck{"skip", detail}) }

	// Package markers
	hasMarker := false
	for _, marker := range []string{"skills", "profiles", "grimoire.toml"} {
		if _, err := os.Stat(filepath.Join(target, marker)); err == nil {
			hasMarker = true
			break
		}
	}
	if !hasMarker {
		fail("no package markers found (expected: skills/, profiles/, or grimoire.toml)")
	} else {
		pass("package structure detected")
	}

	// skills/ structure
	skillsDir := filepath.Join(target, "skills")
	if _, err := os.Stat(skillsDir); err == nil {
		allSkills, _ := skills.ListAllSkills(skillsDir)
		if len(allSkills) == 0 {
			warn("skills/ found but no skills detected (expected: skills/<domain>/<name>/SKILL.md)")
		} else {
			missing := 0
			for i := range allSkills {
				sk := allSkills[i]
				if sk.Path == "" {
					missing++
				}
			}
			if missing > 0 {
				warn(fmt.Sprintf("%d skills, %d missing SKILL.md", len(allSkills), missing))
			} else {
				pass(fmt.Sprintf("%d skill(s), all have SKILL.md", len(allSkills)))
			}
		}
	} else {
		skip("no skills/ directory")
	}

	// profiles/ structure
	profilesDir := filepath.Join(target, "profiles")
	if _, err := os.Stat(profilesDir); err == nil {
		entries, _ := os.ReadDir(profilesDir)
		total, invalid := 0, 0
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".toml") {
				continue
			}
			total++
			if _, err := config.ParseFile(filepath.Join(profilesDir, e.Name())); err != nil {
				invalid++
			}
		}
		switch {
		case total == 0:
			warn("profiles/ found but empty (expected: profiles/<name>.toml)")
		case invalid > 0:
			fail(fmt.Sprintf("%d/%d profile TOML file(s) failed to parse", invalid, total))
		default:
			pass(fmt.Sprintf("%d profile(s), all valid TOML", total))
		}
	} else {
		skip("no profiles/ directory")
	}

	// root grimoire.toml
	rootConfig := filepath.Join(target, "grimoire.toml")
	if _, err := os.Stat(rootConfig); err == nil {
		if _, err := config.ParseFile(rootConfig); err != nil {
			fail(fmt.Sprintf("grimoire.toml parse error: %v", err))
		} else {
			pass("grimoire.toml is valid TOML")
		}
	} else {
		skip("no grimoire.toml")
	}

	fmt.Printf("\nValidating package at %s\n\n", target)
	for _, c := range checks {
		icon := tui.IconOK
		switch c.status {
		case "warn":
			icon = tui.IconWarn
		case "error":
			icon = tui.IconFail
		case "skip":
			icon = tui.IconSkip
		}
		fmt.Printf("  %s  %s\n", icon, c.detail)
	}
	fmt.Println()
	if allOK {
		fmt.Printf("  %s  Package structure is valid — ready to publish.\n\n", tui.IconOK)
		return nil
	}
	fmt.Printf("  Package has issues — fix them before publishing.\n\n")
	return fmt.Errorf("package validation failed")
}
