package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/agent"
	"github.com/jeffreytse/grimoire/internal/config"
	"github.com/jeffreytse/grimoire/internal/lock"
	"github.com/jeffreytse/grimoire/internal/manifest"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var (
	flagUninstallDomain    string
	flagUninstallSubdomain string
	flagUninstallSkill     string
	flagUninstallTarget    string
	flagUninstallScope     string // "project" | "global"
	flagUninstallGlobal    bool   // shorthand for --scope global
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [<grimoire-ref>]",
	Short: "Remove grimoire skills from agent directories",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runUninstall,
}

func init() {
	uninstallCmd.Flags().StringVar(&flagUninstallDomain, "domain", "", "uninstall all skills for a domain")
	uninstallCmd.Flags().StringVar(&flagUninstallSubdomain, "subdomain", "", "restrict to one sub-domain")
	uninstallCmd.Flags().StringVar(&flagUninstallSkill, "skill", "", "uninstall one skill (domain/subdomain/name or domain/name)")
	uninstallCmd.Flags().StringVar(&flagUninstallTarget, "target", "", "target agent (default: all detected)")
	uninstallCmd.Flags().StringVar(&flagUninstallScope, "scope", "project", `uninstall scope: "project" or "global"`)
	uninstallCmd.Flags().BoolVar(&flagUninstallGlobal, "global", false, `global-scoped uninstall — shorthand for --scope global`)
}

// uninstallSkillsDir returns the agent skills dir based on scope flags.
func uninstallSkillsDir(ag string) string {
	if flagUninstallGlobal || flagUninstallScope == "global" {
		return agent.SkillsDir(ag)
	}
	return agent.ProjectSkillsDir(ag, getProjectDir())
}

// uninstallSettingsPath returns the config file to modify based on scope flags.
func uninstallSettingsPath() string {
	if flagUninstallGlobal || flagUninstallScope == "global" {
		return config.GlobalPath()
	}
	return config.ProjectPath(getProjectDir())
}

func runUninstall(cmd *cobra.Command, args []string) error {
	// grimoire uninstall <package-ref> — remove by package URL
	if len(args) > 0 {
		return runUninstallPackage(args[0])
	}

	root := skills.SkillsRoot()
	targets := resolveUninstallTargets(flagUninstallTarget)
	count := 0

	switch {
	case flagUninstallSkill != "":
		skillPath, err := skills.ResolveSkillPath(root, flagUninstallSkill)
		if err != nil {
			// allow uninstalling even if source skill no longer exists
			skillPath = flagUninstallSkill
		}
		name := skillNameFromPath(skillPath, flagUninstallSkill)
		fmt.Printf("Uninstalling skill: %s\n", name)
		for _, ag := range targets {
			ok, err := skills.UninstallSkill(name, agent.SkillsDir(ag))
			if err != nil {
				fmt.Fprintf(os.Stderr, "  warn: %v\n", err)
				continue
			}
			if ok {
				fmt.Printf("  %s %s from %s\n",
					tui.StyleDim.Render("removed:"), name, agent.DisplayName(ag))
				count++
			}
		}

	case flagUninstallDomain != "":
		for _, ag := range targets {
			n, err := uninstallDomainFromAgent(root, flagUninstallDomain, flagUninstallSubdomain, ag)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  error: %v\n", err)
			}
			count += n
		}

	default:
		// uninstall everything
		domains, err := skills.ListDomains(root)
		if err != nil {
			return err
		}
		for _, d := range domains {
			for _, ag := range targets {
				n, err := uninstallDomainFromAgent(root, d.Name, "", ag)
				if err != nil {
					fmt.Fprintf(os.Stderr, "  error: %v\n", err)
				}
				count += n
			}
		}
	}

	// clean broken symlinks
	for _, ag := range targets {
		_, _ = skills.CleanBrokenSymlinks(agent.SkillsDir(ag))
	}

	// remove agent MD config if no skills remain
	for _, ag := range targets {
		if agent.SkillCount(ag) == 0 {
			_ = agent.RemoveAgentMDConfig(ag)
		}
	}

	unique := count / len(targets)
	switch {
	case len(targets) > 1:
		fmt.Printf("\n%s  %d skills uninstalled × %d agents (%d total)\n",
			tui.IconOK, unique, len(targets), count)
	case count > 0:
		fmt.Printf("\n%s  %d skills uninstalled\n", tui.IconOK, count)
	default:
		fmt.Printf("\n%s  nothing to uninstall\n", tui.IconOK)
	}
	return nil
}

