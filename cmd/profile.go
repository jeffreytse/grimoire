package cmd

import (
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
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileShowCmd)
	profileCmd.AddCommand(profileInitCmd)
}

func runProfileList(_ *cobra.Command, _ []string) error {
	cwd, _ := os.Getwd()
	grimoireHome := skills.GrimoireHome()

	searchDirs := []struct {
		label string
		dir   string
	}{
		{"project", filepath.Join(cwd, ".grimoire", "profiles")},
		{"user", filepath.Join(grimoireHome, "profiles")},
	}

	printed := false
	for _, sd := range searchDirs {
		entries, err := os.ReadDir(sd.dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".toml") {
				continue
			}
			name := strings.TrimSuffix(e.Name(), ".toml")
			fmt.Printf("  %s %s  %s\n", tui.IconOK, name, tui.StyleDim.Render("("+sd.label+")"))
			printed = true
		}
	}

	// also show profiles active in settings but not yet listed
	r, err := settings.Load(cwd)
	if err == nil && len(r.Core.Profiles) > 0 {
		fmt.Println()
		fmt.Printf("  %s\n", tui.StyleDim.Render("active in [standards] profiles:"))
		for _, name := range r.Core.Profiles {
			fmt.Printf("    %s\n", name)
		}
		printed = true
	}

	if !printed {
		fmt.Printf("  %s no profiles found\n", tui.IconWarn)
		fmt.Printf("  run: grimoire profile init <name>\n")
	}
	return nil
}

func runProfileShow(_ *cobra.Command, args []string) error {
	name := args[0]
	cwd, _ := os.Getwd()

	p, err := profiles.ResolveWithOptions(name, cwd, profiles.ResolveOptions{
		Sources: skills.AllSkillsSources(),
	})
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

	fmt.Println()
	if len(p.Skills) == 0 {
		fmt.Printf("  %s\n", tui.StyleDim.Render("(no skills defined)"))
	} else {
		for _, sk := range p.Skills {
			fmt.Printf("    %s %s\n", tui.StyleCyan.Render("→"), sk.Name)
		}
	}
	fmt.Println()
	return nil
}

func runProfileInit(_ *cobra.Command, args []string) error {
	name := args[0]
	cwd, _ := os.Getwd()
	dir := filepath.Join(cwd, ".grimoire", "profiles")
	path := filepath.Join(dir, name+".toml")

	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("profile file already exists: %s", path)
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating profiles dir: %w", err)
	}

	content := fmt.Sprintf(`name = "%s"
description = ""

# Add skills to this profile:
# [[skills]]
# name = "apply-solid-principles"
#
# [[skills]]
# name = "apply-tdd"
`, name)

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing profile file: %w", err)
	}

	rel, _ := filepath.Rel(cwd, path)
	fmt.Printf("  %s created %s\n", tui.IconOK, rel)
	return nil
}
