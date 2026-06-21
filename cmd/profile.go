package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/profiles"
	"github.com/jeffreytse/grimoire/internal/settings"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage grimoire profiles",
	Long: `Profiles group related skills under a named paradigm.

  grimoire profile list          List available profiles
  grimoire profile show <name>   Show skills in a profile
  grimoire profile init <name>   Create a profile file in .grimoire/profiles/`,
}

var flagProfileListJSON bool

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available profiles",
	RunE:  runProfileList,
}

var profileShowCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show skills in a named profile",
	Args:  cobra.ExactArgs(1),
	RunE:  runProfileShow,
}

var profileInitCmd = &cobra.Command{
	Use:   "init <name>",
	Short: "Create a profile file in .grimoire/profiles/",
	Args:  cobra.ExactArgs(1),
	RunE:  runProfileInit,
}

func init() {
	profileListCmd.Flags().BoolVar(&flagProfileListJSON, "json", false, "output as JSON")
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileShowCmd)
	profileCmd.AddCommand(profileInitCmd)
}

type profileListEntry struct {
	Name   string `json:"name"`
	Source string `json:"source"` // "project" | "user" | "active"
	File   string `json:"file,omitempty"`
}

func listProfileEntries(cwd, grimoireHome string) []profileListEntry {
	searchDirs := []struct {
		label string
		dir   string
	}{
		{"project", filepath.Join(cwd, ".grimoire", "profiles")},
		{"user", filepath.Join(grimoireHome, "profiles")},
	}
	var entries []profileListEntry
	for _, sd := range searchDirs {
		dirEntries, err := os.ReadDir(sd.dir)
		if err != nil {
			continue
		}
		for _, e := range dirEntries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".toml") {
				continue
			}
			name := strings.TrimSuffix(e.Name(), ".toml")
			entries = append(entries, profileListEntry{
				Name:   name,
				Source: sd.label,
				File:   filepath.Join(sd.dir, e.Name()),
			})
		}
	}
	if r, err := settings.Load(cwd); err == nil {
		for _, name := range r.Core.Profiles {
			entries = append(entries, profileListEntry{Name: name, Source: "active"})
		}
	}
	return entries
}

func runProfileList(_ *cobra.Command, _ []string) error {
	cwd := getProjectDir()
	grimoireHome := skills.OfficialRegistryHome()
	entries := listProfileEntries(cwd, grimoireHome)
	r, settingsErr := settings.Load(cwd)

	if flagProfileListJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(entries)
	}

	if len(entries) == 0 {
		if settingsErr != nil {
			fmt.Printf("  %s no profiles found (settings load error: %v)\n", tui.IconWarn, settingsErr)
			fmt.Printf("  check: %s\n", tui.StyleDim.Render(settings.GlobalPath()))
		} else {
			fmt.Printf("  %s no profiles found\n", tui.IconWarn)
			fmt.Printf("  run: grimoire profile init <name>\n")
		}
		return nil
	}

	opts := resolveOpts(cwd)
	for _, e := range entries {
		if e.Source == "active" {
			continue
		}
		fmt.Printf("  %s %s  %s\n", tui.IconOK, e.Name, tui.StyleDim.Render("("+e.Source+")"))
		p, err := profiles.ResolveWithOptions(e.Name, cwd, opts)
		if err != nil || (len(p.Skills) == 0 && p.Source == "") {
			continue
		}
		for _, sk := range p.Skills {
			fmt.Printf("      %s %s\n", tui.StyleCyan.Render("→"), sk.Name)
		}
		if len(p.Skills) == 0 {
			fmt.Printf("      %s\n", tui.StyleDim.Render("(no installed skills match — AI applies semantically)"))
		}
	}
	if settingsErr == nil && len(r.Core.Profiles) > 0 {
		fmt.Println()
		fmt.Printf("  %s\n", tui.StyleDim.Render("active in [standards] profiles:"))
		for _, name := range r.Core.Profiles {
			p, err := profiles.ResolveWithOptions(name, cwd, opts)
			if err != nil {
				fmt.Printf("    %s\n", tui.StyleDim.Render(name+": (error resolving)"))
				continue
			}
			src := p.Source
			if home, e := os.UserHomeDir(); e == nil && src != "" {
				src = strings.Replace(src, home, "~", 1)
			}
			if src != "" {
				fmt.Printf("    %s  %s\n", tui.StyleBold.Render(name), tui.StyleDim.Render(src))
			} else {
				fmt.Printf("    %s\n", tui.StyleBold.Render(name))
			}
			for _, sk := range p.Skills {
				fmt.Printf("      %s %s\n", tui.StyleCyan.Render("→"), sk.Name)
			}
			if len(p.Skills) == 0 {
				fmt.Printf("      %s\n", tui.StyleDim.Render("(no installed skills match — AI applies semantically)"))
			}
		}
	}
	return nil
}

