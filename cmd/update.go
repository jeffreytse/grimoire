package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

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
	home := skills.OfficialRegistryHome()
	url := skills.GrimoireRepoURL()

	// Local path registry — nothing to pull, just update extends
	if filepath.IsAbs(url) {
		fmt.Printf("  %s  Official registry is a local path — no update needed.\n", tui.IconOK)
		updateCustomRegistries()
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
	if err := gitops.PullWithForceFallback(home); err != nil {
		return fmt.Errorf("updating: %w", err)
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
	changes, err := gitops.RegistryChangesSince(home, oldCommit)
	if err == nil {
		printRegistryChanges(changes, home, oldCommit, os.Stdout)
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
		for _, sk := range allSkills {
			_, _ = skills.InstallSkill(sk.Path, dir, true)
		}
		_, _ = skills.CleanBrokenSymlinks(dir)
	}

	fmt.Println()
}

func printRegistryChanges(c gitops.RegistryChanges, dir, oldCommit string, w io.Writer) {
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
	printChangeSection("presets", c.PresetsAdded, c.PresetsUpdated, w)
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

// updateCustomRegistries pulls or clones all extends targets from resolved settings concurrently.
func updateCustomRegistries() {
	r, err := settings.Load(getProjectDir())
	if err != nil {
		return
	}
	if len(r.StandardsExtends) == 0 {
		return
	}
	fmt.Println()

	n := len(r.StandardsExtends)

	names := make([]string, n)
	for i, ref := range r.StandardsExtends {
		u, _ := settings.ParseRef(ref)
		names[i] = settings.DeriveRegistryName(u)
	}
	board := tui.NewStatusBoard(names)
	stopSpinner := board.StartSpinner()

	limit := 8 // default
	if r.Core.UpdateConcurrency != nil {
		if *r.Core.UpdateConcurrency == 0 {
			limit = n // unlimited
		} else {
			limit = *r.Core.UpdateConcurrency
		}
	}
	if limit > n {
		limit = n
	}
	sem := make(chan struct{}, limit)

	bufs := make([]*bytes.Buffer, n)
	var wg sync.WaitGroup
	for i, ref := range r.StandardsExtends {
		wg.Add(1)
		sem <- struct{}{}
		board.SetUpdating(i)
		go func(i int, ref string) {
			defer wg.Done()
			defer func() { <-sem }()
			buf := &bytes.Buffer{}
			u, ver := settings.ParseRef(ref)
			name := settings.DeriveRegistryName(u)
			dest := skills.ExtendsHome(name)
			ok := true
			if _, err := os.Stat(dest); err != nil {
				if err := os.MkdirAll(dest, 0o755); err != nil {
					fmt.Fprintf(os.Stderr, "  warn: %s: mkdir: %v\n", name, err)
					ok = false
				} else if err := gitops.Clone(u, dest); err != nil {
					fmt.Fprintf(os.Stderr, "  warn: %s: %v\n", name, err)
					ok = false
				} else {
					if ver != "" {
						if err := gitops.CheckoutTag(dest, ver); err != nil {
							fmt.Fprintf(os.Stderr, "  warn: %s: checkout %s: %v\n", name, ver, err)
						}
					}
					fmt.Fprintf(buf, "  %s  %s cloned\n", tui.IconOK, name)
				}
			} else if ver != "" {
				// pinned — ensure correct tag checked out, skip pull
				if err := gitops.CheckoutTag(dest, ver); err != nil {
					fmt.Fprintf(os.Stderr, "  warn: %s: checkout %s: %v\n", name, ver, err)
					ok = false
				} else {
					fmt.Fprintf(buf, "  %s  %s pinned at %s\n", tui.IconOK, name, ver)
				}
			} else {
				oldState, _ := gitops.CurrentState(dest)
				if err := gitops.PullWithForceFallback(dest); err != nil {
					fmt.Fprintf(os.Stderr, "  warn: %s: %v\n", name, err)
					ok = false
				} else {
					fmt.Fprintf(buf, "  %s  %s up to date\n", tui.IconOK, name)
					if changes, err := gitops.RegistryChangesSince(dest, oldState.Commit); err == nil {
						printRegistryChanges(changes, dest, oldState.Commit, buf)
					}
				}
			}
			if ok {
				board.SetDone(i, tui.IconDone)
			} else {
				board.SetDone(i, tui.IconError)
			}
			bufs[i] = buf
		}(i, ref)
	}
	wg.Wait()
	stopSpinner()
	board.Finish()

	for _, buf := range bufs {
		if buf != nil {
			os.Stdout.Write(buf.Bytes())
		}
	}
}

func errNotGit(dir string) error {
	return fmt.Errorf("%s is not a git repository", dir)
}
