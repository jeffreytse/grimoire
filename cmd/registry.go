package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"

	gitops "github.com/jeffreytse/grimoire/internal/git"
	"github.com/jeffreytse/grimoire/internal/settings"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var registryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Manage grimoire skill registries",
	Long: `Manage grimoire skill registries — official and user-defined.

Multiple registries are searched in priority order (highest first).
User registries do not need to follow the official STANDARD.md.

  grimoire registry list                        list all registries
  grimoire registry add <name> <url>            add a registry and clone it
  grimoire registry remove <name>               remove a registry
  grimoire registry enable <name>               enable a disabled registry
  grimoire registry disable <name>              disable a registry without removing it
  grimoire registry set <ref>                   set the official registry URL
  grimoire registry reset                       revert official registry to built-in default
  grimoire registry update [<name>]             pull latest from all (or named) registries
  grimoire registry validate [<path-or-name>]   validate a registry's structure`,
}

var (
	flagRegistryListJSON    bool
	flagRegistryUpdateRef   string
	flagRegistryAddPriority int
)

var registryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured registries",
	RunE:  runRegistryList,
}

var registrySetCmd = &cobra.Command{
	Use:   "set <ref>",
	Short: "Set the official registry URL (owner/repo[@version] or full URL)",
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
	Use:   "add <name> <url>",
	Short: "Add a named registry and clone it",
	Long: `Add a named registry to the [[registry]] list and clone it locally.

<name>  short identifier for this registry (e.g. "my-team", "plugins")
<url>   git URL, owner/repo[@version] shorthand, or absolute local path

User registries do not need to follow the official STANDARD.md — any
directory with a skills/ folder (or flat .md files) is accepted.`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runRegistryAdd,
}

var registryRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a registry by name",
	Args:  cobra.ExactArgs(1),
	RunE:  runRegistryRemove,
}

var registryEnableCmd = &cobra.Command{
	Use:   "enable <name>",
	Short: "Enable a previously disabled registry",
	Args:  cobra.ExactArgs(1),
	RunE:  runRegistryEnable,
}

var registryDisableCmd = &cobra.Command{
	Use:   "disable <name>",
	Short: "Disable a registry without removing it",
	Args:  cobra.ExactArgs(1),
	RunE:  runRegistryDisable,
}

var registryUpdateCmd = &cobra.Command{
	Use:   "update [<name>]",
	Short: "Pull latest from all (or one named) registry",
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
	registryAddCmd.Flags().IntVar(&flagRegistryAddPriority, "priority", 0, "skill resolution priority (higher wins; default 50 for user registries)")
	registryCmd.AddCommand(registryListCmd)
	registryCmd.AddCommand(registryAddCmd)
	registryCmd.AddCommand(registryRemoveCmd)
	registryCmd.AddCommand(registryEnableCmd)
	registryCmd.AddCommand(registryDisableCmd)
	registryCmd.AddCommand(registrySetCmd)
	registryCmd.AddCommand(registryResetCmd)
	registryCmd.AddCommand(registryUpdateCmd)
	registryCmd.AddCommand(registryValidateCmd)
}

type registryListEntry struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Version     string `json:"version,omitempty"`
	Priority    int    `json:"priority"`
	SkillsCount   int    `json:"skills_count"`
	ProfilesCount int    `json:"profiles_count"`
	PresetsCount  int    `json:"presets_count"`
	Cloned        bool   `json:"cloned"`
	Enabled       bool   `json:"enabled"`
	Kind          string `json:"kind"` // "official" | "user" | "local"
}

