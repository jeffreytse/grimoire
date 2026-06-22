package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/agent"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var wizardCmd = &cobra.Command{
	Use:   "wizard",
	Short: "Open the interactive TUI wizard",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWizard()
	},
}

func runWizard() error {
	home := skills.OfficialRegistryHome()
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

	// Download guard: check official registry exists
	if _, err := os.Stat(home); err != nil {
		choice, ok := tui.RunSelect(
			"⚠️  Skills not downloaded yet. Download now?",
			[]string{"Yes, download now", "Cancel"},
		)
		if !ok || choice == "Cancel" {
			fmt.Println("Run: grimoire update")
			return nil
		}
		flagUpdateYes = true
		err2 := runUpdate(nil, nil)
		flagUpdateYes = false
		if err2 != nil {
			return fmt.Errorf("update failed: %w", err2)
		}
	}

	// Pick registries (skip step when only 1 available)
	regs := skills.AllRegistries()
	selectedRegs := regs
	if len(regs) > 1 {
		regNames := make([]string, len(regs))
		for i, r := range regs {
			regNames[i] = r.Name
		}
		presel := make([]bool, len(regs))
		if len(presel) > 0 {
			presel[0] = true // official is always first; pre-select it
		}
		chosenRegs, ok := tui.RunMultiselect("📦 Which registries?", regNames, presel)
		if !ok || len(chosenRegs) == 0 {
			fmt.Println("No registries selected. Exiting.")
			return nil
		}
		chosenSet := make(map[string]bool, len(chosenRegs))
		for _, n := range chosenRegs {
			chosenSet[n] = true
		}
		var filtered []skills.RegistryEntry
		for _, r := range regs {
			if chosenSet[r.Name] {
				filtered = append(filtered, r)
			}
		}
		selectedRegs = filtered
	}

	// Collect domains from all selected registries
	type domainItem struct {
		domain     skills.Domain
		regName    string
		skillsRoot string
	}
	var domainItems []domainItem
	for _, reg := range selectedRegs {
		sr := filepath.Join(reg.Home, "skills")
		if _, err := os.Stat(sr); err != nil {
			continue
		}
		domains, err := skills.ListDomains(sr)
		if err != nil {
			continue
		}
		for _, d := range domains {
			domainItems = append(domainItems, domainItem{domain: d, regName: reg.Name, skillsRoot: sr})
		}
	}

	// Build display names, annotating with registry when >1 selected
	displayNames := make([]string, len(domainItems))
	for i, di := range domainItems {
		if len(selectedRegs) > 1 {
			displayNames[i] = fmt.Sprintf("%s  %s", di.domain.Name, tui.StyleDim.Render("["+di.regName+"]"))
		} else {
			displayNames[i] = di.domain.Name
		}
	}
	displayToItem := make(map[string]domainItem, len(domainItems))
	for i, di := range domainItems {
		displayToItem[displayNames[i]] = di
	}

	selectedDisplays, ok := tui.RunMultiselect("📚 Which domains?", displayNames, nil)
	if !ok || len(selectedDisplays) == 0 {
		fmt.Println("No domains selected. Exiting.")
		return nil
	}
	var selectedItems []domainItem
	for _, d := range selectedDisplays {
		selectedItems = append(selectedItems, displayToItem[d])
	}

	// For nested domains, pick subdomains
	type domainSelection struct {
		name    string
		root    string   // skillsRoot for this domain's registry
		regName string   // for summary display
		subs    []string // empty = all
	}
	var domainSelections []domainSelection

	for _, di := range selectedItems {
		if di.domain.Nested {
			subs, err := skills.ListSubdomains(di.domain.Path)
			if err != nil || len(subs) == 0 {
				domainSelections = append(domainSelections, domainSelection{name: di.domain.Name, root: di.skillsRoot, regName: di.regName})
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
			domainLabel := di.domain.Name
			if len(selectedRegs) > 1 {
				domainLabel = fmt.Sprintf("%s %s", di.domain.Name, tui.StyleDim.Render("["+di.regName+"]"))
			}
			chosen, ok := tui.RunMultiselect(
				fmt.Sprintf("📂 %s: which sub-domains?", domainLabel), subNames, presel)
			if !ok {
				fmt.Println("Cancelled.")
				return nil
			}
			domainSelections = append(domainSelections, domainSelection{name: di.domain.Name, root: di.skillsRoot, regName: di.regName, subs: chosen})
		} else {
			domainSelections = append(domainSelections, domainSelection{name: di.domain.Name, root: di.skillsRoot, regName: di.regName})
		}
	}

	// Summary preview
	fmt.Println()
	fmt.Printf("  %s  %s\n", tui.StyleBold.Render("Agents:"), joinAgentNames(selectedAgents))
	for _, ds := range domainSelections {
		regAnnotation := ""
		if len(selectedRegs) > 1 {
			regAnnotation = "  " + tui.StyleDim.Render("["+ds.regName+"]")
		}
		if len(ds.subs) > 0 {
			fmt.Printf("  %s  %s [%s]%s\n", tui.StyleBold.Render("Domain:"),
				tui.StyleGold.Render(ds.name),
				tui.StyleDim.Render(joinStrings(ds.subs)),
				regAnnotation)
		} else {
			fmt.Printf("  %s  %s%s\n", tui.StyleBold.Render("Domain:"),
				tui.StyleGold.Render(ds.name),
				regAnnotation)
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
						n, err := uninstallDomainFromAgent(ds.root, ds.name, s, ag)
						if err != nil {
							fmt.Fprintf(os.Stderr, "  error: %v\n", err)
						}
						totalCount += n
					}
				} else {
					n, err := uninstallDomainFromAgent(ds.root, ds.name, sub, ag)
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
						n, err := installDomainToAgent(ds.root, ds.name, s, ag, symlink)
						if err != nil {
							fmt.Fprintf(os.Stderr, "  error: %v\n", err)
						}
						totalCount += n
					}
				} else {
					n, err := installDomainToAgent(ds.root, ds.name, "", ag, symlink)
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
