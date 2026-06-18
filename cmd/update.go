package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/agent"
	gitops "github.com/jeffreytse/grimoire/internal/git"
	"github.com/jeffreytse/grimoire/internal/settings"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var (
	flagUpdateStable bool
	flagUpdateYes    bool
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"upgrade"},
	Short:   "Pull the latest grimoire skills and relink",
	RunE:    runUpdate,
}

func init() {
	updateCmd.Flags().BoolVar(&flagUpdateStable, "stable", false, "check out the latest tagged release instead of HEAD")
	updateCmd.Flags().BoolVar(&flagUpdateYes, "yes", false, "skip confirmation prompts")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	home := skills.GrimoireHome()

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
		if err := gitops.Clone(skills.GrimoireRepoURL(), home); err != nil {
			return fmt.Errorf("cloning grimoire: %w", err)
		}
		fmt.Printf("%s  Grimoire installed to %s\n", tui.IconOK, home)
		return nil
	}

	if flagUpdateStable {
		return updateStable(home)
	}
	return updateUnstable(home)
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

	fmt.Printf("Checking out %s...\n", latest)
	if err := gitops.CheckoutTag(home, latest); err != nil {
		return fmt.Errorf("checking out %s: %w", latest, err)
	}

	printUpgradeResult(current, tagState)
	relinkNewSkills(home, current.Commit)
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

	fmt.Printf("Pulling latest grimoire at %s...\n", home)
	if err := gitops.Pull(home); err != nil {
		return fmt.Errorf("pulling: %w", err)
	}

	newState, _ := gitops.CurrentState(home)
	printUpgradeResult(current, newState)
	relinkNewSkills(home, current.Commit)
	updateCustomRegistries()
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
	added, updated, err := gitops.NewSkillsSince(home, oldCommit)
	if err == nil {
		if added > 0 {
			fmt.Printf("  New skills:     %d\n", added)
		}
		if updated > 0 {
			fmt.Printf("  Updated skills: %d\n", updated)
		}
		if added == 0 && updated == 0 {
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
		for _, sk := range allSkills {
			_, _ = skills.InstallSkill(sk.Path, dir, true)
		}
		_, _ = skills.CleanBrokenSymlinks(dir)
	}

	fmt.Println()
}

// updateCustomRegistries pulls or clones all custom registries configured in global settings.
func updateCustomRegistries() {
	fs, err := settings.LoadGlobal()
	if err != nil || len(fs.Registries) == 0 {
		return
	}
	fmt.Println()
	for name, rc := range fs.Registries {
		if name == skills.OfficialRegistryName {
			continue
		}
		dest := skills.RegistryHome(name)
		if _, err := os.Stat(dest); err != nil {
			fmt.Printf("  cloning registry %s...\n", name)
			if err := gitops.Clone(rc.URL, dest); err != nil {
				fmt.Fprintf(os.Stderr, "  warn: %s: %v\n", name, err)
			} else {
				fmt.Printf("  %s  %s cloned\n", tui.IconOK, name)
			}
		} else {
			if err := gitops.Pull(dest); err != nil {
				fmt.Fprintf(os.Stderr, "  warn: %s: %v\n", name, err)
			} else {
				fmt.Printf("  %s  %s up to date\n", tui.IconOK, name)
			}
		}
	}
}

func errNotGit(dir string) error {
	return fmt.Errorf("%s is not a git repository", dir)
}
