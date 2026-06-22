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
	flagListDomain   string
	flagListJSON     bool
	flagListRegistry string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available domains, sub-domains, and skills",
	RunE:  runList,
}

func init() {
	listCmd.Flags().StringVar(&flagListDomain, "domain", "", "filter to a specific domain")
	listCmd.Flags().BoolVar(&flagListJSON, "json", false, "output JSON array")
	listCmd.Flags().StringVar(&flagListRegistry, "registry", "", "list skills from a specific registry only")
}

func runList(cmd *cobra.Command, args []string) error {
	sources := skills.AllSkillsSources()
	if flagListRegistry != "" {
		sources = filterSources(sources, flagListRegistry)
		if len(sources) == 0 {
			return fmt.Errorf("registry %q not found or not cloned", flagListRegistry)
		}
	}

	if len(sources) == 0 {
		root := skills.SkillsRoot()
		return fmt.Errorf("skills not found at %s — run: grimoire update", root)
	}

	allSkills, conflicts, err := skills.ListAllSkillsFromSources(sources)
	if err != nil {
		return fmt.Errorf("listing skills: %w", err)
	}
	for _, c := range conflicts {
		fmt.Fprintf(os.Stderr, "  %s  %s: %s wins over %s\n",
			tui.IconWarn, c.CanonicalPath, c.WinnerRegistry, c.LoserRegistry)
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

	printSkillTreeMulti(allSkills, len(sources) > 1)
	return nil
}

func printSkillTreeMulti(all []skills.Skill, showRegistry bool) {
	// group by registry → domain → subdomain → skills
	registryOrder := []string{}
	byRegistry := map[string][]skills.Skill{}
	for _, s := range all {
		reg := s.Registry
		if reg == "" {
			reg = skills.OfficialRegistryName
		}
		if _, ok := byRegistry[reg]; !ok {
			registryOrder = append(registryOrder, reg)
		}
		byRegistry[reg] = append(byRegistry[reg], s)
	}

	total := 0
	for _, reg := range registryOrder {
		skillList := byRegistry[reg]
		if showRegistry {
			fmt.Printf("\n%s\n", tui.StyleGold.Render("◈ "+reg))
		}
		printSkillTree(skillList)
		total += len(skillList)
	}

	if showRegistry {
		fmt.Printf("\n%s\n", tui.StyleDim.Render(fmt.Sprintf("%d skills total across %d registries", total, len(registryOrder))))
	}
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

func filterSources(sources []skills.SkillsSource, name string) []skills.SkillsSource {
	for _, s := range sources {
		if s.Name == name {
			return []skills.SkillsSource{s}
		}
	}
	return nil
}
