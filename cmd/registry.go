package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	gitops "github.com/jeffreytse/grimoire/internal/git"
	"github.com/jeffreytse/grimoire/internal/settings"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var registryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Manage the grimoire skill registry",
	Long: `View and manage the grimoire skill registry (official source or company mirror).

  grimoire registry list               show current registry and extends targets
  grimoire registry add <ref>          add a registry to standards.extends and clone it
  grimoire registry remove <name>      remove a registry from standards.extends
  grimoire registry set <ref>          set core.registry (owner/repo[@version] or URL)
  grimoire registry reset              clear core.registry (revert to built-in default)
  grimoire registry update [<name>]    pull latest from registry and extends targets
  grimoire registry validate [<path>]  validate a registry's structure before publishing`,
}

var (
	flagRegistryListJSON  bool
	flagRegistryUpdateRef string
)

var registryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List current registry and extends targets",
	RunE:  runRegistryList,
}

var registrySetCmd = &cobra.Command{
	Use:   "set <ref>",
	Short: "Set the core registry (owner/repo[@version] or full URL)",
	Args:  cobra.ExactArgs(1),
	RunE:  runRegistrySet,
}

var registryResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Clear core.registry (revert to built-in default)",
	Args:  cobra.NoArgs,
	RunE:  runRegistryReset,
}

var registryAddCmd = &cobra.Command{
	Use:   "add <ref>",
	Short: "Add a registry to standards.extends and clone it",
	Args:  cobra.ExactArgs(1),
	RunE:  runRegistryAdd,
}

var registryRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a registry from standards.extends",
	Args:  cobra.ExactArgs(1),
	RunE:  runRegistryRemove,
}

var registryUpdateCmd = &cobra.Command{
	Use:   "update [<name>]",
	Short: "Pull latest from registry and extends targets",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runRegistryUpdate,
}

var registryValidateCmd = &cobra.Command{
	Use:   "validate [<path-or-name>]",
	Short: "Validate a registry's structure (for registry authors)",
	Long: `Validate a registry's directory structure before publishing.

Without arguments, validates the current directory (useful when authoring a registry).
With an argument, validates an installed registry by name or an absolute path.

Checks:
  - Registry markers present (skills/, profiles/, presets/, settings.toml)
  - skills/: domain directories with SKILL.md files
  - profiles/: .toml files are valid TOML
  - presets/: each preset directory has settings.toml
  - settings.toml: valid TOML if present`,
	Args: cobra.MaximumNArgs(1),
	RunE: runRegistryValidate,
}

func init() {
	registryListCmd.Flags().BoolVar(&flagRegistryListJSON, "json", false, "output as JSON")
	registryUpdateCmd.Flags().StringVar(&flagRegistryUpdateRef, "ref", "", "change pinned version/branch/commit")
	registryCmd.AddCommand(registryListCmd)
	registryCmd.AddCommand(registryAddCmd)
	registryCmd.AddCommand(registryRemoveCmd)
	registryCmd.AddCommand(registrySetCmd)
	registryCmd.AddCommand(registryResetCmd)
	registryCmd.AddCommand(registryUpdateCmd)
	registryCmd.AddCommand(registryValidateCmd)
}

type registryListEntry struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Version     string `json:"version,omitempty"`
	SkillsCount int    `json:"skills_count"`
	Cloned      bool   `json:"cloned"`
	Kind        string `json:"kind"` // "core" | "extends"
}

