package cmd

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/agent"
	"github.com/jeffreytse/grimoire/internal/config"
	gitops "github.com/jeffreytse/grimoire/internal/git"
	"github.com/jeffreytse/grimoire/internal/lock"
	"github.com/jeffreytse/grimoire/internal/manifest"
	"github.com/jeffreytse/grimoire/internal/profiles"
	"github.com/jeffreytse/grimoire/internal/resolver"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var (
	flagInstallDomain    string
	flagInstallSubdomain string
	flagInstallSkill     string
	flagInstallTarget    string
	flagInstallCopy      bool
	flagInstallYes       bool
	flagInstallNoCfg     bool
	flagInstallFrom      string
	flagInstallPackage   string
	flagInstallScope     string // "project" (default) | "global"
	flagInstallGlobal    bool   // shorthand for --scope global
)

var installCmd = &cobra.Command{
	Use:   "install [<grimoire-ref>]",
	Short: "Install skills from grimoire.toml or a grimoire-ref",
	Long: `Install skills from grimoire.toml [dependencies] (default) or a grimoire-ref.

A grimoire-ref identifies skills to install:

  <grimoire-ref>  =  [host/][owner/repo[@version]][:glob_path]

Examples:
  grimoire install apply-solid                              # skill in official package
  grimoire install engineering/development                  # all skills under a path
  grimoire install acmecorp/practices                       # all skills from a package
  grimoire install acmecorp/practices:engineering/tdd       # one skill from a package
  grimoire install acmecorp/practices:engineering/**        # glob — multiple skills
  grimoire install github.com/acmecorp/practices@v2         # explicit host + version`,
	RunE: runInstall,
}

func init() {
	installCmd.Flags().StringVar(&flagInstallDomain, "domain", "", "install all skills for a domain")
	installCmd.Flags().StringVar(&flagInstallSubdomain, "subdomain", "", "restrict to one sub-domain")
	installCmd.Flags().StringVar(&flagInstallSkill, "skill", "", "install one skill (domain/subdomain/name or domain/name)")
	installCmd.Flags().StringVar(&flagInstallTarget, "target", "", "target agent: claude, codex, gemini, antigravity, openclaw, opencode, agent, all, auto")
	installCmd.Flags().BoolVar(&flagInstallCopy, "copy", false, "copy files instead of symlinking")
	installCmd.Flags().BoolVar(&flagInstallYes, "yes", false, "non-interactive: install all skills to all detected agents")
	installCmd.Flags().BoolVar(&flagInstallNoCfg, "no-configure", false, "skip writing start-best-practice trigger")
	installCmd.Flags().StringVar(&flagInstallFrom, "from", "", "install from a local path or git URL (persisted to ~/.config/grimoire/grimoire.toml)")
	installCmd.Flags().StringVar(&flagInstallPackage, "package", "", "install skills from a specific package only")
	installCmd.Flags().StringVar(&flagInstallScope, "scope", "project", `install scope: "project" (./.claude/skills/) or "global" (~/.claude/skills/)`)
	installCmd.Flags().BoolVar(&flagInstallGlobal, "global", false, `global-scoped install — shorthand for --scope global`)
}

// installSkillsDir returns the target skills directory for an agent based on the current scope flags.
func installSkillsDir(ag string) string {
	if flagInstallGlobal || flagInstallScope == "global" {
		return agent.SkillsDir(ag)
	}
	return agent.ProjectSkillsDir(ag, getProjectDir())
}

// depKeyForSkill returns the dep key to write to grimoire.toml.
// Official-package skills use a bare key; third-party skills use "<pkgName>:<skillName>".
func depKeyForSkill(skillName, pkgName string) string {
	if pkgName == "" {
		return skillName
	}
	return pkgName + ":" + skillName
}

// ensureGlobalManifest returns the global grimoire.toml path, creating it (and its parent dir) if absent.
func ensureGlobalManifest() string {
	path := manifest.GlobalPath()
	if _, err := os.Stat(path); err != nil {
		_ = os.MkdirAll(filepath.Dir(path), 0o755)
		_, _ = scaffoldManifest(path, "global", "", 0, -1)
	}
	return path
}