func uninstallDomainFromAgent(root, domain, subdomain, ag string) (int, error) {
	domainDir := fmt.Sprintf("%s/%s", root, domain)
	count := 0
	destDir := agent.SkillsDir(ag)
	fmt.Printf("Uninstalling domain: %s from %s\n", domain, agent.DisplayName(ag))

	if _, err := os.Stat(domainDir); err != nil {
		return 0, nil // domain dir gone, nothing to do
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
				ok, err := skills.UninstallSkill(sk.Name, destDir)
				if err != nil {
					fmt.Fprintf(os.Stderr, "  warn: %v\n", err)
					continue
				}
				if ok {
					fmt.Printf("  %s %s\n", tui.StyleDim.Render("removed:"), sk.Name)
					count++
				}
			}
		}
	} else {
		skillList, err := skills.ListSkillsInDir(domainDir, domain, "")
		if err != nil {
			return 0, err
		}
		for i := range skillList {
			sk := skillList[i]
			ok, err := skills.UninstallSkill(sk.Name, destDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  warn: %v\n", err)
				continue
			}
			if ok {
				fmt.Printf("  %s %s\n", tui.StyleDim.Render("removed:"), sk.Name)
				count++
			}
		}
	}
	return count, nil
}

// runUninstallPackage implements `grimoire uninstall <package-ref>`:
// removes matching skills from scope-appropriate agent skill dirs
// and removes the ref from [dependencies] skills in settings.
func runUninstallPackage(ref string) error {
	pkgRef := config.ParsePackageRef(ref)
	targets := resolveUninstallTargets(flagUninstallTarget)
	count := 0

	for _, ag := range targets {
		skillDir := uninstallSkillsDir(ag)
		var packageHome string
		switch {
		case pkgRef.IsLocal():
			packageHome = filepath.FromSlash(pkgRef.LocalPath)
		case pkgRef.IsOfficialRepoPath():
			packageHome = skills.OfficialPackageHome()
		default:
			packageHome = skills.PackageHome(pkgRef.PackageName)
		}
		n, err := removeSkillsFromDir(skillDir, packageHome, pkgRef.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  %s  %s: %v\n", tui.IconWarn, agent.DisplayName(ag), err)
		}
		count += n
	}

	// Clean broken symlinks after removal.
	for _, ag := range targets {
		_, _ = skills.CleanBrokenSymlinks(uninstallSkillsDir(ag))
	}

	// Remove matching refs from settings.toml (legacy path).
	settingsPath := uninstallSettingsPath()
	fs, _ := config.ParseFile(settingsPath)
	before := len(fs.Dependencies.Skills)
	fs.Dependencies.Skills = removeMatchingRefs(fs.Dependencies.Skills, &pkgRef)
	if len(fs.Dependencies.Skills) < before {
		if err := config.WriteFile(settingsPath, fs); err != nil {
			fmt.Fprintf(os.Stderr, "  %s  updating config: %v\n", tui.IconWarn, err)
		}
	}

	// Remove matching deps from grimoire.toml when present.
	projectDir := getProjectDir()
	manifestPath := manifest.ProjectPath(projectDir)
	if _, err := os.Stat(manifestPath); err == nil {
		mf, _ := manifest.ParseFile(manifestPath)
		changed := false
		for key := range mf.Deps {
			kRef := config.ParsePackageRef(key)
			if refsMatch(&kRef, &pkgRef) {
				delete(mf.Deps, key)
				changed = true
			}
		}
		if changed {
			if writeErr := manifest.WriteFile(manifestPath, &mf); writeErr != nil {
				fmt.Fprintf(os.Stderr, "  %s  updating grimoire.toml: %v\n", tui.IconWarn, writeErr)
			}
			// Remove from lock file too.
			lockPath := manifest.LockPath(manifestPath)
			if lf, err := lock.ParseFile(lockPath); err == nil {
				lf.Remove(pkgRef.Path)
				_ = lock.WriteFile(lockPath, lf)
			}
		}
	}

	switch {
	case count > 0:
		fmt.Printf("\n%s  %d skills uninstalled\n", tui.IconOK, count)
	default:
		fmt.Printf("\n%s  nothing to uninstall\n", tui.IconOK)
	}
	return nil
}