func runRegistryList(cmd *cobra.Command, args []string) error {
	cfg, err := settings.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}

	regs := skills.AllRegistries()
	entries := make([]registryListEntry, 0, len(regs))
	officialSeen := false

	for _, reg := range regs {
		var url, ver string
		enabled := true
		for _, rd := range cfg.Registries {
			if rd.Name == reg.Name {
				url, ver = settings.ParseRef(rd.URL)
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
		entries = append(entries, registryListEntry{
			Name:          reg.Name,
			URL:           url,
			Version:       ver,
			Priority:      reg.Priority,
			SkillsCount:   countSkills(filepath.Join(reg.Home, "skills")),
			ProfilesCount: countProfiles(filepath.Join(reg.Home, "profiles")),
			PresetsCount:  len(skills.ListPresets(reg.Home)),
			Cloned:        dirExists(reg.Home),
			Enabled:       enabled,
			Kind:          kind,
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
			fmt.Printf("         disabled — run: grimoire registry enable %s\n\n", e.Name)
		case e.Cloned:
			fmt.Printf("         %d skills · %d profiles · %d presets\n\n", e.SkillsCount, e.ProfilesCount, e.PresetsCount)
		default:
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

	cfg, err := settings.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}
	// Update or add the official=true [[registry]] entry.
	for i, rd := range cfg.Registries {
		if !rd.Official {
			continue
		}
		cfg.Registries[i].URL = ref
		if err := settings.SaveGlobal(cfg); err != nil {
			return fmt.Errorf("saving settings: %w", err)
		}
		fmt.Printf("%s  official registry = %s\n", tui.IconOK, ref)
		if filepath.IsAbs(u) {
			fmt.Printf("   local registry set as official\n")
		} else {
			fmt.Printf("   run: grimoire registry update official  to apply\n")
		}
		return nil
	}
	cfg.Registries = append(cfg.Registries, settings.RegistryDef{
		Name:     "official",
		URL:      ref,
		Official: true,
		Priority: 100,
		Enabled:  true,
	})
	if err := settings.SaveGlobal(cfg); err != nil {
		return fmt.Errorf("saving settings: %w", err)
	}
	fmt.Printf("%s  official registry = %s\n", tui.IconOK, ref)
	if filepath.IsAbs(u) {
		fmt.Printf("   local registry set as official\n")
	} else {
		fmt.Printf("   run: grimoire registry update official  to apply\n")
	}
	return nil
}

func runRegistryReset(cmd *cobra.Command, args []string) error {
	cfg, err := settings.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}
	// Remove the official=true [[registry]] entry so GrimoireRepo constant takes effect.
	var kept []settings.RegistryDef
	for _, rd := range cfg.Registries {
		if !rd.Official {
			kept = append(kept, rd)
		}
	}
	cfg.Registries = kept
	if err := settings.SaveGlobal(cfg); err != nil {
		return fmt.Errorf("saving settings: %w", err)
	}
	fmt.Printf("%s  official registry reset — using built-in default (%s)\n", tui.IconOK, skills.GrimoireRepo)
	return nil
}

func runRegistryEnable(cmd *cobra.Command, args []string) error {
	return setRegistryEnabled(args[0], true)
}

func runRegistryDisable(cmd *cobra.Command, args []string) error {
	return setRegistryEnabled(args[0], false)
}

func setRegistryEnabled(name string, enabled bool) error {
	cfg, err := settings.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}
	for i, rd := range cfg.Registries {
		if rd.Name != name {
			continue
		}
		cfg.Registries[i].Enabled = enabled
		if err := settings.SaveGlobal(cfg); err != nil {
			return fmt.Errorf("saving settings: %w", err)
		}
		verb := "enabled"
		if !enabled {
			verb = "disabled"
		}
		fmt.Printf("%s  %s %s\n", tui.IconOK, verb, name)
		return nil
	}
	return fmt.Errorf("registry %q not found in [[registry]] — check: grimoire registry list", name)
}

// updateNamedRegistry clones or pulls a named [[registry]] entry.
func updateNamedRegistry(name, ref, forceVer string, w io.Writer) error {
	u, ver := settings.ParseRef(ref)
	if u == "" {
		u = ref
	}
	if forceVer != "" {
		ver = forceVer
	}
	if flagRegistryUpdateRef != "" {
		ver = flagRegistryUpdateRef
	}

	if filepath.IsAbs(u) {
		if _, err := os.Stat(u); err != nil {
			return fmt.Errorf("local registry %q not found", u)
		}
		fmt.Fprintf(w, "  %s  %s is a local path — no update needed\n", tui.IconOK, name)
		return nil
	}

	dest := filepath.Join(skills.RegistriesRoot(), name)
	if !dirExists(dest) {
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return fmt.Errorf("creating dir: %w", err)
		}
		if err := gitops.Clone(u, dest); err != nil {
			return fmt.Errorf("cloning %s: %w", name, err)
		}
		if ver != "" {
			if err := gitops.CheckoutVersion(dest, ver); err != nil {
				return fmt.Errorf("checkout %s: %w", name, err)
			}
		}
		fmt.Fprintf(w, "  %s  %s cloned\n", tui.IconOK, name)
		return nil
	}

	oldState, _ := gitops.CurrentState(dest)
	if ver != "" {
		if err := gitops.CheckoutVersion(dest, ver); err != nil {
			return fmt.Errorf("checking out version for %s: %w", name, err)
		}
	} else {
		if err := gitops.PullWithForceFallback(dest); err != nil {
			return fmt.Errorf("updating %s: %w", name, err)
		}
	}
	fmt.Fprintf(w, "  %s  %s up to date\n", tui.IconOK, name)
	if changes, err := gitops.RegistryChangesSince(dest, oldState.Commit); err == nil {
		printRegistryChanges(changes, dest, oldState.Commit, w)
	}
	return nil
}