func runInstall(cmd *cobra.Command, args []string) error {
	// grimoire install <ref> — dispatch based on ref type
	if len(args) > 0 {
		return runInstallArg(args[0])
	}

	// --from overrides package search entirely
	if flagInstallFrom != "" {
		resolved, err := resolveAndPersistSource(flagInstallFrom)
		if err != nil {
			return err
		}
		if resolved == "" {
			return nil
		}
		return runInstallFromRoot(resolved)
	}

	// grimoire.toml present → manifest-driven install path
	if flagInstallGlobal || flagInstallScope == "global" {
		if _, err := os.Stat(manifest.GlobalPath()); err == nil {
			return runInstallFromManifest(filepath.Dir(manifest.GlobalPath()))
		}
	} else if _, err := os.Stat(manifest.ProjectPath(getProjectDir())); err == nil {
		return runInstallFromManifest(getProjectDir())
	}

	// Load settings once — used for dep-clone and install pool resolution.
	r, _ := config.Load(getProjectDir())

	// Warn about missing extends refs; auto-clone any declared [dependencies] skills packages.
	for _, ref := range r.MissingExtends {
		fmt.Fprintf(os.Stderr, "%s  standards.extends: package %q is not installed\n", tui.IconWarn, ref)
		fmt.Fprintf(os.Stderr, "   Install it: grimoire package add <name> <grimoire-ref>\n")
	}
	for _, dep := range r.DepSkills {
		pkgRef := config.ParsePackageRef(dep)
		if pkgRef.IsLocal() || pkgRef.IsOfficialRepoPath() || pkgRef.PackageName == "" {
			continue
		}
		dest := skills.PackageHome(pkgRef.PackageName)
		if dirExists(dest) {
			continue
		}
		fmt.Printf("  Cloning %s…\n", pkgRef.PackageName)
		buf := &bytes.Buffer{}
		if cloneErr := updateNamedPackage(pkgRef.PackageName, pkgRef.PackageURL, pkgRef.Tag, buf); cloneErr != nil {
			fmt.Fprintf(os.Stderr, "%s  dependency %q: %v\n", tui.IconWarn, dep, cloneErr)
		} else {
			os.Stdout.Write(buf.Bytes())
		}
	}

	// determine which packages to install from
	regs := skills.AllSkillsPackages()

	// --package filters to one package
	if flagInstallPackage != "" {
		regs = filterPackages(regs, flagInstallPackage)
		if len(regs) == 0 {
			return fmt.Errorf("package %q not found or not cloned — run: grimoire package update %s",
				flagInstallPackage, flagInstallPackage)
		}
	}

	if len(regs) == 0 {
		return fmt.Errorf("skills not found at %s — run: grimoire update", skills.SkillsRoot())
	}

	// parse optional package: prefix from --skill flag  (e.g. "my-org:engineering/tdd")
	skillRef := flagInstallSkill
	if regName, ref, ok := splitPackagePrefix(skillRef); ok {
		skillRef = ref
		regs = filterPackages(regs, regName)
		if len(regs) == 0 {
			return fmt.Errorf("package %q not found or not cloned", regName)
		}
	}

	symlink := !flagInstallCopy
	if !flagInstallCopy && r.Core.InstallMode == "copy" {
		symlink = false
	}
	if err := checkWindowsSymlinkSupport(symlink); err != nil {
		return err
	}
	target := flagInstallTarget
	if flagInstallYes && target == "" {
		target = "auto"
	}
	targets := resolveTargets(target)

	perAgent := make(map[string]int)

	switch {
	case skillRef != "":
		skillPath, src, err := resolveSkillFromPackages(regs, skillRef)
		if err != nil {
			return err
		}
		fmt.Printf("  Resolving skill: %s (from %s)…\n\n", skillRef, src.Name)
		var linkedTo []string
		for _, ag := range targets {
			n, err := installSkillToAgent(skillPath, ag, symlink)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  %s  %s — %s: %v\n", tui.IconWarn, skillRef, agent.DisplayName(ag), err)
				continue
			}
			perAgent[ag] += n
			if n > 0 {
				linkedTo = append(linkedTo, agent.DisplayName(ag))
			}
		}
		if len(linkedTo) > 0 {
			printSkillInstalled(skillRef, linkedTo)
		} else {
			fmt.Printf("  %s  %s already up to date\n", tui.IconSkip, skillRef)
		}

	case flagInstallDomain != "":
		fmt.Print("  Loading skills…")
		all, conflicts, err := skills.ListAllSkillsFromPackagesMeta(regs)
		if err != nil {
			return err
		}
		printConflicts(conflicts)
		var domainSkills []skills.Skill
		for i := range all {
			sk := all[i]
			if sk.Domain != flagInstallDomain {
				continue
			}
			if flagInstallSubdomain != "" && sk.Subdomain != flagInstallSubdomain {
				continue
			}
			domainSkills = append(domainSkills, sk)
		}
		if tui.IsTTY() {
			fmt.Printf("\r\033[K  Loading skills… %d found\n\n", len(domainSkills))
		} else {
			fmt.Printf(" %d found\n\n", len(domainSkills))
		}
		var skippedDomain int
		for i := range domainSkills {
			sk := domainSkills[i]
			var linkedTo []string
			for _, ag := range targets {
				n, installErr := installSkillToAgent(sk.Path, ag, symlink)
				if installErr != nil {
					fmt.Fprintf(os.Stderr, "  %s  %s — %s: %v\n", tui.IconWarn, sk.Name, agent.DisplayName(ag), installErr)
					continue
				}
				perAgent[ag] += n
				if n > 0 {
					linkedTo = append(linkedTo, agent.DisplayName(ag))
				}
			}
			if len(linkedTo) > 0 {
				printSkillInstalled(sk.Name, linkedTo)
			} else {
				skippedDomain++
			}
		}
		if skippedDomain > 0 {
			fmt.Printf("  %s  %d skills already up to date\n", tui.IconSkip, skippedDomain)
		}

	default:
		// No explicit target — require a grimoire.toml or declared deps.
		// Silently installing everything when nothing is declared is a footgun.
		if flagInstallPackage == "" && len(r.DepSkills) == 0 {
			fmt.Println("Nothing to install: no grimoire.toml found and no dependencies declared.")
			fmt.Println("  Create a project manifest:  grimoire init")
			fmt.Println("  Install a specific skill:   grimoire install <skill>")
			fmt.Println("  Install from a ref:         grimoire install <grimoire-ref>")
			return nil
		}
		// install all skills — use [dependencies] skills pool when declared, else all from --package
		fmt.Print("  Loading skills…")
		var all []skills.Skill
		var conflicts []skills.SkillConflict
		var err error
		if flagInstallPackage == "" && len(r.DepSkills) > 0 {
			all, conflicts, err = resolveInstallPool(&r)
		} else {
			all, conflicts, err = skills.ListAllSkillsFromPackagesMeta(regs)
		}
		if err != nil {
			return err
		}
		printConflicts(conflicts)
		if tui.IsTTY() {
			fmt.Printf("\r\033[K  Loading skills… %d found\n\n", len(all))
		} else {
			fmt.Printf(" %d found\n\n", len(all))
		}
		var skipped int
		for i := range all {
			sk := all[i]
			var linkedTo []string
			for _, ag := range targets {
				n, installErr := installSkillToAgent(sk.Path, ag, symlink)
				if installErr != nil {
					fmt.Fprintf(os.Stderr, "  %s  %s — %s: %v\n", tui.IconWarn, sk.Name, agent.DisplayName(ag), installErr)
					continue
				}
				perAgent[ag] += n
				if n > 0 {
					linkedTo = append(linkedTo, agent.DisplayName(ag))
				}
			}
			if len(linkedTo) > 0 {
				printSkillInstalled(sk.Name, linkedTo)
			} else {
				skipped++
			}
		}
		if skipped > 0 {
			fmt.Printf("  %s  %d skills already up to date\n", tui.IconSkip, skipped)
		}
	}

	// clean broken symlinks
	for _, ag := range targets {
		_, _ = skills.CleanBrokenSymlinks(installSkillsDir(ag))
	}

	// configure agent MD files
	if !flagInstallNoCfg {
		for _, ag := range targets {
			if err := agent.ConfigureAgentMD(ag); err != nil {
				fmt.Fprintf(os.Stderr, "  warn: configuring %s: %v\n", ag, err)
			}
		}
	}

	printInstallSummary(perAgent, targets)
	return nil
}