func runProfileShow(_ *cobra.Command, args []string) error {
	name := args[0]
	cwd := getProjectDir()

	p, err := profiles.ResolveWithOptions(name, cwd, resolveOpts(cwd))
	if err != nil {
		return fmt.Errorf("resolving profile %q: %w", name, err)
	}

	fmt.Println()
	fmt.Printf("  %s profile: %s\n", tui.StyleBold.Render(tui.IconOK), tui.StyleBold.Render(name))

	if p.Source == "" {
		fmt.Printf("  %s\n", tui.StyleDim.Render("(no profile file found and no tagged skills)"))
		return nil
	}

	src := p.Source
	if home, err := os.UserHomeDir(); err == nil {
		src = strings.Replace(src, home, "~", 1)
	}
	fmt.Printf("  source: %s\n", tui.StyleDim.Render(src))

	if p.Description != "" {
		fmt.Printf("  %s\n", p.Description)
	}

	if len(p.Extends) > 0 {
		fmt.Printf("  extends: %s\n", strings.Join(p.Extends, ", "))
	}
	if len(p.Tags) > 0 {
		fmt.Printf("  tags: %s\n", strings.Join(p.Tags, ", "))
	}
	if len(p.Exclude) > 0 {
		fmt.Printf("  exclude: %s\n", strings.Join(p.Exclude, ", "))
	}

	fmt.Println()
	if len(p.Skills) == 0 {
		fmt.Printf("  %s\n", tui.StyleDim.Render("(no skills defined)"))
	} else {
		for _, sk := range p.Skills {
			line := fmt.Sprintf("    %s %s", tui.StyleCyan.Render("→"), sk.Name)
			if sk.Priority != 0 && sk.Priority != 50 {
				line += tui.StyleDim.Render(fmt.Sprintf("  (priority %d)", sk.Priority))
			}
			fmt.Println(line)
		}
	}
	fmt.Println()
	return nil
}

func buildProfileTemplate(name string) string {
	return fmt.Sprintf(`name = "%s"
description = ""

# Inherit all skills from other profiles:
# extends = ["oop", "tdd"]

# Bulk-activate all skills matching these tags:
# tags = ["go", "backend"]

# Remove specific skills after all inclusions:
# exclude = ["apply-law-of-demeter"]

# Explicit skill additions (priority: lower = higher priority, default = 50):
# [[skills]]
# name = "apply-solid-principles"
# priority = 10
#
# [[skills]]
# name = "apply-kiss-principle"
`, name)
}

func runProfileInit(_ *cobra.Command, args []string) error {
	name := args[0]
	cwd := getProjectDir()
	dir := filepath.Join(cwd, ".grimoire", "profiles")
	path := filepath.Join(dir, name+".toml")

	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("profile file already exists: %s", path)
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating profiles dir: %w", err)
	}

	if err := os.WriteFile(path, []byte(buildProfileTemplate(name)), 0o644); err != nil {
		return fmt.Errorf("writing profile file: %w", err)
	}

	rel, _ := filepath.Rel(cwd, path)
	fmt.Printf("  %s created %s\n", tui.IconOK, rel)
	return nil
}