func runRegistryUpdate(cmd *cobra.Command, args []string) error {
	cfg, err := settings.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}
	r, _ := settings.Load(getProjectDir())

	if len(args) == 1 {
		name := args[0]
		board := tui.NewStatusBoard([]string{name})
		stopSpinner := board.StartSpinner()
		board.SetUpdating(0)
		buf := &bytes.Buffer{}
		var updateErr error
		if name == skills.OfficialRegistryName || isOfficialByName(name) {
			updateErr = updateCoreRegistry(buf)
		} else if ref := findRegistryRef(name, &cfg); ref != "" {
			updateErr = updateNamedRegistry(name, ref, "", buf)
		} else {
			updateErr = fmt.Errorf("registry %q not found — check: grimoire registry list", name)
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

	// Build work list: new-model registries + old-model extends targets.
	type workItem struct {
		name   string
		ref    string
		isCore bool
	}
	var items []workItem
	var pinnedSkipped []string

	if len(cfg.Registries) > 0 {
		for _, rd := range cfg.Registries {
			if !rd.Enabled {
				continue
			}
			u, ver := settings.ParseRef(rd.URL)
			if u == "" {
				u = rd.URL
			}
			if filepath.IsAbs(u) {
				continue // local registry — user manages checkout
			}
			if ver != "" {
				pinnedSkipped = append(pinnedSkipped, rd.Name)
				continue // pinned registry — immutable by intent; skip bulk update
			}
			items = append(items, workItem{
				name:   rd.Name,
				ref:    rd.URL,
				isCore: rd.Official,
			})
		}
	} else {
		// No [[registry]] configured — update the implicit official registry.
		items = append(items, workItem{name: skills.OfficialRegistryName, isCore: true})
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
				err = updateCoreRegistry(buf)
			default:
				err = updateNamedRegistry(item.name, item.ref, "", buf)
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
		fmt.Printf("       to force-update: grimoire registry update <name>\n")
	}

	for i, res := range results {
		if res.err != nil {
			fmt.Fprintf(os.Stderr, "  warn: %s: %v\n", items[i].name, res.err)
		}
		os.Stdout.Write(res.buf.Bytes())
	}
	return nil
}

// isOfficialByName checks if a named registry is the (demote-resolved) official registry.
// Uses AllRegistries() so demoted lower-priority official entries return false.
func isOfficialByName(name string) bool {
	for _, reg := range skills.AllRegistries() {
		if reg.Official && reg.Name == name {
			return true
		}
	}
	return false
}

// findRegistryRef returns the URL ref for a named registry in [[registry]], or "".
func findRegistryRef(name string, cfg *settings.FileSettings) string {
	for _, rd := range cfg.Registries {
		if rd.Name == name {
			return rd.URL
		}
	}
	return ""
}

func updateCoreRegistry(w io.Writer) error {
	dest := skills.OfficialRegistryHome()
	url := skills.GrimoireRepoURL()

	// Local path: verify it exists, no git ops
	if filepath.IsAbs(url) {
		if _, err := os.Stat(dest); err != nil {
			return fmt.Errorf("local registry %q not found", dest)
		}
		fmt.Fprintf(w, "  %s  official registry is a local path — no update needed\n", tui.IconOK)
		return nil
	}

	if !dirExists(dest) {
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return fmt.Errorf("creating dir: %w", err)
		}
		if err := gitops.Clone(url, dest); err != nil {
			return fmt.Errorf("cloning: %w", err)
		}
		fmt.Fprintf(w, "  %s  official cloned\n", tui.IconOK)
		return nil
	}

	ver := ""
	cfg, _ := settings.LoadGlobal()
	for _, rd := range cfg.Registries {
		if rd.Official {
			_, ver = settings.ParseRef(rd.URL)
			break
		}
	}
	if flagRegistryUpdateRef != "" {
		ver = flagRegistryUpdateRef
		for i, rd := range cfg.Registries {
			if rd.Official {
				cfg.Registries[i].URL = url + "@" + ver
				_ = settings.SaveGlobal(cfg)
				break
			}
		}
	}

	oldState, _ := gitops.CurrentState(dest)
	if ver != "" {
		if err := gitops.CheckoutVersion(dest, ver); err != nil {
			return fmt.Errorf("checking out version: %w", err)
		}
	} else {
		if err := gitops.PullWithForceFallback(dest); err != nil {
			return fmt.Errorf("updating: %w", err)
		}
	}
	fmt.Fprintf(w, "  %s  official up to date\n", tui.IconOK)
	if changes, err := gitops.RegistryChangesSince(dest, oldState.Commit); err == nil {
		printRegistryChanges(changes, dest, oldState.Commit, w)
	}
	return nil
}