func runRegistryList(cmd *cobra.Command, args []string) error {
	fs, err := settings.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}
	r, _ := settings.Load(getProjectDir())

	var entries []registryListEntry

	// Core registry entry
	coreURL := skills.GrimoireRepoURL()
	coreVersion := skills.GrimoireVersion()
	if fs.Core.Registry != "" {
		_, v := settings.ParseRef(fs.Core.Registry)
		if v != "" {
			coreVersion = v
		}
	}
	officialHome := skills.OfficialRegistryHome()
	coreKind := "core"
	if filepath.IsAbs(coreURL) {
		coreKind = "local"
	}
	entries = append(entries, registryListEntry{
		Name:        settings.DeriveRegistryName(coreURL),
		URL:         coreURL,
		Version:     coreVersion,
		SkillsCount: countSkills(skills.SkillsRoot()),
		Cloned:      dirExists(officialHome),
		Kind:        coreKind,
	})

	// Extends targets from resolved settings
	for _, ref := range r.StandardsExtends {
		u, ver := settings.ParseRef(ref)
		name := settings.DeriveRegistryName(u)
		extHome := skills.ExtendsHome(name)
		kind := "extends"
		if filepath.IsAbs(u) {
			kind = "local"
		}
		entries = append(entries, registryListEntry{
			Name:        name,
			URL:         u,
			Version:     ver,
			SkillsCount: countSkills(filepath.Join(extHome, "skills")),
			Cloned:      dirExists(extHome),
			Kind:        kind,
		})
	}

	if flagRegistryListJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(entries)
	}

	for _, e := range entries {
		icon := tui.IconOK
		if !e.Cloned {
			icon = tui.IconWarn
		}
		ref := e.URL
		if e.Version != "" {
			ref += "@" + e.Version
		}
		tag := tui.StyleDim.Render("[" + e.Kind + "]")
		fmt.Printf("  %s  %-30s %s %s\n", icon, e.Name, tag, tui.StyleDim.Render(ref))
		if e.Cloned {
			fmt.Printf("         %d skills\n\n", e.SkillsCount)
		} else {
			fmt.Printf("         not cloned — run: grimoire registry update %s\n\n", e.Name)
		}
	}
	return nil
}

func runRegistrySet(cmd *cobra.Command, args []string) error {
	ref := args[0]
	u, _ := settings.ParseRef(ref)
	if !skills.IsGitURL(u) && !filepath.IsAbs(u) {
		return fmt.Errorf("invalid ref %q — expected owner/repo[@version], git URL, or absolute path", ref)
	}
	if filepath.IsAbs(u) {
		if _, err := os.Stat(u); err != nil {
			return fmt.Errorf("local path %q not found", u)
		}
	}

	fs, err := settings.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}
	fs.Core.Registry = ref
	if err := settings.SaveGlobal(fs); err != nil {
		return fmt.Errorf("saving settings: %w", err)
	}
	fmt.Printf("%s  core.registry = %s\n", tui.IconOK, ref)
	if filepath.IsAbs(u) {
		fmt.Printf("   local registry set as official — grimoire install will symlink from %s\n", u)
	} else {
		fmt.Printf("   run: grimoire update  to apply\n")
	}
	return nil
}

func runRegistryReset(cmd *cobra.Command, args []string) error {
	fs, err := settings.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}
	fs.Core.Registry = ""
	if err := settings.SaveGlobal(fs); err != nil {
		return fmt.Errorf("saving settings: %w", err)
	}
	fmt.Printf("%s  core.registry cleared — using built-in default (%s)\n", tui.IconOK, skills.GrimoireRepo)
	return nil
}

func runRegistryUpdate(cmd *cobra.Command, args []string) error {
	r, err := settings.Load(getProjectDir())
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}

	if len(args) == 1 {
		name := args[0]
		if name == "official" {
			return updateCoreRegistry()
		}
		return updateExtendsTarget(name, r)
	}

	// update all
	if err := updateCoreRegistry(); err != nil {
		fmt.Fprintf(os.Stderr, "  warn: official: %v\n", err)
	}
	for _, ref := range r.StandardsExtends {
		u, _ := settings.ParseRef(ref)
		name := settings.DeriveRegistryName(u)
		if filepath.IsAbs(name) {
			continue // local registry — user manages the checkout
		}
		if err := updateExtendsTarget(name, r); err != nil {
			fmt.Fprintf(os.Stderr, "  warn: %s: %v\n", name, err)
		}
	}
	return nil
}

func updateCoreRegistry() error {
	dest := skills.OfficialRegistryHome()
	url := skills.GrimoireRepoURL()

	// Local path: verify it exists, no git ops
	if filepath.IsAbs(url) {
		if _, err := os.Stat(dest); err != nil {
			return fmt.Errorf("local registry %q not found", dest)
		}
		fmt.Printf("  %s  official registry is a local path — no update needed\n", tui.IconOK)
		return nil
	}

	if !dirExists(dest) {
		fmt.Printf("  cloning official registry...\n")
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return fmt.Errorf("creating dir: %w", err)
		}
		if err := gitops.Clone(url, dest); err != nil {
			return fmt.Errorf("cloning: %w", err)
		}
		fmt.Printf("  %s  official cloned\n", tui.IconOK)
		return nil
	}

	fs, _ := settings.LoadGlobal()
	_, ver := settings.ParseRef(fs.Core.Registry)
	if flagRegistryUpdateRef != "" {
		ver = flagRegistryUpdateRef
		fs.Core.Registry = url + "@" + ver
		_ = settings.SaveGlobal(fs)
	}

	fmt.Printf("  updating official registry...\n")
	if ver != "" {
		if err := gitops.CheckoutVersion(dest, ver); err != nil {
			return fmt.Errorf("checking out version: %w", err)
		}
	} else {
		if err := gitops.Pull(dest); err != nil {
			return fmt.Errorf("pulling: %w", err)
		}
	}
	fmt.Printf("  %s  official up to date\n", tui.IconOK)
	return nil
}