func printInstallSummary(perAgent map[string]int, targets []string) {
	var totalAdded int
	var agentCounts []string
	var upToDate []string
	for _, ag := range targets {
		n := perAgent[ag]
		totalAdded += n
		if n > 0 {
			agentCounts = append(agentCounts, fmt.Sprintf("%s (%d)", agent.DisplayName(ag), n))
		} else {
			upToDate = append(upToDate, agent.DisplayName(ag))
		}
	}

	fmt.Printf("\n%s  grimoire installed\n\n", tui.IconOK)

	if totalAdded > 0 {
		fmt.Printf("  %d skills added — %s\n", totalAdded, strings.Join(agentCounts, ", "))
	} else {
		fmt.Printf("  up to date — %s\n", strings.Join(upToDate, ", "))
	}

	fmt.Println("\n  start any AI session — grimoire skills activate automatically")
	fmt.Println("  or run /start-best-practice in Claude Code to trigger manually")
	fmt.Println("\n  uninstall: grimoire uninstall")
	fmt.Println()
}

// runInstallArg dispatches a `grimoire install <ref>` invocation based on the ref type.
// Any owner/repo ref (with or without path glob) → manifest-aware one-step install.
// Git URL or local path → legacy package URL install.
func runInstallArg(ref string) error {
	pkgRef := config.ParsePackageRef(ref)

	// Any owner/repo ref, or official package path → manifest-aware install
	if pkgRef.IsOfficialRepoPath() || pkgRef.Owner != "" {
		return runInstallSkillRef(ref, &pkgRef)
	}

	// Git URL or local path
	return runInstallFromPackageURL(ref)
}

