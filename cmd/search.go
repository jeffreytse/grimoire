package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search installed packages for skills matching a query",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func runSearch(_ *cobra.Command, args []string) error {
	regs := skills.AllSkillsPackages()
	if len(regs) == 0 {
		return fmt.Errorf("no packages installed — run: grimoire update")
	}

	all, _, err := skills.ListAllSkillsFromPackages(regs)
	if err != nil {
		return err
	}

	query := ""
	if len(args) > 0 {
		query = strings.ToLower(strings.TrimSpace(args[0]))
	}

	var matches []skills.Skill
	for i := range all {
		if query == "" || matchesQuery(&all[i], query) {
			matches = append(matches, all[i])
		}
	}

	if len(matches) == 0 {
		if query != "" {
			fmt.Printf("%s  no skills match %q\n", tui.IconWarn, query)
		} else {
			fmt.Printf("%s  no skills found\n", tui.IconWarn)
		}
		return nil
	}

	fmt.Printf("Found %d skill(s):\n\n", len(matches))
	for i := range matches {
		printSkillSummary(&matches[i])
	}
	return nil
}

func matchesQuery(sk *skills.Skill, query string) bool {
	if strings.Contains(strings.ToLower(sk.Name), query) {
		return true
	}
	if strings.Contains(strings.ToLower(sk.Description), query) {
		return true
	}
	if strings.Contains(strings.ToLower(sk.Domain), query) {
		return true
	}
	for _, tag := range sk.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}
	return false
}

func printSkillSummary(sk *skills.Skill) {
	name := sk.Name
	if sk.Version != "" {
		name += " v" + sk.Version
	}
	fmt.Printf("  %s\n", tui.StyleBold.Render(name))
	if sk.Description != "" {
		fmt.Printf("    %s\n", sk.Description)
	}
	if sk.Domain != "" {
		domain := sk.Domain
		if sk.Subdomain != "" {
			domain += "/" + sk.Subdomain
		}
		fmt.Printf("    domain: %s\n", domain)
	}
	if len(sk.Tags) > 0 {
		fmt.Printf("    tags: %s\n", strings.Join(sk.Tags, ", "))
	}
	fmt.Println()
}
