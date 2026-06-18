package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/agent"
	gitops "github.com/jeffreytse/grimoire/internal/git"
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
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install grimoire skills to AI agent directories",
	RunE:  runInstall,
}

func init() {
	installCmd.Flags().StringVar(&flagInstallDomain, "domain", "", "install all skills for a domain")
	installCmd.Flags().StringVar(&flagInstallSubdomain, "subdomain", "", "restrict to one sub-domain")
	installCmd.Flags().StringVar(&flagInstallSkill, "skill", "", "install one skill (domain/subdomain/name or domain/name)")
	installCmd.Flags().StringVar(&flagInstallTarget, "target", "", "target agent: claude, codex, gemini, openclaw, opencode, all, auto")
	installCmd.Flags().BoolVar(&flagInstallCopy, "copy", false, "copy files instead of symlinking")
	installCmd.Flags().BoolVar(&flagInstallYes, "yes", false, "non-interactive: install all skills to all detected agents")
	installCmd.Flags().BoolVar(&flagInstallNoCfg, "no-configure", false, "skip writing start-best-practice trigger")
	installCmd.Flags().StringVar(&flagInstallFrom, "from", "", "install from a local path or git URL (persisted to ~/.config/grimoire/settings.toml)")
}

func runInstall(cmd *cobra.Command, args []string) error {
	root := skills.SkillsRoot()

	if flagInstallFrom != "" {
		resolved, err := resolveAndPersistSource(flagInstallFrom)
		if err != nil {
			return err
		}
		if resolved == "" {
			return nil // user cancelled
		}
		root = resolved
	}

	if _, err := os.Stat(root); err != nil {
		return fmt.Errorf("skills not found at %s — run: grimoire update", root)
	}

	symlink := !flagInstallCopy
	target := flagInstallTarget
	if flagInstallYes && target == "" {
		target = "auto" // install to all detected agents without prompting
	}
	targets := resolveTargets(target)

	perAgent := make(map[string]int) // ag → newly installed skill count

	switch {
	case flagInstallSkill != "":
		skillPath, err := skills.ResolveSkillPath(root, flagInstallSkill)
		if err != nil {
			return err
		}
		fmt.Printf("Installing skill: %s\n", flagInstallSkill)
		for _, ag := range targets {
			n, err := installSkillToAgent(skillPath, ag, symlink)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  %s  %s: %v\n", tui.IconWarn, agent.DisplayName(ag), err)
			}
			perAgent[ag] += n
		}

	case flagInstallDomain != "":
		for _, ag := range targets {
			n, err := installDomainToAgent(root, flagInstallDomain, flagInstallSubdomain, ag, symlink)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  error: %v\n", err)
			}
			perAgent[ag] += n
		}

	default:
		// install everything
		domains, err := skills.ListDomains(root)
		if err != nil {
			return err
		}
		for _, d := range domains {
			for _, ag := range targets {
				n, err := installDomainToAgent(root, d.Name, "", ag, symlink)
				if err != nil {
					fmt.Fprintf(os.Stderr, "  error: %v\n", err)
				}
				perAgent[ag] += n
			}
		}
	}

	// clean broken symlinks
	for _, ag := range targets {
		_, _ = skills.CleanBrokenSymlinks(agent.SkillsDir(ag))
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
	var installed, skipped []string
	for _, ag := range targets {
		if perAgent[ag] > 0 {
			installed = append(installed, fmt.Sprintf("%s — %d skills",
				agent.DisplayName(ag), perAgent[ag]))
		} else {
			skipped = append(skipped, fmt.Sprintf("%s — already up to date",
				agent.DisplayName(ag)))
		}
	}

	fmt.Printf("\n%s  grimoire installed\n", tui.IconOK)

	if len(installed) > 0 {
		fmt.Println("\n  installed:")
		for _, s := range installed {
			fmt.Printf("    • %s\n", s)
		}
	}
	if len(skipped) > 0 {
		fmt.Println("\n  skipped:")
		for _, s := range skipped {
			fmt.Printf("    • %s\n", s)
		}
	}

	fmt.Println("\n  to get started:")
	fmt.Println("    start any AI session — grimoire skills activate automatically")
	fmt.Println("    or run /start-best-practice in Claude Code to trigger manually")
	fmt.Println("\n  uninstall: grimoire uninstall")
	fmt.Println()
}

// resolveAndPersistSource resolves a --from value (local path or git URL),
// clones if needed, persists to global config, and returns the skills root path.
// Returns ("", nil) when the user cancels.
func resolveAndPersistSource(from string) (string, error) {
	if skills.IsGitURL(from) {
		home := skills.GrimoireHome()
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

func installDomainToAgent(root, domain, subdomain, ag string, symlink bool) (int, error) {
	domainDir := fmt.Sprintf("%s/%s", root, domain)
	if _, err := os.Stat(domainDir); err != nil {
		return 0, fmt.Errorf("domain not found: %s", domain)
	}
	fmt.Printf("Installing domain: %s → %s\n", domain, agent.DisplayName(ag))
	count := 0
	destDir := agent.SkillsDir(ag)
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
			for _, sk := range skillList {
				ok, err := skills.InstallSkill(sk.Path, destDir, symlink)
				if err != nil {
					fmt.Fprintf(os.Stderr, "  warn: %v\n", err)
					continue
				}
				if ok {
					fmt.Printf("  %s %s\n", tui.StyleDim.Render("linked:"), sk.Name)
					count++
				}
			}
		}
	} else {
		skillList, err := skills.ListSkillsInDir(domainDir, domain, "")
		if err != nil {
			return 0, err
		}
		for _, sk := range skillList {
			ok, err := skills.InstallSkill(sk.Path, destDir, symlink)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  warn: %v\n", err)
				continue
			}
			if ok {
				fmt.Printf("  %s %s\n", tui.StyleDim.Render("linked:"), sk.Name)
				count++
			}
		}
	}
	return count, nil
}

func installSkillToAgent(skillPath, ag string, symlink bool) (int, error) {
	destDir := agent.SkillsDir(ag)
	ok, err := skills.InstallSkill(skillPath, destDir, symlink)
	if err != nil {
		return 0, err
	}
	if ok {
		fmt.Printf("  %s %s → %s\n",
			tui.StyleDim.Render("linked:"),
			fmt.Sprintf("skills/%s", splitLast(skillPath, '/')),
			agent.DisplayName(ag))
		return 1, nil
	}
	return 0, nil
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