// runInstallSkillRef installs skills by PackageRef and writes to grimoire.toml [dependencies].
// Three cases:
//   - pkgRef.Owner == "" : bare skill name or domain path — single skill from official package
//   - pkgRef.Owner != "" && pkgRef.Path != "" : package with path glob — install matching skills
//   - pkgRef.Owner != "" && pkgRef.Path == "" : package only — install ALL skills from package
func runInstallSkillRef(rawKey string, pkgRef *config.PackageRef) error {
	projectDir := getProjectDir()
	symlink := !flagInstallCopy
	if !flagInstallCopy {
		r, _ := config.Load(projectDir)
		if r.Core.InstallMode == "copy" {
			symlink = false
		}
	}
	if err := checkWindowsSymlinkSupport(symlink); err != nil {
		return err
	}
	targets := resolveTargets(flagInstallTarget)

	// Ensure non-official package is cloned
	if !pkgRef.IsOfficialRepoPath() && pkgRef.PackageURL != "" {
		dest := skills.PackageHome(pkgRef.PackageName)
		if !dirExists(dest) {
			buf := &bytes.Buffer{}
			if err := updateNamedPackage(pkgRef.PackageName, pkgRef.PackageURL, pkgRef.Tag, buf); err != nil {
				return fmt.Errorf("cloning package %s: %w", pkgRef.PackageName, err)
			}
			os.Stdout.Write(buf.Bytes())
		}
	}

	depKey, version := splitInstallRef(rawKey)
	if version == "" {
		version = "*"
	}

	regs := skills.AllSkillsPackages()
	perAgent := make(map[string]int)

	switch {
	case pkgRef.Owner == "":
		// Bare skill name or domain path — single skill from official package
		skillPath, src, err := resolveSkillFromPackages(regs, pkgRef.Path)
		if err != nil {
			return err
		}
		fmt.Printf("  Resolving skill: %s (from %s)…\n\n", pkgRef.Path, src.Name)
		var linkedTo []string
		for _, ag := range targets {
			n, instErr := installSkillToAgent(skillPath, ag, symlink)
			if instErr != nil {
				fmt.Fprintf(os.Stderr, "  %s  %s — %s: %v\n", tui.IconWarn, pkgRef.Path, agent.DisplayName(ag), instErr)
				continue
			}
			perAgent[ag] += n
			if n > 0 {
				linkedTo = append(linkedTo, agent.DisplayName(ag))
			}
		}
		if len(linkedTo) > 0 {
			printSkillInstalled(pkgRef.Path, linkedTo)
		} else {
			fmt.Printf("  %s  %s already up to date\n", tui.IconSkip, pkgRef.Path)
		}

	case pkgRef.Path != "":
		// Package ref with path glob — install all matching skills
		all, conflicts, err := skills.SkillsMatchingAnyMeta([]config.PackageRef{*pkgRef}, skills.AllPackages())
		if err != nil {
			return err
		}
		printConflicts(conflicts)
		fmt.Printf("  Loading skills… %d found\n\n", len(all))
		var skipped int
		for i := range all {
			sk := all[i]
			var linkedTo []string
			for _, ag := range targets {
				n, instErr := installSkillToAgent(sk.Path, ag, symlink)
				if instErr != nil {
					fmt.Fprintf(os.Stderr, "  %s  %s — %s: %v\n", tui.IconWarn, sk.Name, agent.DisplayName(ag), instErr)
					continue
				}
				perAgent[ag] += n
				if n > 0 {
					linkedTo = append(linkedTo, agent.DisplayName(ag))
				}
			}
			if len(linkedTo) > 0 {
				printSkillInstalled(sk.Name, linkedTo)
			} else {
				skipped++
			}
		}
		if skipped > 0 {
			fmt.Printf("  %s  %d skills already up to date\n", tui.IconSkip, skipped)
		}

	default:
		// Package ref with no path — install ALL skills from this package
		var pkgName string
		if pkgRef.IsOfficialRepoPath() {
			pkgName = skills.OfficialPackageDerivedName()
		} else {
			pkgName = pkgRef.PackageName
		}
		filteredRegs := filterPackages(regs, pkgName)
		if len(filteredRegs) == 0 {
			return fmt.Errorf("package %q not found or not cloned — run: grimoire update", pkgName)
		}
		fmt.Print("  Loading skills…")
		all, conflicts, err := skills.ListAllSkillsFromPackagesMeta(filteredRegs)
		if err != nil {
			return err
		}
		printConflicts(conflicts)
		if tui.IsTTY() {
			fmt.Printf("\r\033[K  Loading skills… %d found\n\n", len(all))
		} else {
			fmt.Printf(" %d found\n\n", len(all))
		}
		var skipped int
		for i := range all {
			sk := all[i]
			var linkedTo []string
			for _, ag := range targets {
				n, instErr := installSkillToAgent(sk.Path, ag, symlink)
				if instErr != nil {
					fmt.Fprintf(os.Stderr, "  %s  %s — %s: %v\n", tui.IconWarn, sk.Name, agent.DisplayName(ag), instErr)
					continue
				}
				perAgent[ag] += n
				if n > 0 {
					linkedTo = append(linkedTo, agent.DisplayName(ag))
				}
			}
			if len(linkedTo) > 0 {
				printSkillInstalled(sk.Name, linkedTo)
			} else {
				skipped++
			}
		}
		if skipped > 0 {
			fmt.Printf("  %s  %d skills already up to date\n", tui.IconSkip, skipped)
		}
	}

	// Save to grimoire.toml [dependencies] if manifest exists
	saveDepToManifest(projectDir, depKey, version)

	// Clean broken symlinks and configure agent MD files
	for _, ag := range targets {
		_, _ = skills.CleanBrokenSymlinks(installSkillsDir(ag))
	}
	if !flagInstallNoCfg {
		for _, ag := range targets {
			if err := agent.ConfigureAgentMD(ag); err != nil {
				fmt.Fprintf(os.Stderr, "  warn: configuring %s: %v\n", ag, err)
			}
		}
	}

	printInstallSummary(perAgent, targets)
	return nil
}

// saveDepToManifest appends depKey=version to [dependencies] in grimoire.toml if not already present.
func saveDepToManifest(projectDir, depKey, version string) {
	var manifestPath string
	if flagInstallGlobal || flagInstallScope == "global" {
		manifestPath = ensureGlobalManifest()
	} else {
		manifestPath = manifest.ProjectPath(projectDir)
	}
	if _, err := os.Stat(manifestPath); err != nil {
		return
	}
	mf, _ := manifest.ParseFile(manifestPath)
	if _, exists := mf.Deps[depKey]; exists {
		return
	}
	if err := manifest.AppendDep(manifestPath, depKey, version); err != nil {
		fmt.Fprintf(os.Stderr, "%s  saving to grimoire.toml: %v\n", tui.IconWarn, err)
		return
	}
	fmt.Printf("  %s  added to grimoire.toml: %s = %q\n", tui.IconOK, depKey, version)
}

// splitInstallRef splits a ref like "pkg:path@version" into (depKey, version).
// Searches for '@' only in the path portion (after ':') to avoid matching git@host SSH URLs.
func splitInstallRef(ref string) (depKey, version string) {
	colonIdx := strings.Index(ref, ":")
	searchFrom := 0
	if colonIdx >= 0 {
		searchFrom = colonIdx
	}
	sub := ref[searchFrom:]
	atIdx := strings.LastIndex(sub, "@")
	if atIdx >= 0 {
		split := searchFrom + atIdx
		return ref[:split], ref[split+1:]
	}
	return ref, ""
}

