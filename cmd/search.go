package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/search"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var searchCmd = &cobra.Command{
	Use:   "search [<query>...]",
	Short: "Search installed packages for skills matching a query",
	Long: `Search skills the way your agent does over MCP: results are ranked by
relevance across skill names, tags, descriptions and full text, matching
word forms ("testing" finds "test"). Several queries are different phrasings
of one need and are ranked together. Your project's declared practices
come first among comparable matches.

With no query, every skill is listed.`,
	Example: `  grimoire search "flaky tests"
  grimoire search "release checklist" "go-live readiness"
  grimoire search incident --domain engineering --top 5 --json`,
	RunE: runSearch,
}

var (
	flagSearchJSON   bool
	flagSearchTop    int
	flagSearchDomain string
)

func init() {
	searchCmd.Flags().BoolVar(&flagSearchJSON, "json", false, "output JSON array, best match first (same shape as grimoire list --json)")
	searchCmd.Flags().IntVar(&flagSearchTop, "top", 10, "how many results to show (1-100)")
	searchCmd.Flags().StringVar(&flagSearchDomain, "domain", "", "only search one domain")
	rootCmd.AddCommand(searchCmd)
}

func runSearch(_ *cobra.Command, args []string) error {
	regs := skills.AllSkillsPackages()
	if len(regs) == 0 {
		return errNoSkillLibrary
	}
	all, _, err := skills.ListAllSkillsFromPackages(regs)
	if err != nil {
		return err
	}
	if flagSearchTop < 1 || flagSearchTop > 100 {
		return fmt.Errorf("--top must be between 1 and 100, got %d", flagSearchTop)
	}

	var queries []string
	for _, a := range args {
		if q := strings.TrimSpace(a); q != "" {
			queries = append(queries, q)
		}
	}
	scope := projectScope(getProjectDir())
	var matches []skills.Skill
	if len(queries) == 0 {
		pool := scope.filter(all)
		for i := range pool {
			if flagSearchDomain == "" || pool[i].Domain == flagSearchDomain {
				matches = append(matches, pool[i])
			}
		}
	} else {
		hits := rankSkills(queries, all, scope, search.Options{Domain: flagSearchDomain}, 0, flagSearchTop)
		for i := range hits {
			matches = append(matches, hits[i].Skill)
		}
	}

	if flagSearchJSON {
		if matches == nil {
			matches = []skills.Skill{}
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(matches)
	}

	if len(matches) == 0 {
		if len(queries) > 0 {
			fmt.Printf("%s  no skills match %q\n", tui.IconWarn, strings.Join(queries, `", "`))
		} else {
			fmt.Printf("%s  no skills found\n", tui.IconWarn)
		}
		return nil
	}

	if len(queries) > 0 {
		fmt.Printf("Top %d match(es), best first:\n\n", len(matches))
	} else {
		fmt.Printf("Found %d skill(s):\n\n", len(matches))
	}
	for i := range matches {
		printSkillSummary(&matches[i], scope.standard(matches[i].Name))
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

// printSkillSummary prints one search result; standard is where the skill
// was declared as a practice ("project", "user", "org"), or "".
func printSkillSummary(sk *skills.Skill, standard string) {
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
	if standard != "" {
		fmt.Printf("    declared standard: %s\n", standard)
	}
	fmt.Println()
}
