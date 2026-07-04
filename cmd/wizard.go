package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func runWizard() error {
	home := skills.OfficialPackageHome()

	for {
		clearScreen()
		tui.PrintBanner(skills.GrimoireVersion())

		var statusParts []string
		for _, ag := range agent.Detected() {
			n := agent.SkillCount(ag)
			if n > 0 {
				statusParts = append(statusParts, fmt.Sprintf("%s (%d skills)", agent.DisplayName(ag), n))
			}
		}
		if len(statusParts) > 0 {
			fmt.Printf("  %s\n\n", tui.StyleDim.Render("Active: "+strings.Join(statusParts, " · ")))
		} else {
			fmt.Printf("  %s\n\n", tui.StyleDim.Render("No skills installed yet — select Install to get started"))
		}

		mode, ok := tui.RunSelect("⚙️  What would you like to do?",
			[]string{"📥 Install", "🚀 Update", "📦 Packages", "🩺 Doctor", "🗑  Uninstall", "📋 Init project", "🚪 Exit"})
		if !ok || mode == "🚪 Exit" {
			return nil
		}

		if mode == "📋 Init project" {
			fmt.Println()
			if err := runInit(nil, nil); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			}
			tui.RunSelect("", []string{"← Back"})
			continue
		}
		if mode == "📦 Packages" {
			if err := runPackageWizard(); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			}
			continue
		}
		if mode == "🩺 Doctor" {
			fmt.Println()
			if err := runDoctor(nil, nil); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			}
			tui.RunSelect("", []string{"← Back"})
			continue
		}
		if mode == "🚀 Update" {
			fmt.Println()
			if err := runUpdate(nil, nil); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			}
			tui.RunSelect("", []string{"← Back"})
			continue
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
			continue
		}
		selectedAgents := make([]string, len(agentDisplays))
		for i, d := range agentDisplays {
			selectedAgents[i] = agent.FromDisplayName(d)
		}

		// Download guard: check official package exists
		if _, err := os.Stat(home); err != nil {
			choice, ok := tui.RunSelect(
				"⚠️  Skills not downloaded yet. Download now?",
				[]string{"Yes, download now", "Cancel"},
			)
			if !ok || choice == "Cancel" {
				continue
			}
			flagUpdateYes = true
			err2 := runUpdate(nil, nil)
			flagUpdateYes = false
			if err2 != nil {
				return fmt.Errorf("update failed: %w", err2)
			}
		}

		// Pick packages (skip step when only 1 available)
		regs := skills.AllPackages()
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
			chosenRegs, ok := tui.RunMultiselect("📦 Which packages?", regNames, presel)
			if !ok || len(chosenRegs) == 0 {
				continue
			}
			chosenSet := make(map[string]bool, len(chosenRegs))
			for _, n := range chosenRegs {
				chosenSet[n] = true
			}
			var filtered []skills.PackageEntry
			for _, r := range regs {
				if chosenSet[r.Name] {
					filtered = append(filtered, r)
				}
			}
			selectedRegs = filtered
		}

		// Collect domains from all selected packages
		type domainItem struct {
			domain     skills.Domain
			regName    string
			skillsRoot string
			pkgName    string // "" for official package, reg.Name for third-party
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
			pkgName := reg.Name
			if reg.Official {
				pkgName = ""
			}
			for _, d := range domains {
				domainItems = append(domainItems, domainItem{domain: d, regName: reg.Name, skillsRoot: sr, pkgName: pkgName})
			}
		}

		// Build display names, annotating with package when >1 selected
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

		// Pre-select all domains by default (consistent with subdomain behavior)
		domainPresel := make([]bool, len(displayNames))
		for i := range domainPresel {
			domainPresel[i] = true
		}
		fmt.Printf("  %s\n\n", tui.StyleDim.Render("Domains group related skills (e.g. engineering, security). All pre-selected — deselect any you don't need."))
		selectedDisplays, ok := tui.RunMultiselect("📚 Which domains?", displayNames, domainPresel)
		if !ok || len(selectedDisplays) == 0 {
			continue
		}
		var selectedItems []domainItem
		for _, d := range selectedDisplays {
			selectedItems = append(selectedItems, displayToItem[d])
		}

		// For nested domains, pick subdomains
		type domainSelection struct {
			name    string
			root    string   // skillsRoot for this domain's package
			regName string   // for summary display
			pkgName string   // "" for official, reg.Name for third-party
			subs    []string // empty = all
		}
		var domainSelections []domainSelection

		for _, di := range selectedItems {
			if di.domain.Nested {
				subs, err := skills.ListSubdomains(di.domain.Path)
				if err != nil || len(subs) == 0 {
					domainSelections = append(domainSelections, domainSelection{name: di.domain.Name, root: di.skillsRoot, regName: di.regName, pkgName: di.pkgName})
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
					continue
				}
				domainSelections = append(domainSelections, domainSelection{name: di.domain.Name, root: di.skillsRoot, regName: di.regName, pkgName: di.pkgName, subs: chosen})
			} else {
				domainSelections = append(domainSelections, domainSelection{name: di.domain.Name, root: di.skillsRoot, regName: di.regName, pkgName: di.pkgName})
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

		// Confirmation before execute
		confirmLabel := "Install"
		if mode == "🗑  Uninstall" {
			confirmLabel = "Uninstall"
		}
		confirm, ok := tui.RunSelect(
			fmt.Sprintf("Proceed with %s?", confirmLabel),
			[]string{"Yes", "Cancel"},
		)
		if !ok || confirm == "Cancel" {
			continue
		}

		// Execute
		flagInstallGlobal = true
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
			tui.RunSelect("", []string{"← Back"})
			continue
		} else {
			for _, ds := range domainSelections {
				for _, ag := range selectedAgents {
					if len(ds.subs) > 0 {
						for _, s := range ds.subs {
							n, err := installDomainToAgent(ds.root, ds.name, s, ag, symlink, ds.pkgName)
							if err != nil {
								fmt.Fprintf(os.Stderr, "  error: %v\n", err)
							}
							totalCount += n
						}
					} else {
						n, err := installDomainToAgent(ds.root, ds.name, "", ag, symlink, ds.pkgName)
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

			fmt.Println("\n  to get started:")
			fmt.Println("    start any AI session — skills activate automatically")
			fmt.Println("    or run /start-best-practice in Claude Code to trigger manually")
			fmt.Printf("\n  project setup: %s\n", tui.StyleCyan.Render("grimoire init"))
			fmt.Printf("    %s\n", tui.StyleDim.Render("creates grimoire.toml — tracks skill deps, enables compliance checks"))

			// marketplace tip (install only)
			fmt.Println()
			fmt.Printf("%s Also available via marketplace:\n", tui.StyleBold.Render("💡"))
			fmt.Printf("   %s  %s\n",
				tui.StyleGold.Render("🐙 Copilot"),
				tui.StyleCyan.Render("copilot plugin marketplace add jeffreytse/grimoire"))
			fmt.Printf("              %s\n",
				tui.StyleCyan.Render("copilot plugin install grimoire@grimoire"))
			fmt.Println()
		}
		flagInstallGlobal = false
		tui.RunSelect("", []string{"← Back"})
	}
}

func runPackageWizard() error {
	r := bufio.NewReader(os.Stdin)
	for {
		clearScreen()
		regs := skills.AllPackages()

		items := []string{tui.ProfileSectionPrefix + "Installed packages"}
		displayToEntry := make(map[string]skills.PackageEntry)
		if len(regs) == 0 {
			items = append(items, tui.ProfileSectionPrefix+"(no packages installed)")
		} else {
			for _, reg := range regs {
				name := reg.Name
				if parts := strings.Split(name, "/"); len(parts) > 1 {
					name = parts[len(parts)-1]
				}
				if reg.Official {
					name += "  (official)"
				}
				items = append(items, name)
				displayToEntry[name] = reg
			}
		}
		items = append(items, tui.ProfileSectionPrefix+" ")
		items = append(items, "Add package", "Done")

		action, ok := tui.RunSelectScrollable("📦 Manage packages", items, 8)
		if !ok || action == "Done" {
			return nil
		}
		if action != "Add package" {
			if reg, found := displayToEntry[action]; found {
				showPackageDetail(reg)
			}
			continue
		}

		fmt.Print("  Package name (e.g. my-team): ")
		name := strings.TrimSpace(readLine(r))
		if name == "" {
			continue
		}
		fmt.Print("  Package URL (git URL or owner/repo): ")
		url := strings.TrimSpace(readLine(r))
		if url == "" {
			continue
		}

		if err := runPackageAdd(nil, []string{name, url}); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
	}
}

func showPackageDetail(reg skills.PackageEntry) {
	skillsRoot := filepath.Join(reg.Home, "skills")

	source := reg.Name
	if !filepath.IsAbs(source) {
		if idx := strings.LastIndex(source, "@"); idx >= 0 {
			source = source[:idx]
		}
		source = "https://" + source
	}

	domains, _ := skills.ListDomains(skillsRoot)
	domainNames := make([]string, len(domains))
	for i, d := range domains {
		domainNames[i] = d.Name
	}
	domainStr := strings.Join(domainNames, ", ")
	if domainStr == "" {
		domainStr = "(none)"
	}

	allSkills, _ := skills.WalkSkills(skillsRoot)

	displayName := reg.Name
	if parts := strings.Split(displayName, "/"); len(parts) > 1 {
		displayName = parts[len(parts)-1]
	}
	official := ""
	if reg.Official {
		official = "  (official)"
	}

	fmt.Printf("\n  %s%s\n", tui.StyleBold.Render(displayName), tui.StyleDim.Render(official))
	fmt.Printf("  %s\n\n", strings.Repeat("─", 40))
	fmt.Printf("  %-10s %s\n", tui.StyleDim.Render("Source:"), source)
	fmt.Printf("  %-10s %d skills across %d domains\n", tui.StyleDim.Render("Skills:"), len(allSkills), len(domains))
	fmt.Printf("  %-10s %s\n", tui.StyleDim.Render("Domains:"), domainStr)
	fmt.Printf("  %-10s %s\n\n", tui.StyleDim.Render("Path:"), tui.StyleDim.Render(reg.Home))

	actions := []string{"Update package", "← Back"}
	if !reg.Official {
		actions = []string{"Update package", "Remove package", "← Back"}
	}
	action, _ := tui.RunSelect("", actions)
	switch action {
	case "Update package":
		fmt.Println()
		if err := runPackageUpdate(nil, []string{reg.Name}); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
		tui.RunSelect("", []string{"← Back"})
	case "Remove package":
		confirm, ok := tui.RunSelect(
			fmt.Sprintf("Remove package %q?", displayName),
			[]string{"Yes, remove", "Cancel"},
		)
		if ok && confirm == "Yes, remove" {
			if err := runPackageRemove(nil, []string{reg.Name}); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			}
		}
	}
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
