package cmd

import (
	"fmt"
	"os"

	"github.com/jeffreytse/grimoire/internal/agent"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

func runInteractive() error {
	home := skills.GrimoireHome()
	tui.PrintBanner(skills.GrimoireVersion())

	mode, ok := tui.RunSelect("⚙️  What would you like to do?",
		[]string{"📥 Install", "🗑  Uninstall", "🚀 Update", "🩺 Doctor", "🚪 Exit"})
	if !ok || mode == "🚪 Exit" {
		return nil
	}

	if mode == "🩺 Doctor" {
		return runDoctor(nil, nil)
	}
	if mode == "🚀 Update" {
		return runUpdate(nil, nil)
	}

	// Install / Uninstall — pick agents
	detected := agent.Detected()
	var agentOptions []string
	if len(detected) == 0 {
		agentOptions = []string{"Claude Code"}
	} else {
		for _, ag := range detected {
			agentOptions = append(agentOptions, agent.DisplayName(ag))
		}
	}

	promptVerb := "install to"
	if mode == "🗑  Uninstall" {
		promptVerb = "uninstall from"
	}
	agentDisplays, ok := tui.RunMultiselect(
		fmt.Sprintf("🤖 Which agents to %s?", promptVerb),
		agentOptions, nil)
	if !ok || len(agentDisplays) == 0 {
		fmt.Println("No agents selected. Exiting.")
		return nil
	}
	selectedAgents := make([]string, len(agentDisplays))
	for i, d := range agentDisplays {
		selectedAgents[i] = agent.FromDisplayName(d)
	}

	// Pick domains
	root := skills.SkillsRoot()
	if _, err := os.Stat(root); err != nil {
		choice, ok := tui.RunSelect(
			"⚠️  Skills not downloaded yet. Download now?",
			[]string{"Yes, download now", "Cancel"},
		)
		if !ok || choice == "Cancel" {
			fmt.Println("Run: grimoire update")
			return nil
		}
		if err := runUpdate(nil, nil); err != nil {
			return fmt.Errorf("update failed: %w", err)
		}
	}
	domains, err := skills.ListDomains(root)
	if err != nil {
		return fmt.Errorf("listing domains: %w", err)
	}
	domainNames := make([]string, len(domains))
	for i, d := range domains {
		domainNames[i] = d.Name
	}
	selectedDomainNames, ok := tui.RunMultiselect(
		"📚 Which domains?", domainNames, nil)
	if !ok || len(selectedDomainNames) == 0 {
		fmt.Println("No domains selected. Exiting.")
		return nil
	}

	// For nested domains, pick subdomains
	type domainSelection struct {
		name string
		subs []string // empty = all
	}
	var domainSelections []domainSelection

	for _, dName := range selectedDomainNames {
		dPath := fmt.Sprintf("%s/%s", root, dName)
		if skills.IsNested(dPath) {
			subs, err := skills.ListSubdomains(dPath)
			if err != nil || len(subs) == 0 {
				domainSelections = append(domainSelections, domainSelection{name: dName})
				continue
			}
			subNames := make([]string, len(subs))
			for i, s := range subs {
				subNames[i] = s.Name
			}
			// pre-select all
			presel := make([]bool, len(subNames))
			for i := range presel {
				presel[i] = true
			}
			chosen, ok := tui.RunMultiselect(
				fmt.Sprintf("📂 %s: which sub-domains?", dName), subNames, presel)
			if !ok {
				fmt.Println("Cancelled.")
				return nil
			}
			domainSelections = append(domainSelections, domainSelection{name: dName, subs: chosen})
		} else {
			domainSelections = append(domainSelections, domainSelection{name: dName})
		}
	}

	// Summary preview
	fmt.Println()
	fmt.Printf("  %s  %s\n", tui.StyleBold.Render("Agents:"), joinAgentNames(selectedAgents))
	for _, ds := range domainSelections {
		if len(ds.subs) > 0 {
			fmt.Printf("  %s  %s [%s]\n", tui.StyleBold.Render("Domain:"),
				tui.StyleGold.Render(ds.name),
				tui.StyleDim.Render(joinStrings(ds.subs)))
		} else {
			fmt.Printf("  %s  %s\n", tui.StyleBold.Render("Domain:"),
				tui.StyleGold.Render(ds.name))
		}
	}
	fmt.Println()

	// Execute
	symlink := true
	totalCount := 0

	if mode == "🗑  Uninstall" {
		for _, ds := range domainSelections {
			for _, ag := range selectedAgents {
				sub := ""
				if len(ds.subs) == 1 {
					sub = ds.subs[0]
				}
				// uninstall each selected sub
				if len(ds.subs) > 1 {
					for _, s := range ds.subs {
						n, err := uninstallDomainFromAgent(root, ds.name, s, ag)
						if err != nil {
							fmt.Fprintf(os.Stderr, "  error: %v\n", err)
						}
						totalCount += n
					}
				} else {
					n, err := uninstallDomainFromAgent(root, ds.name, sub, ag)
					if err != nil {
						fmt.Fprintf(os.Stderr, "  error: %v\n", err)
					}
					totalCount += n
				}
			}
		}
		for _, ag := range selectedAgents {
			_, _ = skills.CleanBrokenSymlinks(agent.SkillsDir(ag))
			if agent.SkillCount(ag) == 0 {
				_ = agent.RemoveAgentMDConfig(ag)
			}
		}
		fmt.Printf("\n%s  %d skills uninstalled → %s\n",
			tui.IconOK, totalCount/max(len(selectedAgents), 1), joinAgentNames(selectedAgents))
	} else {
		for _, ds := range domainSelections {
			for _, ag := range selectedAgents {
				if len(ds.subs) > 0 {
					for _, s := range ds.subs {
						n, err := installDomainToAgent(root, ds.name, s, ag, symlink)
						if err != nil {
							fmt.Fprintf(os.Stderr, "  error: %v\n", err)
						}
						totalCount += n
					}
				} else {
					n, err := installDomainToAgent(root, ds.name, "", ag, symlink)
					if err != nil {
						fmt.Fprintf(os.Stderr, "  error: %v\n", err)
					}
					totalCount += n
				}
			}
		}
		for _, ag := range selectedAgents {
			_, _ = skills.CleanBrokenSymlinks(agent.SkillsDir(ag))
			_ = agent.ConfigureAgentMD(ag)
		}

		_ = home
		fmt.Printf("\n%s  %d skills installed → %s\n",
			tui.IconOK, totalCount/max(len(selectedAgents), 1), joinAgentNames(selectedAgents))
	}

	// marketplace tip
	fmt.Println()
	fmt.Printf("%s Also available via marketplace:\n", tui.StyleBold.Render("💡"))
	fmt.Printf("   %s  %s\n",
		tui.StyleGold.Render("🐙 Copilot"),
		tui.StyleCyan.Render("copilot plugin marketplace add jeffreytse/grimoire"))
	fmt.Printf("              %s\n",
		tui.StyleCyan.Render("copilot plugin install grimoire@grimoire"))
	fmt.Println()

	return nil
}

func joinStrings(ss []string) string {
	result := ""
	for i, s := range ss {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}