// runInstallFromManifest implements manifest-driven `grimoire install` (no args, grimoire.toml present).
// It reads [dependencies] from grimoire.toml, ensures packages are cloned, installs matching skills,
// and writes/updates grimoire.lock.
func runInstallFromManifest(projectDir string) error {
	r, err := manifest.Load(projectDir)
	if err != nil {
		return fmt.Errorf("loading grimoire.toml: %w", err)
	}

	// Project install: deps come only from the project file, not merged with global/system.
	// This prevents global [dependencies] from bleeding into project installs.
	// Global install: use merged deps (global + system layers both contribute).
	allDeps := r.Deps
	devDeps := r.DevDeps
	if !flagInstallGlobal && flagInstallScope != "global" {
		projFile, projErr := manifest.ParseFile(manifest.ProjectPath(projectDir))
		if projErr != nil {
			return fmt.Errorf("reading project grimoire.toml: %w", projErr)
		}
		allDeps = projFile.Deps
		devDeps = projFile.DevDeps
	}
	if len(allDeps) == 0 && len(devDeps) == 0 {
		fmt.Println("grimoire.toml has no [dependencies] — nothing to install")
		return nil
	}

	fmt.Println("  Installing from grimoire.toml…")

	symlink := r.Core.InstallMode != "copy"
	if err := checkWindowsSymlinkSupport(symlink); err != nil {
		return err
	}
	targets := resolveTargets(flagInstallTarget)

	// Ensure non-official packages are cloned
	allKeys := make([]string, 0, len(allDeps)+len(devDeps))
	for key := range allDeps {
		allKeys = append(allKeys, key)
	}
	for key := range devDeps {
		allKeys = append(allKeys, key)
	}

	for _, key := range allKeys {
		pkgRef := config.ParsePackageRef(key)
		if pkgRef.IsLocal() || pkgRef.IsOfficialRepoPath() || pkgRef.PackageURL == "" {
			continue
		}
		dest := skills.PackageHome(pkgRef.PackageName)
		if dirExists(dest) {
			continue
		}
		fmt.Printf("  Cloning %s…\n", pkgRef.PackageName)
		buf := &bytes.Buffer{}
		if cloneErr := updateNamedPackage(pkgRef.PackageName, pkgRef.PackageURL, pkgRef.Tag, buf); cloneErr != nil {
			fmt.Fprintf(os.Stderr, "%s  dependency %q: %v\n", tui.IconWarn, key, cloneErr)
		} else {
			os.Stdout.Write(buf.Bytes())
		}
	}

	// Build PackageRef list and resolve skills
	refs := make([]config.PackageRef, 0, len(allKeys))
	for _, key := range allKeys {
		refs = append(refs, config.ParsePackageRef(key))
	}
	fmt.Print("  Loading skills…")
	all, conflicts, listErr := skills.SkillsMatchingAnyMeta(refs, skills.AllPackages())
	if listErr != nil {
		return listErr
	}
	printConflicts(conflicts)
	if tui.IsTTY() {
		fmt.Printf("\r\033[K  Loading skills… %d found\n\n", len(all))
	} else {
		fmt.Printf(" %d found\n\n", len(all))
	}

	// Install skills
	perAgent := make(map[string]int)
	var skipped int
	for i := range all {
		sk := all[i]
		var linkedTo []string
		for _, ag := range targets {
			n, installErr := installSkillToAgent(sk.Path, ag, symlink)
			if installErr != nil {
				fmt.Fprintf(os.Stderr, "  %s  %s — %s: %v\n", tui.IconWarn, sk.Name, agent.DisplayName(ag), installErr)
				continue
			}
			perAgent[ag] += n
			if n > 0 {
				linkedTo = append(linkedTo, agent.DisplayName(ag))
			}
		}
		if len(linkedTo) > 0 {
			printSkillInstalled(sk.Name, linkedTo)
		} else {
			skipped++
		}
	}
	if skipped > 0 {
		fmt.Printf("  %s  %d skills already up to date\n", tui.IconSkip, skipped)
	}

	// Clean broken symlinks
	for _, ag := range targets {
		_, _ = skills.CleanBrokenSymlinks(installSkillsDir(ag))
	}
	if !flagInstallNoCfg {
		for _, ag := range targets {
			if err := agent.ConfigureAgentMD(ag); err != nil {
				fmt.Fprintf(os.Stderr, "  warn: configuring %s: %v\n", ag, err)
			}
		}
	}

	// Build resolver metadata from installed packages + write grimoire.lock
	if writeErr := updateLockFile(projectDir, allDeps, all); writeErr != nil {
		fmt.Fprintf(os.Stderr, "%s  updating grimoire.lock: %v\n", tui.IconWarn, writeErr)
	}

	printInstallSummary(perAgent, targets)
	return nil
}