// removeSkillsFromDir removes symlinks in skillDir whose targets are under packageHome
// and whose relative path from packageHome matches pathFilter (empty = all).
func removeSkillsFromDir(skillDir, packageHome, pathFilter string) (int, error) {
	entries, err := os.ReadDir(skillDir)
	if err != nil {
		return 0, nil // dir absent = nothing to remove
	}
	count := 0
	for _, e := range entries {
		entryPath := filepath.Join(skillDir, e.Name())
		target, err := os.Readlink(entryPath)
		if err != nil {
			continue // not a symlink
		}
		if !strings.HasPrefix(target, packageHome) {
			continue
		}
		if pathFilter != "" {
			rel, _ := filepath.Rel(packageHome, target)
			rel = filepath.ToSlash(rel)
			if !skills.GlobMatch(pathFilter, rel) && !strings.HasPrefix(rel, filepath.ToSlash(pathFilter)+"/") {
				continue
			}
		}
		if err := os.Remove(entryPath); err != nil {
			fmt.Fprintf(os.Stderr, "  warn: removing %s: %v\n", entryPath, err)
			continue
		}
		fmt.Printf("  %s %s\n", tui.StyleDim.Render("removed:"), e.Name())
		count++
	}
	return count, nil
}

// removeMatchingRefs filters out package refs that match pkgRef's package (and path, if set).
func removeMatchingRefs(refs []string, pkgRef *config.PackageRef) []string {
	var result []string
	for _, r := range refs {
		parsed := config.ParsePackageRef(r)
		sameSource := false
		switch {
		case pkgRef.IsLocal():
			sameSource = parsed.IsLocal() && parsed.LocalPath == pkgRef.LocalPath
		case pkgRef.IsOfficialRepoPath():
			sameSource = parsed.IsOfficialRepoPath()
		default:
			sameSource = parsed.PackageName == pkgRef.PackageName
		}
		if !sameSource {
			result = append(result, r)
			continue
		}
		// Same source: keep only if paths don't overlap with pkgRef.Path.
		if pkgRef.Path == "" {
			continue // remove all from this source
		}
		if parsed.Path == "" {
			continue // existing "all from package" ref — remove it
		}
		// Keep if parsed.Path is NOT a sub-path of pkgRef.Path
		if parsed.Path != pkgRef.Path && !strings.HasPrefix(parsed.Path, pkgRef.Path+"/") {
			result = append(result, r)
		}
	}
	return result
}

// refsMatch reports whether two PackageRefs refer to the same package source.
func refsMatch(a, b *config.PackageRef) bool {
	if a.IsOfficialRepoPath() && b.IsOfficialRepoPath() {
		return b.Path == "" || a.Path == b.Path
	}
	if a.IsLocal() && b.IsLocal() {
		return a.LocalPath == b.LocalPath
	}
	return a.PackageName == b.PackageName && (b.Path == "" || a.Path == b.Path)
}

// resolveUninstallTargets is like resolveTargets but uses DetectedOrInstalled for the
// auto case, so uninstall covers agents whose binary is no longer in PATH but whose
// skills directory is non-empty (e.g. agy removed from PATH after install).
func resolveUninstallTargets(target string) []string {
	switch target {
	case "", "auto":
		r, _ := config.Load(getProjectDir())
		if len(r.Core.Agents) > 0 {
			return r.Core.Agents
		}
		found := agent.DetectedOrInstalled()
		if len(found) == 0 {
			return []string{"claude"}
		}
		return found
	default:
		return resolveTargets(target)
	}
}

func skillNameFromPath(skillPath, ref string) string {
	// if skillPath is a proper path, get basename
	if len(skillPath) > len(ref) {
		parts := splitLast(skillPath, '/')
		return parts
	}
	// fall back to last segment of ref
	return splitLast(ref, '/')
}

func splitLast(s string, sep byte) string {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == sep {
			return s[i+1:]
		}
	}
	return s
}