func updateExtendsTarget(name string, r settings.Resolved) error {
	if filepath.IsAbs(name) {
		if _, err := os.Stat(name); err != nil {
			return fmt.Errorf("local registry %q not found", name)
		}
		fmt.Printf("  %s  %s is a local path — no update needed\n", tui.IconOK, name)
		return nil
	}

	var extURL, extVer string
	for _, ref := range r.StandardsExtends {
		u, ver := settings.ParseRef(ref)
		if settings.DeriveRegistryName(u) == name {
			extURL, extVer = u, ver
			break
		}
	}
	if extURL == "" {
		return fmt.Errorf("extends target %q not declared in any settings file", name)
	}

	dest := skills.ExtendsHome(name)
	if !dirExists(dest) {
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return fmt.Errorf("creating dir: %w", err)
		}
		fmt.Printf("  cloning %s...\n", name)
		if err := gitops.Clone(extURL, dest); err != nil {
			return fmt.Errorf("cloning: %w", err)
		}
		if err := gitops.CheckoutVersion(dest, extVer); err != nil {
			return fmt.Errorf("checkout: %w", err)
		}
		fmt.Printf("  %s  %s cloned\n", tui.IconOK, name)
		return nil
	}

	fmt.Printf("  updating %s...\n", name)
	if flagRegistryUpdateRef != "" {
		extVer = flagRegistryUpdateRef
	}
	if extVer != "" {
		if err := gitops.CheckoutVersion(dest, extVer); err != nil {
			return fmt.Errorf("checking out version: %w", err)
		}
	} else {
		if err := gitops.Pull(dest); err != nil {
			return fmt.Errorf("pulling: %w", err)
		}
	}
	fmt.Printf("  %s  %s up to date\n", tui.IconOK, name)
	return nil
}

func runRegistryAdd(cmd *cobra.Command, args []string) error {
	ref := args[0]

	if filepath.IsAbs(ref) {
		return runRegistryAddLocal(ref)
	}

	u, _ := settings.ParseRef(ref)
	if !skills.IsGitURL(u) {
		return fmt.Errorf("invalid ref %q — expected owner/repo[@version], git URL, or absolute path", ref)
	}
	name := settings.DeriveRegistryName(u)

	fs, err := settings.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}

	// idempotent: skip if already present
	for _, existing := range fs.StandardsExtends {
		eu, _ := settings.ParseRef(existing)
		if settings.DeriveRegistryName(eu) == name {
			fmt.Printf("%s  %s already in standards.extends\n", tui.IconOK, name)
			r, _ := settings.Load(getProjectDir())
			return updateExtendsTarget(name, r)
		}
	}

	fs.StandardsExtends = append(fs.StandardsExtends, ref)
	if err := settings.SaveGlobal(fs); err != nil {
		return fmt.Errorf("saving settings: %w", err)
	}
	fmt.Printf("%s  added %s to standards.extends\n", tui.IconOK, name)

	r, _ := settings.Load(getProjectDir())
	if err := updateExtendsTarget(name, r); err != nil {
		return err
	}

	extHome := skills.ExtendsHome(name)
	if sc := countSkills(filepath.Join(extHome, "skills")); sc > 0 {
		fmt.Printf("   %d skills available from %s\n", sc, name)
	}
	if pc := countProfiles(filepath.Join(extHome, "profiles")); pc > 0 {
		fmt.Printf("   %d profiles available from %s\n", pc, name)
	}
	return nil
}

func runRegistryAddLocal(path string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("local path %q not found", path)
	}

	fs, err := settings.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}
	for _, existing := range fs.StandardsExtends {
		eu, _ := settings.ParseRef(existing)
		if eu == path {
			fmt.Printf("%s  %s already in standards.extends\n", tui.IconOK, path)
			return nil
		}
	}

	fs.StandardsExtends = append(fs.StandardsExtends, path)
	if err := settings.SaveGlobal(fs); err != nil {
		return fmt.Errorf("saving settings: %w", err)
	}
	fmt.Printf("%s  added %s to standards.extends\n", tui.IconOK, path)

	if sc := countSkills(filepath.Join(path, "skills")); sc > 0 {
		fmt.Printf("   %d skills available\n", sc)
	}
	if pc := countProfiles(filepath.Join(path, "profiles")); pc > 0 {
		fmt.Printf("   %d profiles available\n", pc)
	}
	return nil
}

