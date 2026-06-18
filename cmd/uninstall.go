package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/agent"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var (
	flagUninstallDomain    string
	flagUninstallSubdomain string
	flagUninstallSkill     string
	flagUninstallTarget    string
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove grimoire skills from agent directories",
	RunE:  runUninstall,
}

func init() {
	uninstallCmd.Flags().StringVar(&flagUninstallDomain, "domain", "", "uninstall all skills for a domain")
	uninstallCmd.Flags().StringVar(&flagUninstallSubdomain, "subdomain", "", "restrict to one sub-domain")
	uninstallCmd.Flags().StringVar(&flagUninstallSkill, "skill", "", "uninstall one skill (domain/subdomain/name or domain/name)")
	uninstallCmd.Flags().StringVar(&flagUninstallTarget, "target", "", "target agent (default: all detected)")
}

func runUninstall(cmd *cobra.Command, args []string) error {
	root := skills.SkillsRoot()
	targets := resolveTargets(flagUninstallTarget)
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
			for _, sk := range skillList {
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
		for _, sk := range skillList {
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
