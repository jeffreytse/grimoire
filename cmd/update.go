package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/agent"
	"github.com/jeffreytse/grimoire/internal/config"
	gitops "github.com/jeffreytse/grimoire/internal/git"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var (
	flagUpdateStable bool
	flagUpdateYes    bool
	flagUpdateDryRun bool
	flagUpdateForce  bool
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"upgrade"},
	Short:   "Update all packages to latest",
	RunE:    runUpdate,
}

func init() {
	updateCmd.Flags().BoolVar(&flagUpdateStable, "stable", false, "check out the latest tagged release instead of HEAD")
	updateCmd.Flags().BoolVar(&flagUpdateYes, "yes", false, "skip confirmation prompts")
	updateCmd.Flags().BoolVar(&flagUpdateDryRun, "dry-run", false, "show what would change without pulling")
	updateCmd.Flags().BoolVar(&flagUpdateForce, "force", false, "discard local package modifications and update")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	home := skills.OfficialPackageHome()
	url := skills.GrimoireRepoURL()

	// Local path package — nothing to pull
	if filepath.IsAbs(url) {
		fmt.Printf("  %s  Official package is a local path — no update needed.\n", tui.IconOK)
		return nil
	}

	// Clone if not present
	if _, err := os.Stat(home); err != nil {
		fmt.Printf("Grimoire not found at %s\n", home)
		if !flagUpdateYes {
			chosen, ok := tui.RunSelect("Clone grimoire to "+home+"?", []string{"Yes", "No"})
			if !ok || chosen == "No" {
				fmt.Println("Cancelled.")
				return nil
			}
		}
		fmt.Printf("Cloning %s...\n", skills.GrimoireRepoURL())
		if err := os.MkdirAll(filepath.Dir(home), 0o755); err != nil {
			return fmt.Errorf("creating dir: %w", err)
		}
		if err := gitops.Clone(skills.GrimoireRepoURL(), home); err != nil {
			return fmt.Errorf("cloning grimoire: %w", err)
		}
		fmt.Printf("%s  Grimoire installed to %s\n", tui.IconOK, home)
		return nil
	}

	var coreErr error
	if flagUpdateStable {
		coreErr = updateStable(home)
	} else {
		coreErr = updateUnstable(home)
	}
	if coreErr != nil {
		return coreErr
	}

	// Also update all non-official, non-pinned, non-local packages from [[package]].
	cfg, err := config.LoadGlobal()
	if err != nil {
		return nil // best-effort
	}
	for _, rd := range cfg.Packages {
		if rd.Official || !rd.Enabled {
			continue
		}
		u, ver := config.ParseRef(rd.URL)
		if u == "" {
			u = rd.URL
		}
		if filepath.IsAbs(u) || ver != "" {
			continue // local or pinned — skip
		}
		if err := updateNamedPackage(rd.Name, rd.URL, "", os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "  warn: %s: %v\n", rd.Name, err)
		}
	}
	return nil
}

func updateStable(home string) error {
	current, err := gitops.CurrentState(home)
	if err != nil {
		return fmt.Errorf("reading current state: %w", err)
	}

	fmt.Println("Fetching release tags...")
	if err := gitops.FetchTags(home); err != nil && !errors.Is(err, errNotGit(home)) {
		fmt.Fprintf(os.Stderr, "  warn: fetch tags: %v\n", err)
	}

	latest, err := gitops.LatestTag(home)
	if err != nil {
		return fmt.Errorf("finding latest tag: %w", err)
	}
	tagState, err := gitops.TagState(home, latest)
	if err != nil {
		return fmt.Errorf("reading tag state: %w", err)
	}

	if tagState.Commit == current.Commit {
		fmt.Printf("  %s  Already on latest stable release. (%s, commit %s, %s)\n\n",
			tui.IconOK, latest, tagState.Commit, tagState.Date)
		return nil
	}

	fmt.Printf("\n  New stable release available:\n")
	fmt.Printf("    Current:  v%s (commit %s, %s)\n", current.Version, current.Commit, current.Date)
	fmt.Printf("    New:      %s  (commit %s, %s)\n", latest, tagState.Commit, tagState.Date)
	fmt.Println()

	if !flagUpdateYes {
		chosen, ok := tui.RunSelect("Upgrade now?", []string{"Yes", "No"})
		if !ok || chosen == "No" {
			fmt.Println("Upgrade cancelled.")
			return nil
		}
	}

	if flagUpdateDryRun {
		fmt.Printf("  (dry run — not checking out)\n")
		return nil
	}

	snapshot := snapshotSkillVersions(home)

	fmt.Printf("Checking out %s...\n", latest)
	if err := gitops.CheckoutTag(home, latest); err != nil {
		return fmt.Errorf("checking out %s: %w", latest, err)
	}
	skills.InvalidateSkillCache(home)

	printUpgradeResult(current, tagState)
	relinkNewSkills(home, current.Commit)
	printVersionDiff(home, snapshot)
	return nil
}