func runRegistryRemove(cmd *cobra.Command, args []string) error {
	target := args[0]

	fs, err := settings.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}

	var kept []string
	removed := false
	for _, existing := range fs.StandardsExtends {
		eu, _ := settings.ParseRef(existing)
		if settings.DeriveRegistryName(eu) == target {
			removed = true
			continue
		}
		kept = append(kept, existing)
	}
	if !removed {
		return fmt.Errorf("registry %q not found in standards.extends", target)
	}

	fs.StandardsExtends = kept
	if err := settings.SaveGlobal(fs); err != nil {
		return fmt.Errorf("saving settings: %w", err)
	}
	fmt.Printf("%s  removed %s from standards.extends\n", tui.IconOK, target)
	fmt.Printf("   local clone at %s preserved — delete manually if no longer needed\n",
		skills.ExtendsHome(target))
	return nil
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

func runRegistryValidate(cmd *cobra.Command, args []string) error {
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
			// Treat as installed registry name or shorthand ref
			u, _ := settings.ParseRef(ref)
			name := settings.DeriveRegistryName(u)
			home := skills.RegistryHome(name)
			if !dirExists(home) {
				return fmt.Errorf("registry %q not installed — run: grimoire registry update %s", ref, ref)
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

	// Registry markers
	hasMarker := false
	for _, marker := range []string{"skills", "profiles", "presets", "settings.toml"} {
		if _, err := os.Stat(filepath.Join(target, marker)); err == nil {
			hasMarker = true
			break
		}
	}
	if !hasMarker {
		fail("no registry markers found (expected: skills/, profiles/, presets/, or settings.toml)")
	} else {
		pass("registry structure detected")
	}

	// skills/ structure
	skillsDir := filepath.Join(target, "skills")
	if _, err := os.Stat(skillsDir); err == nil {
		allSkills, _ := skills.ListAllSkills(skillsDir)
		if len(allSkills) == 0 {
			warn("skills/ found but no skills detected (expected: skills/<domain>/<name>/SKILL.md)")
		} else {
			missing := 0
			for _, sk := range allSkills {
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
			if _, err := settings.ParseFile(filepath.Join(profilesDir, e.Name())); err != nil {
				invalid++
			}
		}
		if total == 0 {
			warn("profiles/ found but empty (expected: profiles/<name>.toml)")
		} else if invalid > 0 {
			fail(fmt.Sprintf("%d/%d profile TOML file(s) failed to parse", invalid, total))
		} else {
			pass(fmt.Sprintf("%d profile(s), all valid TOML", total))
		}
	} else {
		skip("no profiles/ directory")
	}

	// presets/ structure
	presetsDir := filepath.Join(target, "presets")
	if _, err := os.Stat(presetsDir); err == nil {
		presets := skills.ListPresets(target)
		if len(presets) == 0 {
			warn("presets/ found but no presets detected (expected: presets/<name>/settings.toml)")
		} else {
			invalid := 0
			for _, p := range presets {
				if _, err := settings.ParseFile(filepath.Join(presetsDir, p, "settings.toml")); err != nil {
					invalid++
				}
			}
			if invalid > 0 {
				fail(fmt.Sprintf("%d/%d preset settings.toml file(s) failed to parse", invalid, len(presets)))
			} else {
				pass(fmt.Sprintf("%d preset(s), all have valid settings.toml", len(presets)))
			}
		}
	} else {
		skip("no presets/ directory")
	}

	// root settings.toml
	rootSettings := filepath.Join(target, "settings.toml")
	if _, err := os.Stat(rootSettings); err == nil {
		if _, err := settings.ParseFile(rootSettings); err != nil {
			fail(fmt.Sprintf("settings.toml parse error: %v", err))
		} else {
			pass("settings.toml is valid TOML")
		}
	} else {
		skip("no settings.toml")
	}

	fmt.Printf("\nValidating registry at %s\n\n", target)
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
		fmt.Printf("  %s  Registry structure is valid — ready to publish.\n\n", tui.IconOK)
		return nil
	}
	fmt.Printf("  Registry has issues — fix them before publishing.\n\n")
	return fmt.Errorf("registry validation failed")
}
