package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var (
	flagListDomain string
	flagListJSON   bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available domains, sub-domains, and skills",
	RunE:  runList,
}

func init() {
	listCmd.Flags().StringVar(&flagListDomain, "domain", "", "filter to a specific domain")
	listCmd.Flags().BoolVar(&flagListJSON, "json", false, "output JSON array")
}

func runList(cmd *cobra.Command, args []string) error {
	root := skills.SkillsRoot()
	if _, err := os.Stat(root); err != nil {
		return fmt.Errorf("skills not found at %s — run: grimoire update", root)
	}

	allSkills, err := skills.ListAllSkills(root)
	if err != nil {
		return fmt.Errorf("listing skills: %w", err)
	}

	if flagListDomain != "" {
		filtered := allSkills[:0]
		for _, s := range allSkills {
			if s.Domain == flagListDomain {
				filtered = append(filtered, s)
			}
		}
		allSkills = filtered
	}

	if flagListJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(allSkills)
	}

	printSkillTree(allSkills)
	return nil
}

func printSkillTree(all []skills.Skill) {
	// group by domain → subdomain → skills
	domains := map[string]map[string][]skills.Skill{}
	domOrder := []string{}
	for _, s := range all {
		if _, ok := domains[s.Domain]; !ok {
			domains[s.Domain] = map[string][]skills.Skill{}
			domOrder = append(domOrder, s.Domain)
		}
		domains[s.Domain][s.Subdomain] = append(domains[s.Domain][s.Subdomain], s)
	}

	for _, d := range domOrder {
		subs := domains[d]
		fmt.Printf("%s\n", tui.StyleGold.Render("▸ "+d))
		for sub, skillList := range subs {
			if sub != "" {
				fmt.Printf("  %s\n", tui.StyleCyan.Render("  "+sub))
			}
			for _, sk := range skillList {
				fmt.Printf("    %s\n", tui.StyleDim.Render(sk.Name))
			}
		}
	}

	total := len(all)
	fmt.Printf("\n%s\n", tui.StyleDim.Render(fmt.Sprintf("%d skills total", total)))
}