// updateLockFile builds lock entries from installed skills and writes grimoire.lock.
func updateLockFile(projectDir string, deps map[string]manifest.DepSpec, installed []skills.Skill) error {
	// Build resolver.SkillMeta from installed skills
	meta := make(map[string]resolver.SkillMeta, len(installed))
	for i := range installed {
		sk := installed[i]
		if _, seen := meta[sk.Name]; seen {
			continue
		}
		// Get package commit from cloned package dir (best-effort)
		commit := ""
		if sk.Package != "" {
			regDir := skills.PackageHome(sk.Package)
			if state, err := gitops.CurrentState(regDir); err == nil {
				commit = state.Commit
			}
		}
		// Derive source (dep key prefix) from the skill's package name
		source := sk.Package
		resolved := ""
		officialPkgName := skills.OfficialPackageDerivedName()
		skillIsOfficial := sk.Package == officialPkgName
		// Find matching dep key to get the resolved URL
		for depKey := range deps {
			pkgRef := config.ParsePackageRef(depKey)
			matched := (sk.Package != "" && pkgRef.PackageName == sk.Package) ||
				(pkgRef.IsOfficialRepoPath() && skillIsOfficial)
			if matched {
				resolved = pkgRef.PackageURL
				if resolved == "" {
					resolved = skills.GrimoireRepoURL()
				}
				break
			}
		}
		meta[sk.Name] = resolver.SkillMeta{
			Name:     sk.Name,
			Version:  sk.Version,
			Source:   source,
			Resolved: resolved,
			Commit:   commit,
			Checksum: skillChecksum(sk.Path),
		}
	}

	r := resolver.New(meta)
	entries, err := r.Resolve(deps)
	if err != nil {
		return err
	}

	lockPath := manifest.LockPath(manifest.ProjectPath(projectDir))
	lf, _ := lock.ParseFile(lockPath)
	for i := range entries {
		lf.Upsert(&entries[i])
	}
	return lock.WriteFile(lockPath, lf)
}

// runInstallFromPackageURL implements the URL/local-path branch of `grimoire install <grimoire-ref>`:
// derives a package name, persists to global settings, clones, then installs skills.
func runInstallFromPackageURL(ref string) error {
	u, ver := config.ParseRef(ref)
	if u == "" {
		u = ref
	}
	if !skills.IsGitURL(u) && !filepath.IsAbs(u) {
		return fmt.Errorf("%q is not a git URL or local path\nhint: grimoire package add <name> <grimoire-ref>", ref)
	}

	name := config.DerivePackageName(u)
	if ver != "" {
		name += "@" + ver
	}

	// Persist to [dependencies] in the scope-appropriate config file (idempotent).
	var settingsPath string
	if flagInstallGlobal || flagInstallScope == "global" {
		settingsPath = config.GlobalPath()
	} else {
		settingsPath = config.ProjectPath(getProjectDir())
	}
	fs, _ := config.ParseFile(settingsPath)
	alreadyDep := false
	for _, d := range fs.Dependencies.Skills {
		if d == ref {
			alreadyDep = true
			break
		}
	}
	if !alreadyDep {
		fs.Dependencies.Skills = append(fs.Dependencies.Skills, ref)
		if err := config.WriteFile(settingsPath, fs); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}
		fmt.Printf("  %s  added dependency: %s\n", tui.IconOK, name)
	}

	// Clone / update.
	buf := &bytes.Buffer{}
	if err := updateNamedPackage(name, ref, "", buf); err != nil {
		return err
	}
	os.Stdout.Write(buf.Bytes())

	// Install skills from this package only.
	regs := skills.AllSkillsPackages()
	regs = filterPackages(regs, name)
	if len(regs) == 0 {
		return fmt.Errorf("package %q has no skills after cloning", name)
	}

	symlink := !flagInstallCopy
	r, _ := config.Load(getProjectDir())
	if r.Core.InstallMode == "copy" {
		symlink = false
	}
	if err := checkWindowsSymlinkSupport(symlink); err != nil {
		return err
	}
	target := flagInstallTarget
	if flagInstallYes && target == "" {
		target = "auto"
	}
	targets := resolveTargets(target)
	perAgent := make(map[string]int)

	fmt.Print("  Loading skills…")
	all, conflicts, err := skills.ListAllSkillsFromPackagesMeta(regs)
	if err != nil {
		return err
	}
	printConflicts(conflicts)
	if tui.IsTTY() {
		fmt.Printf("\r\033[K  Loading skills… %d found\n\n", len(all))
	} else {
		fmt.Printf(" %d found\n\n", len(all))
	}
	var skippedURL int
	for i := range all {
		sk := all[i]
		var linkedTo []string
		for _, ag := range targets {
			n, installErr := installSkillToAgent(sk.Path, ag, symlink)
			if installErr != nil {
				fmt.Fprintf(os.Stderr, "  %s  %s — %s: %v\n", tui.IconWarn, sk.Name, agent.DisplayName(ag), installErr)
				continue
			}
			perAgent[ag] += n
			if n > 0 {
				linkedTo = append(linkedTo, agent.DisplayName(ag))
			}
		}
		if len(linkedTo) > 0 {
			printSkillInstalled(sk.Name, linkedTo)
		} else {
			skippedURL++
		}
	}
	if skippedURL > 0 {
		fmt.Printf("  %s  %d skills already up to date\n", tui.IconSkip, skippedURL)
	}
	for _, ag := range targets {
		_, _ = skills.CleanBrokenSymlinks(installSkillsDir(ag))
	}
	if !flagInstallNoCfg {
		for _, ag := range targets {
			if err := agent.ConfigureAgentMD(ag); err != nil {
				fmt.Fprintf(os.Stderr, "  warn: configuring %s: %v\n", ag, err)
			}
		}
	}
	printInstallSummary(perAgent, targets)
	return nil
}