func updateUnstable(home string) error {
	current, err := gitops.CurrentState(home)
	if err != nil {
		return fmt.Errorf("reading current state: %w", err)
	}

	fmt.Println("Checking for updates...")
	upToDate, _, remote, err := gitops.IsUpToDate(home)
	if err != nil {
		return fmt.Errorf("checking upstream: %w", err)
	}
	if upToDate {
		fmt.Printf("  %s  Already up to date. (v%s, commit %s, %s)\n\n",
			tui.IconOK, current.Version, current.Commit, current.Date)
		return nil
	}

	fmt.Printf("\n  New version available:\n")
	fmt.Printf("    Current:  v%s (commit %s, %s)\n", current.Version, current.Commit, current.Date)
	if remote.Version != "" && remote.Version != "unknown" {
		fmt.Printf("    New:      v%s (commit %s, %s)\n", remote.Version, remote.Commit, remote.Date)
	} else {
		fmt.Printf("    New:      commit %s (%s)\n", remote.Commit, remote.Date)
	}
	fmt.Println()

	if !flagUpdateYes {
		chosen, ok := tui.RunSelect("Pull now?", []string{"Yes", "No"})
		if !ok || chosen == "No" {
			fmt.Println("Update cancelled.")
			return nil
		}
	}

	if flagUpdateDryRun {
		fmt.Printf("  (dry run — not pulling)\n")
		return nil
	}

	// Snapshot criteria before pull so we can diff after.
	snapshot := snapshotSkillVersions(home)

	fmt.Printf("Pulling latest grimoire at %s...\n", home)
	if err := gitops.PullWithForceFallback(home, flagUpdateForce); err != nil {
		return fmt.Errorf("updating: %w", err)
	}
	skills.InvalidateSkillCache(home)

	newState, _ := gitops.CurrentState(home)
	printUpgradeResult(current, newState)
	relinkNewSkills(home, current.Commit)
	printVersionDiff(home, snapshot)
	return nil
}

func printUpgradeResult(old, updated gitops.State) {
	fmt.Printf("\n  %s  Grimoire upgraded to latest.\n\n", tui.IconOK)
	fmt.Printf("    Previous: v%s (commit %s, %s)\n", old.Version, old.Commit, old.Date)
	fmt.Printf("    Current:  v%s (commit %s, %s)\n", updated.Version, updated.Commit, updated.Date)
	fmt.Println()
}

func relinkNewSkills(home, oldCommit string) {
	root := skills.SkillsRoot()
	changes, err := gitops.PackageChangesSince(home, oldCommit)
	if err == nil {
		printPackageChanges(changes, home, oldCommit, os.Stdout)
		if !changes.HasSkillChanges() {
			fmt.Printf("  %s  No skill changes — skipping relink.\n", tui.IconOK)
			fmt.Println()
			return
		}
	}

	// relink for all agents that already have a skills dir
	for _, ag := range agent.All {
		dir := agent.SkillsDir(ag)
		if _, err := os.Stat(dir); err != nil {
			continue
		}
		allSkills, err := skills.ListAllSkills(root)
		if err != nil {
			continue
		}
		for i := range allSkills {
			_, _ = skills.InstallSkill(allSkills[i].Path, dir, true)
		}
		_, _ = skills.CleanBrokenSymlinks(dir)
	}

	fmt.Println()
}

func printPackageChanges(c gitops.PackageChanges, dir, oldCommit string, w io.Writer) { //nolint:gocritic // callers in two files; value semantics avoids pointer threading
	if dir != "" && oldCommit != "" {
		if commits, _ := gitops.CommitsSince(dir, oldCommit); len(commits) > 0 {
			fmt.Fprintf(w, "    %s\n", tui.StyleBold.Render("commits:"))
			for _, line := range commits {
				parts := strings.SplitN(line, " ", 2)
				if len(parts) == 2 {
					fmt.Fprintf(w, "      %s %s\n", tui.StyleDim.Render(parts[0]), parts[1])
				} else {
					fmt.Fprintf(w, "      %s\n", line)
				}
			}
		}
	}
	printChangeSection("skills", c.SkillsAdded, c.SkillsUpdated, w)
	printChangeSection("profiles", c.ProfilesAdded, c.ProfilesUpdated, w)
}

func printChangeSection(label string, added, updated []string, w io.Writer) {
	if len(added)+len(updated) == 0 {
		return
	}
	fmt.Fprintf(w, "    %s\n", tui.StyleBold.Render(label+":"))
	for _, name := range added {
		fmt.Fprintf(w, "      %s %s\n", tui.StyleGreen.Render("+"), name)
	}
	for _, name := range updated {
		fmt.Fprintf(w, "      %s %s\n", tui.StyleGold.Render("~"), name)
	}
}

func errNotGit(dir string) error {
	return fmt.Errorf("%s is not a git repository", dir)
}

// versionSnapshot records skill versions before an update so changes can be shown after.
type versionSnapshot map[string]string // skill name → version (empty string if unversioned)

// snapshotSkillVersions walks the package and captures skill versions.
func snapshotSkillVersions(home string) versionSnapshot {
	snap := make(versionSnapshot)
	all, err := skills.ListAllSkills(home)
	if err != nil {
		return snap
	}
	for i := range all {
		sk := all[i]
		v := sk.Version
		if v == "" {
			v = "unversioned"
		}
		snap[sk.Name] = v
	}
	return snap
}

// printVersionDiff compares the current package skills against the pre-update snapshot
// and prints any version changes.
func printVersionDiff(home string, before versionSnapshot) {
	if len(before) == 0 {
		return
	}
	all, err := skills.ListAllSkills(home)
	if err != nil {
		return
	}
	var diffs []string
	for i := range all {
		sk := all[i]
		oldVer, exists := before[sk.Name]
		if !exists {
			continue // new skill — already shown in package changes
		}
		newVer := sk.Version
		if newVer == "" {
			newVer = "unversioned"
		}
		if oldVer == newVer {
			continue
		}
		diffs = append(diffs, fmt.Sprintf("    %s %s  %s → %s",
			tui.StyleGold.Render("~"), sk.Name,
			tui.StyleDim.Render(oldVer), newVer))
	}
	if len(diffs) == 0 {
		return
	}
	fmt.Printf("  %s  Skills updated — re-run `grimoire check` to verify compliance:\n", tui.IconWarn)
	for _, line := range diffs {
		fmt.Println(line)
	}
	fmt.Println()
}