func runRegistryAdd(cmd *cobra.Command, args []string) error {
	if len(args) == 1 {
		return fmt.Errorf("usage: grimoire registry add <name> <url>\n\nExample:\n  grimoire registry add my-team https://github.com/acme/grimoire.git")
	}

	name := args[0]
	ref := args[1]

	if strings.ContainsAny(name, "/\\") {
		return fmt.Errorf("registry name %q must not contain path separators — use a short identifier like %q", name, strings.ReplaceAll(name, "/", "-"))
	}

	u, _ := settings.ParseRef(ref)
	if u == "" {
		u = ref
	}
	if !skills.IsGitURL(u) && !filepath.IsAbs(u) {
		return fmt.Errorf("invalid url %q — expected git URL, owner/repo[@version], or absolute path", ref)
	}
	if filepath.IsAbs(u) {
		if _, err := os.Stat(u); err != nil {
			return fmt.Errorf("local path %q not found", u)
		}
	}

	cfg, err := settings.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}

	// Idempotent: if name already exists, update URL/priority.
	for i, existing := range cfg.Registries {
		if existing.Name != name {
			continue
		}
		cfg.Registries[i].URL = ref
		if flagRegistryAddPriority > 0 {
			cfg.Registries[i].Priority = flagRegistryAddPriority
		}
		if err := settings.SaveGlobal(cfg); err != nil {
			return fmt.Errorf("saving settings: %w", err)
		}
		fmt.Printf("%s  updated registry %s → %s\n", tui.IconOK, name, u)
		if filepath.IsAbs(u) {
			return nil
		}
		return updateNamedRegistry(name, ref, "", os.Stdout)
	}

	rd := settings.RegistryDef{
		Name:    name,
		URL:     ref,
		Enabled: true,
	}
	if flagRegistryAddPriority > 0 {
		rd.Priority = flagRegistryAddPriority
	}
	cfg.Registries = append(cfg.Registries, rd)
	if err := settings.SaveGlobal(cfg); err != nil {
		return fmt.Errorf("saving settings: %w", err)
	}
	fmt.Printf("%s  added registry %s\n", tui.IconOK, name)

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

	if err := updateNamedRegistry(name, ref, "", os.Stdout); err != nil {
		return err
	}
	home := filepath.Join(skills.RegistriesRoot(), name)
	if sc := countSkills(filepath.Join(home, "skills")); sc > 0 {
		fmt.Printf("   %d skills available from %s\n", sc, name)
	}
	if pc := countProfiles(filepath.Join(home, "profiles")); pc > 0 {
		fmt.Printf("   %d profiles available from %s\n", pc, name)
	}
	return nil
}

func runRegistryRemove(cmd *cobra.Command, args []string) error {
	target := args[0]

	cfg, err := settings.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}

	var kept []settings.RegistryDef
	removed := false
	for _, rd := range cfg.Registries {
		if rd.Name == target {
			removed = true
			continue
		}
		kept = append(kept, rd)
	}
	if !removed {
		return fmt.Errorf("registry %q not found — check: grimoire registry list", target)
	}
	cfg.Registries = kept
	if err := settings.SaveGlobal(cfg); err != nil {
		return fmt.Errorf("saving settings: %w", err)
	}
	home := filepath.Join(skills.RegistriesRoot(), target)
	fmt.Printf("%s  removed %s from [[registry]]\n", tui.IconOK, target)
	fmt.Printf("   local clone at %s preserved — delete manually if no longer needed\n", home)
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