// runInstallFromRoot runs a single-root install (--from path), bypassing multi-package logic.
func runInstallFromRoot(root string) error {
	if _, err := os.Stat(root); err != nil {
		return fmt.Errorf("skills not found at %s", root)
	}
	symlink := !flagInstallCopy
	if !flagInstallCopy {
		r, _ := config.Load(getProjectDir())
		if r.Core.InstallMode == "copy" {
			symlink = false
		}
	}
	if err := checkWindowsSymlinkSupport(symlink); err != nil {
		return err
	}
	target := flagInstallTarget
	if flagInstallYes && target == "" {
		target = "auto"
	}
	targets := resolveTargets(target)
	perAgent := make(map[string]int)

	domains, err := skills.ListDomains(root)
	if err != nil {
		return err
	}
	for _, d := range domains {
		for _, ag := range targets {
			n, err := installDomainToAgent(root, d.Name, "", ag, symlink, "")
			if err != nil {
				fmt.Fprintf(os.Stderr, "  error: %v\n", err)
			}
			perAgent[ag] += n
		}
	}

	for _, ag := range targets {
		_, _ = skills.CleanBrokenSymlinks(installSkillsDir(ag))
	}
	if !flagInstallNoCfg {
		for _, ag := range targets {
			if err := agent.ConfigureAgentMD(ag); err != nil {
				fmt.Fprintf(os.Stderr, "  warn: configuring %s: %v\n", ag, err)
			}
		}
	}
	printInstallSummary(perAgent, targets)
	return nil
}

// splitPackagePrefix parses "my-org:engineering/tdd" into ("my-org", "engineering/tdd", true).
func splitPackagePrefix(ref string) (pkg, skill string, ok bool) {
	for i, ch := range ref {
		if ch == ':' {
			return ref[:i], ref[i+1:], true
		}
		if ch == '/' {
			break // no package prefix
		}
	}
	return "", "", false
}

// resolveSkillFromPackages finds a skill by ref across multiple packages; first match wins.
func resolveSkillFromPackages(regs []skills.SkillsPackage, ref string) (path string, reg skills.SkillsPackage, err error) {
	for _, r := range regs {
		p, e := skills.ResolveSkillPath(r.Root, ref)
		if e == nil {
			return p, r, nil
		}
	}
	// Fallback: glob-based search (handles bare names like "apply-dry-principle" and
	// .gitignore-style patterns). r.Root is the skills/ subdir; Dir gives the package home.
	for _, r := range regs {
		matches, e := skills.SkillsMatchingGlob(filepath.Dir(r.Root), ref)
		if e == nil && len(matches) > 0 {
			return matches[0].Path, r, nil
		}
	}
	return "", skills.SkillsPackage{}, fmt.Errorf("skill %q not found in any configured package", ref)
}

// resolveAndPersistSource resolves a --from value (local path or git URL),
// clones if needed, persists to global config, and returns the skills root path.
// Returns ("", nil) when the user cancels.
func resolveAndPersistSource(from string) (string, error) {
	if skills.IsGitURL(from) {
		home := skills.OfficialPackageHome()
		if _, err := os.Stat(home); err == nil {
			chosen, ok := tui.RunSelect(
				fmt.Sprintf("Replace existing grimoire at %s with %s?", home, from),
				[]string{"Yes", "Cancel"},
			)
			if !ok || chosen == "Cancel" {
				fmt.Println("Cancelled.")
				return "", nil
			}
			if err := os.RemoveAll(home); err != nil {
				return "", fmt.Errorf("removing %s: %w", home, err)
			}
		}
		fmt.Printf("Cloning %s → %s...\n", from, home)
		if err := os.MkdirAll(filepath.Dir(home), 0o755); err != nil {
			return "", fmt.Errorf("creating dir: %w", err)
		}
		if err := gitops.Clone(from, home); err != nil {
			return "", fmt.Errorf("cloning: %w", err)
		}
	}

	if skills.IsGitURL(from) {
		return skills.SkillsRoot(), nil
	}
	abs, err := filepath.Abs(from)
	if err != nil {
		return "", fmt.Errorf("resolving path %s: %w", from, err)
	}
	return filepath.Join(abs, "skills"), nil
}

func installDomainToAgent(root, domain, subdomain, ag string, symlink bool, pkgName string) (int, error) {
	domainDir := fmt.Sprintf("%s/%s", root, domain)
	if _, err := os.Stat(domainDir); err != nil {
		return 0, fmt.Errorf("domain not found: %s", domain)
	}
	count := 0
	destDir := installSkillsDir(ag)

	// Pre-load global manifest once — avoids per-skill file reads and handles
	// already-linked skills (InstallSkill ok=false) that still need dep recording.
	var globalManifestPath string
	existingDeps := make(map[string]manifest.DepSpec)
	if flagInstallGlobal || flagInstallScope == "global" {
		globalManifestPath = ensureGlobalManifest()
		mf, _ := manifest.ParseFile(globalManifestPath)
		if mf.Deps != nil {
			existingDeps = mf.Deps
		}
	}
	recordDep := func(skillName string) {
		if globalManifestPath == "" {
			return
		}
		key := depKeyForSkill(skillName, pkgName)
		if _, exists := existingDeps[key]; !exists {
			_ = manifest.AppendDep(globalManifestPath, key, "*")
			existingDeps[key] = manifest.DepSpec{Version: "*"}
		}
	}

	if skills.IsNested(domainDir) {
		subs, err := skills.ListSubdomains(domainDir)
		if err != nil {
			return 0, err
		}
		for _, sub := range subs {
			if subdomain != "" && sub.Name != subdomain {
				continue
			}
			skillList, err := skills.ListSkillsInDir(sub.Path, domain, sub.Name)
			if err != nil {
				continue
			}
			for i := range skillList {
				sk := skillList[i]
				ok, err := skills.InstallSkill(sk.Path, destDir, symlink)
				if err != nil {
					fmt.Fprintf(os.Stderr, "  warn: %v\n", err)
					continue
				}
				if ok {
					fmt.Printf("  %s  %s\n", tui.IconDone, sk.Name)
					count++
				}
				recordDep(sk.Name)
			}
		}
	} else {
		skillList, err := skills.ListSkillsInDir(domainDir, domain, "")
		if err != nil {
			return 0, err
		}
		for i := range skillList {
			sk := skillList[i]
			ok, err := skills.InstallSkill(sk.Path, destDir, symlink)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  warn: %v\n", err)
				continue
			}
			if ok {
				fmt.Printf("  %s %s\n", tui.StyleDim.Render("linked:"), sk.Name)
				count++
			}
			recordDep(sk.Name)
		}
	}
	return count, nil
}

