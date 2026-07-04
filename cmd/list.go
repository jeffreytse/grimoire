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
	flagListDomain  string
	flagListJSON    bool
	flagListPackage string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available domains, sub-domains, and skills",
	RunE:  runList,
}

func init() {
	listCmd.Flags().StringVar(&flagListDomain, "domain", "", "filter to a specific domain")
	listCmd.Flags().BoolVar(&flagListJSON, "json", false, "output JSON array")
	listCmd.Flags().StringVar(&flagListPackage, "package", "", "list skills from a specific package only")
}

func runList(cmd *cobra.Command, args []string) error {
	regs := skills.AllSkillsPackages()
	if flagListPackage != "" {
		regs = filterPackages(regs, flagListPackage)
		if len(regs) == 0 {
			return fmt.Errorf("package %q not found or not cloned", flagListPackage)
		}
	}

	if len(regs) == 0 {
		root := skills.SkillsRoot()
		return fmt.Errorf("skills not found at %s — run: grimoire update", root)
	}

	allSkills, conflicts, err := skills.ListAllSkillsFromPackages(regs)
	if err != nil {
		return fmt.Errorf("listing skills: %w", err)
	}
	for _, c := range conflicts {
		fmt.Fprintf(os.Stderr, "  %s  %s: %s wins over %s\n",
			tui.IconWarn, c.CanonicalPath, c.WinnerPackage, c.LoserPackage)
	}

	if flagListDomain != "" {
		filtered := allSkills[:0]
		for i := range allSkills {
			s := allSkills[i]
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

	printSkillTreeMulti(allSkills, len(regs) > 1)
	return nil
}

func printSkillTreeMulti(all []skills.Skill, showPackage bool) {
	// group by package → domain → subdomain → skills
	packageOrder := []string{}
	byPackage := map[string][]skills.Skill{}
	for i := range all {
		s := all[i]
		pkg := s.Package
		if pkg == "" {
			pkg = skills.OfficialPackageDerivedName()
		}
		if _, ok := byPackage[pkg]; !ok {
			packageOrder = append(packageOrder, pkg)
		}
		byPackage[pkg] = append(byPackage[pkg], s)
	}

	total := 0
	for _, pkg := range packageOrder {
		skillList := byPackage[pkg]
		if showPackage {
			fmt.Printf("\n%s\n", tui.StyleGold.Render("◈ "+pkg))
		}
		printSkillTree(skillList)
		total += len(skillList)
	}

	if showPackage {
		fmt.Printf("\n%s\n", tui.StyleDim.Render(fmt.Sprintf("%d skills total across %d packages", total, len(packageOrder))))
	}
}

func printSkillTree(all []skills.Skill) {
	// group by domain → subdomain → skills
	domains := map[string]map[string][]skills.Skill{}
	domOrder := []string{}
	for i := range all {
		s := all[i]
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
			for i := range skillList {
				sk := skillList[i]
				fmt.Printf("    %s\n", tui.StyleDim.Render(sk.Name))
			}
		}
	}

	total := len(all)
	fmt.Printf("\n%s\n", tui.StyleDim.Render(fmt.Sprintf("%d skills total", total)))
}

func filterPackages(regs []skills.SkillsPackage, name string) []skills.SkillsPackage {
	for _, r := range regs {
		if r.Name == name {
			return []skills.SkillsPackage{r}
		}
	}
	return nil
}