func installSkillToAgent(skillPath, ag string, symlink bool) (int, error) {
	destDir := installSkillsDir(ag)
	ok, err := skills.InstallSkill(skillPath, destDir, symlink)
	if err != nil {
		return 0, err
	}
	if ok {
		return 1, nil
	}
	return 0, nil
}

// printSkillInstalled prints one ✓ line per skill, listing agent display names when >1.
func printSkillInstalled(name string, toAgents []string) {
	if len(toAgents) <= 1 {
		fmt.Printf("  %s  %s\n", tui.IconDone, name)
	} else {
		fmt.Printf("  %s  %-48s %s\n", tui.IconDone, name, strings.Join(toAgents, ", "))
	}
}

func printConflicts(conflicts []skills.SkillConflict) {
	for _, c := range conflicts {
		fmt.Fprintf(os.Stderr, "  %s  %s: %s wins over %s (override: grimoire install --package %s --skill %s)\n",
			tui.IconWarn, c.CanonicalPath, c.WinnerPackage, c.LoserPackage, c.LoserPackage, c.CanonicalPath)
	}
}

// resolveInstallPool builds the candidate skill pool from [dependencies] skills,
// then filters by active profiles when standards.profiles is declared.
func resolveInstallPool(r *config.Config) ([]skills.Skill, []skills.SkillConflict, error) {
	var all []skills.Skill
	var conflicts []skills.SkillConflict
	var err error
	if len(r.DepSkills) > 0 {
		refs := make([]config.PackageRef, len(r.DepSkills))
		for i, s := range r.DepSkills {
			refs[i] = config.ParsePackageRef(s)
		}
		all, conflicts, err = skills.SkillsMatchingAnyMeta(refs, skills.AllPackages())
	} else {
		all, conflicts, err = skills.ListAllSkillsFromPackagesMeta(skills.AllSkillsPackages())
	}
	if err != nil {
		return nil, nil, err
	}

	// Filter candidate pool to skills declared in active profiles.
	if len(r.Core.Profiles) > 0 {
		projectDir := getProjectDir()
		effectiveSkills, resolveErr := profiles.ResolveEffectiveSkills(r.Core.Profiles, projectDir, nil, nil)
		if resolveErr == nil && len(effectiveSkills) > 0 {
			allowed := make(map[string]bool, len(effectiveSkills))
			for _, sk := range effectiveSkills {
				allowed[sk.Name] = true
			}
			filtered := all[:0]
			for i := range all {
				sk := all[i]
				if allowed[sk.Name] {
					filtered = append(filtered, sk)
				}
			}
			all = filtered
		}
	}

	return all, conflicts, nil
}

// skillChecksum returns sha256:<hex> of the skill's SKILL.md content, or "" on error.
func skillChecksum(skillPath string) string {
	data, err := os.ReadFile(filepath.Join(skillPath, "SKILL.md"))
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func joinAgentNames(agents []string) string {
	names := make([]string, len(agents))
	for i, ag := range agents {
		names[i] = agent.DisplayName(ag)
	}
	result := ""
	for i, n := range names {
		if i > 0 {
			result += ", "
		}
		result += n
	}
	return result
}

// checkWindowsSymlinkSupport tests symlink capability on Windows before install.
// Returns nil on non-Windows or when copy mode is active.
func checkWindowsSymlinkSupport(symlink bool) error {
	if runtime.GOOS != "windows" || !symlink {
		return nil
	}
	tmp, err := os.MkdirTemp("", "grimoire-symlink-test-*")
	if err != nil {
		return nil // can't test; let the real install surface errors
	}
	defer func() { _ = os.RemoveAll(tmp) }()
	src := filepath.Join(tmp, "src")
	_ = os.WriteFile(src, nil, 0o600)
	if err := os.Symlink(src, filepath.Join(tmp, "link")); err != nil {
		return fmt.Errorf("symlinks are not available on this system\n\n" +
			"Windows requires one of:\n" +
			"  • Enable Developer Mode: Settings → System → Developer Mode → On\n" +
			"  • Run grimoire as Administrator\n\n" +
			"Or use copy mode instead:\n" +
			"  grimoire config set core.install-mode copy --global\n" +
			"  grimoire install")
	}
	return nil
}
